# User Story: US-036 - Tag Management for Semantic Versioning

## Title
Create and Manage Semantic Version Tags

## Description
As a developer, I want to create and manage tags with semantic version information, so that I can mark specific commits for releases and provide meaningful version identifiers.

## Priority
**Medium** - Important for release management

## User Story
**As a** developer
**I want to** create and manage semantic version tags
**So that** I can mark specific commits for releases and provide meaningful version identifiers

## Acceptance Criteria

### AC1: Create Tags
- [ ] Create lightweight tags (just a name)
- [ ] Create annotated tags with message and metadata
- [ ] Support semantic versioning format (MAJOR.MINOR.PATCH)
- [ ] Validate tag names against version format

### AC2: Tag Management
- [ ] List all tags with filtering options
- [ ] View tag details (message, creator, date)
- [ ] Delete tags
- [ ] Rename tags (with caution)

### AC3: Tag Operations
- [ ] Push tags to remote repositories
- [ ] Pull tags from remote repositories
- [ ] Checkout commits by tag
- [ ] Compare tags (differences between tagged commits)

### AC4: Semantic Versioning Support
- [ ] Enforce semantic versioning rules
- [ ] Automatic version increment suggestions
- [ ] Version constraint checking
- [ ] Changelog generation from tags

### AC5: Tag Metadata in RDF
- [ ] Store tag information in RDF format
- [ ] Link tags to commit entities
- [ ] Support querying tags via SPARQL
- [ ] Track tag creation and modification history

## Implementation Notes

### RDF Data Model for Tags
```turtle
@prefix semsrc: <https://semsrc.org/ontology/> .
@prefix prov: <http://www.w3.org/ns/prov#> .

:tag-v1.0.0 a semsrc:Tag ;
    prov:wasGeneratedBy :repository-main ;
    prov:wasGeneratedBy :user-developer ;
    prov:startedAtTime "2026-06-14T14:00:00Z"^^xsd:dateTime ;
    semsrc:tagName "v1.0.0" ;
    semsrc:tagType "annotated" ;
    semsrc:targetCommit "abc123..." ;
    semsrc:tagMessage "Initial release" ;
    semsrc:versionMajor 1 ;
    semsrc:versionMinor 0 ;
    semsrc:versionPatch 0 .
```

### SPARQL Queries Required
```sparql
# List all tags with version information
PREFIX semsrc: <https://semsrc.org/ontology/>

SELECT ?tag ?name ?commit ?message ?date
WHERE {
    ?tag a semsrc:Tag ;
        semsrc:tagName ?name ;
        semsrc:targetCommit ?commit ;
        semsrc:tagMessage ?message ;
        prov:startedAtTime ?date .
}
ORDER BY DESC(?date)
```

### Interface Design
```go
type TagManager interface {
    // CreateLightweightTag creates a simple tag
    CreateLightweightTag(name string, commitHash string) error
    
    // CreateAnnotatedTag creates tag with message and metadata
    CreateAnnotatedTag(name string, commitHash string, message string) error
    
    // ListTags returns all tags with optional filtering
    ListTags(filter TagFilter) ([]Tag, error)
    
    // GetTagDetails returns detailed tag information
    GetTagDetails(name string) (Tag, error)
    
    // DeleteTag removes a tag
    DeleteTag(name string) error
    
    // PushTags pushes tags to remote
    PushTags(remote string, tags []string) error
}

type Tag struct {
    Name        string
    Type        string  // "lightweight", "annotated"
    TargetHash  string
    Message     string
    CreatedAt   time.Time
    Version     *SemanticVersion
}

type SemanticVersion struct {
    Major int
    Minor int
    Patch int
}
```

### Semantic Version Validation
```go
func ValidateSemanticVersion(version string) error {
    // Validate MAJOR.MINOR.PATCH format
    // Check for valid numeric components
    // Ensure no leading zeros (except 0 itself)
    return nil
}
```

## Edge Cases

### Edge Case 1: Duplicate Tag Names
**Scenario**: User tries to create tag with existing name
**Handling**: Error message, suggest alternative name

### Edge Case 2: Invalid Version Format
**Scenario**: Tag name doesn't follow semantic versioning
**Handling**: Validation error with clear guidance

### Edge Case 3: Tagging Non-Commit Objects
**Scenario**: User tries to tag something that's not a commit
**Handling**: Restrict tagging to commits only

### Edge Case 4: Remote Tag Conflicts
**Scenario**: Local and remote tags have same name but different targets
**Handling**: Conflict resolution strategy, user notification

## Dependencies

### US-017: Commit Changes Message
- Required for understanding commit structure

### US-018: View Commit History
- Required for understanding commit relationships

### US-033: Interactive History Rewriting
- May affect tagged commits

## Test Cases

### Test Case 1: Create Semantic Version Tag
**Input**: `tag create v1.0.0 abc123... "Initial release"`
**Expected**: Tag created with semantic version metadata

### Test Case 2: List Tags
**Input**: `tag list`
**Expected**: Display all tags with version information

### Test Case 3: Validate Version Format
**Input**: `tag create invalid-version abc123...`
**Expected**: Error with semantic versioning guidance

### Test Case 4: Push Tags to Remote
**Input**: `tag push origin v1.0.0`
**Expected**: Tag pushed to remote repository

## Related User Stories
- **US-017: Commit Changes Message** - Related commit operations
- **US-018: View Commit History** - Related history viewing
- **US-033: Interactive History Rewriting** - May affect tagged commits

## Implementation Priority
**Medium** - Important for release management but can follow core functionality.