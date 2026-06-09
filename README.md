# Semsrc - Semantic Source Control

[![Go Reference](https://pkg.go.dev/badge/github.com/djha-skin/semsrc.svg)](https://pkg.go.dev/github.com/djha-skin/semsrc)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Semsrc is an AI-native, holistic version control system based on object store and RDF triple store technologies.

## Vision

Traditional version control systems like Git focus narrowly on source code changes. Semsrc extends this to encompass the entire software development lifecycle using semantic web technologies, making it understandable and extensible by AI agents while remaining familiar to human developers.

## Key Features

- **AI-Native**: Semantic triples provide structure that AIs can understand and extend
- **Holistic**: Unified management of code, tickets, builds, configurations, and documentation
- **Portable**: Works locally, on-premise, or in the cloud with the same interface
- **Semantic**: Built on RDF/SPARQL with PROV-O and NATMUK ontologies
- **Git-Compatible**: Familiar commands with semantic enhancements

## Quick Start

### Installation

```bash
# Build from source
git clone https://github.com/djha-skin/semsrc.git
cd semsrc
go build ./cmd/semsrc

# Or install via go install
go install github.com/djha-skin/semsrc/cmd/semsrc@latest
```

### Basic Usage

```bash
# Initialize a new repository
semsrc init my-project
cd my-project

# Add files and commit (just like Git)
semsrc add README.md
semsrc commit -m "Initial commit"

# View semantic history
semsrc log --semantic

# Query relationships
semsrc query "SELECT ?commit WHERE { ?commit semsrc-ticket:fixedBy :ticket-123 }"
```

## Architecture

Semsrc is built on two fundamental layers:

1. **Object Store**: Content-addressable storage for blobs (files, binaries, artifacts)
2. **Triple Store**: RDF-based semantic database for relationships and metadata

These layers are abstracted behind clean Go interfaces, allowing multiple implementations:

- **Object Store Backends**: Filesystem, S3, Azure Blob, memory
- **Triple Store Backends**: SQLite, PostgreSQL/RDF, Blazegraph, in-memory

## Design Philosophy

### 1. Pragmatic Semantics
Use RDF where it adds value, not everywhere. Start with Git compatibility, add semantics where they help.

### 2. Progressive Enhancement
Opt-in to advanced features. The basic Git workflow works immediately; semantic features are available when needed.

### 3. Open Standards
Build on existing standards (PROV-O, NATMUK) rather than inventing new ones.

### 4. AI First
Design for AI consumption: structured data, clear semantics, extensible without breaking changes.

## Documentation

- [Golden Thread](./design/golden-thread.md) - Vision and overview
- [User Stories](./design/user-stories/) - Phased development approach
- [Ontology Design](./design/ontologies/semsrc-ontology.md) - Semantic data model
- [System Architecture](./design/architecture/system-architecture.md) - Technical design

## Development Status

⚠️ **Early Development** - Semsrc is currently in design and prototyping phase.

### Current Focus
- Core object store interface and filesystem implementation
- Triple store interface with SQLite backend
- Basic Git command compatibility layer

### Roadmap
See [User Stories](./design/user-stories/01-foundation-stories.md) for detailed phased approach.

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines.

1. Check open issues or propose new features
2. Fork the repository
3. Create a feature branch
4. Submit a pull request

## License

MIT License - see [LICENSE](./LICENSE) for details.

## Acknowledgments

- **PROV-O Ontology**: For provenance tracking standards
- **NATMUK Filesystem Ontology**: For filesystem semantics
- **Git**: For inspiration and compatibility target
- **RDF/SPARQL Community**: For semantic web foundations

## Contact

- GitHub Issues: [Issue Tracker](https://github.com/djha-skin/semsrc/issues)
- Discussions: [GitHub Discussions](https://github.com/djha-skin/semsrc/discussions)

---

*"The future of version control is semantic"*