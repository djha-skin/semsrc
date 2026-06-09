# US-030: Full-Text Search Across Blobs and Metadata

## Description
Implement full-text search capabilities across file content (blobs) and semantic metadata. This enables users to find content and relationships quickly using keyword searches.


## Test Cases
1. **Content search**: Search file content by keywords
2. **Metadata search**: Search semantic metadata by keywords
3. **Combined search**: Search both content and metadata simultaneously
4. **Ranking**: Results ranked by relevance
5. **Large dataset**: Search performance scales with repository size

## Dependencies
- US-001: Store Files with SHA-256 Content Addressing
- US-009: Store RDF Triples in SQLite

## Implementation Notes
- Use full-text search index (e.g., SQLite FTS5)
- Index both file content and semantic metadata
- Implement efficient query processing
- Support search ranking and relevance scoring
- Provide search filters (file type, date range, etc.)
- Support search across multiple repositories

## Edge Cases
- **Large files**: Should index large files efficiently
- **Binary files**: Should skip or handle binary files appropriately
- **Special characters**: Should handle special characters in search queries
- **Multiple languages**: Should support multiple languages
- **Search performance**: Should maintain performance with large datasets