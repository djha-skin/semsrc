# Semsrc Golden Thread

## Vision Statement

Semsrc is a semantic, AI-native version control system that extends beyond
traditional source code management to encompass the entire software development
lifecycle using semantic web technologies.

## Core Problem Statement

Traditional version control systems like Git focus narrowly on source code
changes, forcing teams to maintain parallel systems for tickets, build results,
CI/CD configurations, and documentation. This fragmentation creates cognitive
overhead, integration challenges, and limits AI assistance capabilities.

## Solution Overview

Semsrc integrates:

1. **Object Store** - Content-addressable storage for blobs (files, binaries,
   artifacts)

2. **Triple Store** - RDF-based semantic database tracking relationships,
   metadata, and provenance

3. **Semantic Understanding** - AI-native interfaces that understand software
   development concepts

4. **Holistic Scope** - Unified management of code, tickets, builds,
   configurations, and documentation

## Key Differentiators

### 1. AI-Native Design

- Semantic triples provide structure that AIs can understand and extend
- Built-in ontologies (PROV-O, NATMUK) give AIs familiar conceptual frameworks
- AI agents can add meaningful metadata without schema modifications
- Queryable semantic relationships enable intelligent analysis

### 2. Holistic Software Management

- Code commits linked to tickets via semantic relationships
- Build results stored as objects with provenance metadata
- CI/CD configurations as versioned semantic specifications
- Documentation as first-class citizens in the triple store
- Everything versioned together with clear dependencies

### 3. Portability and Deployment

- Object store abstraction supports filesystem, S3, or custom backends
- Triple store portable across **Amazon Neptune** (cloud) and **Oxigraph** (local)
- **Council Decision**: Native RDF stores preferred over SQLite (see ADR-001)
- Single-binary deployment for local development
- Scalable to enterprise multi-tenant clouds
- Works offline with eventual consistency capabilities

## Technical Foundation

### Bounded Context Architecture

Semsrc follows a strict bounded context architecture to ensure clear separation of concerns and maintainability:

```
┌─────────────────────────────────────────────────────────┐
│                    Application Layer                    │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐    │
│  │   CLI/API   │  │  Web UI     │  │   Plugins   │    │
│  └─────────────┘  └─────────────┘  └─────────────┘    │
└─────────────────────────────────────────────────────────┘
                    │ Interface Contracts
┌─────────────────────────────────────────────────────────┐
│                  Business Logic Layer                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐    │
│  │  Version    │  │  Semantic   │  │   Query     │    │
│  │  Control    │  │ Annotation  │  │  & Reason   │    │
│  └─────────────┘  └─────────────┘  └─────────────┘    │
│  Bounded Context: Core Domain Logic                    │
└─────────────────────────────────────────────────────────┘
                    │ Interface Contracts
┌─────────────────────────────────────────────────────────┐
│                   Data Access Layer                     │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐    │
│  │ Object Store│  │ Triple Store│  │   Index &   │    │
│  │  Interface  │  │  Interface  │  │   Cache     │    │
│  └─────────────┘  └─────────────┘  └─────────────┘    │
│  Bounded Context: Data Persistence                     │
└─────────────────────────────────────────────────────────┘
                    │ Interface Contracts
┌─────────────────────────────────────────────────────────┐
│                Storage Implementation Layer             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐    │
│  │  Filesystem │  │     S3      │  │ Blazegraph/ │    │
│  │Object Store │  │Object Store │  │   Oxigraph  │    │
│  └─────────────┘  └─────────────┘  └─────────────┘    │
│  Bounded Context: Storage Backends                     │
└─────────────────────────────────────────────────────────┘
```

**Key Bounded Context Principles:**
- Each layer has well-defined interfaces
- No cross-layer dependencies except through interfaces
- Data flows unidirectionally downward
- Each context can be tested independently

### Object Store Layer

- Content-addressable storage (SHA-256 hashes)
- Pluggable storage backends (Filesystem, S3, Azure Blob, etc.)
- Efficient deduplication of identical content
- Streaming support for large blobs
- Atomic operations with rollback capability
- **Strategy Pattern**: Backend selection via configuration
- **Factory Pattern**: Object store creation via dependency injection

### Triple Store Layer

- RDF N-Triples/Turtle serialization
- SPARQL query endpoint
- Support for OWL inferencing
- Built-in ontologies:
  - PROV-O for provenance tracking
  - NATMUK Filesystem for file operations
  - Custom Semsrc ontology extending both
- Transaction support with ACID semantics
- **Repository Pattern**: Abstract access to semantic data

### Semantic Ontology Strategy

1. **Reuse existing standards** where possible (PROV-O, NATMUK)
2. **Extend carefully** with domain-specific concepts
3. **Validate with OWL** for consistency checking
4. **Document thoroughly** with human and AI readable descriptions
5. **Version ontologies** alongside the data they describe
6. **Bounded Contexts**: Core ontologies vs. domain-specific extensions

## Core User Value Propositions

### For Developers

- Unified view of code, tickets, and builds
- Semantic search across entire project history
- Easy navigation of relationships (e.g., "show me all commits that fixed ticket X")
- Automated documentation from structured metadata
- Better code understanding through semantic annotations

### For Teams

- Shared understanding encoded in triples
- Reduced tribal knowledge through explicit relationships
- Automated compliance and audit trails
- Better onboarding with semantically rich project context
- Cross-project analysis via common ontology

### For AI Agents

- Structured data suitable for reasoning
- Ability to add meaningful annotations
- Understandable software development concepts
- Queryable knowledge base for planning and analysis
- Extensible without breaking changes

## Non-Goals

- **Not** just Git with RDF metadata
- **Not** requiring all users to understand RDF/SPARQL
- **Not** replacing all existing tools immediately
- **Not** sacrificing performance for semantic richness
- **Not** limited to specific cloud providers

## Success Metrics

1. **Developer Experience**: As productive as Git for daily workflows
2. **AI Integration**: AIs can successfully query and extend the system
3. **Performance**: Comparable to Git for common operations
4. **Adoption**: Teams choose Semsrc for new projects
5. **Extensibility**: Third-party integrations and plugins thrive
6. **Evolutionary Progress**: Track metrics for migration success:
   - Time to onboard new developer
   - AI query success rate
   - Semantic feature adoption rate
   - Migration completion percentage

## Implementation Principles

1. **Pragmatic Semantics**: Use RDF where it adds value, not everywhere
2. **Gradual Migration**: Interoperate with existing Git repositories
3. **Progressive Enhancement**: Start with Git-like features, add semantics
4. **Backward Compatibility**: Never break existing queries/API
5. **Open Standards**: Use and contribute to standard ontologies
6. **Evolutionary Patterns**: Apply Strangler Fig pattern for migration

## Evolutionary Migration Pattern

### Strangler Fig Pattern for Semsrc Adoption

The Strangler Fig pattern provides a proven approach for migrating from Git to Semsrc:

**Phase 1: Parallel Operation**
- Deploy Semsrc alongside existing Git repositories
- Semsrc reads from Git, writes to both Git and Semsrc
- No disruption to existing workflows
- **Success Metric**: 0% workflow interruption

**Phase 2: Gradual Feature Introduction**
- Enable semantic features for new projects only
- Existing projects continue using Git-only workflows
- Developers opt-in to semantic features gradually
- **Success Metric**: 30% of new projects using semantic features

**Phase 3: Semantic Enhancement**
- Migrate critical projects to Semsrc with full semantics
- Maintain Git compatibility for external collaboration
- Enable AI-assisted features for migrated projects
- **Success Metric**: 50% of active projects using Semsrc

**Phase 4: Full Migration**
- Complete migration of all projects
- Retire Git compatibility layer for internal projects
- **Success Metric**: 90% project migration completion

### Migration Scenarios

**Scenario 1: Small Team (1-10 developers)**
- Direct migration with minimal parallel operation
- Estimated timeline: 2-4 weeks
- Risk level: Low

**Scenario 2: Medium Team (10-50 developers)**
- Gradual migration with parallel operation for 1-2 months
- Estimated timeline: 2-3 months
- Risk level: Medium

**Scenario 3: Large Organization (50+ developers)**
- Phased migration with dedicated migration team
- Estimated timeline: 6-12 months
- Risk level: Medium-High

## Next Steps

See the user stories in `user-stories/` for phased development approach,
starting with core object/triple store foundations and building toward full Git
compatibility plus semantic extensions.

### Immediate Actions

1. **Implement Bounded Context Architecture**
   - Define interface contracts between layers
   - Create dependency injection framework
   - Establish testing strategy for each context

2. **Add Design Pattern Implementations**
   - Repository Pattern for data access
   - Strategy Pattern for storage backend selection
   - Factory Pattern for object creation
   - **Worktree Management**: New pattern for parallel development workflows
   - **Staging Area**: Advanced staging with hunk-by-hunk control
   - **Semantic Merge**: RDF-aware conflict resolution

3. **Establish Evolutionary Migration Framework**
   - Create migration tools for Git repositories
   - Define success metrics and tracking
   - Document migration scenarios
   - **Worktree Support**: Multiple working trees for parallel development
   - **Reflog & Recovery**: Safety net for developer mistakes
   - **Tag Management**: Semantic versioning support