## Description
Brief description of the changes in this PR.

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Performance improvement
- [ ] Code refactoring

## Mathematical Changes
If this PR involves mathematical operations:
- [ ] Maintains compatibility with original KAS library behavior
- [ ] Includes mathematical correctness verification
- [ ] Adds appropriate test cases for edge cases
- [ ] Updates LaTeX formatting if applicable

## Testing
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] I have tested the mathematical accuracy of any new operations

## Checklist
- [ ] My code follows the style guidelines of this project
- [ ] I have performed a self-review of my code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] Any dependent changes have been merged and published

## Examples
If this adds new functionality, provide usage examples:

```go
// Example of new functionality
expr, _ := parser.Parse("example")
result := newFunction(expr)
fmt.Println(result.String())
```

## Related Issues
Fixes #(issue number)

## Additional Notes
Add any additional notes about the PR here.