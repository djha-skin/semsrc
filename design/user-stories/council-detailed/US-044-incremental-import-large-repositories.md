# US-044: Incremental Import for Large Repositories

## Description
Support incremental import of large Git repositories, allowing the import process to be paused and resumed. This makes migration of very large repositories feasible by breaking the process into manageable chunks.


## Test Cases
1. **Incremental import**: Import repository in multiple stages
2. **Pause/resume**: Pause import and resume from checkpoint
3. **Large repositories**: Handle repositories with hundreds of thousands of commits
4. **Import consistency**: Maintain consistency across multiple import sessions
5. **Error recovery**: Recover from import failures

## Dependencies
- US-040: Import Existing Git Repositories
- US-041: Read Git Pack Files Directly

## Implementation Notes
- Implement checkpoint system for import progress
- Track imported commits and objects
- Support pause/resume functionality
- Maintain import consistency across sessions
- Provide progress reporting and estimation
- Handle partial imports gracefully

## Edge Cases
- **Repository changes during import**: Should detect and handle concurrent modifications
- **Disk space**: Should handle running out of disk space during import
- **Import corruption**: Should detect and recover from corrupted import state
- **Very long imports**: Should handle imports that take days to complete
- **Network repositories**: Should handle importing from remote repositories