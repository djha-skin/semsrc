# User Story: US-032 - Staging Area Management

## Title
Stage and Unstage Changes Before Commit

## Description
As a developer, I want to stage and unstage changes selectively before committing, so that I can control exactly which changes are included in each commit and maintain clean, focused commit histories.

## Priority
**High** - Fundamental Git workflow requirement

## User Story
**As a** developer
**I want to** stage and unstage changes before committing
**So that** I can create clean, focused commits with only relevant changes

## Acceptance Criteria

### AC1: Stage Individual Files
- [ ] Stage specific files or directories
- [ ] Support staging multiple files at once
- [ ] Preserve file modification states
- [ ] Update staging area metadata in RDF

### AC2: Unstage Specific Changes
- [ ] Unstage individual files from staging area
- [ ] Support unstage multiple files at once
- [ ] Maintain file state history
- [ ] Update staging area metadata in RDF

### AC3: View Staged vs Unstaged Differences
- [ ] Show differences between working directory and staging area
- [ ] Show differences between staging area and last commit
- [ ] Support viewing diffs for individual files
- [ ] Color-coded diff output for readability

### AC4: Interactive Staging
- [ ] Stage individual hunks of changes
- [ ] Stage individual lines of changes
- [ ] Interactive hunk selection interface
- [ ] Preview changes before staging

### AC5: Staging Area Persistence
- [ ] Store staging area state in RDF
- [ ] Preserve staging across process restarts
- [ ] Sync staging area with working directory changes
- [ ] Support partial staging scenarios

## Implementation Notes

### RDF Data Model for Staging Area
```turtle
@prefix semsrc: <https://semsrc.org/ontology/> .
@prefix prov: <http://www.w3.org/ns/prov#> .

:staging-area-123 a semsrc:StagingArea ;
    prov:wasGeneratedBy :repository-main ;
    semsrc:stagingState "active" ;
    semsrc:stagedFiles :file-list-123 .

:file-list-123 a semsrc:FileList ;
    semsrc:containsFile :file-1, :file-2, :file-3 .

:file-1 a semsrc:StagedFile ;
    semsrc:filePath "src/main.go" ;
    semsrc:stagedAt "2026-06-14T10:30:00Z"^^xsd:dateTime ;
    semsrc:stagingState "full" ;
    semsrc:originalHash "abc123..." ;
    semsrc:stagedHash "def456..." .
```

### SPARQL Queries Required
```sparql
# List all staged files
PREFIX semsrc: <https://semsrc.org/ontology/>

SELECT ?filePath ?stagingState ?stagedAt
WHERE {
    ?file a semsrc:StagedFile ;
        semsrc:filePath ?filePath ;
        semsrc:stagingState ?stagingState ;
        semsrc:stagedAt ?stagedAt .
}
```

### Interface Design
```go
type StagingArea interface {
    // StageFile stages a single file
    StageFile(path string) error
    
    // StageDirectory stages all files in a directory
    StageDirectory(path string) error
    
    // UnstageFile removes a file from staging area
    UnstageFile(path string) error
    
    // GetStagedFiles returns list of staged files
    GetStagedFiles() ([]StagedFile, error)
    
    // GetDiffWorkingToStaging shows diff between working and staged
    GetDiffWorkingToStaging(path string) (string, error)
    
    // GetDiffStagingToCommit shows diff between staged and commit
    GetDiffStagingToCommit(path string) (string, error)
    
    // InteractiveStage allows hunk-by-hunk staging
    InteractiveStage(path string) error
}

type StagedFile struct {
    Path          string
    StagingState  string  // "full", "partial", "hunk"
    OriginalHash  string
    StagedHash    string
    StagedAt      time.Time
    HunkSelection []HunkRange
}
```

### Integration with Existing Systems
1. **Repository Pattern**: Staging area tied to repository
2. **Object Store**: Staged files stored in object store
3. **SPARQL Endpoint**: Query staging area state
4. **Git Compatibility**: Mirror Git staging workflow

## Edge Cases

### Edge Case 1: File Modified After Staging
**Scenario**: User stages file, then modifies it again
**Handling**: Show both staged and unstaged versions, allow restaging

### Edge Case 2: Large File Staging
**Scenario**: User stages very large files
**Handling**: Efficient diff computation, progress indicators

### Edge Case 3: Conflicting Staging States
**Scenario**: File marked as both staged and deleted
**Handling**: Clear conflict resolution, user notification

### Edge Case 4: Partial Staging Complexities
**Scenario**: User stages only some hunks of a file
**Handling**: Track hunk ranges precisely, validate on commit

### Edge Case 5: Staging Across Repositories
**Scenario**: User works with multiple repositories
**Handling**: Separate staging areas per repository

## Dependencies

### US-001: Store Files with SHA256 Hashing
- Required for file hashing in staging area

### US-004: Atomic Writes with Rollback
- Required for staging area integrity

### US-016: Add Files Staging Area
- Related functionality (may be duplicate or extension)

### US-019: View Differences Between Commits
- Related diff viewing functionality

## Test Cases

### Test Case 1: Stage Single File
**Input**: `stage src/main.go`
**Expected**: File added to staging area, RDF metadata updated

### Test Case 2: Stage Directory
**Input**: `stage src/`
**Expected**: All files in directory staged, metadata for each file

### Test Case 3: Unstage File
**Input**: `unstage src/main.go`
**Expected**: File removed from staging area, metadata updated

### Test Case 4: View Working vs Staged Diff
**Input**: `diff --working --staged src/main.go`
**Expected**: Shows differences between working directory and staged version

### Test Case 5: Interactive Hunk Staging
**Input**: `stage -i src/main.go`
**Expected**: Interactive interface for selecting hunks to stage

## Related User Stories
- **US-016: Add Files Staging Area** - Basic staging functionality (US-032 extends with advanced features)
- **US-019: View Differences Between Commits** - Related diff functionality
- **US-033: Staging Area Persistence** - Related storage functionality

## Note: Relationship to US-016
US-032 extends US-016 with advanced staging features:
- **US-016**: Basic file/directory staging
- **US-032**: Interactive hunk staging, advanced diff viewing, staging area persistence

## Implementation Priority
**High** - Fundamental to Git workflow, should be implemented early.

## Council Review
This user story was identified by Junio C Hamano (Git Expert) as a critical gap in the current user story set. It addresses the lack of staging area functionality in the Git competitor implementation.