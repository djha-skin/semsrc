package store

import (
	"errors"
	"time"
)

// Triple represents an RDF triple (subject, predicate, object)
type Triple struct {
	Subject   Resource
	Predicate Resource
	Object    Resource
}

// Resource represents an RDF resource (IRI, literal, or blank node)
type Resource interface {
	// String returns the string representation
	String() string
	// Type returns the resource type
	Type() ResourceType
	// Equal checks if two resources are equal
	Equal(other Resource) bool
}

// ResourceType indicates the type of RDF resource
type ResourceType int

const (
	ResourceIRI ResourceType = iota
	ResourceLiteral
	ResourceBlank
)

// IRI represents an Internationalized Resource Identifier
type IRI string

func (i IRI) String() string  { return string(i) }
func (i IRI) Type() ResourceType { return ResourceIRI }
func (i IRI) Equal(other Resource) bool {
	if o, ok := other.(IRI); ok {
		return i == o
	}
	return false
}

// Literal represents an RDF literal with optional datatype and language
type Literal struct {
	Value    string
	Datatype IRI
	Language string
}

func (l Literal) String() string {
	if l.Language != "" {
		return `"` + l.Value + `"@` + l.Language
	}
	if l.Datatype != "" {
		return `"` + l.Value + `"^^` + string(l.Datatype)
	}
	return `"` + l.Value + `"`
}

func (l Literal) Type() ResourceType { return ResourceLiteral }
func (l Literal) Equal(other Resource) bool {
	if o, ok := other.(Literal); ok {
		return l.Value == o.Value && l.Datatype == o.Datatype && l.Language == o.Language
	}
	return false
}

// BlankNode represents an RDF blank node
type BlankNode struct {
	ID string
}

func (b BlankNode) String() string { return "_:" + b.ID }
func (b BlankNode) Type() ResourceType { return ResourceBlank }
func (b BlankNode) Equal(other Resource) bool {
	if o, ok := other.(BlankNode); ok {
		return b.ID == o.ID
	}
	return false
}

// TripleStore defines the interface for RDF triple storage and querying
type TripleStore interface {
	// Add adds a triple to the store
	Add(triple Triple) error
	
	// Remove removes a triple from the store
	Remove(triple Triple) error
	
	// Query executes a SPARQL query
	Query(sparql string) (QueryResult, error)
	
	// QueryOne executes a SPARQL query expecting exactly one result
	QueryOne(sparql string) (map[string]Resource, error)
	
	// QueryExists executes a SPARQL ASK query
	QueryExists(sparql string) (bool, error)
	
	// Begin starts a transaction
	Begin() (Transaction, error)
	
	// Stats returns statistics about the store
	Stats() (TripleStoreStats, error)
	
	// Close releases resources
	Close() error
}

// QueryResult represents the result of a SPARQL query
type QueryResult struct {
	Variables []string
	Rows      []map[string]Resource
	Boolean   *bool // For ASK queries
}

// Transaction represents a triple store transaction
type Transaction interface {
	// Add adds a triple within the transaction
	Add(triple Triple) error
	
	// Remove removes a triple within the transaction
	Remove(triple Triple) error
	
	// Query executes a SPARQL query within transaction context
	Query(sparql string) (QueryResult, error)
	
	// Commit commits the transaction
	Commit() error
	
	// Rollback rolls back the transaction
	Rollback() error
}

// TripleStoreStats contains statistics about the triple store
type TripleStoreStats struct {
	TripleCount    int64
	DistinctSubjects int64
	DistinctPredicates int64
	DistinctObjects  int64
	LastUpdate     time.Time
}

// NamedGraphStore extends TripleStore with named graph support
type NamedGraphStore interface {
	TripleStore
	
	// AddToGraph adds a triple to a specific named graph
	AddToGraph(graph string, triple Triple) error
	
	// RemoveFromGraph removes a triple from a specific named graph
	RemoveFromGraph(graph string, triple Triple) error
	
	// QueryGraph executes a SPARQL query against a specific graph
	QueryGraph(graph string, sparql string) (QueryResult, error)
	
	// ListGraphs returns all named graphs in the store
	ListGraphs() ([]string, error)
}

// InferenceStore extends TripleStore with OWL inference capabilities
type InferenceStore interface {
	TripleStore
	
	// Infer performs OWL reasoning and adds inferred triples
	Infer() error
	
	// GetInferred returns all inferred triples
	GetInferred() ([]Triple, error)
	
	// ClearInferred removes all inferred triples
	ClearInferred() error
	
	// IsInferred checks if a triple is inferred (not asserted)
	IsInferred(triple Triple) (bool, error)
}

// BatchTripleStore extends TripleStore with batch operations
type BatchTripleStore interface {
	TripleStore
	
	// AddBatch adds multiple triples in a single operation
	AddBatch(triples []Triple) error
	
	// RemoveBatch removes multiple triples in a single operation
	RemoveBatch(triples []Triple) error
	
	// QueryBatch executes multiple queries in a single operation
	QueryBatch(queries []string) ([]QueryResult, error)
}

// Common errors
var (
	ErrTripleExists    = errors.New("triple already exists")
	ErrTripleNotFound  = errors.New("triple not found")
	ErrInvalidSPARQL   = errors.New("invalid SPARQL query")
	ErrTransaction     = errors.New("transaction error")
	ErrInference       = errors.New("inference error")
	ErrGraphNotFound   = errors.New("named graph not found")
)

// Helper functions for creating resources

// NewIRI creates a new IRI resource
func NewIRI(s string) IRI {
	return IRI(s)
}

// NewLiteral creates a new literal resource
func NewLiteral(value string) Literal {
	return Literal{Value: value}
}

// NewLiteralWithType creates a new typed literal
func NewLiteralWithType(value string, datatype IRI) Literal {
	return Literal{Value: value, Datatype: datatype}
}

// NewLiteralWithLanguage creates a new language-tagged literal
func NewLiteralWithLanguage(value, language string) Literal {
	return Literal{Value: value, Language: language}
}

// NewBlankNode creates a new blank node
func NewBlankNode(id string) BlankNode {
	return BlankNode{ID: id}
}

// NewTriple creates a new triple
func NewTriple(subject, predicate, object Resource) Triple {
	return Triple{
		Subject:   subject,
		Predicate: predicate,
		Object:    object,
	}
}

// Standard RDF/OWL IRIs for convenience
var (
	RDFType      = NewIRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type")
	RDFSLabel    = NewIRI("http://www.w3.org/2000/01/rdf-schema#label")
	RDFSComment  = NewIRI("http://www.w3.org/2000/01/rdf-schema#comment")
	OWLClass     = NewIRI("http://www.w3.org/2002/07/owl#Class")
	OWLObjectProperty = NewIRI("http://www.w3.org/2002/07/owl#ObjectProperty")
	OWLDatatypeProperty = NewIRI("http://www.w3.org/2002/07/owl#DatatypeProperty")
	
	// XSD datatypes
	XSDString    = NewIRI("http://www.w3.org/2001/XMLSchema#string")
	XSDInteger   = NewIRI("http://www.w3.org/2001/XMLSchema#integer")
	XSDBoolean   = NewIRI("http://www.w3.org/2001/XMLSchema#boolean")
	XSDDateTime  = NewIRI("http://www.w3.org/2001/XMLSchema#dateTime")
)