# US-018: View Commit History

## Description

View the commit history of the repository, showing a log of commits with their
messages, authors, and timestamps. This allows users to track project changes
over time.


## Test Cases

1. **Basic history**: View recent commits
2. **Limit results**: Limit number of commits shown
3. **Pagination**: Navigate through commit history
4. **Filter by author**: Show commits by specific author
5. **Filter by date**: Show commits within date range

## Dependencies

- US-017: Commit Changes with a Message
- US-009: Store RDF Triples in Native RDF Store

## Implementation Notes

- Query commits from triple store using SPARQL
- Support pagination for large histories
- Implement filtering by author, date, message pattern
- Support different output formats (oneline, full, email)
- Provide reverse chronological ordering
- Include commit metadata in query results

## Edge Cases

- **Empty repository**: Should handle no commits gracefully
- **Large history**: Should handle repositories with thousands of commits
- **Filter with no results**: Should return empty list, not error
- **Invalid date range**: Should handle invalid date formats
- **Concurrent modification**: Should handle history changes during query