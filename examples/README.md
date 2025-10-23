# CAS Examples

This directory contains comprehensive examples demonstrating the capabilities of the CAS (Computer Algebra System) library.

## Running the Examples

Each example is a standalone Go program. To run them:

```bash
cd /Users/paramjeet/Projects/cas
go run examples/basic_usage.go
go run examples/calculus_operations.go
go run examples/polynomial_operations.go
go run examples/latex_formatting.go
```

## Example Files

### 1. `basic_usage.go`
Demonstrates fundamental operations:
- Expression parsing from strings
- Variable evaluation with specific values
- Working with mathematical constants (π, e)
- Expression simplification
- Basic function evaluation

**Key concepts covered:**
- `parser.Parse()` for parsing expressions
- `expr.Eval()` for numerical evaluation
- `expr.Variables()` to get variable names
- `expr.Simplify()` for algebraic simplification
- Working with `big.Float` for precision

### 2. `calculus_operations.go`
Showcases calculus functionality:
- Basic differentiation rules
- Chain rule applications
- Product rule examples
- Higher order derivatives
- Multivariable calculus (gradients)

**Key concepts covered:**
- `calculus.Derivative()` for first derivatives
- `calculus.NthDerivative()` for higher order derivatives
- `calculus.Gradient()` for partial derivatives
- `latex.FormatDerivative()` for derivative notation

### 3. `polynomial_operations.go`
Focuses on polynomial algebra:
- Polynomial expansion (basic and full)
- Linear equation solving
- Quadratic equation solving
- Equation verification
- Different equation forms (expr = 0 vs lhs = rhs)

**Key concepts covered:**
- `expand.Expand()` and `expand.ExpandFully()`
- `solve.Solve()` for equation solving
- `solve.SolveEquation()` for lhs = rhs form
- Solution verification techniques
- Working with exact vs approximate solutions

### 4. `latex_formatting.go`
Demonstrates LaTeX output generation:
- Basic expression formatting
- Function notation
- Complex mathematical expressions
- Formatting options and customization
- Special formatting for equations, derivatives, and integrals

**Key concepts covered:**
- `latex.Format()` with custom options
- `latex.FormatEquation()` for equations
- `latex.FormatDerivative()` for derivatives
- `latex.FormatIntegral()` for integrals
- Greek letters and special symbols
- Formatting options control

## Common Patterns

### Error Handling
```go
expr, err := parser.Parse("x^2 + 1")
if err != nil {
    log.Fatal(err)
}
```

### Variable Evaluation
```go
vars := map[string]*big.Float{
    "x": big.NewFloat(3.14159),
    "y": big.NewFloat(2.71828),
}
result, err := expr.Eval(vars)
```

### Chaining Operations
```go
// Parse -> Expand -> Differentiate -> Format
expr, _ := parser.Parse("(x+1)^2")
expanded := expand.Expand(expr)
derivative, _ := calculus.Derivative(expanded, "x")
latexOutput := latex.Format(derivative)
```

### Solution Verification
```go
solutions := solve.Solve(expr)
for _, sol := range solutions.Solutions {
    // Substitute back to verify
    vars := map[string]*big.Float{"x": solValue}
    if result, err := expr.Eval(vars); err == nil {
        isValid := math.Abs(result) < 1e-10
    }
}
```

## Expected Output Examples

### Basic Usage
```
Expression: x^2 + 3*x + 2
LaTeX: x^{2} + 3x + 2
When x = 2: 16.00
```

### Calculus Operations
```
d/dx(sin(x^2)) = cos(x^2) * 2*x
LaTeX: \cos\left(x^{2}\right) \cdot 2x
```

### Polynomial Operations
```
Solve: x^2 - 4 = 0
Solutions:
  x = 2 (exact)
  x = -2 (exact)
```

### LaTeX Formatting
```
Expression: sin(x^2) * cos(x) + ln(x + 1)
LaTeX: \sin\left(x^{2}\right)\cos(x) + \ln(x + 1)
```

## Building Your Own Examples

When creating your own examples, follow this pattern:

1. **Import required packages**
   ```go
   import (
       "github.com/quizizz/cas/pkg/parser"
       "github.com/quizizz/cas/pkg/calculus"
       "github.com/quizizz/cas/pkg/latex"
       // ... other packages as needed
   )
   ```

2. **Parse expressions safely**
   ```go
   expr, err := parser.Parse("your expression")
   if err != nil {
       // handle error
   }
   ```

3. **Use appropriate data types**
   - Use `big.Float` for high precision
   - Check for errors from all operations
   - Handle cases where expressions contain undefined variables

4. **Format output appropriately**
   - Use `latex.Format()` for mathematical notation
   - Use `expr.String()` for debugging
   - Consider using custom formatting options for specific needs

## Tips for Exploration

- **Start simple**: Begin with basic expressions and gradually add complexity
- **Check variables**: Use `expr.Variables()` to see what variables an expression contains
- **Verify results**: When solving equations, substitute solutions back to verify correctness
- **Use LaTeX**: The LaTeX output is often more readable than the string representation
- **Handle errors**: Many operations can fail if expressions are malformed or contain undefined variables
- **Experiment with options**: Try different `FormatOptions` to customize LaTeX output

## Common Issues and Solutions

1. **Parse errors**: Check for unsupported functions or malformed syntax
2. **Evaluation errors**: Ensure all variables have values before evaluation
3. **Complex expressions**: Some operations may not be implemented for all expression types
4. **Precision**: Use `big.Float` for high precision arithmetic

## Next Steps

After running these examples:
1. Try modifying the expressions to see different behaviors
2. Combine operations in new ways
3. Create your own mathematical problems to solve
4. Explore the CLI interface for interactive experimentation
5. Look at the test files for more usage patterns