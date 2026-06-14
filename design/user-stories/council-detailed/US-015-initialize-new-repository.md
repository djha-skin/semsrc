# US-015: Initialize a New Repository

## Description

Initialize a new repository with the necessary structure and metadata. This
sets up the foundational directory structure, creates initial configuration
files, and establishes the semantic graph with default contexts.


## Test Cases

1. **Basic initialization**: Create repository with default structure
2. **Custom path**: Initialize repository at specified path
3. **Configuration files**: Create proper config.yaml with defaults
4. **Graph initialization**: Create default named graphs
5. **Repository ID**: Generate unique repository identifier

## Dependencies

- US-001: Store Files with SHA-256 Content Addressing
- US-009: Store RDF Triples in Native RDF Store
- US-013: Named Graph Support

## Implementation Notes

- Create `.semsrc` directory for configuration
- Initialize object store with configured backend
- Initialize triple store with SQLite database
- Create default graphs: provenance, annotations, validation
- Generate unique repository ID (UUID or similar)
- Store repository metadata in triple store
- Create initial `.gitignore` equivalent for semsrc
- Support initialization from existing Git repository

## Edge Cases

- **Existing directory**: Should handle initialization in non-empty directory
- **Permission errors**: Should provide clear error messages
- **Already initialized**: Should detect and report if already initialized
- **Custom backends**: Should support custom storage backends during init
- **Failed initialization**: Should rollback partially created structure