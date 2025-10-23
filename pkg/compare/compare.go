// Package compare provides utilities for comparing mathematical expressions.
package compare

import (
	"fmt"
	"math/big"

	"github.com/quizizz/cas/pkg/ast"
	"github.com/quizizz/cas/pkg/simplify"
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
		Tolerance:        1e-10,
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

	// 2. Check if expressions are identical in structure
	if expr1.String() == expr2.String() {
		return ComparisonResult{
			Equal:   true,
			Message: "Expressions are structurally identical",
		}
	}

	// 3. Check semantic equivalence by simplifying both
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

	// 4. Check numeric equivalence by evaluation
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

	// 5. Optional form check
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

	// 6. Optional simplification check
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
	simplified1 := simplify.Simplify(expr1)
	simplified2 := simplify.Simplify(expr2)
	return simplified1.String() == simplified2.String()
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
	// Test with multiple variable assignments
	testValues := []float64{0, 1, -1, 2, -2, 0.5, -0.5, 10}

	for _, testVal := range testValues {
		varMap := make(map[string]*big.Float)
		for _, varName := range vars {
			varMap[varName] = big.NewFloat(testVal)
		}

		val1, err1 := expr1.Eval(varMap)
		val2, err2 := expr2.Eval(varMap)

		if err1 != nil || err2 != nil {
			continue // Skip this test value
		}

		diff := new(big.Float).Sub(val1, val2)
		diff.Abs(diff)
		if diff.Cmp(big.NewFloat(tolerance)) > 0 {
			return ComparisonResult{
				Equal:   false,
				Message: fmt.Sprintf("Expressions differ at %v = %f", vars, testVal),
				Details: map[string]interface{}{
					"test_value":   testVal,
					"expr1_result": val1.Text('g', -1),
					"expr2_result": val2.Text('g', -1),
					"difference":   diff.Text('g', -1),
				},
			}
		}
	}

	return ComparisonResult{
		Equal:   true,
		Message: "Expressions are numerically equivalent for tested values",
		Details: map[string]interface{}{
			"test_values": testValues,
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
