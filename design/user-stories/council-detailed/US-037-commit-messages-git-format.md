# US-037: Commit Messages in Git Format

## Description
Store commit messages in standard Git format to ensure compatibility with
external tools, CI/CD systems, and Git clients. This enables seamless
integration with the existing Git ecosystem.


## Test Cases
1. **Standard format**: Commit messages follow Git format conventions
2. **Multi-line messages**: Handle commit messages with multiple paragraphs
3. **Ticket references**: Preserve ticket references in messages
4. **Encoding**: Handle Unicode and special characters correctly
5. **Message validation**: Validate message format and length

## Dependencies
- US-017: Commit Changes with a Message

## Implementation Notes
- Follow Git commit message format specification
- Support subject line + body format
- Handle line wrapping automatically
- Preserve ticket references and metadata
- Support signing commits with GPG
- Provide commit message validation

## Edge Cases
- **Very long messages**: Should handle messages exceeding line length limits
- **Unicode characters**: Should handle international characters correctly
- **Special characters**: Should handle special characters without corruption
- **Message templates**: Support commit message templates
- **Validation errors**: Provide clear error messages for invalid formats