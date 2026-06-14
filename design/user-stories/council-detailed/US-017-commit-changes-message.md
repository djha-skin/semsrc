# US-017: Commit Changes with a Message

## Description

Commit staged changes to the repository with a descriptive message. This
creates a permanent snapshot of the repository state that can be referenced
later.


## Test Cases

1. **Basic commit**: Commit staged changes with message
2. **Empty commit**: Handle commit with no staged changes (error)
3. **Commit metadata**: Author, timestamp, and message are recorded
4. **Parent reference**: Commit points to previous commit
5. **Tree creation**: Commit creates tree object for file structure

## Dependencies

- US-016: Add Files to Staging Area
- US-001: Store Files with SHA-256 Content Addressing
- US-009: Store RDF Triples in SQLite

## Implementation Notes

- Create commit object with metadata
- Generate tree object representing file structure
- Store commit and tree in object store
- Update repository HEAD to point to new commit
- Record commit in triple store for semantic queries
- Validate commit message format
- Support commit signing (optional)

## Edge Cases

- **Invalid message**: Should validate commit message format
- **Concurrent commits**: Should handle race conditions safely
- **Large commits**: Should handle commits with many files efficiently
- **Failed commit**: Should rollback on failure
- **Repository corruption**: Should detect and report corruption