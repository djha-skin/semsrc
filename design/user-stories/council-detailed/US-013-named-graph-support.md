# US-013: Named Graph Support

## Description

Implement named graph support in the triple store to organize RDF triples by
context. This allows different parts of the repository to have separate graphs
for different purposes (e.g., provenance, annotations, metadata).


## Test Cases

1. **Graph isolation**: Triples in different graphs are isolated
2. **Graph queries**: Query specific graphs or all graphs
3. **Graph management**: Create, list, and delete graphs
4. **Cross-graph queries**: Query across multiple graphs
5. **Graph metadata**: Store and retrieve graph metadata

## Dependencies

- US-009: Store RDF Triples in SQLite

## Implementation Notes

- Store graph URI with each triple
- Create index on graph for efficient queries
- Support default graph for ungraphed triples
- Implement graph-level operations (clear, export, import)
- Support graph visibility control (public/private graphs)
- Provide SPARQL FROM/FROM NAMED support

## Edge Cases

- **Default graph**: Should handle triples without explicit graph
- **Graph URI validation**: Should validate graph URIs
- **Large graphs**: Should handle graphs with millions of triples
- **Graph deletion**: Should handle graph deletion with referential integrity
- **Concurrent graph access**: Should handle concurrent graph operations