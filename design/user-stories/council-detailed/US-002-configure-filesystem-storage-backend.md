# US-002: Configure Filesystem Storage Backend

## Description

Configure the object store to use a local filesystem backend. This allows
storing object data directly on the local disk with optional compression for
space efficiency. The backend should be easily configurable.

## Test Cases

1. **Basic storage and retrieval**: Store an object and retrieve it
   successfully
2. **Compression verification**: Compressed objects should be transparently
   decompressed on retrieval
3. **Concurrent access**: Multiple processes should be able to access the
   backend safely
4. **Disk space management**: Backend should handle disk space exhaustion
   gracefully

## Dependencies

- US-001: Store Files with SHA-256 Content Addressing

## Implementation Notes

- Support multiple compression algorithms (gzip, zstd, none)
- Use file locking for concurrent access safety
- Store object metadata (size, compression, hash) in separate index file
- Implement garbage collection for orphaned objects

## Edge Cases

- **Disk full**: Should fail gracefully with clear error message
- **Permission errors**: Should handle read/write permission issues
- **Corrupted objects**: Should detect and handle corrupted storage
- **Concurrent writes**: Should handle race conditions safely
- **Symlink attacks**: Should validate paths to prevent security issues
