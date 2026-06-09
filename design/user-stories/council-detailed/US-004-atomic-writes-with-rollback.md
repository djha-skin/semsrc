# US-004: Atomic Writes with Rollback Capability

## Description
Implement atomic write operations for the object store. When writing objects, the operation should be atomic - either the entire object is written successfully, or nothing is written. If a write fails mid-operation, the system should automatically rollback to a consistent state.


## Test Cases
1. **Successful write**: Object is written atomically and available for retrieval
2. **Write failure mid-operation**: System remains consistent, no partial writes
3. **Concurrent writes**: Multiple processes writing same object should not corrupt data
4. **Disk full during write**: Should rollback cleanly without leaving temp files
5. **Power loss during write**: System should recover to consistent state on restart

## Dependencies
- US-001: Store Files with SHA-256 Content Addressing
- US-002: Configure Filesystem Storage Backend

## Implementation Notes
- Use filesystem atomic rename operation for final commit
- Write to temporary file first, verify hash, then rename
- Use file locking to prevent concurrent writes to same object
- Implement transaction log for rollback capability
- On startup, scan for orphaned temp files and clean up
- Use fsync() or equivalent to ensure data is written to disk

## Edge Cases
- **Temp file cleanup**: Orphaned temp files from crashes should be cleaned
- **Concurrent access**: Multiple processes should coordinate via file locking
- **Filesystem quirks**: Some filesystems may not support atomic rename
- **Network filesystems**: Atomic operations may not be guaranteed over NFS
- **Symbolic links**: Should not follow symlinks for security