# ADR-002: Copy-Editor Review of RDF Store Decision

## Review Date
June 14, 2026

## Reviewer
**Copy-Editor Persona** - Ensuring documentation consistency and clarity

## Task
Review all documents affected by the RDF store decision (ADR-001) and ensure consistency across the documentation.

## Documents Reviewed

### 1. User Story US-009
**File**: `user-stories/council-detailed/US-009-store-rdf-triples-native-rdf-store.md`
**Status**: ✅ Updated
**Changes Made**:
- Title changed from "Store RDF Triples in SQLite" to "Store RDF Triples in Native RDF Store"
- Description updated to reference Blazegraph/Oxigraph instead of SQLite
- Implementation notes updated for native RDF store features
- Edge cases updated for SPARQL optimization considerations

### 2. User Story References
**Files Updated**: 15 files in `user-stories/council-detailed/`
**Status**: ✅ Updated
**Changes Made**:
- All references to "US-009: Store RDF Triples in SQLite" changed to "US-009: Store RDF Triples in Native RDF Store"
- Verified with `rg` search for "US-009" - all instances updated

### 3. README.md
**File**: `user-stories/council-detailed/README.md`
**Status**: ✅ Updated
**Changes Made**:
- File reference updated from `US-009-store-rdf-triples-sqlite.md` to `US-009-store-rdf-triples-native-rdf-store.md`
- Description updated from "RDF triple storage in SQLite" to "RDF triple storage in native RDF store"

### 4. System Architecture
**File**: `architecture/system-architecture.md`
**Status**: ✅ Updated
**Changes Made**:
- Configuration example changed from `type: "sqlite"` to `type: "blazegraph"`
- Added comments about alternative RDF store options (Oxigraph, PostgreSQL)
- Updated deployment scenario descriptions
- Extended Go factory code to support multiple RDF store types

### 5. Bounded Contexts
**File**: `bounded-contexts.md`
**Status**: ✅ Updated
**Changes Made**:
- Triple store configuration examples updated
- Backend selection logic documented

### 6. Council Design Review
**File**: `council/council-design-review-updated.md`
**Status**: ✅ Updated
**Changes Made**:
- Added Kelsey Hightower as new council member
- Updated council ratings table
- Documented critical architectural decisions including RDF store selection

### 7. New Documentation Created
**Files Created**:
- `persona-kelsey-hightower.md` - Cloud deployment expertise
- `council-debate-rdf-stores.md` - Structured council debate
- `council-debate-personas-limited-context.md` - Limited context simulation
- `ADR-001-rdf-store-selection.md` - Architecture Decision Record

## Inconsistencies Found

### Minor Inconsistency: Configuration Comments
**Location**: `architecture/system-architecture.md` lines 1180-1210
**Issue**: Comments mention "sqlite" as alternative but code uses "blazegraph"
**Resolution**: ✅ Added clarification that SQLite remains as lighter embedded option

### Minor Inconsistency: Go Factory Code
**Location**: `architecture/system-architecture.md` lines 700-720
**Issue**: Factory code shows both "sqlite" and "blazegraph" options
**Resolution**: ✅ Updated to show Blazegraph as default, SQLite as alternative

## Recommendations

### 1. Add Decision Summary to Project README
**Location**: Project root `README.md`
**Action**: Add section summarizing RDF store decision and migration path

### 2. Update Quick Start Guide
**Location**: `quick-start.md` (if exists)
**Action**: Update installation instructions for Oxigraph local development

### 3. Add Migration Notes
**Location**: `migration-guide.md` (create if needed)
**Action**: Document steps for users migrating from SQLite to native RDF store

### 4. Consistent Naming Convention
**Recommendation**: Use "Native RDF Store" consistently instead of mixing "Blazegraph/Oxigraph"
**Exception**: Configuration examples should show specific store names

## Copy-Editor Approval
✅ **Approved** - All documents reviewed and updated for consistency with ADR-001.

**Next**: Git expert persona review for user story gaps.