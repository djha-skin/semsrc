# US-001: Store Files with SHA-256 Content Addressing

## Description

Store files using content-addressable hashing with SHA-256. This ensures that
identical content is automatically deduplicated, reducing storage overhead and
improving cache efficiency. When storing a file, the system calculates its
SHA-256 hash and uses this hash as the file's identifier in storage.

## Test Cases

1. **Identical content deduplication**: Two files with identical content
   should be stored once with the same hash
2. **Unique content creates unique hashes**: Files with different content
   should have different hashes
3. **Empty file handling**: Empty files should generate a valid SHA-256 hash
4. **Large file support**: Files larger than available memory should be hashed
   correctly using streaming
5. **Binary file support**: Binary files should be hashed correctly without
   encoding issues

## Dependencies

- None (foundational feature)

## Implementation Notes

- Use SHA-256 hashing algorithm (standard, secure, collision-resistant)
- Hash format: hexadecimal string representation (64 characters)
- Storage path pattern: `objects/sha256/ab/cdef1234...` (first 2 chars
  as directory, remainder as filename)
- Implement streaming hash calculation to handle large files without loading
  entire file into memory
- Store hash metadata alongside file content for quick verification
- Hash calculation should be deterministic across platforms

## Edge Cases

- **Hash collisions**: Extremely unlikely with SHA-256, but system should
   handle gracefully
- **Unicode filenames**: Hash is based on content, not filename
- **File modification time**: Should not affect hash calculation
- **Permission bits**: Should not affect hash calculation (content only)
- **Line endings**: Should not affect hash (raw bytes are hashed)
