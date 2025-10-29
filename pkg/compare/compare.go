// Package compare provides utilities for comparing mathematical expressions.
package compare

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"time"

	"github.com/quizizz/cas/pkg/ast"
	"github.com/quizizz/cas/pkg/simplify"
)

const (
	// ITERATIONS is the number of test points to use for comparison
	ITERATIONS = 12
	// TOLERANCE_EXP is the exponent for tolerance (10^-TOLERANCE_EXP)
	TOLERANCE_EXP = 9
	// MAX_EXPONENT_BITS limits exponent size to prevent slow evaluations (2^10 = 1024)
	MAX_EXPONENT_BITS = 10
)

// ComparisonResult contains the result of comparing two expressions
type ComparisonResult struct {
	Equal   bool
	Message string
	Details map[string]interface{}
}

// Options controls comparison behavior
type Options struct {
	// CheckForm verifies that expressions have the same structural form
	CheckForm bool
	// CheckSimplified ensures expressions are in simplified form
	CheckSimplified bool
	// Tolerance for numeric comparisons
	Tolerance float64
	// Variables to check for consistency
	RequireVariables []string
}

// DefaultOptions returns the default comparison options
func DefaultOptions() Options {
	return Options{
		CheckForm:        false,
		CheckSimplified:  false,
		Tolerance:        math.Pow(10, -TOLERANCE_EXP), // 1e-9, matching Node.js
		RequireVariables: nil,
	}
}

// Compare compares two mathematical expressions for equivalence
func Compare(expr1, expr2 ast.Expr, opts ...Options) ComparisonResult {
	options := DefaultOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	// 1. Check variables are consistent
	vars1 := expr1.Variables()
	vars2 := expr2.Variables()

	if !sameVariables(vars1, vars2) {
		return ComparisonResult{
			Equal: false,
			Message: fmt.Sprintf("Different variables: expr1 has %v, expr2 has %v",
				vars1, vars2),
			Details: map[string]interface{}{
				"expr1_vars": vars1,
				"expr2_vars": vars2,
			},
		}
	}

	// Check for specific required variables
	if options.RequireVariables != nil {
		for _, requiredVar := range options.RequireVariables {
			if !containsVariable(vars1, requiredVar) {
				return ComparisonResult{
					Equal:   false,
					Message: fmt.Sprintf("Missing required variable: %s", requiredVar),
					Details: map[string]interface{}{
						"missing_variable": requiredVar,
					},
				}
			}
		}
	}

	// 2. Early check for expressions with extremely large exponents before simplification
	if hasLargeExponents(expr1) || hasLargeExponents(expr2) {
		return ComparisonResult{
			Equal:   false,
			Message: "Expressions contain very large exponents - cannot evaluate safely",
			Details: map[string]interface{}{
				"expr1": expr1.String(),
				"expr2": expr2.String(),
			},
		}
	}

	// 3. Check if expressions are identical in structure
	if expr1.String() == expr2.String() {
		return ComparisonResult{
			Equal:   true,
			Message: "Expressions are structurally identical",
		}
	}

	// 4. Check semantic equivalence by simplifying both
	simplified1 := simplify.Simplify(expr1)
	simplified2 := simplify.Simplify(expr2)

	if simplified1.String() == simplified2.String() {
		return ComparisonResult{
			Equal:   true,
			Message: "Expressions are semantically equivalent",
			Details: map[string]interface{}{
				"simplified1": simplified1.String(),
				"simplified2": simplified2.String(),
			},
		}
	}

	// 5. Check numeric equivalence by evaluation
	if len(vars1) > 0 {
		result := checkNumericEquivalence(expr1, expr2, vars1, options.Tolerance)
		if result.Equal {
			return result
		}
	} else {
		// No variables - direct evaluation
		val1, err1 := expr1.Eval(make(map[string]*big.Float))
		val2, err2 := expr2.Eval(make(map[string]*big.Float))

		if err1 == nil && err2 == nil {
			diff := new(big.Float).Sub(val1, val2)
			diff.Abs(diff)
			tolerance := big.NewFloat(options.Tolerance)

			if diff.Cmp(tolerance) <= 0 {
				return ComparisonResult{
					Equal:   true,
					Message: "Expressions are numerically equivalent",
					Details: map[string]interface{}{
						"value1":     val1.Text('g', -1),
						"value2":     val2.Text('g', -1),
						"difference": diff.Text('g', -1),
					},
				}
			}
		}
	}

	// 6. Optional form check
	if options.CheckForm {
		if !checkSameForm(expr1, expr2) {
			return ComparisonResult{
				Equal:   false,
				Message: "Expressions do not have the same form",
				Details: map[string]interface{}{
					"form1": getExpressionForm(expr1),
					"form2": getExpressionForm(expr2),
				},
			}
		}
	}

	// 7. Optional simplification check
	if options.CheckSimplified {
		if !isSimplified(expr1) {
			return ComparisonResult{
				Equal:   false,
				Message: "First expression is not in simplified form",
				Details: map[string]interface{}{
					"simplified_form": simplified1.String(),
				},
			}
		}
		if !isSimplified(expr2) {
			return ComparisonResult{
				Equal:   false,
				Message: "Second expression is not in simplified form",
				Details: map[string]interface{}{
					"simplified_form": simplified2.String(),
				},
			}
		}
	}

	return ComparisonResult{
		Equal:   false,
		Message: "Expressions are not equivalent",
		Details: map[string]interface{}{
			"expr1":       expr1.String(),
			"expr2":       expr2.String(),
			"simplified1": simplified1.String(),
			"simplified2": simplified2.String(),
		},
	}
}

// StructurallyEqual checks if two expressions have exactly the same structure
func StructurallyEqual(expr1, expr2 ast.Expr) bool {
	return expr1.String() == expr2.String()
}

// SemanticallyEqual checks if two expressions are mathematically equivalent
func SemanticallyEqual(expr1, expr2 ast.Expr) bool {
	// First check variables are consistent
	vars1 := expr1.Variables()
	vars2 := expr2.Variables()

	if !sameVariables(vars1, vars2) {
		return false
	}

	// Try simplification first
	simplified1 := simplify.Simplify(expr1)
	simplified2 := simplify.Simplify(expr2)
	if simplified1.String() == simplified2.String() {
		return true
	}

	// Fall back to numeric comparison like the Node.js implementation
	if len(vars1) == 0 {
		// No variables - direct comparison
		val1, err1 := expr1.Eval(make(map[string]*big.Float))
		val2, err2 := expr2.Eval(make(map[string]*big.Float))

		if err1 != nil || err2 != nil {
			return false
		}

		diff := new(big.Float).Sub(val1, val2)
		diff.Abs(diff)
		tolerance := big.NewFloat(math.Pow(10, -TOLERANCE_EXP))
		return diff.Cmp(tolerance) <= 0
	}

	// Test with multiple variable values (similar to Node.js compare method)
	result := checkNumericEquivalence(expr1, expr2, vars1, math.Pow(10, -TOLERANCE_EXP))
	return result.Equal
}

// NumericallyEqual checks if expressions evaluate to the same value(s)
func NumericallyEqual(expr1, expr2 ast.Expr, tolerance float64) bool {
	vars1 := expr1.Variables()
	vars2 := expr2.Variables()

	if !sameVariables(vars1, vars2) {
		return false
	}

	if len(vars1) == 0 {
		// No variables - direct comparison
		val1, err1 := expr1.Eval(make(map[string]*big.Float))
		val2, err2 := expr2.Eval(make(map[string]*big.Float))

		if err1 != nil || err2 != nil {
			return false
		}

		diff := new(big.Float).Sub(val1, val2)
		diff.Abs(diff)
		return diff.Cmp(big.NewFloat(tolerance)) <= 0
	}

	// Test with multiple variable values
	result := checkNumericEquivalence(expr1, expr2, vars1, tolerance)
	return result.Equal
}

// Helper functions

func sameVariables(vars1, vars2 []string) bool {
	if len(vars1) != len(vars2) {
		return false
	}

	// Convert to sets for comparison
	set1 := make(map[string]bool)
	set2 := make(map[string]bool)

	for _, v := range vars1 {
		set1[v] = true
	}
	for _, v := range vars2 {
		set2[v] = true
	}

	for v := range set1 {
		if !set2[v] {
			return false
		}
	}
	for v := range set2 {
		if !set1[v] {
			return false
		}
	}

	return true
}

func containsVariable(vars []string, target string) bool {
	for _, v := range vars {
		if v == target {
			return true
		}
	}
	return false
}

func checkNumericEquivalence(expr1, expr2 ast.Expr, vars []string, tolerance float64) ComparisonResult {
	// Use a seeded random generator for reproducible results within a test run
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Compare at ITERATIONS number of points to determine equality
	// Similar to the Node.js implementation
	for i := 0; i < ITERATIONS; i++ {
		varMap := make(map[string]*big.Float)

		// One third total iterations each with range 10, 100, and 1000
		rangeExp := 1 + int(math.Floor(3*float64(i)/float64(ITERATIONS)))
		valueRange := math.Pow(10, float64(rangeExp))

		// Half of the iterations should only use integer values
		// This is important for expressions like (-2)^x that result in NaN
		// with non-integer values of x in JavaScript
		useFloats := i%2 == 0

		for _, varName := range vars {
			var value float64
			if useFloats {
				// Generate random float in range [-valueRange, valueRange]
				value = (rng.Float64()*2 - 1) * valueRange
			} else {
				// Generate random integer in range [-valueRange, valueRange]
				value = float64(rng.Intn(int(2*valueRange)+1) - int(valueRange))
			}
			varMap[varName] = big.NewFloat(value)
		}

		val1, err1 := expr1.Eval(varMap)
		val2, err2 := expr2.Eval(varMap)

		// Handle evaluation errors - if both expressions fail to evaluate
		// at the same point, they might still be equivalent
		if err1 != nil && err2 != nil {
			continue // Both failed - skip this test point
		}
		if err1 != nil || err2 != nil {
			// Only one failed - expressions are different
			return ComparisonResult{
				Equal:   false,
				Message: fmt.Sprintf("Expressions differ in evaluation success"),
				Details: map[string]interface{}{
					"iteration":   i,
					"variables":   varMap,
					"expr1_error": err1 != nil,
					"expr2_error": err2 != nil,
				},
			}
		}

		// Check for NaN values - if both are NaN, continue
		val1Float, _ := val1.Float64()
		val2Float, _ := val2.Float64()
		if math.IsNaN(val1Float) && math.IsNaN(val2Float) {
			continue
		}
		if math.IsNaN(val1Float) || math.IsNaN(val2Float) {
			return ComparisonResult{
				Equal:   false,
				Message: fmt.Sprintf("One expression evaluates to NaN while the other doesn't"),
				Details: map[string]interface{}{
					"iteration": i,
					"variables": varMap,
					"expr1_nan": math.IsNaN(val1Float),
					"expr2_nan": math.IsNaN(val2Float),
				},
			}
		}

		// Check for infinity values
		if math.IsInf(val1Float, 0) && math.IsInf(val2Float, 0) {
			if math.IsInf(val1Float, 1) == math.IsInf(val2Float, 1) {
				continue // Both are same type of infinity
			}
		}

		// Perform numeric comparison with tolerance
		diff := new(big.Float).Sub(val1, val2)
		diff.Abs(diff)

		// Use relative comparison for large numbers, absolute for small ones
		var toleranceValue *big.Float
		if math.Abs(val1Float) < 1 || math.Abs(val2Float) < 1 {
			toleranceValue = big.NewFloat(tolerance)
		} else {
			// Relative tolerance
			max := big.NewFloat(math.Max(math.Abs(val1Float), math.Abs(val2Float)))
			toleranceValue = new(big.Float).Mul(max, big.NewFloat(tolerance))
		}

		if diff.Cmp(toleranceValue) > 0 {
			return ComparisonResult{
				Equal:   false,
				Message: fmt.Sprintf("Expressions differ at test point %d", i),
				Details: map[string]interface{}{
					"iteration":    i,
					"variables":    formatVarMap(varMap),
					"expr1_result": val1.Text('g', -1),
					"expr2_result": val2.Text('g', -1),
					"difference":   diff.Text('g', -1),
					"tolerance":    toleranceValue.Text('g', -1),
				},
			}
		}
	}

	return ComparisonResult{
		Equal:   true,
		Message: "Expressions are numerically equivalent for all test points",
		Details: map[string]interface{}{
			"iterations_tested": ITERATIONS,
		},
	}
}

func checkSameForm(expr1, expr2 ast.Expr) bool {
	return getExpressionForm(expr1) == getExpressionForm(expr2)
}

func getExpressionForm(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Add:
		return fmt.Sprintf("Add[%d]", len(e.Terms()))
	case *ast.Mul:
		return fmt.Sprintf("Mul[%d]", len(e.Terms()))
	case *ast.Pow:
		return "Pow"
	case *ast.Var:
		return "Var"
	case *ast.Func:
		return fmt.Sprintf("Func[%s]", e.Name())
	case ast.Numeric:
		return "Numeric"
	default:
		return "Unknown"
	}
}

func isSimplified(expr ast.Expr) bool {
	simplified := simplify.Simplify(expr)
	return expr.String() == simplified.String()
}

// formatVarMap converts a variable map to a readable string representation
func formatVarMap(varMap map[string]*big.Float) map[string]string {
	result := make(map[string]string)
	for k, v := range varMap {
		result[k] = v.Text('g', -1)
	}
	return result
}

// isLargeExponent checks if a big.Float represents an exponent that's too large to evaluate safely
func isLargeExponent(val *big.Float) bool {
	intVal := new(big.Int)
	if val.IsInt() {
		val.Int(intVal)
		return intVal.BitLen() > MAX_EXPONENT_BITS
	}
	return false
}

// isLargeComputedExponent tries to evaluate an expression and checks if the result is too large for exponentiation
func isLargeComputedExponent(expr ast.Expr) bool {
	emptyVars := make(map[string]*big.Float)
	if result, err := expr.Eval(emptyVars); err == nil {
		return isLargeExponent(result)
	}
	return false
}

// hasLargeExponents checks if an expression contains power operations with extremely large exponents
// that would be prohibitively slow to evaluate numerically
func hasLargeExponents(expr ast.Expr) bool {
	if expr == nil {
		return false
	}

	switch e := expr.(type) {
	case *ast.Add:
		for _, term := range e.Terms() {
			if hasLargeExponents(term) {
				return true
			}
		}
		return false

	case *ast.Mul:
		for _, term := range e.Terms() {
			if hasLargeExponents(term) {
				return true
			}
		}
		return false

	case *ast.Pow:
		exp := e.Exponent()

		// Check if the exponent is a very large integer
		if intExp, ok := exp.(*ast.Int); ok {
			if isLargeExponent(intExp.Value()) {
				return true
			}
		}

		// Check if exponent is a power or multiplication expression that evaluates to a large number
		if _, ok := exp.(*ast.Pow); ok {
			if isLargeComputedExponent(exp) {
				return true
			}
		}
		if _, ok := exp.(*ast.Mul); ok {
			if isLargeComputedExponent(exp) {
				return true
			}
		}

		// Recursively check base and exponent
		return hasLargeExponents(e.Base()) || hasLargeExponents(exp)

	case *ast.Func:
		for _, arg := range e.Args() {
			if hasLargeExponents(arg) {
				return true
			}
		}
		return false

	default:
		return false
	}
}
