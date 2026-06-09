# Semsrc System Architecture

## Architectural Principles

### 1. Separation of Concerns
- **Object Store**: Content-addressable blob storage
- **Triple Store**: Semantic relationship storage
- **Business Logic**: Version control operations
- **Interface Layers**: CLI, API, Protocols

### 2. Pluggable Components

- Storage backends interchangeable via interfaces
- Triple store implementations swappable
- Protocol handlers pluggable
- Authentication/authorization providers

### 3. Semantic First

- All operations generate semantic metadata
- Data structures designed for AI consumption
- Standard ontologies preferred over custom solutions
- Inference capabilities built-in

### 4. Progressive Enhancement

- Start with Git compatibility
- Add semantic features incrementally
- Opt-in advanced features
- Graceful degradation when semantics unavailable

## High-Level Architecture

**Note**: For detailed bounded context definitions, see [bounded-contexts.md](bounded-contexts.md)

```
┌─────────────────────────────────────────────────────────────┐
│                     Application Layer                       │
│              (Bounded Context: Presentation)                │
├───────────────┬───────────────┬─────────────────────────────┤
│    Git CLI    │  Semsrc CLI   │      HTTP/API/SSH           │
│  Compatibility│  Native Commands│     Protocols             │
└───────────────┴───────────────┴───────────────┬─────────────┘
                                                 │
                    Interface Contracts                  │
┌─────────────────────────────────────────────────────────────┐
│              Business Logic Layer                           │
│         (Bounded Context: Core Domain Logic)               │
├───────────────┬───────────────┬─────────────────────────────┤
│  Version Control│  Semantic     │   Query & Inference       │
│  Operations    │  Annotation   │   Engine                   │
└───────────────┴───────────────┴───────────────┬─────────────┘
                                                 │
                    Interface Contracts                  │
┌─────────────────────────────────────────────────────────────┐
│                  Data Access Layer                          │
│          (Bounded Context: Data Persistence)               │
├───────────────┬───────────────┬─────────────────────────────┤
│  Object Store │  Triple Store │      Index & Cache         │
│   Interface   │   Interface   │                             │
└───────────────┴───────────────┴───────────────┬─────────────┘
                                                 │
                    Interface Contracts                  │
┌─────────────────────────────────────────────────────────────┐
│                 Storage Implementation Layer                │
│            (Bounded Context: Storage Backends)             │
├───────────────┬───────────────┬─────────────────────────────┤
│  Filesystem   │      S3       │     SQLite/Embedded        │
│  Object Store │  Object Store │     Triple Store           │
└───────────────┴───────────────┴─────────────────────────────┘
```

### Key Architectural Principles

1. **Strict Layering**: Each layer only communicates with adjacent layers
2. **Interface Contracts**: Dependencies flow downward via interfaces only
3. **Dependency Injection**: Storage implementations are injected, not hardcoded
4. **Testability**: Each bounded context can be tested independently
5. **Evolutionary Design**: Easy to add new storage backends or features

## Component Details

### Interface Layer

#### Git CLI Compatibility
- Binary named `semsrc` (or `git` alias)
- Subset of Git commands with same flags/behavior
- Translates Git operations to semantic operations
- Produces same stdout/stderr format where possible
- **Pattern**: Adapter Pattern (Git CLI → Semsrc operations)

### Git Pack File Compatibility

#### Import Strategy
- **Git Pack Format**: Support reading Git pack files directly
- **Object Mapping**: Map Git SHA-1 objects to SHA-256 hashes
- **Delta Compression**: Handle Git's delta compression efficiently
- **Incremental Import**: Import repositories in phases for large repos

**Implementation**:
```go
type GitPackImporter struct {
    objectStore ObjectStore
    tripleStore TripleStore
}

func (i *GitPackImporter) ImportPack(packPath string) error {
    // 1. Read pack file index
    index, err := readPackIndex(packPath)
    if err != nil {
        return err
    }
    
    // 2. Iterate through objects
    for _, obj := range index.Objects {
        // 3. Resolve deltas
        data, err := i.resolveDelta(obj)
        if err != nil {
            return err
        }
        
        // 4. Store in Object Store (SHA-256)
        sha256Hash := computeSHA256(data)
        err = i.objectStore.Put(sha256Hash, bytes.NewReader(data))
        if err != nil {
            return err
        }
        
        // 5. Convert to RDF triples
        triples := i.convertGitObject(obj, sha256Hash)
        err = i.tripleStore.AddBatch(triples)
        if err != nil {
            return err
        }
    }
    
    return nil
}

func (i *GitPackImporter) resolveDelta(obj GitObject) ([]byte, error) {
    if obj.Type != "delta" {
        return obj.Data, nil
    }
    
    // Resolve delta against base object
    base, err := i.resolveDelta(obj.BaseObject)
    if err != nil {
        return nil, err
    }
    
    return applyDelta(base, obj.DeltaData), nil
}
```

#### Export Strategy
- **Git Pack Format**: Generate Git-compatible pack files
- **Object Mapping**: Map SHA-256 back to Git SHA-1 if needed
- **Delta Compression**: Apply Git's delta compression for efficiency
- **Incremental Export**: Export changes since last sync

**Implementation**:
```go
type GitPackExporter struct {
    objectStore ObjectStore
}

func (e *GitPackExporter) ExportPack(repoID string, outputPath string) error {
    // 1. List all objects in repository
    objects, err := e.listObjects(repoID)
    if err != nil {
        return err
    }
    
    // 2. Generate pack file
    packFile, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer packFile.Close()
    
    // 3. Write pack header
    _, err = packFile.Write([]byte("PACK"))
    if err != nil {
        return err
    }
    
    // 4. Write objects with Git format
    for _, obj := range objects {
        data, err := e.objectStore.Get(obj.Hash)
        if err != nil {
            return err
        }
        
        // Convert to Git object format and write
        gitObj := e.convertToGitObject(obj, data)
        _, err = packFile.Write(gitObj)
        if err != nil {
            return err
        }
    }
    
    // 5. Write index file
    return e.writeIndexFile(outputPath + ".idx")
}
```

### Git Hook Compatibility

#### Hook Translation Layer
- **Pre-commit**: Run semantic validation before commit
- **Post-commit**: Add semantic annotations after commit
- **Pre-push**: Validate semantic integrity before push
- **Post-merge**: Update semantic relationships after merge

**Implementation**:
```go
type GitHookAdapter struct {
    semsrc VersionControl
}

func (h *GitHookAdapter) PreCommit() error {
    // Run SHACL validation on staged triples
    validationErrors := h.semsrc.ValidateStaged()
    if len(validationErrors) > 0 {
        return fmt.Errorf("validation failed: %v", validationErrors)
    }
    return nil
}

func (h *GitHookAdapter) PostCommit(commitHash string) error {
    // Add semantic annotations to commit
    commit, err := h.semsrc.GetCommit(commitHash)
    if err != nil {
        return err
    }
    
    // Extract and link tickets
    ticketIDs := extractTicketIDs(commit.Message)
    for _, ticketID := range ticketIDs {
        h.semsrc.LinkCommitToTicket(commitHash, ticketID)
    }
    
    return nil
}

func (h *GitHookAdapter) PrePush() error {
    // Validate semantic integrity before push
    return h.semsrc.ValidateIntegrity()
}

func (h *GitHookAdapter) PostMerge(mergeCommit string) error {
    // Update semantic relationships after merge
    return h.semsrc.UpdateMergeRelationships(mergeCommit)
}

// Concrete Git Hook Translation Example
type GitPreCommitHook struct {
    hookDir string
}

func (g *GitPreCommitHook) Install() error {
    hookPath := filepath.Join(g.hookDir, "pre-commit")
    
    // Shell script that calls semsrc validation
    hookScript := `#!/bin/sh
# Generated Semsrc Hook
semsrc validate-staged
if [ $? -ne 0 ]; then
    echo "Semantic validation failed. Commit aborted."
    exit 1
fi
`
    
    return os.WriteFile(hookPath, []byte(hookScript), 0755)
}

// Example: Pre-commit hook that extracts and links tickets
type TicketLinkingHook struct {
    semsrc VersionControl
}

func (t *TicketLinkingHook) Execute(commitMsg string) error {
    // Extract ticket IDs from commit message
    // Pattern: "Fixes PROJ-123: Description"
    re := regexp.MustCompile(`PROJ-\d+`)
    ticketIDs := re.FindAllString(commitMsg, -1)
    
    // For each ticket, create semantic relationship
    for _, ticketID := range ticketIDs {
        err := t.semsrc.CreateTicketRelationship(ticketID)
        if err != nil {
            return fmt.Errorf("failed to link ticket %s: %w", ticketID, err)
        }
    }
    
    return nil
}
```

#### Hook Installation
- **Automatic**: Install hooks when repository is initialized
- **Optional**: Allow users to disable hooks if needed
- **Custom**: Support custom hooks alongside semantic hooks

**Installation Example**:
```go
func (r *Repository) InitializeHooks() error {
    hooks := []GitHook{
        &GitPreCommitHook{hookDir: r.hookDir},
        &TicketLinkingHook{semsrc: r.semsrc},
    }
    
    for _, hook := range hooks {
        if err := hook.Install(); err != nil {
            return err
        }
    }
    
    return nil
}
```

### Git Ref Handling

#### Reference Types
- **Branches**: Stored as semantic properties
- **Tags**: Stored as semantic entities with versioning
- **HEAD**: Stored as current branch reference

**Implementation**:
```go
type GitRefHandler struct {
    tripleStore TripleStore
}

func (h *GitRefHandler) UpdateBranch(branchName string, commitHash string) error {
    // Store branch as semantic triple
    triple := Triple{
        Subject:   fmt.Sprintf("semsrc:branch/%s", branchName),
        Predicate: "semsrc:pointsTo",
        Object:    commitHash,
    }
    return h.tripleStore.Add(triple)
}

func (h *GitRefHandler) GetBranchCommit(branchName string) (string, error) {
    // Query branch HEAD
    sparql := fmt.Sprintf(
        "SELECT ?commit WHERE { semsrc:branch/%s semsrc:pointsTo ?commit }",
        branchName,
    )
    result, err := h.tripleStore.Query(sparql)
    if err != nil {
        return "", err
    }
    
    if len(result) == 0 {
        return "", fmt.Errorf("branch not found: %s", branchName)
    }
    
    return result[0].Object, nil
}
```

#### Semsrc Native CLI
- Extended commands for semantic features
- `semsrc semantic-query` - SPARQL interface
- `semsrc annotate` - Add semantic annotations
- `semsrc provenance` - Show provenance chains
- `semsrc inference` - Run OWL reasoning

#### Protocol Handlers
- **Smart HTTP**: Git-compatible push/pull
- **SSH**: Secure shell access
- **gRPC**: High-performance API
- **WebSocket**: Real-time notifications
- **GraphQL**: Flexible querying
- **Pattern**: Strategy Pattern (protocol selection via config)

### Semantic API (HATEOAS Implementation)

#### RESTful API Design
- **Resource-Oriented**: Everything is a resource with URI
- **Hypermedia Controls**: Include links for state transitions
- **Content Negotiation**: Support multiple serialization formats
- **Semantic Metadata**: Include RDF annotations in responses

**API Endpoints**:
```
GET  /api/repos/{repo-id}              → Repository resource
GET  /api/repos/{repo-id}/commits      → Collection of commits
GET  /api/repos/{repo-id}/commits/{id} → Specific commit
GET  /api/repos/{repo-id}/query        → SPARQL endpoint
POST /api/repos/{repo-id}/commits      → Create new commit
GET  /api/repos/{repo-id}/ontology     → Repository ontology
```

**HATEOAS Response Example**:
```json
{
  "@context": "http://semsrc.dev/api/v1/context",
  "@id": "http://api.semsrc.dev/repos/my-project/commits/abc123",
  "@type": "semsrc:Commit",
  "message": "Fix login bug",
  "author": "http://api.semsrc.dev/users/alice",
  "timestamp": "2024-01-15T10:30:00Z",
  "_links": {
    "self": {"href": "/repos/my-project/commits/abc123"},
    "parent": {"href": "/repos/my-project/commits/def456"},
    "tree": {"href": "/repos/my-project/trees/xyz789"},
    "children": {"href": "/repos/my-project/commits?parent=abc123"},
    "diff": {"href": "/repos/my-project/commits/abc123/diff"},
    "semantic": {"href": "/repos/my-project/commits/abc123/semantic"}
  }
}
```

#### Content Negotiation
- **Turtle (.ttl)**: For semantic web tools
- **JSON-LD**: For web APIs
- **N-Triples**: For bulk data transfer
- **HTML**: For browser viewing

**Implementation**:
```go
type ContentNegotiator struct {
    formats map[string]Formatter
}

func (c *ContentNegotiator) Format(resource interface{}, accept string) ([]byte, error) {
    formatter, ok := c.formats[accept]
    if !ok {
        // Default to JSON-LD
        formatter = c.formats["application/ld+json"]
    }
    
    return formatter.Format(resource)
}

// Usage in HTTP handler
func handleGetCommit(w http.ResponseWriter, r *http.Request) {
    commit := getCommit(r)
    
    negotiator := &ContentNegotiator{
        formats: map[string]Formatter{
            "text/turtle":          &TurtleFormatter{},
            "application/ld+json":  &JSONLDFormatter{},
            "application/n-triples": &NTriplesFormatter{},
            "text/html":            &HTMLFormatter{},
        },
    }
    
    data, err := negotiator.Format(commit, r.Header.Get("Accept"))
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    w.Write(data)
}
```

#### HATEOAS Navigation
**Purpose**: Clients discover API capabilities through hypermedia links

**Implementation**:
```go
type HALResource struct {
    Data  map[string]interface{} `json:"_embedded,omitempty"`
    Links map[string]Link        `json:"_links"`
}

type Link struct {
    Href  string `json:"href"`
    Title string `json:"title,omitempty"`
}

func createCommitResource(commit *Commit) HALResource {
    return HALResource{
        Data: map[string]interface{}{
            "id":      commit.ID,
            "message": commit.Message,
            "author":  commit.Author,
        },
        Links: map[string]Link{
            "self":   {Href: fmt.Sprintf("/commits/%s", commit.ID)},
            "parent": {Href: fmt.Sprintf("/commits/%s", commit.ParentID)},
            "tree":   {Href: fmt.Sprintf("/trees/%s", commit.TreeID)},
        },
    }
}
```

#### Semantic Metadata in API Responses
- **RDFa**: Embed semantic triples in HTML responses
- **JSON-LD @context**: Define semantic vocabulary
- **Link headers**: Include semantic relationships in HTTP headers

**Example Link Header**:
```
Link: </commits/abc123>; rel="parent",
      </trees/xyz789>; rel="tree",
      </tickets/42>; rel="semsrc:fixedBy"
```

### Business Logic Layer

#### Version Control Operations
```go
// Core operations matching Git
type VersionControl interface {
    InitRepo(path string) error
    Add(paths []string) error
    Commit(message string) error
    Checkout(ref string) error
    Merge(source string) error
    Log(options LogOptions) ([]Commit, error)
    Diff(a, b string) ([]FileDiff, error)
}

// Repository Pattern Implementation
type CommitRepository interface {
    Save(commit *Commit) error
    FindByID(id string) (*Commit, error)
    FindByRange(start, end int) ([]*Commit, error)
    FindByBranch(branch string) ([]*Commit, error)
    FindByTicket(ticketID string) ([]*Commit, error)
}

// Factory Pattern for Repository Creation
type RepositoryFactory interface {
    CreateCommitRepository(store TripleStore) CommitRepository
    CreateBlobRepository(store ObjectStore) BlobRepository
    CreateTreeRepository(store ObjectStore) TreeRepository
}
```

#### Semantic Annotation Engine
- Automatically adds PROV-O relationships
- Extracts ticket references from commits
- Links builds to commits
- Tracks file renames/moves semantically
- Generates RDF from operational metadata
- **Pattern**: Observer Pattern - Reacts to version control events

#### Query & Inference Engine
- SPARQL endpoint over triple store
- OWL reasoning for inferred triples
- SHACL validation for data quality
- Full-text search across blobs and metadata
- Graph traversal for relationship chains
- **Pattern**: Strategy Pattern - Different query optimization strategies

### Design Pattern Implementations

#### 1. Repository Pattern (Data Access)

**Purpose**: Abstract data access and provide a collection-like interface

**Implementation**:
```go
// Generic repository interface
type Repository[T any] interface {
    Save(entity *T) error
    FindByID(id string) (*T, error)
    FindAll() ([]*T, error)
    Delete(id string) error
}

// Commit-specific repository
type CommitRepository struct {
    tripleStore TripleStore
    objectStore ObjectStore
}

func (r *CommitRepository) Save(commit *Commit) error {
    // 1. Store commit object in Object Store
    // 2. Convert to RDF triples
    // 3. Store triples in Triple Store
    // 4. Add provenance relationships
    
    objectHash, err := r.objectStore.Put(commit.Data)
    if err != nil {
        return err
    }
    
    triples := commit.ToRDF()
    triples = append(triples, prov.OwasGeneratedBy(commit.ID, objectHash))
    
    return r.tripleStore.AddBatch(triples)
}

// Usage in Business Logic
func (vc *VersionControl) Commit(message string) error {
    repo := vc.repoFactory.CreateCommitRepository()
    commit := NewCommit(message)
    return repo.Save(commit)
}
```

**Benefits**:
- Decouples business logic from data access details
- Enables easy swapping of storage implementations
- Simplifies testing with mock repositories

#### 2. Strategy Pattern (Storage Backend Selection)

**Purpose**: Select storage backend based on configuration

**Implementation**:
```go
// Strategy interface
type ObjectStoreStrategy interface {
    Create(config Config) (ObjectStore, error)
    Supports(config Config) bool
}

// Concrete strategies
type FilesystemStrategy struct{}

func (s *FilesystemStrategy) Create(config Config) (ObjectStore, error) {
    return NewFilesystemStore(config.Path, config.Compression), nil
}

func (s *FilesystemStrategy) Supports(config Config) bool {
    return config.Type == "filesystem"
}

type S3Strategy struct{}

func (s *S3Strategy) Create(config Config) (ObjectStore, error) {
    return NewS3Store(config.Bucket, config.Region), nil
}

func (s *S3Strategy) Supports(config Config) bool {
    return config.Type == "s3"
}

// Strategy registry
type ObjectStoreRegistry struct {
    strategies []ObjectStoreStrategy
}

func (r *ObjectStoreRegistry) Create(config Config) (ObjectStore, error) {
    for _, strategy := range r.strategies {
        if strategy.Supports(config) {
            return strategy.Create(config)
        }
    }
    return nil, fmt.Errorf("no strategy found for config: %v", config)
}

// Usage
registry := &ObjectStoreRegistry{
    strategies: []ObjectStoreStrategy{
        &FilesystemStrategy{},
        &S3Strategy{},
    },
}

store, err := registry.Create(config)
```

**Benefits**:
- Easy to add new storage backends
- Configuration-driven backend selection
- Clean separation of backend implementations

#### 3. Factory Pattern (Object Creation)

**Purpose**: Centralize object creation with dependency injection

**Implementation**:
```go
type ObjectStoreFactory struct {
    registry ObjectStoreRegistry
}

func (f *ObjectStoreFactory) Create(config Config) (ObjectStore, error) {
    store, err := f.registry.Create(config)
    if err != nil {
        return nil, err
    }
    
    // Add caching layer if configured
    if config.CacheEnabled {
        store = NewCachedObjectStore(store, config.CacheSize)
    }
    
    // Add compression if configured
    if config.Compression != "none" {
        store = NewCompressedObjectStore(store, config.Compression)
    }
    
    return store, nil
}

// Triple Store Factory
type TripleStoreFactory struct {
    inferenceEnabled bool
}

func (f *TripleStoreFactory) Create(config Config) (TripleStore, error) {
    var store TripleStore
    
    switch config.Type {
    case "sqlite":
        store = NewSQLiteTripleStore(config.Path)
    case "postgres":
        store = NewPostgresTripleStore(config.ConnectionString)
    default:
        return nil, fmt.Errorf("unsupported triple store type: %s", config.Type)
    }
    
    // Wrap with inference if enabled
    if f.inferenceEnabled {
        store = NewInferenceTripleStore(store)
    }
    
    return store, nil
}
```

**Benefits**:
- Centralized object creation logic
- Easy to add cross-cutting concerns (caching, compression)
- Simplifies dependency injection

#### 4. Adapter Pattern (Git CLI Compatibility)

**Purpose**: Translate Git CLI commands to Semsrc operations

**Implementation**:
```go
type GitAdapter struct {
    semsrc VersionControl
}

func (g *GitAdapter) Exec(args []string) (string, error) {
    switch args[0] {
    case "commit":
        return g.handleCommit(args[1:])
    case "add":
        return g.handleAdd(args[1:])
    case "log":
        return g.handleLog(args[1:])
    default:
        return "", fmt.Errorf("unknown command: %s", args[0])
    }
}

func (g *GitAdapter) handleCommit(args []string) (string, error) {
    // Parse Git-style arguments
    message := extractMessage(args)
    
    // Execute Semsrc operation
    commit, err := g.semsrc.Commit(message)
    if err != nil {
        return "", err
    }
    
    // Return Git-compatible output
    return fmt.Sprintf("[%s] %s", commit.Hash, message), nil
}
```

**Benefits**:
- Git compatibility without modifying core logic
- Clean separation between Git interface and Semsrc operations
- Easy to maintain and extend

#### 5. Observer Pattern (Semantic Annotation Engine)

**Purpose**: React to version control events and add semantic annotations

**Implementation**:
```go
type CommitObserver interface {
    OnCommit(commit *Commit) error
}

type SemanticAnnotationEngine struct {
    observers []CommitObserver
}

func (e *SemanticAnnotationEngine) Register(observer CommitObserver) {
    e.observers = append(e.observers, observer)
}

func (e *SemanticAnnotationEngine) ProcessCommit(commit *Commit) error {
    for _, observer := range e.observers {
        if err := observer.OnCommit(commit); err != nil {
            return err
        }
    }
    return nil
}

// Concrete observers
type TicketExtractor struct{}

func (t *TicketExtractor) OnCommit(commit *Commit) error {
    // Extract ticket IDs from commit message
    ticketIDs := extractTicketIDs(commit.Message)
    
    // Add semantic relationships
    for _, ticketID := range ticketIDs {
        // Add triple: commit semsrc-ticket:fixedBy ticketID
    }
    
    return nil
}

type ProvenanceGenerator struct{}

func (p *ProvenanceGenerator) OnCommit(commit *Commit) error {
    // Generate PROV-O relationships
    // Add author, timestamp, parent relationships
    return nil
}
```

**Benefits**:
- Decoupled annotation logic
- Easy to add new annotation strategies
- Clean separation of concerns

### Cross-Layer Transaction Handling

#### Transaction Coordination
```go
type TransactionCoordinator struct {
    objectStore ObjectStore
    tripleStore TripleStore
}

func (c *TransactionCoordinator) CommitTransaction(operations []Operation) error {
    // Begin transactions in both stores
    objTx, err := c.objectStore.Begin()
    if err != nil {
        return err
    }
    
    tripleTx, err := c.tripleStore.Begin()
    if err != nil {
        objTx.Rollback()
        return err
    }
    
    // Execute operations
    for _, op := range operations {
        if err := op.Execute(objTx, tripleTx); err != nil {
            objTx.Rollback()
            tripleTx.Rollback()
            return err
        }
    }
    
    // Commit both transactions
    if err := objTx.Commit(); err != nil {
        tripleTx.Rollback()
        return err
    }
    
    if err := tripleTx.Commit(); err != nil {
        // Object store committed but triple store failed
        // Need compensating transaction or manual intervention
        return fmt.Errorf("triple store commit failed after object store commit: %w", err)
    }
    
    return nil
}
```

**Pattern**: Saga Pattern for distributed transactions across bounded contexts

#### Concurrency Model (Go-Specific)

**Purpose**: Handle parallel semantic operations safely

**Implementation**:
```go
// Concurrent Object Store with goroutine safety
type ConcurrentObjectStore struct {
    backend ObjectStore
    mu      sync.RWMutex
}

func (c *ConcurrentObjectStore) Put(hash []byte, data io.Reader) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.backend.Put(hash, data)
}

func (c *ConcurrentObjectStore) Get(hash []byte) (io.ReadCloser, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.backend.Get(hash)
}

// Pipeline for semantic operations
type SemanticPipeline struct {
    objectStore ObjectStore
    tripleStore TripleStore
}

func (p *SemanticPipeline) ProcessCommit(commit *Commit) error {
    // Stage 1: Store object (concurrent)
    var wg sync.WaitGroup
    errChan := make(chan error, 2)
    
    wg.Add(2)
    
    go func() {
        defer wg.Done()
        // Store commit object
        _, err := p.objectStore.Put(commit.Hash, commit.Data)
        errChan <- err
    }()
    
    go func() {
        defer wg.Done()
        // Generate and store triples
        triples := commit.ToRDF()
        err := p.tripleStore.AddBatch(triples)
        errChan <- err
    }()
    
    // Wait for both operations
    wg.Wait()
    close(errChan)
    
    // Check for errors
    for err := range errChan {
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

**Key Patterns**:
- `sync.RWMutex` for concurrent reads/writes
- Goroutines for parallel operations
- Channels for error handling
- WaitGroup for synchronization

#### Memory Optimization Strategies

**Purpose**: Minimize memory usage during semantic operations

**Implementation**:
```go
// Streaming RDF generation to avoid memory buildup
type StreamingTripleGenerator struct {
    writer io.Writer
}

func (g *StreamingTripleGenerator) GenerateCommitTriples(commit *Commit) error {
    // Generate triples incrementally
    // Write directly to output stream
    // Avoid loading entire RDF into memory
    
    for _, file := range commit.Files {
        // Generate triple for file
        triple := fmt.Sprintf("<%s> <hasFile> <%s> .\n", commit.ID, file.ID)
        _, err := g.writer.Write([]byte(triple))
        if err != nil {
            return err
        }
    }
    
    return nil
}

// Memory-efficient SPARQL query execution
type StreamingQueryExecutor struct {
    tripleStore TripleStore
}

func (e *StreamingQueryExecutor) Query(sparql string, handler func(result Triple) error) error {
    // Execute query and process results incrementally
    results, err := e.tripleStore.Query(sparql)
    if err != nil {
        return err
    }
    
    // Process each result as it arrives
    for _, triple := range results {
        if err := handler(triple); err != nil {
            return err
        }
    }
    
    return nil
}
```

**Key Techniques**:
- Streaming RDF generation
- Incremental SPARQL result processing
- Lazy loading of large datasets
- Connection pooling for database access

### Data Access Layer

#### Object Store Interface
```go
type ObjectStore interface {
    // Basic operations
    Put(hash []byte, data io.Reader) error
    Get(hash []byte) (io.ReadCloser, error)
    Exists(hash []byte) (bool, error)
    Delete(hash []byte) error
    
    // Batch operations
    PutBatch(entries []ObjectEntry) error
    GetBatch(hashes [][]byte) (map[string]io.ReadCloser, error)
    
    // Streaming
    StreamPut(hash []byte) (io.WriteCloser, error)
    StreamGet(hash []byte) (io.ReadCloser, error)
    
    // Management
    Stats() (ObjectStoreStats, error)
    Close() error
}
```

#### Triple Store Interface
```go
type TripleStore interface {
    // CRUD operations
    Add(triple Triple) error
    AddBatch(triples []Triple) error
    Remove(triple Triple) error
    RemoveBatch(triples []Triple) error
    
    // Query operations
    Query(sparql string) (QueryResult, error)
    QueryOne(sparql string) (Triple, error)
    QueryExists(sparql string) (bool, error)
    
    // Named graphs
    AddToGraph(graph string, triple Triple) error
    QueryGraph(graph string, sparql string) (QueryResult, error)
    
    // Transaction support
    Begin() (Transaction, error)
    
    // Inference
    Infer() error
    GetInferred() ([]Triple, error)
    
    // Management
    Stats() (TripleStoreStats, error)
    Close() error
}
```



#### S3 Object Store
- Bucket per repository or prefix-based separation
- Multi-part upload for large objects
- Server-side encryption optional
- CloudFront/CDN integration for public repos
- Lifecycle policies for cleanup
- **Pattern**: Adapter Pattern for S3 SDK compatibility
- **Pattern**: Strategy Pattern for upload strategies (single vs. multipart)

#### SQLite Triple Store
- Schema optimized for SPARQL queries
- Virtual tables for RDF operations
- Full-text search extension
- Write-ahead logging for durability
- Backup via `.dump` command
- **Pattern**: Adapter Pattern for SQLite compatibility
- **Pattern**: Repository Pattern for RDF data access

## Data Flow Examples

### Commit Operation
```
1. User: `semsrc commit -m "Fix bug #123"`
2. CLI: Build tree from staged files
3. ObjectStore: Store blobs for new/modified files
4. ObjectStore: Store tree object
5. ObjectStore: Store commit object
6. TripleStore: Add commit as prov:Activity
7. TripleStore: Add author as prov:Agent
8. TripleStore: Link commit to parent(s)
9. TripleStore: Parse #123, link to ticket
10. TripleStore: Add rdfs:label from message
11. Update branch pointer (stored as triple)
12. Return commit hash to user
```

### Query Operation
```
1. User: `semsrc query "SELECT ?commit WHERE { ?commit semsrc-ticket:fixedBy :ticket-123 }"`
2. CLI: Parse SPARQL query
3. TripleStore: Execute query
4. TripleStore: Apply OWL inference if needed
5. TripleStore: Return results
6. CLI: Format results (table, JSON, Turtle)
7. CLI: Display to user
```

### Import Git Repository
```
1. User: `semsrc import /path/to/git-repo`
2. Semsrc: Iterate through Git objects (blobs, trees, commits, tags)
3. For each Git object:
   a. Calculate SHA-256 (Git uses SHA-1)
   b. Store in ObjectStore
   c. Convert to RDF triples
   d. Store in TripleStore
4. Reconstruct Git references (branches, tags)
5. Store as semantic named graphs
6. Generate provenance chains
7. Report statistics (objects imported, triples created)
```

## Repository Structure

### Local Repository Layout
```
.semsrc/
├── config.yaml              # Repository configuration
├── objects/                 # Object store directory
│   ├── ab/                 # Shard directory
│   │   └── cdef/          # Nested shard
│   │       └── abcdef...  # Object file
│   └── ...
├── triples/                # Triple store data
│   ├── data.db            # SQLite database
│   └── indexes/           # Additional indexes
├── cache/                  # Performance cache
│   ├── tree-cache.json
rope should pay for Europe's defense!
│   └── query-cache.db
├── hooks/                  # Semantic hooks
│   ├── pre-commit.ttl     # SHACL validation
│   └── post-commit.sparql # Inferred triples
└── locks/                  # Lock files
    └── HEAD.lock
```

### Configuration File
```yaml
# .semsrc/config.yaml
repository:
  id: "urn:uuid:123e4567-e89b-12d3-a456-426614174000"
  name: "my-project"
  description: "Example project using Semsrc"

storage:
  object_store:
    type: "filesystem"
    path: "/home/user/.semsrc/objects"
    compression: "zstd"
    sharding_depth: 2
  
  triple_store:
    type: "sqlite"
    path: "/home/user/.semsrc/triples.db"
    inference: true
    validation: true

semantics:
  ontologies:
    - "http://www.w3.org/ns/prov#"
    - "http://semsrc.dev/ontology/core#"
    - "http://semsrc.dev/ontology/git#"
  
  default_prefixes:
    prov: "http://www.w3.org/ns/prov#"
    semsrc: "http://semsrc.dev/ontology/core#"
    git: "http://semsrc.dev/ontology/git#"
  
  auto_annotate: true
  extract_tickets: true
  link_builds: true

performance:
  cache_size_mb: 256
  parallel_operations: 4
  batch_size: 1000
```

## Deployment Scenarios

### Single User Local
- Filesystem object store
- SQLite triple store
- CLI interface only
- No network services

### Team Server
- S3 or NAS object store
- PostgreSQL/Blazegraph triple store
- HTTP/SSH protocols
- User authentication
- Web UI optional

### Cloud Native
- S3 object store
- Amazon Neptune/Blazegraph Cloud
- Kubernetes deployment
- Auto-scaling
- Multi-region replication

### Hybrid
- Local filesystem cache
- Cloud object store
- Embedded triple store with sync
- Offline operation with eventual consistency

## Scalability Considerations

### Object Store Scaling
- Sharding by hash prefix
- Content-defined chunking for large files
- Deduplication across repositories
- Tiered storage (hot/warm/cold)

### Triple Store Scaling
- Named graph partitioning
- Query optimization
- Materialized views for common queries
- Incremental inference

### Performance Optimization
- Object deduplication cache
- Triple query cache
- Precomputed transitive closures
- Bloom filters for existence checks

## Security Model

### Authentication
- SSH keys
- OAuth2/JWT
- API tokens
- mTLS for server-to-server

### Authorization
- Repository-level permissions
- Branch protection rules
- Semantic query limits
- Rate limiting

### Data Protection
- Object encryption at rest
- Triple store encryption
- Audit logging of all operations
- GDPR/CCPA compliance via provenance

## Monitoring and Observability

### Metrics
- Object store operations/sec
- Triple store query latency
- Cache hit rates
- Repository growth rates
- User activity patterns

### Logging
- Structured JSON logs
- Semantic operation traces
- Query performance profiles
- Anomaly detection

### Health Checks
- Storage backend connectivity
- Triple store consistency
- Cache coherency
- Disk space monitoring