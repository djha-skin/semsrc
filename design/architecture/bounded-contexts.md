# Bounded Contexts Architecture

## Overview

This document defines the bounded contexts for the Semsrc semantic version control system, following Domain-Driven Design principles to ensure clear separation of concerns, maintainability, and testability.

## Context Map

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           Application Layer                                 │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐          │
│  │     CLI/API      │  │    Web UI        │  │     Plugins      │          │
│  │  (Presentation)  │  │  (Presentation)  │  │ (Presentation)   │          │
│  └──────────────────┘  └──────────────────┘  └──────────────────┘          │
│                      │ Interface Contracts                                 │
└─────────────────────────────────────────────────────────────────────────────┘
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Business Logic Layer                                   │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐          │
│  │ Version Control  │  │ Semantic         │  │ Query &          │          │
│  │  Operations      │  │ Annotation       │  │ Reasoning        │          │
│  └──────────────────┘  └──────────────────┘  └──────────────────┘          │
│                      │ Interface Contracts                                 │
│  Bounded Context: Core Domain Logic                                        │
└─────────────────────────────────────────────────────────────────────────────┘
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Data Access Layer                                     │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐          │
│  │  Object Store    │  │  Triple Store    │  │   Index &        │          │
│  │   Interface      │  │   Interface      │  │   Cache          │          │
│  └──────────────────┘  └──────────────────┘  └──────────────────┘          │
│                      │ Interface Contracts                                 │
│  Bounded Context: Data Persistence                                         │
└─────────────────────────────────────────────────────────────────────────────┘
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                  Storage Implementation Layer                               │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐          │
│  │   Filesystem     │  │      S3          │  │    SQLite/       │          │
│  │  Object Store    │  │  Object Store    │  │  Embedded        │          │
│  └──────────────────┘  └──────────────────┘  └──────────────────┘          │
│  Bounded Context: Storage Backends                                         │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Bounded Context Definitions

### 1. Application Layer (Presentation Context)

**Purpose**: User interaction and presentation logic

**Components**:
- **CLI**: Command-line interface for terminal users
- **API**: HTTP/REST API for programmatic access
- **Web UI**: Browser-based interface (optional)
- **Plugins**: Extension system for third-party integrations

**Key Principles**:
- No business logic in this layer
- Delegates all operations to Business Logic Layer
- Formats responses for different presentation formats
- Handles input validation and command parsing

**Interface Contracts**:
```go
type CommandHandler interface {
    Execute(args []string) (Response, error)
}

type APIHandler interface {
    HandleRequest(req *http.Request) (*http.Response, error)
}
```

### 2. Business Logic Layer (Core Domain Context)

**Purpose**: Core version control and semantic operations

**Sub-Contexts**:

#### 2.1 Version Control Operations
**Responsibilities**:
- Commit operations
- Branch management
- Merge operations
- Diff calculations
- History traversal

**Key Patterns**:
- **Repository Pattern**: Abstract data access
- **Command Pattern**: Encapsulate operations

**Code Example (Repository Pattern)**:
```go
// CommitRepository interface abstracts data access
type CommitRepository interface {
    Save(commit *Commit) error
    FindByID(id string) (*Commit, error)
    FindByRange(start, end int) ([]*Commit, error)
}

// Concrete implementation
type SqlCommitRepository struct {
    db *sql.DB
}

func (r *SqlCommitRepository) Save(commit *Commit) error {
    // SQL implementation details
    _, err := r.db.Exec("INSERT INTO commits (...) VALUES (...)", ...)
    return err
}

// Usage in Business Logic Layer
func (vc *VersionControl) Commit(message string) error {
    repo := vc.repoFactory.CreateCommitRepository()
    commit := NewCommit(message)
    return repo.Save(commit)
}
```

#### 2.2 Semantic Annotation Engine
**Responsibilities**:
- PROV-O relationship generation
- Ticket extraction from commit messages
- Build link creation

**Code Example (Observer Pattern)**:
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

// Ticket Extraction Observer
type TicketExtractor struct{}

func (t *TicketExtractor) OnCommit(commit *Commit) error {
    ticketIDs := extractTicketIDs(commit.Message)
    for _, ticketID := range ticketIDs {
        // Link commit to ticket in triple store
    }
    return nil
}

// Usage
e := &SemanticAnnotationEngine{}
e.Register(&TicketExtractor{})
e.Register(&ProvenanceGenerator{})
// When commit happens:
e.ProcessCommit(commit)
```
- File rename tracking
- Metadata inference

**Key Patterns**:
- **Observer Pattern**: React to version control events
- **Template Method**: Standard annotation pipeline

#### 2.3 Query & Reasoning Engine
**Responsibilities**:
- SPARQL query execution
- OWL inference
- SHACL validation
- Full-text search
- Graph traversal

**Key Patterns**:
- **Strategy Pattern**: Query optimization strategies
- **Factory Pattern**: Query plan generation

### 3. Data Access Layer (Persistence Context)

**Purpose**: Abstract storage operations

#### 3.1 Object Store Interface
**Responsibilities**:
- Content-addressable storage operations
- Streaming support
- Batch operations
- Compression management

**Interface Definition**:
```go
type ObjectStore interface {
    Put(hash []byte, data io.Reader) error
    Get(hash []byte) (io.ReadCloser, error)
    Exists(hash []byte) (bool, error)
    Delete(hash []byte) error
    PutBatch(entries []ObjectEntry) error
    GetBatch(hashes [][]byte) (map[string]io.ReadCloser, error)
    StreamPut(hash []byte) (io.WriteCloser, error)
    StreamGet(hash []byte) (io.ReadCloser, error)
    Stats() (ObjectStoreStats, error)
    Close() error
}
```

#### 3.2 Triple Store Interface
**Responsibilities**:
- RDF triple storage
- SPARQL query execution
- Named graph support
- Transaction support
- Inference execution

**Interface Definition**:
```go
type TripleStore interface {
    Add(triple Triple) error
    AddBatch(triples []Triple) error
    Remove(triple Triple) error
    RemoveBatch(triples []Triple) error
    Query(sparql string) (QueryResult, error)
    QueryOne(sparql string) (Triple, error)
    QueryExists(sparql string) (bool, error)
    AddToGraph(graph string, triple Triple) error
    QueryGraph(graph string, sparql string) (QueryResult, error)
    Begin() (Transaction, error)
    Infer() error
    GetInferred() ([]Triple, error)
    Stats() (TripleStoreStats, error)
    Close() error
}
```

### 4. Storage Implementation Layer (Infrastructure Context)

**Purpose**: Concrete storage backend implementations

#### 4.1 Filesystem Object Store
**Characteristics**:
- Directory-per-repository layout
- Hash-based sharding (2 levels)
- File locking via flock()
- Atomic writes via temp files
- Compression support (zstd/lz4)

**Implementation Pattern**:
- **Factory Pattern**: Create filesystem store with config
- **Strategy Pattern**: Compression algorithm selection

#### 4.2 S3 Object Store
**Characteristics**:
- Bucket-per-repository or prefix-based
- Multi-part upload for large objects
- Server-side encryption
- CDN integration support

**Implementation Pattern**:
- **Adapter Pattern**: Wrap S3 SDK for ObjectStore interface

#### 4.3 SQLite Triple Store
**Characteristics**:
- Schema optimized for SPARQL queries
- Virtual tables for RDF operations
- Full-text search extension
- Write-ahead logging for durability

**Implementation Pattern**:
- **Adapter Pattern**: Wrap SQLite for TripleStore interface

## Context Relationships

### Downward Dependencies
```
Application Layer → Business Logic Layer
    ↓ (delegates operations)
Business Logic Layer → Data Access Layer
    ↓ (uses interfaces)
Data Access Layer → Storage Implementation Layer
    ↓ (concrete implementations)
```

### Upward Flow
```
Storage Implementation Layer → Data Access Layer
    ↓ (returns data)
Data Access Layer → Business Logic Layer
    ↓ (processes data)
Business Logic Layer → Application Layer
    ↓ (formats responses)
```

### Cross-Context Communication
- **Application → Business**: Via interface contracts only
- **Business → Data Access**: Via repository interfaces
- **Data Access → Storage**: Via dependency injection

## Dependency Injection

### Factory Pattern for Storage Creation
```go
type ObjectStoreFactory interface {
    Create(config Config) (ObjectStore, error)
}

type TripleStoreFactory interface {
    Create(config Config) (TripleStore, error)
}
```

### Configuration-Driven Creation
```yaml
storage:
  object_store:
    type: "filesystem"  # or "s3", "azure"
    path: "/path/to/objects"
    compression: "zstd"
  
  triple_store:
    type: "blazegraph"  # Native RDF store preferred
    # Alternative options:
    # type: "sqlite"    # Embedded, lighter weight
    # type: "oxigraph"   # Rust-based embedded
    # type: "postgres"   # For team server scenarios
    path: "/path/to/triples.db"
```

## Testing Strategy

### Unit Tests
- **Application Layer**: Command parsing, API request/response handling
- **Business Logic Layer**: All business rules and operations
- **Data Access Layer**: Interface contract verification
- **Storage Implementation Layer**: Backend-specific operations

### Integration Tests
- **Cross-layer**: Full operation flow from CLI to storage
- **Storage backends**: Verify all backends work correctly
- **Semantic operations**: End-to-end semantic annotation and queries

### Contract Tests
- **Interface verification**: Ensure implementations satisfy interfaces
- **Backward compatibility**: Verify no breaking changes

## Evolutionary Migration Support

### Strangler Fig Implementation
- **Phase 1**: Semsrc reads from Git, writes to both
- **Phase 2**: Gradual enablement of semantic features
- **Phase 3**: Full migration with Git compatibility
- **Phase 4**: Complete migration with optional Git compatibility

### Git Compatibility Layer
- **Bounded Context**: Separate context for Git operations
- **Interface**: Translates Git commands to Semsrc operations
- **Data Flow**: Bidirectional synchronization during migration

## Benefits of Bounded Context Architecture

1. **Clear Separation of Concerns**: Each layer has distinct responsibilities
2. **Testability**: Each context can be tested independently
3. **Maintainability**: Changes in one layer don't affect others
4. **Scalability**: Each layer can scale independently
5. **Evolutionary Design**: Easy to add new storage backends or features
