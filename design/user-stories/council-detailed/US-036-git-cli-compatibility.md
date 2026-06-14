# US-036: Git CLI Compatibility

## Description
Implement a `semsrc` binary with Git-compatible commands, allowing users to
continue using familiar Git workflows while leveraging semsrc's semantic
features. This provides a seamless migration path from Git.


## Test Cases
1. **Basic Git commands**: `status`, `add`, `commit`, `log` work correctly
2. **Branch operations**: Create, switch, merge branches
3. **Remote operations**: Push, pull from remote repositories
4. **Command compatibility**: Commands match Git semantics and output
5. **Error messages**: Git-compatible error messages and status codes

## Dependencies
- US-015: Initialize New Repository
- US-016: Add Files to Staging Area
- US-017: Commit Changes with a Message

## Implementation Notes
- Parse Git-style command line arguments
- Implement Git-compatible command interface
- Provide Git-compatible output formatting
- Support common Git options and flags
- Maintain compatibility with Git hooks
- Provide migration path for existing Git users

## Edge Cases
- **Unsupported commands**: Handle commands not yet implemented
- **Version differences**: Handle differences between Git versions
- **Configuration**: Support Git-compatible configuration
- **Performance**: Should match Git performance expectations
- **Migration**: Should support gradual migration from Git