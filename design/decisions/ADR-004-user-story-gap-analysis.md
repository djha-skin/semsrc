# ADR-004: User Story Gap Analysis and New User Stories

## Status
**Completed** - Gap analysis by Junio C Hamano (Git Expert) and Copy-Editor Persona

## Review Date
June 14, 2026

## Context
After reviewing the complete user story set (60 existing stories), the Git Expert persona (Junio C Hamano) identified significant gaps in Git-like functionality, particularly around worktree management and intermediate workflow steps.

## Gap Analysis Summary

### Existing User Stories Coverage
**Total Stories**: 60 (US-001 through US-060)

**Coverage Areas**:
- ✅ **Object Storage**: US-001-008 (8 stories)
- ✅ **RDF/SPARQL**: US-009-014 (6 stories)
- ✅ **Core Git Operations**: US-015-020 (6 stories)
- ✅ **PROV-O History**: US-021-030 (10 stories)
- ✅ **Git Compatibility**: US-066-044 (9 stories)
- ✅ **Hooks**: US-045-049 (5 stories)
- ✅ **Semantic Features**: US-050-058 (9 stories)
- ✅ **Architecture**: US-059-060 (2 stories)

### Identified Gaps

#### 1. Worktree Management ❌ **MISSING**
**Priority**: High
**Git Equivalent**: `git worktree add`, `git worktree list`, `git worktree remove`
**Impact**: Essential for parallel development workflows
**Solution**: Created US-061

#### 2. Advanced Staging Area ❌ **MISSING** (Partial)
**Priority**: High
**Git Equivalent**: Interactive staging, hunk-by-hunk staging
**Impact**: Critical for clean commit histories
**Solution**: Created US-062 (extends US-016)

#### 3. Interactive History Rewriting ❌ **MISSING**
**Priority**: Medium
**Git Equivalent**: `git rebase -i`, `git commit --amend`
**Impact**: Important for code cleanup
**Solution**: Created US-063

#### 4. Reflog and Recovery ❌ **MISSING**
**Priority**: Medium
**Git Equivalent**: `git reflog`, recovery operations
**Impact**: Safety net for developer mistakes
**Solution**: Created US-064

#### 5. Semantic Merge Conflicts ❌ **MISSING**
**Priority**: High
**Git Equivalent**: Three-way merge with semantic understanding
**Impact**: Critical for semantic data workflows
**Solution**: Created US-065

#### 6. Tag Management ❌ **MISSING**
**Priority**: Medium
**Git Equivalent**: `git tag`, version management
**Impact**: Important for release management
**Solution**: Created US-066

#### 7. Worktree Sharing ❌ **NOTED**
**Priority**: Medium
**Git Equivalent**: Remote worktree operations
**Impact**: Collaboration workflows
**Status**: Future consideration

#### 8. Worktree Cleanup ❌ **NOTED**
**Priority**: Low
**Git Equivalent**: `git worktree prune`
**Impact**: Maintenance tasks
**Status**: Future consideration

## New User Stories Created

### US-061: Manage Multiple Working Trees
**File**: `US-061-manage-multiple-working-trees.md`
**Priority**: High
**Coverage**: Worktree creation, listing, removal, conflict prevention
**RDF Integration**: Uses PROV-O for worktree relationships

### US-062: Staging Area Management
**File**: `US-062-staging-area-management.md`
**Priority**: High
**Coverage**: Interactive staging, hunk-by-hunk staging, diff viewing
**Note**: Extends US-016 (basic staging) with advanced features

### US-063: Interactive History Rewriting
**File**: `US-063-interactive-history-rewriting.md`
**Priority**: Medium
**Coverage**: Interactive rebase, commit amendment, commit reordering
**RDF Integration**: Tracks rewrite operations in RDF

### US-064: Reflog and Recovery
**File**: `US-064-reflog-and-recovery.md`
**Priority**: Medium
**Coverage**: Reflog viewing, commit recovery, worktree state recovery
**RDF Integration**: Stores operation history in RDF

### US-065: Semantic Merge Conflict Resolution
**File**: `US-065-semantic-merge-conflict-resolution.md`
**Priority**: High
**Coverage**: RDF conflict detection, semantic merge suggestions, validation
**RDF Integration**: Uses SPARQL for conflict detection

### US-066: Tag Management for Semantic Versioning
**File**: `US-066-tag-management-semantic-versioning.md`
**Priority**: Medium
**Coverage**: Tag creation, management, semantic versioning support
**RDF Integration**: Stores tag metadata in RDF

## User Story Count Update
- **Before**: 60 user stories (US-001 through US-060)
- **After**: 66 user stories (US-001 through US-066 + US-037-060)
- **New Stories**: 6 (US-061 through US-066)

## Document Updates Required

### Updated Documents
1. ✅ `decisions/ADR-004-user-story-gap-analysis.md` (this document)
2. ✅ `user-stories/council-detailed/US-061-manage-multiple-working-trees.md`
3. ✅ `user-stories/council-detailed/US-062-staging-area-management.md`
4. ✅ `user-stories/council-detailed/US-063-interactive-history-rewriting.md`
5. ✅ `user-stories/council-detailed/US-064-reflog-and-recovery.md`
6. ✅ `user-stories/council-detailed/US-065-semantic-merge-conflict-resolution.md`
7. ✅ `user-stories/council-detailed/US-066-tag-management-semantic-versioning.md`

### Documents to Update
1. `golden-thread.md` - Add worktree management workflow
2. `system-architecture.md` - Reference new user stories
3. `council/council-design-review-updated.md` - Update with new stories
4. `council/persona-research-summary.md` - Note gap analysis completion

## Copy-Editor Review Findings

### Documentation Consistency
✅ All new user stories follow existing patterns
✅ RDF integration consistent across new stories
✅ SPARQL queries properly formatted
✅ Interface designs follow Go conventions

### Recommendations
1. **Update Golden Thread**: Add worktree management to overall workflow
2. **Update Architecture Docs**: Reference new functionality
3. **Consistent Naming**: Use "Native RDF Store" consistently

## Git Expert (Junio C Hamano) Recommendations

### Immediate Actions
1. **Implement US-061** (Worktree Management) - Highest priority
2. **Implement US-062** (Staging Area) - High priority
3. **Implement US-065** (Semantic Merge) - High priority

### Short-term Actions
4. Implement US-063 (Interactive History)
5. Implement US-064 (Reflog/Recovery)
6. Implement US-066 (Tag Management)

### Future Considerations
7. Worktree sharing across machines
8. Worktree cleanup automation
9. Advanced conflict resolution strategies

## Council Consensus

### Gap Analysis Completion
✅ **Approved** - All identified gaps addressed with new user stories
✅ **Documentation** - New stories follow existing patterns
✅ **Priority** - High-priority stories identified for early implementation

### Next Steps
1. Update golden-thread.md with new workflow
2. Update architecture documentation
3. Prioritize implementation based on council recommendations
4. Review remaining user stories for additional gaps

## Related Documents
- `decisions/ADR-001-rdf-store-selection.md` - RDF store decision
- `decisions/ADR-002-copy-editor-review.md` - Documentation review
- `decisions/ADR-003-git-expert-review.md` - Gap analysis details
- `council/council-debate-personas-limited-context.md` - Council debate simulation

## Review Cycle
- **Next Review**: When implementing new user stories
- **Review Criteria**: Gap coverage, implementation priority, documentation completeness