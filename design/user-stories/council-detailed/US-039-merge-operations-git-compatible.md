# US-039: Merge Operations Compatible with Git

## Description
Implement merge operations that are compatible with Git semantics, allowing
teams to continue using familiar merge workflows while leveraging semantic
features.


## Test Cases
1. **Fast-forward merge**: Merge with fast-forward when possible
2. **Non-fast-forward merge**: Create merge commit explicitly
3. **Squash merge**: Combine multiple commits into one
4. **Conflict resolution**: Handle merge conflicts appropriately
5. **Merge metadata**: Record merge information semantically

## Dependencies
- US-038: Branch Management Compatible with Git
- US-017: Commit Changes with a Message

## Implementation Notes
- Support standard Git merge strategies
- Handle merge conflicts gracefully
- Record merge metadata in semantic graph
- Support three-way merge algorithm
- Provide merge conflict resolution tools
- Maintain merge history for semantic queries

## Edge Cases
- **No changes to merge**: Handle already-merged branches
- **Conflicting changes**: Provide clear conflict resolution guidance
- **Large merges**: Handle merges with many changed files efficiently
- **Binary files**: Handle binary file merges appropriately
- **Merge strategies**: Support different merge strategies (recursive,
   octopus, etc.)