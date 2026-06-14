# US-040: Import Existing Git Repositories

## Description
Import existing Git repositories into Semsrc, preserving commit history,
branches, and tags. This enables gradual migration from Git to Semsrc without
losing historical data.


## Test Cases
1. **Full repository import**: Import complete Git repository with history
2. **Large repository**: Handle repositories with thousands of commits
3. **Incremental import**: Import repository in stages for large repos
4. **Preserve metadata**: Maintain branches, tags, and commit metadata
5. **Error handling**: Handle import failures gracefully

## Dependencies
- US-009: Store RDF Triples in Native RDF Store
- US-017: Commit Changes with a Message
- US-036: Git CLI Compatibility

## Implementation Notes
- Use Git library to read repository data
- Import commits in chronological order
- Preserve commit metadata (author, timestamp, message)
- Import branches and tags as references
- Store commit relationships in semantic graph
- Provide progress reporting for long imports

## Edge Cases
- **Corrupted Git repository**: Should detect and report corruption
- **Very large repositories**: Should handle repos with millions of commits
- **Failed imports**: Should rollback or provide resume capability
- **Permission issues**: Should handle file permission problems
- **Cross-platform**: Should work on different operating systems