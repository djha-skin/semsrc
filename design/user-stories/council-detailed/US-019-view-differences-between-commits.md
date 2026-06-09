# US-019: View Differences Between Commits

## Description
View the differences (diff) between two commits to understand what changed. This helps users review changes, understand the evolution of the codebase, and identify potential issues.


## Test Cases
1. **Basic diff**: View changes between two commits
2. **Diff with no changes**: Handle commits with no differences
3. **Large diff**: Handle diffs with many changed files
4. **Binary files**: Handle binary file differences appropriately
5. **Renamed files**: Detect and show file renames

## Dependencies
- US-017: Commit Changes with a Message
- US-001: Store Files with SHA-256 Content Addressing

## Implementation Notes
- Compare tree objects between commits
- Use efficient diff algorithms (e.g., Myers diff)
- Support different diff formats (unified, context, side-by-side)
- Detect file renames by content similarity
- Handle binary files (show "binary files differ")
- Provide statistics (files changed, insertions, deletions)

## Edge Cases
- **Identical commits**: Should show no differences
- **Very large diffs**: Should handle memory constraints
- **Binary file comparison**: Should not attempt text diff on binary files
- **Missing files**: Handle files deleted between commits
- **Encoding issues**: Handle different file encodings gracefully