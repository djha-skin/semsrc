# US-022: Ticket ID Extraction from Commit Messages

## Description
Automatically extract ticket IDs from commit messages and link commits to issue
tracking systems. This enables traceability between code changes and the issues
they address.


## Test Cases
1. **Basic ticket reference**: Extract ticket ID from commit message
2. **Multiple tickets**: Extract multiple ticket references
3. **Different formats**: Handle various ticket reference formats
4. **No tickets**: Handle commit messages without ticket references
5. **Integration with issue tracker**: Link to external issue tracking system

## Dependencies
- US-017: Commit Changes with a Message
- US-009: Store RDF Triples in Native RDF Store

## Implementation Notes
- Support multiple ticket reference formats (#12345, TICKET-12345, etc.)
- Configure ticket patterns per repository
- Store ticket links in semantic graph
- Support integration with external issue trackers
- Handle ticket ID validation
- Provide ticket lookup functionality

## Edge Cases
- **False positives**: Should avoid extracting non-ticket strings
- **Multiple formats**: Should handle mixed formats in same message
- **Private tickets**: Should handle access control for ticket links
- **Ticket deletion**: Should handle deleted/removed tickets
- **Format changes**: Should support customizable extraction patterns