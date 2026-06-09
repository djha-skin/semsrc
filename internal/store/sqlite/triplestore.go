package sqlite

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/djha-skin/semsrc/internal/store"
	_ "github.com/mattn/go-sqlite3"
)

// Config holds configuration for the SQLite triple store
type Config struct {
	Path        string // Path to SQLite database file
	JournalMode string // SQLite journal mode (WAL, DELETE, TRUNCATE, etc.)
	SyncMode    string // SQLite synchronous mode
	CacheSize   int    // SQLite cache size in pages
}

// TripleStore implements store.TripleStore using SQLite
type TripleStore struct {
	db *sql.DB
}

// New creates a new SQLite triple store
func New(config Config) (*TripleStore, error) {
	// Build connection string with pragmas
	pragmas := []string{
		"busy_timeout=5000",
		fmt.Sprintf("journal_mode=%s", config.JournalMode),
		fmt.Sprintf("synchronous=%s", config.SyncMode),
		fmt.Sprintf("cache_size=%d", config.CacheSize),
		"foreign_keys=ON",
	}

	connStr := fmt.Sprintf("file:%s?%s", config.Path, strings.Join(pragmas, "&"))
	
	db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	store := &TripleStore{
		db: db,
	}

	// Initialize schema
	if err := store.initializeSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return store, nil
}

// initializeSchema creates the necessary tables and indexes
func (s *TripleStore) initializeSchema() error {
	// Create tables
	schema := []string{
		// Resources table for normalized subject/predicate/object storage
		`CREATE TABLE IF NOT EXISTS resources (
			id INTEGER PRIMARY KEY,
			type TEXT NOT NULL CHECK(type IN ('iri', 'literal', 'blank')),
			value TEXT NOT NULL,
			datatype TEXT,
			language TEXT,
			UNIQUE(type, value, datatype, language)
		)`,

		// Triples table
		`CREATE TABLE IF NOT EXISTS triples (
			id INTEGER PRIMARY KEY,
			subject_id INTEGER NOT NULL,
			predicate_id INTEGER NOT NULL,
			object_id INTEGER NOT NULL,
			graph TEXT DEFAULT '',
			inferred BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(subject_id) REFERENCES resources(id),
			FOREIGN KEY(predicate_id) REFERENCES resources(id),
			FOREIGN KEY(object_id) REFERENCES resources(id),
			UNIQUE(subject_id, predicate_id, object_id, graph)
		)`,

		// Named graphs table
		`CREATE TABLE IF NOT EXISTS graphs (
			name TEXT PRIMARY KEY,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	// Create indexes
	indexes := []string{
		// Resource lookups
		`CREATE INDEX IF NOT EXISTS idx_resources_type_value ON resources(type, value)`,
		`CREATE INDEX IF NOT EXISTS idx_resources_value ON resources(value)`,
		
		// Triple lookups
		`CREATE INDEX IF NOT EXISTS idx_triples_subject ON triples(subject_id)`,
		`CREATE INDEX IF NOT EXISTS idx_triples_predicate ON triples(predicate_id)`,
		`CREATE INDEX IF NOT EXISTS idx_triples_object ON triples(object_id)`,
		`CREATE INDEX IF NOT EXISTS idx_triples_graph ON triples(graph)`,
		`CREATE INDEX IF NOT EXISTS idx_triples_sp ON triples(subject_id, predicate_id)`,
		`CREATE INDEX IF NOT EXISTS idx_triples_po ON triples(predicate_id, object_id)`,
		`CREATE INDEX IF NOT EXISTS idx_triples_so ON triples(subject_id, object_id)`,
		
		// Graph lookups
		`CREATE INDEX IF NOT EXISTS idx_graphs_created ON graphs(created_at)`,
	}

	// Execute schema creation
	for _, stmt := range schema {
		if _, err := s.db.Exec(stmt); err != nil {
			return fmt.Errorf("failed to execute schema statement: %w\nStatement: %s", err, stmt)
		}
	}

	// Execute index creation
	for _, stmt := range indexes {
		if _, err := s.db.Exec(stmt); err != nil {
			return fmt.Errorf("failed to create index: %w\nStatement: %s", err, stmt)
		}
	}

	return nil
}

// Add adds a triple to the store
func (s *TripleStore) Add(triple store.Triple) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Will be ignored if committed

	// Get or create resource IDs
	subjectID, err := s.getOrCreateResource(tx, triple.Subject)
	if err != nil {
		return fmt.Errorf("failed to get/create subject: %w", err)
	}

	predicateID, err := s.getOrCreateResource(tx, triple.Predicate)
	if err != nil {
		return fmt.Errorf("failed to get/create predicate: %w", err)
	}

	objectID, err := s.getOrCreateResource(tx, triple.Object)
	if err != nil {
		return fmt.Errorf("failed to get/create object: %w", err)
	}

	// Insert triple
	query := `INSERT INTO triples (subject_id, predicate_id, object_id, graph) 
	          VALUES (?, ?, ?, '') 
	          ON CONFLICT(subject_id, predicate_id, object_id, graph) DO NOTHING`
	
	result, err := tx.Exec(query, subjectID, predicateID, objectID)
	if err != nil {
		return fmt.Errorf("failed to insert triple: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return store.ErrTripleExists
	}

	return tx.Commit()
}

// Remove removes a triple from the store
func (s *TripleStore) Remove(triple store.Triple) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get resource IDs
	subjectID, err := s.getResourceID(tx, triple.Subject)
	if err != nil {
		return fmt.Errorf("failed to get subject ID: %w", err)
	}
	if subjectID == 0 {
		return store.ErrTripleNotFound
	}

	predicateID, err := s.getResourceID(tx, triple.Predicate)
	if err != nil {
		return fmt.Errorf("failed to get predicate ID: %w", err)
	}
	if predicateID == 0 {
		return store.ErrTripleNotFound
	}

	objectID, err := s.getResourceID(tx, triple.Object)
	if err != nil {
		return fmt.Errorf("failed to get object ID: %w", err)
	}
	if objectID == 0 {
		return store.ErrTripleNotFound
	}

	// Delete triple
	query := `DELETE FROM triples WHERE subject_id = ? AND predicate_id = ? AND object_id = ? AND graph = ''`
	result, err := tx.Exec(query, subjectID, predicateID, objectID)
	if err != nil {
		return fmt.Errorf("failed to delete triple: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return store.ErrTripleNotFound
	}

	// Clean up orphaned resources (optional, could be done periodically)
	// s.cleanupOrphanedResources(tx)

	return tx.Commit()
}

// Query executes a SPARQL query
func (s *TripleStore) Query(sparql string) (store.QueryResult, error) {
	// Note: Full SPARQL implementation is complex. This is a simplified version.
	// For production, we'd need a proper SPARQL parser and query planner.
	
	// Simple SELECT * WHERE { ?s ?p ?o } implementation for now
	if strings.Contains(sparql, "SELECT") && strings.Contains(sparql, "WHERE") {
		return s.executeBasicSelect(sparql)
	}
	
	// ASK query
	if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(sparql)), "ASK") {
		return s.executeAskQuery(sparql)
	}
	
	return store.QueryResult{}, store.ErrInvalidSPARQL
}

// executeBasicSelect executes a basic SELECT query
func (s *TripleStore) executeBasicSelect(sparql string) (store.QueryResult, error) {
	// Parse SELECT clause to get variables
	// For now, just handle SELECT * WHERE { ?s ?p ?o }
	
	query := `
		SELECT 
			s.value as subject,
			s.type as subject_type,
			s.datatype as subject_datatype,
			s.language as subject_language,
			p.value as predicate,
			p.type as predicate_type,
			p.datatype as predicate_datatype,
			p.language as predicate_language,
			o.value as object,
			o.type as object_type,
			o.datatype as object_datatype,
			o.language as object_language
		FROM triples t
		JOIN resources s ON t.subject_id = s.id
		JOIN resources p ON t.predicate_id = p.id
		JOIN resources o ON t.object_id = o.id
		WHERE t.graph = ''
		LIMIT 1000
	`
	
	rows, err := s.db.Query(query)
	if err != nil {
		return store.QueryResult{}, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()
	
	var result store.QueryResult
	result.Variables = []string{"s", "p", "o"}
	
	for rows.Next() {
		var sVal, sType, sDatatype, sLang sql.NullString
		var pVal, pType, pDatatype, pLang sql.NullString
		var oVal, oType, oDatatype, oLang sql.NullString
		
		if err := rows.Scan(
			&sVal, &sType, &sDatatype, &sLang,
			&pVal, &pType, &pDatatype, &pLang,
			&oVal, &oType, &oDatatype, &oLang,
		); err != nil {
			return store.QueryResult{}, fmt.Errorf("failed to scan row: %w", err)
		}
		
		row := make(map[string]store.Resource)
		
		// Build subject resource
		row["s"] = s.buildResource(sType.String, sVal.String, sDatatype.String, sLang.String)
		
		// Build predicate resource
		row["p"] = s.buildResource(pType.String, pVal.String, pDatatype.String, pLang.String)
		
		// Build object resource
		row["o"] = s.buildResource(oType.String, oVal.String, oDatatype.String, oLang.String)
		
		result.Rows = append(result.Rows, row)
	}
	
	return result, nil
}

// executeAskQuery executes an ASK query
func (s *TripleStore) executeAskQuery(sparql string) (store.QueryResult, error) {
	// Simple ASK implementation: just check if any triples exist
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM triples WHERE graph = ''").Scan(&count)
	if err != nil {
		return store.QueryResult{}, fmt.Errorf("failed to count triples: %w", err)
	}
	
	exists := count > 0
	return store.QueryResult{
		Boolean: &exists,
	}, nil
}

// QueryOne executes a SPARQL query expecting exactly one result
func (s *TripleStore) QueryOne(sparql string) (map[string]store.Resource, error) {
	result, err := s.Query(sparql)
	if err != nil {
		return nil, err
	}
	
	if len(result.Rows) == 0 {
		return nil, store.ErrTripleNotFound
	}
	
	if len(result.Rows) > 1 {
		return nil, fmt.Errorf("expected one result, got %d", len(result.Rows))
	}
	
	return result.Rows[0], nil
}

// QueryExists executes a SPARQL ASK query
func (s *TripleStore) QueryExists(sparql string) (bool, error) {
	result, err := s.Query(sparql)
	if err != nil {
		return false, err
	}
	
	if result.Boolean == nil {
		return false, fmt.Errorf("not an ASK query")
	}
	
	return *result.Boolean, nil
}

// Begin starts a transaction
func (s *TripleStore) Begin() (store.Transaction, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	
	return &transaction{
		tx: tx,
		s:  s,
	}, nil
}

// Stats returns statistics about the store
func (s *TripleStore) Stats() (store.TripleStoreStats, error) {
	var stats store.TripleStoreStats
	
	// Get triple count
	err := s.db.QueryRow("SELECT COUNT(*) FROM triples WHERE graph = ''").Scan(&stats.TripleCount)
	if err != nil {
		return stats, fmt.Errorf("failed to count triples: %w", err)
	}
	
	// Get distinct counts
	err = s.db.QueryRow(`
		SELECT 
			COUNT(DISTINCT subject_id),
			COUNT(DISTINCT predicate_id),
			COUNT(DISTINCT object_id)
		FROM triples
		WHERE graph = ''
	`).Scan(&stats.DistinctSubjects, &stats.DistinctPredicates, &stats.DistinctObjects)
	if err != nil {
		return stats, fmt.Errorf("failed to count distinct resources: %w", err)
	}
	
	// Get last update time
	var lastUpdateStr string
	err = s.db.QueryRow("SELECT MAX(created_at) FROM triples").Scan(&lastUpdateStr)
	if err == nil && lastUpdateStr != "" {
		stats.LastUpdate, _ = time.Parse("2006-01-02 15:04:05", lastUpdateStr)
	}
	
	return stats, nil
}

// Close releases resources
func (s *TripleStore) Close() error {
	return s.db.Close()
}

// Helper methods

// getOrCreateResource gets or creates a resource and returns its ID
func (s *TripleStore) getOrCreateResource(tx *sql.Tx, resource store.Resource) (int64, error) {
	// First try to get existing ID
	id, err := s.getResourceID(tx, resource)
	if err != nil {
		return 0, err
	}
	if id > 0 {
		return id, nil
	}
	
	// Create new resource
	var rType, value, datatype, language string
	
	switch r := resource.(type) {
	case store.IRI:
		rType = "iri"
		value = string(r)
	case store.Literal:
		rType = "literal"
		value = r.Value
		datatype = string(r.Datatype)
		language = r.Language
	case store.BlankNode:
		rType = "blank"
		value = r.ID
	default:
		return 0, fmt.Errorf("unknown resource type: %T", resource)
	}
	
	result, err := tx.Exec(`
		INSERT INTO resources (type, value, datatype, language)
		VALUES (?, ?, ?, ?)
	`, rType, value, datatype, language)
	if err != nil {
		return 0, fmt.Errorf("failed to insert resource: %w", err)
	}
	
	return result.LastInsertId()
}

// getResourceID gets the ID of a resource, returns 0 if not found
func (s *TripleStore) getResourceID(tx *sql.Tx, resource store.Resource) (int64, error) {
	var query string
	var args []interface{}
	
	switch r := resource.(type) {
	case store.IRI:
		query = "SELECT id FROM resources WHERE type = 'iri' AND value = ?"
		args = []interface{}{string(r)}
	case store.Literal:
		query = "SELECT id FROM resources WHERE type = 'literal' AND value = ? AND datatype = ? AND language = ?"
		args = []interface{}{r.Value, string(r.Datatype), r.Language}
	case store.BlankNode:
		query = "SELECT id FROM resources WHERE type = 'blank' AND value = ?"
		args = []interface{}{r.ID}
	default:
		return 0, fmt.Errorf("unknown resource type: %T", resource)
	}
	
	var id int64
	err := tx.QueryRow(query, args...).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to query resource: %w", err)
	}
	
	return id, nil
}

// buildResource builds a store.Resource from database fields
func (s *TripleStore) buildResource(rType, value, datatype, language string) store.Resource {
	switch rType {
	case "iri":
		return store.IRI(value)
	case "literal":
		lit := store.Literal{
			Value: value,
		}
		if datatype != "" {
			lit.Datatype = store.IRI(datatype)
		}
		if language != "" {
			lit.Language = language
		}
		return lit
	case "blank":
		return store.BlankNode{ID: value}
	default:
		// Fallback to literal
		return store.Literal{Value: value}
	}
}

// transaction implements store.Transaction for SQLite
type transaction struct {
	tx *sql.Tx
	s  *TripleStore
}

func (t *transaction) Add(triple store.Triple) error {
	subjectID, err := t.s.getOrCreateResource(t.tx, triple.Subject)
	if err != nil {
		return err
	}
	
	predicateID, err := t.s.getOrCreateResource(t.tx, triple.Predicate)
	if err != nil {
		return err
	}
	
	objectID, err := t.s.getOrCreateResource(t.tx, triple.Object)
	if err != nil {
		return err
	}
	
	_, err = t.tx.Exec(`
		INSERT INTO triples (subject_id, predicate_id, object_id, graph)
		VALUES (?, ?, ?, '')
		ON CONFLICT(subject_id, predicate_id, object_id, graph) DO NOTHING
	`, subjectID, predicateID, objectID)
	
	if err != nil {
		return fmt.Errorf("failed to add triple: %w", err)
	}
	
	return nil
}

func (t *transaction) Remove(triple store.Triple) error {
	subjectID, err := t.s.getResourceID(t.tx, triple.Subject)
	if err != nil || subjectID == 0 {
		return store.ErrTripleNotFound
	}
	
	predicateID, err := t.s.getResourceID(t.tx, triple.Predicate)
	if err != nil || predicateID == 0 {
		return store.ErrTripleNotFound
	}
	
	objectID, err := t.s.getResourceID(t.tx, triple.Object)
	if err != nil || objectID == 0 {
		return store.ErrTripleNotFound
	}
	
	_, err = t.tx.Exec(`
		DELETE FROM triples 
		WHERE subject_id = ? AND predicate_id = ? AND object_id = ? AND graph = ''
	`, subjectID, predicateID, objectID)
	
	if err != nil {
		return fmt.Errorf("failed to remove triple: %w", err)
	}
	
	return nil
}

func (t *transaction) Query(sparql string) (store.QueryResult, error) {
	// Note: Transaction queries should use the transaction context
	// For simplicity, we reuse the main store's query method
	return t.s.Query(sparql)
}

func (t *transaction) Commit() error {
	return t.tx.Commit()
}

func (t *transaction) Rollback() error {
	return t.tx.Rollback()
}