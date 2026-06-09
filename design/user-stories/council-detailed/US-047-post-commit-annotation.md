# US-047: Post-commit Annotation

## Description
Automatically annotate commits with semantic metadata after commit creation. This enriches the knowledge graph with relationships and context derived from the commit operation.


## Test Cases
1. **Automatic annotation**: Annotations are generated automatically after commit
2. **Annotation accuracy**: Annotations are accurate and relevant
3. **Performance**: Annotation does not significantly slow commit process
4. **Custom strategies**: Support custom annotation strategies
5. **Error handling**: Handle annotation failures gracefully

## Dependencies
- US-045: Automatic Hook Installation
- US-021: Automatic PROV-O Relationship Generation
- US-022: Ticket ID Extraction from Commit Messages

## Implementation Notes
- Implement annotation as post-commit hook
- Support multiple annotation strategies
- Store annotations in semantic graph
- Handle annotation failures without breaking commits
- Provide annotation configuration options
- Support incremental annotation for large commits

## Edge Cases
- **Large commits**: Should handle commits with many files efficiently
- **Annotation failures**: Should not prevent commit from completing
- **Custom strategies**: Should support user-defined annotation strategies
- **Performance**: Should not significantly impact commit speed
- **Annotation conflicts**: Should handle conflicting annotations