# US-020: Checkout Different Commits

## Description
Checkout a different commit to explore the repository at a specific point in history. This allows users to view or work with an older version of the repository.


## Test Cases
1. **Checkout specific commit**: Update working directory to specific commit
2. **Checkout branch**: Switch to a different branch
3. **Checkout tag**: Switch to a tagged commit
4. **Checkout with changes**: Handle uncommitted changes appropriately
5. **Checkout previous commit**: Navigate back in history

## Dependencies
- US-017: Commit Changes with a Message
- US-018: View Commit History
- US-001: Store Files with SHA-256 Content Addressing

## Implementation Notes
- Update working directory to match commit tree
- Update HEAD reference to point to checked-out commit
- Handle uncommitted changes (error, stash, or discard)
- Support checkout by commit ID, branch, or tag
- Validate commit exists before checkout
- Maintain repository state consistency

## Edge Cases
- **Uncommitted changes**: Should warn or error when changes would be lost
- **Checkout to same commit**: Should handle gracefully
- **Corrupted commit**: Should detect and report invalid commits
- **Large checkout**: Should handle checkouts with many files efficiently
- **Permission issues**: Should handle file permission changes