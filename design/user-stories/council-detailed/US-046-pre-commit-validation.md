# US-046: Pre-commit Validation

## Description
Implement pre-commit validation hooks to ensure semantic data quality before
commits are created. This prevents invalid data from entering the repository and
maintains data integrity.


## Test Cases
1. **Valid commit**: Pre-commit validation passes for valid commits
2. **Invalid commit**: Pre-commit validation rejects invalid commits
3. **Custom rules**: Support custom validation rules
4. **Performance**: Validation should not significantly slow commits
5. **Error reporting**: Clear error messages for validation failures

## Dependencies
- US-045: Automatic Hook Installation
- US-029: SHACL Validation for Data Quality

## Implementation Notes
- Implement validation rules as separate functions
- Run validation before commit creation
- Provide clear error messages for validation failures
- Support custom validation rules per repository
- Cache validation results for performance
- Allow skipping validation for specific commits

## Edge Cases
- **Large commits**: Should handle commits with many files efficiently
- **Custom rules**: Should support user-defined validation rules
- **Validation errors**: Should provide actionable error messages
- **Performance**: Should not significantly impact commit speed
- **False positives**: Should minimize false validation errors