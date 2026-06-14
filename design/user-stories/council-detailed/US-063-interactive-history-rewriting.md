# User Story: US-033 - Interactive History Rewriting

## Title
Interactive Rebase and Commit Amendment

## Description
As a developer, I want to rewrite commit history interactively, so that I can clean up my commit history before sharing it with others, fix commit messages, and combine related commits.

## Priority
**Medium** - Important for code quality and collaboration

## User Story
**As a** developer
**I want to** rewrite commit history interactively
**So that** I can clean up my commit history before sharing it with others

## Acceptance Criteria

### AC1: Interactive Rebase
- [ ] Rebase onto a different commit with interactive interface
- [ ] Support squash, reword, edit, drop actions
- [ ] Preview changes before applying rebase
- [ ] Preserve commit relationships during rebase

### AC2: Commit Amendment
- [ ] Amend the most recent commit message
- [ ] Amend commit content (add/remove files)
- [ ] Preserve commit timestamp option
- [ ] Validate amended commit integrity

### AC3: Commit Reordering
- [ ] Reorder commits in history
- [ ] Drag-and-drop interface for reordering
- [ ] Validate reordering doesn't break dependencies
- [ ] Preview reordering results

### AC4: Interactive Commit Selection
- [ ] Select specific commits for rewriting
- [ ] Filter commits by message, author, date
- [ ] Bulk operations on selected commits
- [ ] Undo/redo history rewriting operations

### AC5: History Integrity
- [ ] Maintain commit hash relationships
- [ ] Track rewriting operations in RDF
- [ ] Provide recovery from rewriting mistakes
- [ ] Validate history integrity after rewriting

## Implementation Notes

### RDF Data Model for History Rewriting
```turtle
@prefix semsrc: <https://semsrc.org/ontology/> .
@prefix prov: <http://www.w3.org/ns/prov#> .

:rewrite-operation-123 a semsrc:HistoryRewrite ;
    prov:wasGeneratedBy :repository-main ;
    prov:startedAtTime "2026-06-14T11:00:00Z"^^xsd:dateTime ;
    prov:endedAtTime "2026-06-14T11:05:00Z"^^xsd:dateTime ;
    semsrc:rewriteType "interactive-rebase" ;
    semsrc:targetCommit "abc123..." ;
    semsrc:operations :operation-list-123 .

:operation-list-123 a semsrc:RewriteOperationList ;
    semsrc:containsOperation :op-reword-1, :op-squash-2 .
```

### Interface Design
```go
type HistoryRewriter interface {
    // InteractiveRebase performs interactive rebase
    InteractiveRebase(targetCommit string, operations []RebaseOperation) error
    
    // AmendCommit modifies the most recent commit
    AmendCommit(newMessage string, addToCommit bool) error
    
    // ReorderCommits changes commit order
    ReorderCommits(commits []string, newOrder []string) error
    
    // GetRebasePreview shows planned changes
    GetRebasePreview(targetCommit string, operations []RebaseOperation) (RebasePreview, error)
    
    // UndoRewrite operation
    UndoRewrite(operationID string) error
}

type RebaseOperation struct {
    CommitHash string
    Action     string  // "pick", "squash", "reword", "edit", "drop"
    NewMessage string  // For reword action
}
```

## Edge Cases

### Edge Case 1: Rewriting Public History
**Scenario**: User tries to rewrite commits that have been pushed
**Handling**: Warning message, suggest force push alternatives

### Edge Case 2: Circular Dependencies
**Scenario**: Reordering creates commit dependency cycles
**Handling**: Detect and prevent, provide clear error messages

### Edge Case 3: Large History Rewrites
**Scenario**: Rewriting hundreds of commits
**Handling**: Progress indicators, batch processing, validation checkpoints

### Edge Case 4: Conflicting Operations
**Scenario**: Multiple rewrite operations conflict
**Handling**: Conflict detection, resolution suggestions

## Dependencies

### US-017: Commit Changes Message
- Required for commit amendment functionality

### US-018: View Commit History
- Required for understanding commit structure

### US-021: Automatic PROV-O Relationship Generation
- Required for tracking rewrite operations

## Test Cases

### Test Case 1: Interactive Rebase
**Input**: `rebase -i HEAD~3`
**Expected**: Interactive interface for selecting rebase actions

### Test Case 2: Amend Commit Message
**Input**: `commit --amend -m "New message"`
**Expected**: Most recent commit message updated

### Test Case 3: Reorder Commits
**Input**: `reorder commit1 commit2 commit3`
**Expected**: Commits reordered, history updated

## Related User Stories
- **US-017: Commit Changes Message** - Related commit operations
- **US-018: View Commit History** - Related history viewing
- **US-034: Reflog and Recovery** - Related recovery functionality

## Implementation Priority
**Medium** - Important for code quality but can follow core functionality.