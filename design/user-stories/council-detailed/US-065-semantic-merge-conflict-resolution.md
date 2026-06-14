# User Story: US-035 - Semantic Merge Conflict Resolution

## Title
Resolve Merge Conflicts with Semantic Understanding

## Description
As a developer working with semantic data, I want intelligent merge conflict resolution that understands RDF structure, so that I can resolve conflicts in semantic triples rather than just line-by-line text conflicts.

## Priority
**High** - Critical for semantic data workflows

## User Story
**As a** developer working with semantic data
**I want to** resolve merge conflicts with semantic understanding
**So that** I can properly handle RDF triple conflicts rather than text conflicts

## Acceptance Criteria

### AC1: Detect RDF Conflicts
- [ ] Identify conflicts in RDF triples rather than text lines
- [ ] Detect conflicting additions/deletions of triples
- [ ] Identify property conflicts on same subject-predicate pairs
- [ ] Classify conflicts by severity and type

### AC2: Semantic Merge Suggestions
- [ ] Provide intelligent merge suggestions based on RDF structure
- [ ] Suggest conflict resolution based on ontology rules
- [ ] Offer automatic merging when safe
- [ ] Explain merge suggestions with reasoning

### AC3: Manual Conflict Resolution
- [ ] Allow manual selection of conflicting triples
- [ ] Support triple-level merge operations
- [ ] Provide preview of merge results
- [ ] Validate merge results against ontology

### AC4: Merge Validation
- [ ] Validate merged RDF against SHACL constraints
- [ ] Check for consistency with existing triples
- [ ] Ensure no orphaned nodes or properties
- [ ] Provide validation report

### AC5: Merge History Tracking
- [ ] Track merge operations in RDF
- [ ] Record conflict resolution decisions
- [ ] Support undo/redo of merge operations
- [ ] Provide audit trail of merge conflicts

## Implementation Notes

### RDF Data Model for Merge Conflicts
```turtle
@prefix semsrc: <https://semsrc.org/ontology/> .
@prefix prov: <http://www.w3.org/ns/prov#> .

:merge-operation-123 a semsrc:MergeOperation ;
    prov:wasGeneratedBy :repository-main ;
    prov:startedAtTime "2026-06-14T13:00:00Z"^^xsd:dateTime ;
    prov:endedAtTime "2026-06-14T13:05:00Z"^^xsd:dateTime ;
    semsrc:sourceBranch "feature/new-feature" ;
    semsrc:targetBranch "main" ;
    semsrc:conflictCount 3 ;
    semsrc:resolutionStrategy "semantic-merge" .

:conflict-123 a semsrc:MergeConflict ;
    semsrc:conflictType "triple-conflict" ;
    semsrc:conflictingTriple :triple-123 ;
    semsrc:resolution :resolution-123 .
```

### Interface Design
```go
type SemanticMergeManager interface {
    // DetectConflicts identifies RDF-specific conflicts
    DetectConflicts(sourceBranch, targetBranch string) ([]MergeConflict, error)
    
    // SuggestResolution provides intelligent merge suggestions
    SuggestResolution(conflict MergeConflict) ([]ResolutionSuggestion, error)
    
    // ApplyResolution applies selected resolution
    ApplyResolution(conflictID string, resolution Resolution) error
    
    // ValidateMerge validates merged RDF
    ValidateMerge(mergedGraph RDFGraph) (ValidationResult, error)
    
    // GetMergePreview shows merge results before applying
    GetMergePreview(sourceBranch, targetBranch string) (MergePreview, error)
}

type MergeConflict struct {
    ID            string
    Type          string  // "triple-conflict", "property-conflict"
    Triple        RDFTriple
    SourceValue   interface
    TargetValue   interface
    Severity      string  // "critical", "warning", "info"
}

type ResolutionSuggestion struct {
    ConflictID    string
    Strategy      string  // "choose-source", "choose-target", "merge"
    Confidence    float64
    Explanation   string
    ApplyFunction func() error
}
```

### SPARQL for Conflict Detection
```sparql
# Detect conflicting triples between branches
PREFIX semsrc: <https://semsrc.org/ontology/>

SELECT ?subject ?predicate ?sourceValue ?targetValue
WHERE {
    {
        GRAPH :source-branch {
            ?subject ?predicate ?sourceValue .
        }
        GRAPH :target-branch {
            ?subject ?predicate ?targetValue .
        }
        FILTER(?sourceValue != ?targetValue)
    }
    UNION
    {
        GRAPH :source-branch {
            ?subject ?predicate ?sourceValue .
        }
        NOT EXISTS {
            GRAPH :target-branch {
                ?subject ?predicate ?targetValue .
            }
        }
    }
    UNION
    {
        GRAPH :target-branch {
            ?subject ?predicate ?targetValue .
        }
        NOT EXISTS {
            GRAPH :source-branch {
                ?subject ?predicate ?sourceValue .
            }
        }
    }
}
```

## Edge Cases

### Edge Case 1: Ontology-Aware Merging
**Scenario**: Merging triples that violate ontology constraints
**Handling**: Detect ontology violations, suggest corrections

### Edge Case 2: Circular References
**Scenario**: Merge creates circular relationships
**Handling**: Detect cycles, prevent or flag for review

### Edge Case 3: Large Graph Merges
**Scenario**: Merging large RDF graphs with thousands of triples
**Handling**: Incremental merging, progress indicators, memory management

### Edge Case 4: Conflicting Property Values
**Scenario**: Same subject-predicate pair with different values
**Handling**: Provide resolution strategies based on ontology rules

## Dependencies

### US-010: SPARQL Query Endpoint
- Required for querying conflicting triples

### US-013: Named Graph Support
- Required for storing branch-specific graphs

### US-029: SHACL Validation
- Required for merge validation

### US-038: Branch Management Git Compatible
- Required for understanding branch structure

## Test Cases

### Test Case 1: Detect Triple Conflicts
**Input**: Merge operation between feature and main branches
**Expected**: List of conflicting triples identified

### Test Case 2: Semantic Merge Suggestion
**Input**: Conflict detection result
**Expected**: Intelligent merge suggestions with explanations

### Test Case 3: Apply Merge Resolution
**Input**: Selected resolution strategy
**Expected**: Merged RDF graph with conflicts resolved

### Test Case 4: Validate Merge Results
**Input**: Merged RDF graph
**Expected**: Validation report showing any issues

## Related User Stories
- **US-013: Named Graph Support** - Related branch graph storage
- **US-029: SHACL Validation** - Related validation functionality
- **US-038: Branch Management Git Compatible** - Related merge operations

## Implementation Priority
**High** - Critical for semantic data workflows, should be implemented early.