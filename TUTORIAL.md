# CAS Library Tutorial

Welcome to the comprehensive tutorial for the CAS (Computer Algebra System) library. This guide will take you from basic concepts to advanced usage patterns.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Core Concepts](#core-concepts)
3. [Expression Parsing](#expression-parsing)
4. [Working with Variables](#working-with-variables)
5. [Mathematical Operations](#mathematical-operations)
6. [Calculus Operations](#calculus-operations)
7. [Polynomial Algebra](#polynomial-algebra)
8. [Equation Solving](#equation-solving)
9. [LaTeX Formatting](#latex-formatting)
10. [Advanced Usage](#advanced-usage)
11. [Performance Considerations](#performance-considerations)
12. [Troubleshooting](#troubleshooting)

## Getting Started

### Installation and Setup

```bash
# Clone the repository
git clone https://github.com/quizizz/cas.git
cd cas

# Install dependencies
go mod tidy

# Build the CLI tool
go build ./cmd/cas

# Run tests to verify installation
go test ./...
```

### Your First Expression

```go
package main

import (
    "fmt"
    "log"

    "github.com/quizizz/cas/pkg/parser"
)

func main() {
    // Parse a simple expression
    expr, err := parser.Parse("x^2 + 2*x + 1")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Expression:", expr.String())
    fmt.Println("LaTeX:", expr.LaTeX())
}
```

## Core Concepts

### Abstract Syntax Tree (AST)

The CAS library represents mathematical expressions as Abstract Syntax Trees (ASTs). Every expression implements the `Expr` interface:

```go
type Expr interface {
    String() string                                         // Human-readable representation
    LaTeX() string                                         // LaTeX representation
    Eval(vars map[string]*big.Float) (*big.Float, error)  // Numerical evaluation
    Variables() []string                                   // List of variables
    Simplify() Expr                                       // Algebraic simplification
    Clone() Expr                                          // Create a copy
}
```

### Expression Types

The library supports several types of expressions:

- **Numbers**: `Int`, `Float`, `Rational`
- **Variables**: `Var` (like x, y, theta)
- **Constants**: `Const` (like π, e)
- **Operations**: `Add`, `Mul`, `Pow`
- **Functions**: `Func` (like sin, cos, ln)

### Precision and Accuracy

The library uses `math/big` package for arbitrary precision arithmetic:

```go
import "math/big"

// High precision evaluation
vars := map[string]*big.Float{
    "x": big.NewFloat(3.141592653589793),
}
result, _ := expr.Eval(vars)
fmt.Println(result.Text('f', 20)) // 20 decimal places
```

## Expression Parsing

### Basic Syntax

The parser supports standard mathematical notation:

```go
examples := []string{
    "x + y",           // Addition
    "x - y",           // Subtraction
    "x * y",           // Multiplication
    "x / y",           // Division
    "x ^ y",           // Exponentiation
    "(x + 1) * 2",     // Parentheses
    "2*x^2 + 3*x + 1", // Complex expressions
}
```

### Functions

Supported mathematical functions:

```go
functions := []string{
    "sin(x)",     // Sine
    "cos(x)",     // Cosine
    "tan(x)",     // Tangent
    "ln(x)",      // Natural logarithm
    "log(x)",     // Base-10 logarithm
    "sqrt(x)",    // Square root
    "abs(x)",     // Absolute value
    "exp(x)",     // Exponential (e^x)
    "sinh(x)",    // Hyperbolic sine
    "cosh(x)",    // Hyperbolic cosine
    "tanh(x)",    // Hyperbolic tangent
}
```

### Constants

Built-in mathematical constants:

```go
constants := []string{
    "pi",   // π ≈ 3.14159...
    "e",    // e ≈ 2.71828...
}

expr, _ := parser.Parse("pi * r^2")  // Circle area
```

### Error Handling

Always check for parse errors:

```go
expr, err := parser.Parse("invalid expression!")
if err != nil {
    fmt.Printf("Parse error: %v\n", err)
    // Handle the error appropriately
}
```

## Working with Variables

### Variable Names

Variables can be single characters or words:

```go
expressions := []string{
    "x",                    // Single character
    "theta",               // Greek letter name
    "variable_name",       // Underscore allowed
    "x1",                  // Numbers allowed
}
```

### Getting Variable Names

```go
expr, _ := parser.Parse("x^2 + y*z + pi")
variables := expr.Variables()
fmt.Println("Variables:", variables) // [x y z]
// Note: constants like 'pi' are not included
```

### Variable Evaluation

```go
expr, _ := parser.Parse("x^2 + 2*x + 1")

// Evaluate with specific values
vars := map[string]*big.Float{
    "x": big.NewFloat(3.0),
}

result, err := expr.Eval(vars)
if err != nil {
    fmt.Printf("Evaluation error: %v\n", err)
} else {
    fmt.Printf("Result: %s\n", result.Text('f', 2))
}
```

### Undefined Variables

```go
expr, _ := parser.Parse("x + y + z")
vars := map[string]*big.Float{
    "x": big.NewFloat(1.0),
    // y and z are undefined
}

_, err := expr.Eval(vars)
if err != nil {
    fmt.Printf("Error: %v\n", err) // Will report missing variables
}
```

## Mathematical Operations

### Basic Arithmetic

```go
// Addition and subtraction
add_expr, _ := parser.Parse("x + y - z")

// Multiplication (various notations)
mul_examples := []string{
    "2*x",        // Explicit multiplication
    "2x",         // Implicit (parsed as 2*x)
    "x*y*z",      // Multiple factors
    "(x+1)(x-1)", // Implicit between parentheses
}

// Division
div_expr, _ := parser.Parse("(x + 1) / (x - 1)")

// Exponentiation
pow_examples := []string{
    "x^2",        // Square
    "x^(1/2)",    // Square root
    "2^(3^4)",    // Right associative
}
```

### Operator Precedence

The parser follows standard mathematical precedence:

1. Parentheses: `()`
2. Functions: `sin()`, `cos()`, etc.
3. Exponentiation: `^` (right associative)
4. Multiplication and Division: `*`, `/` (left associative)
5. Addition and Subtraction: `+`, `-` (left associative)

```go
expr, _ := parser.Parse("2 + 3 * 4^2")
// Parsed as: 2 + (3 * (4^2)) = 2 + (3 * 16) = 50
```

### Simplification

```go
expr, _ := parser.Parse("x + 0 + x*1 + 0*y")
simplified := expr.Simplify()
fmt.Println("Original:", expr.String())     // x + 0 + x*1 + 0*y
fmt.Println("Simplified:", simplified.String()) // 2*x
```

## Calculus Operations

### First Derivatives

```go
import "github.com/quizizz/cas/pkg/calculus"

expr, _ := parser.Parse("x^3 + 2*x^2 + x + 1")
derivative, err := calculus.Derivative(expr, "x")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("f(x) = %s\n", expr.String())
fmt.Printf("f'(x) = %s\n", derivative.String())
// Output: f'(x) = 3*x^2 + 4*x + 1
```

### Higher Order Derivatives

```go
// Second derivative
second_deriv, err := calculus.NthDerivative(expr, "x", 2)

// Third derivative
third_deriv, err := calculus.NthDerivative(expr, "x", 3)

fmt.Printf("f''(x) = %s\n", second_deriv.String())
fmt.Printf("f'''(x) = %s\n", third_deriv.String())
```

### Chain Rule

The library automatically applies the chain rule:

```go
expr, _ := parser.Parse("sin(x^2)")
derivative, _ := calculus.Derivative(expr, "x")
fmt.Printf("d/dx[sin(x²)] = %s\n", derivative.String())
// Output: cos(x^2) * 2*x
```

### Product Rule

```go
expr, _ := parser.Parse("x * sin(x)")
derivative, _ := calculus.Derivative(expr, "x")
fmt.Printf("d/dx[x·sin(x)] = %s\n", derivative.String())
// Output: sin(x) + x*cos(x)
```

### Multivariable Calculus

```go
expr, _ := parser.Parse("x^2 + x*y + y^2")
variables := []string{"x", "y"}

gradient, err := calculus.Gradient(expr, variables)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("f(x,y) = %s\n", expr.String())
for _, variable := range variables {
    partial := gradient[variable]
    fmt.Printf("∂f/∂%s = %s\n", variable, partial.String())
}
```

## Polynomial Algebra

### Expansion

```go
import "github.com/quizizz/cas/pkg/expand"

// Basic expansion
expr, _ := parser.Parse("(x + 1)^2")
expanded := expand.Expand(expr)
fmt.Printf("(%s) = %s\n", expr.String(), expanded.String())
// Output: (x + 1)^2 = x^2 + 2*x + 1

// Full expansion for complex expressions
complex_expr, _ := parser.Parse("(x + 1)^2 * (x - 2)")
fully_expanded := expand.ExpandFully(complex_expr)
fmt.Println("Fully expanded:", fully_expanded.String())
```

### Expansion Options

```go
options := expand.ExpandOptions{
    MaxDepth:   10,     // Maximum recursion depth
    ExpandTrig: true,   // Expand trigonometric identities
    ExpandLog:  true,   // Expand logarithmic properties
}

expanded := expand.ExpandWithOptions(expr, options)
```

## Equation Solving

### Linear Equations

```go
import "github.com/quizizz/cas/pkg/solve"

// Solve: 2x - 4 = 0
expr, _ := parser.Parse("2*x - 4")
solutions := solve.Solve(expr)

if solutions.HasSolutions {
    for _, sol := range solutions.Solutions {
        fmt.Printf("%s = %s\n", sol.Variable, sol.Value.String())
    }
}
```

### Quadratic Equations

```go
// Solve: x² - 5x + 6 = 0
expr, _ := parser.Parse("x^2 - 5*x + 6")
solutions := solve.Solve(expr)

fmt.Printf("Status: %s\n", solutions.Message)
for i, sol := range solutions.Solutions {
    fmt.Printf("x_%d = %s", i+1, sol.Value.String())
    if sol.IsExact {
        fmt.Print(" (exact)")
    }
    fmt.Println()
}
```

### Equation Form (lhs = rhs)

```go
// Solve: x + 1 = 3
lhs, _ := parser.Parse("x + 1")
rhs, _ := parser.Parse("3")
solutions := solve.SolveEquation(lhs, rhs)
```

### Custom Solving Options

```go
options := solve.SolveOptions{
    Variable:         "x",     // Variable to solve for
    AllowComplex:     false,   // Allow complex solutions
    AllowApproximate: true,    // Allow numerical approximations
    MaxDegree:        4,       // Maximum polynomial degree
}

solutions := solve.Solve(expr, options)
```

## LaTeX Formatting

### Basic Formatting

```go
import "github.com/quizizz/cas/pkg/latex"

expr, _ := parser.Parse("x^2 + 2*x + 1")
latexStr := latex.Format(expr)
fmt.Printf("LaTeX: %s\n", latexStr)
// Output: x^{2} + 2x + 1
```

### Custom Formatting Options

```go
options := latex.FormatOptions{
    UseFractions:     true,   // Convert divisions to \frac{}{}
    UseSymbols:       true,   // Use \pi, \alpha, etc.
    UseParentheses:   true,   // Use \left( \right)
    SimplifyRoots:    true,   // Convert x^{1/2} to \sqrt{x}
    MaxDecimalPlaces: 6,      // Decimal precision
}

formatted := latex.Format(expr, options)
```

### Special Formatting Functions

```go
// Equation formatting
lhs, _ := parser.Parse("x^2")
rhs, _ := parser.Parse("4")
equation := latex.FormatEquation(lhs, rhs)
fmt.Println(equation) // x^{2} = 4

// Derivative notation
expr, _ := parser.Parse("x^3 + x")
derivative_latex := latex.FormatDerivative(expr, "x", 1)
fmt.Println(derivative_latex) // \frac{d}{dx}\left(x^{3} + x\right)

// Integral notation
integral_latex := latex.FormatIntegral(expr, "x", false, nil, nil)
fmt.Println(integral_latex) // \int x^{3} + x \, dx
```

### Function and Symbol Formatting

```go
// Functions are properly formatted
expr, _ := parser.Parse("sin(x) + ln(y)")
fmt.Println(latex.Format(expr))
// Output: \sin(x) + \ln(y)

// Greek letters
expr, _ := parser.Parse("alpha + beta + pi")
fmt.Println(latex.Format(expr))
// Output: \alpha + \beta + \pi
```

## Advanced Usage

### Cloning Expressions

```go
original, _ := parser.Parse("x^2 + 1")
clone := original.Clone()

// Modify clone without affecting original
// (expressions are generally immutable anyway)
```

### Custom Expression Building

```go
import "github.com/quizizz/cas/pkg/ast"

// Build expressions programmatically
x := ast.NewVar("x")
two := ast.NewInt(2)
power := ast.NewPow(x, two)
one := ast.NewInt(1)
expr := ast.NewAdd(power, one) // x^2 + 1
```

### Combining Operations

```go
// Chain operations together
original, _ := parser.Parse("(x + 1)^2")
expanded := expand.Expand(original)
derivative, _ := calculus.Derivative(expanded, "x")
latex_output := latex.Format(derivative)

fmt.Printf("Original: %s\n", original.String())
fmt.Printf("Expanded: %s\n", expanded.String())
fmt.Printf("Derivative: %s\n", derivative.String())
fmt.Printf("LaTeX: %s\n", latex_output)
```

### Working with Multiple Variables

```go
expr, _ := parser.Parse("x^2 + 2*x*y + y^2")

// Get all variables
vars := expr.Variables()
fmt.Printf("Variables: %v\n", vars)

// Compute partial derivatives for all variables
for _, variable := range vars {
    partial, err := calculus.Derivative(expr, variable)
    if err == nil {
        fmt.Printf("∂/∂%s = %s\n", variable, partial.String())
    }
}
```

## Performance Considerations

### Expression Caching

For frequently used expressions, consider caching parsed results:

```go
type ExprCache struct {
    cache map[string]ast.Expr
}

func (ec *ExprCache) Parse(exprStr string) (ast.Expr, error) {
    if expr, exists := ec.cache[exprStr]; exists {
        return expr, nil
    }

    expr, err := parser.Parse(exprStr)
    if err == nil {
        ec.cache[exprStr] = expr
    }
    return expr, err
}
```

### Numerical vs Symbolic

- Use symbolic operations when you need exact results
- Use numerical evaluation when you need speed and approximation is acceptable
- Consider the trade-off between precision and performance

### Memory Usage

- Large expressions with many terms can use significant memory
- Consider simplifying expressions when possible
- Clone expressions only when necessary

### Optimization Tips

```go
// Simplify expressions early
expr, _ := parser.Parse("x + 0 + x*1")
expr = expr.Simplify() // Reduces to 2*x

// Reuse variable maps
vars := map[string]*big.Float{
    "x": big.NewFloat(1.0),
}
// Reuse 'vars' for multiple evaluations

// Pre-compute constants
pi_val, _ := big.NewFloat(0).SetString("3.141592653589793")
vars["pi"] = pi_val
```

## Troubleshooting

### Common Parse Errors

```go
// Invalid function names
expr, err := parser.Parse("sine(x)")  // Should be "sin(x)"

// Mismatched parentheses
expr, err := parser.Parse("(x + 1")   // Missing closing parenthesis

// Invalid operators
expr, err := parser.Parse("x ** y")   // Should be "x ^ y"
```

### Evaluation Issues

```go
// Undefined variables
expr, _ := parser.Parse("x + y")
vars := map[string]*big.Float{"x": big.NewFloat(1.0)}
// Missing "y" will cause evaluation error

// Domain errors
expr, _ := parser.Parse("ln(-1)")  // Natural log of negative number
result, err := expr.Eval(make(map[string]*big.Float))
// Will return an error
```

### Differentiation Limitations

```go
// Non-differentiable functions
expr, _ := parser.Parse("abs(x)")
derivative, err := calculus.Derivative(expr, "x")
// May not handle discontinuous derivatives properly

// Implicit functions
expr, _ := parser.Parse("x^y")  // Both x and y are variables
derivative, err := calculus.Derivative(expr, "x")
// May require treating y as constant
```

### Solving Limitations

```go
// High-degree polynomials
expr, _ := parser.Parse("x^5 + x + 1")
solutions := solve.Solve(expr)
// May not find exact solutions for degree > 4

// Transcendental equations
expr, _ := parser.Parse("sin(x) - x")
solutions := solve.Solve(expr)
// Currently not supported
```

### Debug Techniques

```go
// Check expression structure
fmt.Printf("String representation: %s\n", expr.String())
fmt.Printf("Variables: %v\n", expr.Variables())

// Verify LaTeX output
fmt.Printf("LaTeX: %s\n", latex.Format(expr))

// Test evaluation with simple values
vars := map[string]*big.Float{"x": big.NewFloat(1.0)}
if result, err := expr.Eval(vars); err == nil {
    fmt.Printf("f(1) = %s\n", result.Text('f', 4))
}
```

## Best Practices

1. **Always check errors** from parsing and evaluation
2. **Use appropriate precision** for your use case
3. **Simplify expressions** when possible to improve performance
4. **Cache frequently used expressions** and computed results
5. **Use LaTeX formatting** for readable mathematical output
6. **Verify solutions** by substituting back into original equations
7. **Handle edge cases** like division by zero or undefined functions
8. **Use meaningful variable names** in your expressions
9. **Document complex mathematical operations** in your code
10. **Test with known mathematical identities** to verify correctness

## Next Steps

- Explore the [examples](examples/) directory for more complex use cases
- Try the CLI interface for interactive exploration
- Read the source code to understand implementation details
- Contribute new features or improvements
- Create your own mathematical applications using the library

This tutorial should give you a solid foundation for using the CAS library effectively. For more specific use cases, refer to the examples and API documentation.