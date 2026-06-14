# User Story: US-034 - Reflog and Recovery

## Title
Recover Lost Commits and Worktree States

## Description
As a developer, I want to view a log of all state changes and recover lost commits, so that I can recover from mistakes and never lose work due to accidental deletion or rebase operations.

## Priority
**Medium** - Safety net for developer mistakes

## User Story
**As a** developer
**I want to** view a log of all state changes and recover lost commits
**So that** I can recover from mistakes and never lose work

## Acceptance Criteria

### AC1: Reflog Viewing
- [ ] View log of all state changes (commits, branches, merges)
- [ ] Filter reflog by operation type
- [ ] Search reflog by commit hash or message
- [ ] Display reflog entries with timestamps and operations

### AC2: Commit Recovery
- [ ] Recover lost commits by reference
- [ ] Restore commit to branch or create new branch
- [ ] Validate recovered commit integrity
- [ ] Preserve commit metadata during recovery

### AC3: Worktree State Recovery
- [ ] Restore previous worktree states
- [ ] Recover from failed merge operations
- [ ] Restore branch states from reflog
- [ ] Preview recovery results before applying

### AC4: Audit Trail
- [ ] Maintain comprehensive operation history
- [ ] Track all repository modifications
- [ ] Provide searchable audit log
- [ ] Export audit trail for compliance

### AC5: Recovery Safety
- [ ] Prevent accidental overwrites during recovery
- [ ] Provide confirmation prompts for dangerous operations
- [ ] Maintain backup of original state
- [ ] Support undo of recovery operations

## Implementation Notes

### RDF Data Model for Reflog
```turtle
@prefix semsrc: <https://semsrc.org/ontology/> .
@prefix prov: <http://www.w3.org/ns/prov#> .

:reflog-entry-123 a semsrc:ReflogEntry ;
    prov:wasGeneratedBy :repository-main ;
    prov:startedAtTime "2026-06-14T12:00:00Z"^^xsd:dateTime ;
    semsrc:operationType "commit" ;
    semsrc:oldReference "refs/heads/main" ;
    semsrc:newReference "abc123..." ;
    semsrc:operationDetails :operation-details-123 .

:operation-details-123 a semsrc:OperationDetails ;
    semsrc:command "git commit -m 'Add feature'" ;
    semsrc:workingDirectory "/home/user/project" ;
    semsrc:exitCode 0 .
```

### SPARQL Queries Required
```sparql
# Get reflog for specific branch
PREFIX semsrc: <https://semsrc.org/ontology/>

SELECT ?entry ?timestamp ?operation ?oldRef ?newRef
WHERE {
    ?entry a semsrc:ReflogEntry ;
        prov:startedAtTime ?timestamp ;
        semsrc:operationType ?operation ;
        semsrc:oldReference ?oldRef ;
        semsrc:newReference ?newRef .
    FILTER(?oldRef = "refs/heads/main" || ?newRef = "refs/heads/main")
}
ORDER BY DESC(?timestamp)
```

### Interface Design
```go
type ReflogManager interface {
    // GetReflog returns reflog entries
    GetReflog(reference string, limit int) ([]ReflogEntry, error)
    
    // RecoverCommit restores a lost commit
    RecoverCommit(commitHash string, targetBranch string) error
    
    // RestoreWorktreeState restores previous worktree state
    RestoreWorktreeState(stateID string) error
    
    // SearchReflog searches reflog by criteria
    SearchReflog(criteria SearchCriteria) ([]ReflogEntry, error)
    
    // ExportAuditTrail exports operation history
    ExportAuditTrail(format string) ([]byte, error)
}

type ReflogEntry struct {
    Reference   string
    Operation   string
    OldHash     string
    NewHash     string
    Timestamp   time.Time
    Command     string
    WorkingDir  string
}
```

## Edge Cases

### Edge Case 1: Reflog Size Management
**Scenario**: Reflog grows very large over time
**Handling**: Automatic pruning, configurable retention periods

### Edge Case 2: Corrupted Reflog Entries
**Scenario**: Reflog contains inconsistent data
**Handling**: Validation on access, repair functionality

### Edge Case 3: Recovering Deleted Branches
**Scenario**: Branch was deleted but commits still exist
**Handling**: Scan reflog for branch deletion, offer recovery

### Edge Case 4: Conflicting Recovery Operations
**Scenario**: Multiple recovery operations conflict
**Handling**: Sequential processing, conflict detection

## Dependencies

### US-018: View Commit History
- Required for understanding commit relationships

### US-020: Checkout Different Commits
- Required for restoring commit states

### US-021: Automatic PROV-O Relationship Generation
- Required for tracking operation history

## Test Cases

### Test Case 1: View Reflog
**Input**: `reflog`
**Expected**: Display recent reflog entries with timestamps

### Test Case 2: Recover Lost Commit
**Input**: `recover-commit abc123...`
**Expected**: Commit restored to specified branch

### Test Case 3: Search Reflog
**Input**: `reflog --search "Add feature"`
**Expected**: Reflog entries matching search criteria

## Related User Stories
- **US-018: View Commit History** - Related history viewing
- **US-020: Checkout Different Commits** - Related commit operations
- **US-033: Interactive History Rewriting** - Related history operations

## Implementation Priority
**Medium** - Important safety feature but can follow core functionality.