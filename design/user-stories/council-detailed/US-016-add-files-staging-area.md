# US-016: Add Files to Staging Area

## Description
Add files to the staging area to prepare for commit. The staging area tracks files that will be included in the next commit, allowing selective commits and partial commits.


## Test Cases
1. **Single file staging**: Stage a single file successfully
2. **Multiple files staging**: Stage multiple files in one operation
3. **Directory staging**: Stage all files in a directory recursively
4. **Already staged**: Handle already staged files gracefully
5. **Non-existent files**: Handle missing files with clear errors

## Dependencies
- US-001: Store Files with SHA-256 Content Addressing
- US-015: Initialize New Repository

## Implementation Notes
- Calculate SHA-256 hash for each file
- Store file content in object store using hash-based addressing
- Maintain staging index separate from committed state
- Track file metadata: path, hash, size, permissions, timestamps
- Support incremental staging (only changed files)
- Use efficient hashing for large files (streaming)

## Edge Cases
- **Modified files**: Should detect if file changed since last staging
- **Binary files**: Should handle binary files correctly
- **Large files**: Should handle files larger than memory
- **Symbolic links**: Should follow or ignore symlinks based on config
- **Permission changes**: Should detect and track permission changes