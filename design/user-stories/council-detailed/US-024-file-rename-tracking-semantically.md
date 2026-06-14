# US-024: File Rename Tracking Semantically

## Description
Track file renames semantically to maintain coherent history across file name
changes. This ensures that operations like `git log --follow` work correctly
when files are renamed.


## Test Cases
1. **Basic rename**: Track a simple file rename
2. **Multiple renames**: Handle multiple files renamed in one commit
3. **Rename chain**: Handle file renamed multiple times across commits
4. **Rename detection**: Automatically detect renames based on content
   similarity
5. **History query**: Query complete history across renames

## Dependencies
- US-009: Store RDF Triples in SQLite
- US-017: Commit Changes with a Message
- US-019: View Differences Between Commits

## Implementation Notes
- Store rename relationships in semantic graph
- Use content similarity for automatic rename detection
- Track rename chains across multiple commits
- Support rename queries for historical tracking
- Include metadata (rename timestamp, author, etc.)
- Handle partial renames (content changes + name change)

## Edge Cases
- **Content changes**: Distinguish renames from content modifications
- **Concurrent renames**: Handle race conditions in rename detection
- **Large file renames**: Handle renames of large files efficiently
- **Cross-directory renames**: Track renames across different directories
- **Rename undo**: Handle renames that are later undone