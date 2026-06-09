# US-002: Configure Filesystem Storage Backend

## Description
Configure the object store to use a local filesystem backend. This allows storing object data directly on the local disk with optional compression for space efficiency. The backend should be easily configurable and support different directory layouts for organizing stored objects.


## Test Cases
1. **Basic storage and retrieval**: Store an object and retrieve it successfully
2. **Compression verification**: Compressed objects should be transparently decompressed on retrieval
3. **Directory layout**: Objects should be organized according to configured layout
4. **Concurrent access**: Multiple processes should be able to access the backend safely
5. **Disk space management**: Backend should handle disk space exhaustion gracefully

## Dependencies
- US-001: Store Files with SHA-256 Content Addressing

## Implementation Notes
- Support multiple compression algorithms (gzip, zstd, none)
- Use file locking for concurrent access safety
- Implement atomic write operations (write to temp file, then rename)
- Support different directory layouts:
  - `flat`: All objects in single directory (simple but may hit filesystem limits)
  - `hierarchical`: Based on hash prefixes (e.g., `ab/cdef...`)
  - `content-hash`: Full hash-based path
- Store object metadata (size, compression, hash) in separate index file
- Implement garbage collection for orphaned objects

## Edge Cases
- **Disk full**: Should fail gracefully with clear error message
- **Permission errors**: Should handle read/write permission issues
- **Corrupted objects**: Should detect and handle corrupted storage
- **Concurrent writes**: Should handle race conditions safely
- **Symlink attacks**: Should validate paths to prevent security issues