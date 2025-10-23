// Package solve implements equation solving algorithms for various types of equations.
package solve

import (
	"fmt"
	"math/big"

	"github.com/quizizz/cas/pkg/ast"
	"github.com/quizizz/cas/pkg/simplify"
)

// Solution represents a solution to an equation
type Solution struct {
	Variable string
	Value    ast.Expr
	IsReal   bool
	IsExact  bool
}

// SolutionSet represents all solutions to an equation
type SolutionSet struct {
	Solutions    []Solution
	Message      string
	HasSolutions bool
}

// SolveOptions controls equation solving behavior
type SolveOptions struct {
	// Variable specifies which variable to solve for
	Variable string
	// AllowComplex enables complex number solutions
	AllowComplex bool
	// AllowApproximate enables approximate numerical solutions
	AllowApproximate bool
	// MaxDegree limits the polynomial degree for solving
	MaxDegree int
}

// DefaultSolveOptions returns default solving options
func DefaultSolveOptions() SolveOptions {
	return SolveOptions{
		Variable:         "x",
		AllowComplex:     false,
		AllowApproximate: true,
		MaxDegree:        4,
	}
}

// Solve attempts to solve an equation expr = 0 for the specified variable
func Solve(expr ast.Expr, opts ...SolveOptions) SolutionSet {
	options := DefaultSolveOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	// Simplify the expression first
	simplified := simplify.Simplify(expr)

	// Determine equation type and solve accordingly
	return solveEquation(simplified, options)
}

// SolveEquation solves equation lhs = rhs for the specified variable
func SolveEquation(lhs, rhs ast.Expr, opts ...SolveOptions) SolutionSet {
	options := DefaultSolveOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	// Convert to standard form: lhs - rhs = 0
	diff := ast.NewAdd(lhs, ast.NewMul(ast.NewInt(-1), rhs))
	simplified := simplify.Simplify(diff)

	return solveEquation(simplified, options)
}

// solveEquation is the main solving dispatcher
func solveEquation(expr ast.Expr, opts SolveOptions) SolutionSet {
	// Check if expression contains the variable
	variables := expr.Variables()
	hasVariable := false
	for _, v := range variables {
		if v == opts.Variable {
			hasVariable = true
			break
		}
	}

	if !hasVariable {
		// No variable present - check if expression equals zero
		result, err := expr.Eval(make(map[string]*big.Float))
		if err != nil {
			return SolutionSet{
				Message:      "Cannot evaluate constant expression",
				HasSolutions: false,
			}
		}

		if result.Sign() == 0 {
			return SolutionSet{
				Message:      "Identity: true for all values",
				HasSolutions: true,
			}
		} else {
			return SolutionSet{
				Message:      "No solution: expression is never zero",
				HasSolutions: false,
			}
		}
	}

	// Classify equation type
	degree := getPolynomialDegree(expr, opts.Variable)

	switch {
	case degree == 0:
		return solveConstant(expr, opts)
	case degree == 1:
		return solveLinear(expr, opts)
	case degree == 2:
		return solveQuadratic(expr, opts)
	case degree == 3:
		return solveCubic(expr, opts)
	case degree == 4:
		return solveQuartic(expr, opts)
	case degree > 4:
		return SolutionSet{
			Message:      fmt.Sprintf("Polynomial degree %d too high for exact solution", degree),
			HasSolutions: false,
		}
	default:
		return solveGeneral(expr, opts)
	}
}

// getPolynomialDegree determines the degree of a polynomial in the given variable
func getPolynomialDegree(expr ast.Expr, variable string) int {
	return getExpressionDegree(expr, variable)
}

func getExpressionDegree(expr ast.Expr, variable string) int {
	switch e := expr.(type) {
	case *ast.Var:
		if e.Name() == variable {
			return 1
		}
		return 0
	case ast.Numeric:
		return 0
	case *ast.Const:
		return 0
	case *ast.Add:
		maxDegree := 0
		for _, term := range e.Terms() {
			degree := getExpressionDegree(term, variable)
			if degree > maxDegree {
				maxDegree = degree
			}
		}
		return maxDegree
	case *ast.Mul:
		totalDegree := 0
		for _, factor := range e.Terms() {
			totalDegree += getExpressionDegree(factor, variable)
		}
		return totalDegree
	case *ast.Pow:
		base := e.Base()
		exp := e.Exponent()

		baseDegree := getExpressionDegree(base, variable)
		if baseDegree == 0 {
			return 0 // Base doesn't contain variable
		}

		// Check if exponent is a constant integer
		if expInt, ok := exp.(*ast.Int); ok {
			val, err := expInt.Eval(make(map[string]*big.Float))
			if err == nil {
				intVal, _ := val.Int64()
				if intVal >= 0 {
					return baseDegree * int(intVal)
				}
			}
		}

		// Non-constant or negative exponent - not a polynomial
		return -1
	default:
		// Functions and other expressions - not polynomial
		return -1
	}
}

// solveConstant handles constant equations (no variable)
func solveConstant(expr ast.Expr, opts SolveOptions) SolutionSet {
	result, err := expr.Eval(make(map[string]*big.Float))
	if err != nil {
		return SolutionSet{
			Message:      "Cannot evaluate expression",
			HasSolutions: false,
		}
	}

	if result.Sign() == 0 {
		return SolutionSet{
			Message:      "Identity: true for all values of " + opts.Variable,
			HasSolutions: true,
		}
	} else {
		return SolutionSet{
			Message:      "No solution: " + result.Text('g', -1) + " ≠ 0",
			HasSolutions: false,
		}
	}
}

// solveLinear solves linear equations ax + b = 0
func solveLinear(expr ast.Expr, opts SolveOptions) SolutionSet {
	// Extract coefficients: ax + b = 0
	a, b := extractLinearCoefficients(expr, opts.Variable)

	// Check if 'a' is zero
	aVal, err := a.Eval(make(map[string]*big.Float))
	if err != nil || aVal.Sign() == 0 {
		return SolutionSet{
			Message:      "Not a linear equation in " + opts.Variable,
			HasSolutions: false,
		}
	}

	// Solution: x = -b/a
	solution := ast.NewMul(ast.NewInt(-1), b, ast.NewPow(a, ast.NewInt(-1)))
	simplified := simplify.Simplify(solution)

	return SolutionSet{
		Solutions: []Solution{{
			Variable: opts.Variable,
			Value:    simplified,
			IsReal:   true,
			IsExact:  true,
		}},
		Message:      "Linear equation solved",
		HasSolutions: true,
	}
}

// extractLinearCoefficients extracts a and b from ax + b
func extractLinearCoefficients(expr ast.Expr, variable string) (a, b ast.Expr) {
	// This is a simplified implementation
	// In a full implementation, we would need to collect all terms
	// and separate those with the variable from constants

	a = ast.NewInt(0)
	b = ast.NewInt(0)

	switch e := expr.(type) {
	case *ast.Add:
		for _, term := range e.Terms() {
			if containsVariable(term, variable) {
				// Extract coefficient of the variable
				coeff := extractCoefficient(term, variable)
				a = ast.NewAdd(a, coeff)
			} else {
				// Constant term
				b = ast.NewAdd(b, term)
			}
		}
	case *ast.Mul:
		if containsVariable(e, variable) {
			a = extractCoefficient(e, variable)
		} else {
			b = e
		}
	case *ast.Var:
		if e.Name() == variable {
			a = ast.NewInt(1)
		} else {
			b = e
		}
	default:
		if containsVariable(e, variable) {
			a = ast.NewInt(1) // Simplified assumption
		} else {
			b = e
		}
	}

	return a, b
}

// extractCoefficient extracts the coefficient of a variable from a term
func extractCoefficient(term ast.Expr, variable string) ast.Expr {
	switch t := term.(type) {
	case *ast.Var:
		if t.Name() == variable {
			return ast.NewInt(1)
		}
		return ast.NewInt(0)
	case *ast.Mul:
		var coeff ast.Expr = ast.NewInt(1)
		hasVar := false
		for _, factor := range t.Terms() {
			if v, ok := factor.(*ast.Var); ok && v.Name() == variable {
				hasVar = true
			} else if !containsVariable(factor, variable) {
				coeff = ast.NewMul(coeff, factor)
			}
		}
		if hasVar {
			return coeff
		}
		return ast.NewInt(0)
	default:
		if containsVariable(term, variable) {
			return ast.NewInt(1) // Simplified
		}
		return ast.NewInt(0)
	}
}

// containsVariable checks if an expression contains a specific variable
func containsVariable(expr ast.Expr, variable string) bool {
	variables := expr.Variables()
	for _, v := range variables {
		if v == variable {
			return true
		}
	}
	return false
}

// solveQuadratic solves quadratic equations ax² + bx + c = 0
func solveQuadratic(expr ast.Expr, opts SolveOptions) SolutionSet {
	// Extract coefficients: ax² + bx + c = 0
	a, b, c := extractQuadraticCoefficients(expr, opts.Variable)

	// Check if 'a' is zero (not actually quadratic)
	aVal, err := a.Eval(make(map[string]*big.Float))
	if err != nil || aVal.Sign() == 0 {
		// Fall back to linear solution
		return solveLinear(expr, opts)
	}

	// Calculate discriminant: b² - 4ac
	discriminant := ast.NewAdd(
		ast.NewPow(b, ast.NewInt(2)),
		ast.NewMul(ast.NewInt(-4), a, c),
	)
	discriminantSimplified := simplify.Simplify(discriminant)

	// Evaluate discriminant
	discVal, err := discriminantSimplified.Eval(make(map[string]*big.Float))
	if err != nil {
		return SolutionSet{
			Message:      "Cannot evaluate discriminant",
			HasSolutions: false,
		}
	}

	if discVal.Sign() < 0 {
		if !opts.AllowComplex {
			return SolutionSet{
				Message:      "No real solutions (discriminant < 0)",
				HasSolutions: false,
			}
		}
		// TODO: Implement complex solutions
		return SolutionSet{
			Message:      "Complex solutions not yet implemented",
			HasSolutions: false,
		}
	}

	// Calculate solutions using quadratic formula: x = (-b ± √discriminant) / (2a)
	sqrt := ast.NewFunc("sqrt", discriminantSimplified)
	twoA := ast.NewMul(ast.NewInt(2), a)

	// Solution 1: (-b + √discriminant) / (2a)
	sol1Num := ast.NewAdd(ast.NewMul(ast.NewInt(-1), b), sqrt)
	sol1 := ast.NewMul(sol1Num, ast.NewPow(twoA, ast.NewInt(-1)))
	sol1Simplified := simplify.Simplify(sol1)

	// Solution 2: (-b - √discriminant) / (2a)
	sol2Num := ast.NewAdd(ast.NewMul(ast.NewInt(-1), b), ast.NewMul(ast.NewInt(-1), sqrt))
	sol2 := ast.NewMul(sol2Num, ast.NewPow(twoA, ast.NewInt(-1)))
	sol2Simplified := simplify.Simplify(sol2)

	solutions := []Solution{
		{
			Variable: opts.Variable,
			Value:    sol1Simplified,
			IsReal:   true,
			IsExact:  true,
		},
		{
			Variable: opts.Variable,
			Value:    sol2Simplified,
			IsReal:   true,
			IsExact:  true,
		},
	}

	// Check for repeated root (discriminant = 0)
	if discVal.Sign() == 0 {
		return SolutionSet{
			Solutions:    solutions[:1], // Only one unique solution
			Message:      "Quadratic equation solved (repeated root)",
			HasSolutions: true,
		}
	}

	return SolutionSet{
		Solutions:    solutions,
		Message:      "Quadratic equation solved",
		HasSolutions: true,
	}
}

// extractQuadraticCoefficients extracts a, b, c from ax² + bx + c
func extractQuadraticCoefficients(expr ast.Expr, variable string) (a, b, c ast.Expr) {
	// Initialize coefficients
	a = ast.NewInt(0)
	b = ast.NewInt(0)
	c = ast.NewInt(0)

	// This is a simplified implementation
	// A complete implementation would need more sophisticated term analysis
	switch e := expr.(type) {
	case *ast.Add:
		for _, term := range e.Terms() {
			degree := getExpressionDegree(term, variable)
			switch degree {
			case 2:
				coeff := extractCoefficientForDegree(term, variable, 2)
				a = ast.NewAdd(a, coeff)
			case 1:
				coeff := extractCoefficientForDegree(term, variable, 1)
				b = ast.NewAdd(b, coeff)
			case 0:
				c = ast.NewAdd(c, term)
			}
		}
	default:
		degree := getExpressionDegree(expr, variable)
		switch degree {
		case 2:
			a = extractCoefficientForDegree(expr, variable, 2)
		case 1:
			b = extractCoefficientForDegree(expr, variable, 1)
		case 0:
			c = expr
		}
	}

	return a, b, c
}

// extractCoefficientForDegree extracts coefficient for a specific degree term
func extractCoefficientForDegree(term ast.Expr, variable string, targetDegree int) ast.Expr {
	degree := getExpressionDegree(term, variable)
	if degree != targetDegree {
		return ast.NewInt(0)
	}

	switch t := term.(type) {
	case *ast.Var:
		if t.Name() == variable && targetDegree == 1 {
			return ast.NewInt(1)
		}
		return ast.NewInt(0)
	case *ast.Pow:
		if base, ok := t.Base().(*ast.Var); ok && base.Name() == variable {
			if expInt, ok := t.Exponent().(*ast.Int); ok {
				val, _ := expInt.Eval(make(map[string]*big.Float))
				intVal, _ := val.Int64()
				if int(intVal) == targetDegree {
					return ast.NewInt(1)
				}
			}
		}
		return ast.NewInt(0)
	case *ast.Mul:
		var coeff ast.Expr = ast.NewInt(1)
		varPower := 0

		for _, factor := range t.Terms() {
			if v, ok := factor.(*ast.Var); ok && v.Name() == variable {
				varPower++
			} else if pow, ok := factor.(*ast.Pow); ok {
				if base, ok := pow.Base().(*ast.Var); ok && base.Name() == variable {
					if expInt, ok := pow.Exponent().(*ast.Int); ok {
						val, _ := expInt.Eval(make(map[string]*big.Float))
						intVal, _ := val.Int64()
						varPower += int(intVal)
					}
				} else if !containsVariable(factor, variable) {
					coeff = ast.NewMul(coeff, factor)
				}
			} else if !containsVariable(factor, variable) {
				coeff = ast.NewMul(coeff, factor)
			}
		}

		if varPower == targetDegree {
			return coeff
		}
		return ast.NewInt(0)
	default:
		return ast.NewInt(0)
	}
}

// Placeholder implementations for higher-order equations
func solveCubic(expr ast.Expr, opts SolveOptions) SolutionSet {
	return SolutionSet{
		Message:      "Cubic equation solving not yet implemented",
		HasSolutions: false,
	}
}

func solveQuartic(expr ast.Expr, opts SolveOptions) SolutionSet {
	return SolutionSet{
		Message:      "Quartic equation solving not yet implemented",
		HasSolutions: false,
	}
}

func solveGeneral(expr ast.Expr, opts SolveOptions) SolutionSet {
	return SolutionSet{
		Message:      "General equation solving not yet implemented",
		HasSolutions: false,
	}
}

// Helper function to check if a solution is valid
func validateSolution(solution ast.Expr, originalExpr ast.Expr, variable string) bool {
	// Substitute solution back into original expression
	vars := make(map[string]*big.Float)

	// Try to evaluate the solution
	solVal, err := solution.Eval(vars)
	if err != nil {
		return false
	}

	vars[variable] = solVal
	result, err := originalExpr.Eval(vars)
	if err != nil {
		return false
	}

	// Check if result is close to zero (within tolerance)
	tolerance := big.NewFloat(1e-10)
	return result.Abs(result).Cmp(tolerance) <= 0
}
