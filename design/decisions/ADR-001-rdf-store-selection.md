# ADR-001: RDF Store Selection for Cloud and Local Deployment

## Status
**Accepted** - Council Decision (June 14, 2026)

## Context
The project requires RDF triple storage with the following competing concerns:
1. **Ease of deployment** in cloud (AWS-backed) environments
2. **Storage efficiency** for local development
3. **Query performance** for triple store operations with large datasets

The original User Story Nine specified SQLite for RDF triple storage, which has limitations for semantic web operations (SPARQL optimization, inference capabilities, graph-specific indexing).

## Decision
The council has decided to replace SQLite with **native RDF stores** for both cloud and local deployments:

### Cloud Deployment (AWS-Backed)
**Primary Choice: Amazon Neptune** (Council Vote: 3-2)
- **Rationale**: Fully managed service, easy deployment, SPARQL 1.1 compliance, AWS ecosystem integration
- **Deployment Time**: < 30 minutes with CloudFormation
- **Cost Model**: Pay-as-you-go, serverless option available
- **Backup**: Automated snapshots, cross-region replication

**Alternative: Blazegraph on EKS**
- **When**: Teams with Kubernetes expertise
- **Benefits**: Full control, open source, portable across clouds
- **Trade-off**: Higher operational burden

### Local Development
**Primary Choice: Oxigraph** (Council Vote: 4-1)
- **Rationale**: Single binary deployment, fast startup, cross-platform compatibility
- **Performance**: Native Rust speed, suitable for 1M-10M triple datasets
- **Integration**: Mirrors Neptune SPARQL endpoint for development consistency
- **Backup**: Git LFS + automated scripts

**Alternative: Blazegraph Local**
- **When**: Need full inference capabilities
- **Benefits**: Same features as cloud deployment
- **Trade-off**: Heavier resource usage (Java-based)

## Rationale
### Query Performance Considerations
Triple stores involve complex joins across large datasets. Native RDF stores provide:
- **Specialized indexing**: SPO, POS, OSP graph-aware indexes
- **SPARQL-specific optimization**: Query planning for graph operations
- **Efficient memory management**: Built for large-scale semantic data

### Council Member Perspectives
- **Kelsey Hightower (Cloud Expert)**: Neptune provides best balance of ease-of-use and managed service
- **Ruben Verborgh (Semantic Web)**: Native stores essential for RDF standards compliance
- **Ian Lance Taylor (Performance)**: Graph-aware optimization crucial for large triple stores
- **Junio C Hamano (Data Integrity)**: Cloud snapshots + Git LFS provide robust backup strategy
- **Martin Fowler (Architecture)**: Strong abstraction enables future evolution

## Consequences
### Positive
- ✅ Improved SPARQL query performance (native optimization)
- ✅ Better inference capabilities (OWL support)
- ✅ Easier cloud deployment (managed Neptune)
- ✅ Consistent local/cloud development experience
- ✅ Standards compliance (SPARQL 1.1, RDF 1.1)

### Negative
- ⚠️ Vendor lock-in to AWS (Neptune)
- ⚠️ Learning curve for native RDF stores
- ⚠️ Migration effort from existing SQLite setup

## Implementation Notes
### Abstraction Layer Required
```go
type TripleStore interface {
    Query(sparql string) (ResultSet, error)
    Insert(triples []Triple) error
    Delete(triples []Triple) error
    Backup() (BackupID, error)
    Restore(BackupID) error
}
```

### Migration Strategy
1. **Phase 1**: Add Oxigraph support for local development
2. **Phase 2**: Add Neptune support for cloud deployment
3. **Phase 3**: Provide migration tools from SQLite
4. **Phase 4**: Deprecate SQLite with feature flags
5. **Phase 5**: Evaluate Blazegraph if Neptune proves limiting

### Configuration Examples
```yaml
# Cloud (Neptune)
triple_store:
  type: "blazegraph"  # Default choice
  endpoint: "https://your-neptune-cluster.cluster-xyz.us-east-1.neptune.amazonaws.com:8182/sparql"

# Local (Oxigraph)
triple_store:
  type: "oxigraph"  # Embedded store
  path: "/home/user/.semsrc/triples"
```

## Related Documents
- `council-debate-rdf-stores.md` - Detailed council debate
- `council-debate-personas-limited-context.md` - Limited context simulation
- `persona-kelsey-hightower.md` - Cloud deployment expertise
- `user-stories/council-detailed/US-009-store-rdf-triples-native-rdf-store.md` - Updated user story

## Review Cycle
- **Next Review**: When migrating from SQLite or evaluating new RDF stores
- **Review Criteria**: Query performance benchmarks, cost analysis, vendor lock-in assessment

## References
- Amazon Neptune Documentation: https://aws.amazon.com/neptune/
- Oxigraph: https://github.com/oxigraph/oxigraph
- Blazegraph: https://blazegraph.com/