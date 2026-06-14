# User Story: US-031 - Manage Multiple Working Trees

## Title
Manage Multiple Working Trees for Parallel Development

## Description
As a developer working on multiple features simultaneously, I want to create and manage multiple working trees from the same repository, so that I can work on different branches in parallel without switching contexts or contaminating my workspace.

## Priority
**High** - Critical for parallel development workflows

## User Story
**As a** developer
**I want to** create, list, and remove working trees for parallel development
**So that** I can work on multiple features simultaneously without switching branches

## Acceptance Criteria

### AC1: Create New Worktree
- [ ] Create a new worktree from an existing commit/branch
- [ ] Specify a unique path for the new worktree
- [ ] Automatically track the relationship to the original repository
- [ ] Preserve the branch/commit state in the new worktree

### AC2: List Active Worktrees
- [ ] List all active worktrees with their paths
- [ ] Show associated branch/commit for each worktree
- [ ] Display worktree metadata (creation time, last activity)
- [ ] Filter worktrees by branch or commit

### AC3: Remove Worktree
- [ ] Remove a worktree and clean up associated resources
- [ ] Prevent removal of the primary worktree
- [ ] Validate worktree exists before removal
- [ ] Clean up any temporary files associated with worktree

### AC4: Worktree Conflict Prevention
- [ ] Prevent creating worktrees with duplicate paths
- [ ] Prevent creating multiple worktrees for the same branch
- [ ] Provide clear error messages for conflict scenarios
- [ ] Allow force removal if conflicts exist

### AC5: Worktree State Tracking
- [ ] Track worktree metadata in RDF using PROV-O
- [ ] Record relationships between worktrees and repositories
- [ ] Support querying worktree information via SPARQL
- [ ] Maintain worktree history for audit purposes

## Implementation Notes

### RDF Data Model
```turtle
@prefix semsrc: <https://semsrc.org/ontology/> .
@prefix prov: <http://www.w3.org/ns/prov#> .

:worktree-123 a semsrc:WorkingTree ;
    prov:wasGeneratedBy :repository-main ;
    prov:generatedAtTime "2026-06-14T10:00:00Z"^^xsd:dateTime ;
    semsrc:worktreePath "/home/user/projects/feature-branch" ;
    semsrc:associatedBranch "feature/new-feature" ;
    semsrc:basedOnCommit "abc123..." .
```

### SPARQL Queries Required
```sparql
# List all worktrees for a repository
PREFIX semsrc: <https://semsrc.org/ontology/>
PREFIX prov: <http://www.w3.org/ns/prov#>

SELECT ?worktree ?path ?branch ?commit
WHERE {
    ?worktree a semsrc:WorkingTree ;
        prov:wasGeneratedBy ?repo ;
        semsrc:worktreePath ?path ;
        semsrc:associatedBranch ?branch ;
        semsrc:basedOnCommit ?commit .
    FILTER(?repo = :repository-main)
}
```

### Interface Design
```go
type WorktreeManager interface {
    // CreateWorktree creates a new worktree from a commit/branch
    CreateWorktree(commitHash string, path string) (Worktree, error)
    
    // ListWorktrees returns all active worktrees
    ListWorktrees() ([]Worktree, error)
    
    // RemoveWorktree deletes a worktree and cleans up resources
    RemoveWorktree(path string) error
    
    // GetWorktree retrieves worktree by path
    GetWorktree(path string) (Worktree, error)
}

type Worktree struct {
    Path          string
    Branch        string
    CommitHash    string
    CreatedAt     time.Time
    LastActivity  time.Time
    RepositoryURL string
}
```

### Integration with Existing Systems
1. **Repository Pattern**: Worktrees stored in repository RDF graph
2. **PROV-O Relationships**: Track creation and modification history
3. **SPARQL Endpoint**: Query worktree information
4. **Git Compatibility**: Mirror Git worktree command interface

## Edge Cases

### Edge Case 1: Worktree Path Conflicts
**Scenario**: User tries to create worktree with existing path
**Handling**: Return error with clear message suggesting alternative path

### Edge Case 2: Branch Already Has Worktree
**Scenario**: User tries to create second worktree for same branch
**Handling**: Prevent creation, suggest using existing worktree or specifying different commit

### Edge Case 3: Primary Worktree Removal
**Scenario**: User tries to remove the main repository worktree
**Handling**: Prevent removal with clear error message

### Edge Case 4: Corrupted Worktree Metadata
**Scenario**: RDF store contains inconsistent worktree metadata
**Handling**: Validate worktree state on access, provide repair functionality

### Edge Case 5: Large Number of Worktrees
**Scenario**: User creates dozens of worktrees
**Handling**: Efficient SPARQL queries, pagination support, cleanup suggestions

## Dependencies

### US-015: Initialize New Repository
- Required for understanding repository structure

### US-021: Automatic PROV-O Relationship Generation
- Required for tracking worktree provenance

### US-027: SPARQL Endpoint Semantic Queries
- Required for querying worktree information

### US-032: Staging Area Management (Future)
- Related functionality for complete Git workflow

## Test Cases

### Test Case 1: Create Worktree from Branch
**Input**: `create-worktree --branch feature/new-feature --path /tmp/feature-worktree`
**Expected**: New worktree created, RDF metadata stored, relationship recorded

### Test Case 2: List All Worktrees
**Input**: `list-worktrees`
**Expected**: Returns list of all active worktrees with metadata

### Test Case 3: Remove Worktree
**Input**: `remove-worktree --path /tmp/feature-worktree`
**Expected**: Worktree removed, resources cleaned up, metadata updated

### Test Case 4: Conflict Prevention
**Input**: `create-worktree --branch feature/new-feature --path /tmp/existing-worktree`
**Expected**: Error returned with clear conflict message

### Test Case 5: Query via SPARQL
**Input**: SPARQL query for worktree information
**Expected**: Returns worktree data in RDF format

## Related User Stories
- **US-038: Branch Management Git Compatible** - Related branch operations
- **US-039: Merge Operations Git Compatible** - Related merge operations
- **US-040: Import Existing Git Repositories** - Related repository operations

## Implementation Priority
**High** - Should be implemented early in the development process as it's fundamental to parallel development workflows.

## Council Review
This user story was identified by Junio C Hamano (Git Expert) as a critical gap in the current user story set. It addresses the lack of worktree management functionality in the Git competitor implementation.