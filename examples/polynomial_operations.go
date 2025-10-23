// Example demonstrating polynomial operations: expansion and solving
package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/quizizz/cas/pkg/ast"

	"github.com/quizizz/cas/pkg/expand"
	"github.com/quizizz/cas/pkg/latex"
	"github.com/quizizz/cas/pkg/parser"
	"github.com/quizizz/cas/pkg/solve"
)

func main() {
	fmt.Println("=== Polynomial Operations Examples ===\n")

	// Example 1: Basic expansion
	fmt.Println("1. Basic Polynomial Expansion")
	expansionExamples := []string{
		"(x + 1)^2",
		"(x - 2)^2",
		"(x + 1)^3",
		"(x + 2)*(x - 3)",
		"(x + 1)*(x + 2)*(x + 3)",
		"(2*x + 3)^2",
	}

	for _, exprStr := range expansionExamples {
		expr, err := parser.Parse(exprStr)
		if err != nil {
			log.Printf("Parse error for %s: %v", exprStr, err)
			continue
		}

		expanded := expand.Expand(expr)
		fmt.Printf("(%s) = %s\n", exprStr, expanded.String())
		fmt.Printf("LaTeX: %s = %s\n", latex.Format(expr), latex.Format(expanded))
		fmt.Println()
	}

	// Example 2: Full expansion vs basic expansion
	fmt.Println("2. Full Expansion vs Basic Expansion")
	complexExpr, _ := parser.Parse("(x + 1)^2 * (x - 2)")

	basicExpanded := expand.Expand(complexExpr)
	fullyExpanded := expand.ExpandFully(complexExpr)

	fmt.Printf("Original: %s\n", complexExpr.String())
	fmt.Printf("Basic Expansion: %s\n", basicExpanded.String())
	fmt.Printf("Full Expansion: %s\n", fullyExpanded.String())
	fmt.Printf("LaTeX: %s\n", latex.Format(fullyExpanded))
	fmt.Println()

	// Example 3: Solving linear equations
	fmt.Println("3. Linear Equation Solving")
	linearEquations := []string{
		"2*x - 4",  // 2x - 4 = 0
		"3*x + 6",  // 3x + 6 = 0
		"x - 1",    // x - 1 = 0
		"5*x + 10", // 5x + 10 = 0
		"-x + 7",   // -x + 7 = 0
	}

	for _, eqStr := range linearEquations {
		expr, _ := parser.Parse(eqStr)
		solutions := solve.Solve(expr)

		fmt.Printf("Solve: %s = 0\n", eqStr)
		if solutions.HasSolutions && len(solutions.Solutions) > 0 {
			for _, sol := range solutions.Solutions {
				fmt.Printf("  %s = %s", sol.Variable, sol.Value.String())
				if sol.IsExact {
					fmt.Printf(" (exact)")
				}
				fmt.Println()

				// Show LaTeX
				fmt.Printf("  LaTeX: %s = %s\n", sol.Variable, latex.Format(sol.Value))

				// Verify the solution
				if verified := verifySolution(sol.Value, expr, sol.Variable); verified {
					fmt.Printf("  ✓ Solution verified\n")
				} else {
					fmt.Printf("  ✗ Solution verification failed\n")
				}
			}
		} else {
			fmt.Printf("  %s\n", solutions.Message)
		}
		fmt.Println()
	}

	// Example 4: Solving quadratic equations
	fmt.Println("4. Quadratic Equation Solving")
	quadraticEquations := []string{
		"x^2 - 4",       // x² - 4 = 0, solutions: ±2
		"x^2 + 3*x + 2", // x² + 3x + 2 = 0, solutions: -1, -2
		"x^2 + 2*x + 1", // (x + 1)² = 0, solution: -1 (repeated)
		"x^2 - 5*x + 6", // x² - 5x + 6 = 0, solutions: 2, 3
		"2*x^2 - 8",     // 2x² - 8 = 0, solutions: ±2
		"x^2 + 1",       // x² + 1 = 0, no real solutions
	}

	for _, eqStr := range quadraticEquations {
		expr, _ := parser.Parse(eqStr)
		solutions := solve.Solve(expr)

		fmt.Printf("Solve: %s = 0\n", eqStr)
		fmt.Printf("Status: %s\n", solutions.Message)

		if solutions.HasSolutions && len(solutions.Solutions) > 0 {
			fmt.Printf("Solutions:\n")
			for i, sol := range solutions.Solutions {
				fmt.Printf("  x_%d = %s", i+1, sol.Value.String())
				if sol.IsExact {
					fmt.Printf(" (exact)")
				}
				if !sol.IsReal {
					fmt.Printf(" (complex)")
				}
				fmt.Println()

				// Show LaTeX
				fmt.Printf("  LaTeX: x_{%d} = %s\n", i+1, latex.Format(sol.Value))

				// Try numerical evaluation
				if numVal, err := sol.Value.Eval(make(map[string]*big.Float)); err == nil {
					fmt.Printf("  Numerical: x_%d ≈ %s\n", i+1, numVal.Text('g', 6))
				}
			}
		} else {
			fmt.Printf("No real solutions found\n")
		}
		fmt.Println()
	}

	// Example 5: Equation form solving (lhs = rhs)
	fmt.Println("5. Equation Form Solving (lhs = rhs)")
	equationPairs := []struct {
		lhs, rhs string
	}{
		{"x + 1", "3"},   // x + 1 = 3
		{"2*x", "x + 5"}, // 2x = x + 5
		{"x^2", "4"},     // x² = 4
		{"x^2 + x", "2"}, // x² + x = 2
	}

	for _, eq := range equationPairs {
		lhs, _ := parser.Parse(eq.lhs)
		rhs, _ := parser.Parse(eq.rhs)
		solutions := solve.SolveEquation(lhs, rhs)

		fmt.Printf("Solve: %s = %s\n", eq.lhs, eq.rhs)
		fmt.Printf("LaTeX: %s = %s\n", latex.Format(lhs), latex.Format(rhs))

		if solutions.HasSolutions && len(solutions.Solutions) > 0 {
			for i, sol := range solutions.Solutions {
				fmt.Printf("  Solution %d: %s = %s\n", i+1, sol.Variable, sol.Value.String())
				fmt.Printf("  LaTeX: %s = %s\n", sol.Variable, latex.Format(sol.Value))

				// Verify by substitution
				if verified := verifyEquationSolution(sol.Value, lhs, rhs, sol.Variable); verified {
					fmt.Printf("  ✓ Verified by substitution\n")
				}
			}
		} else {
			fmt.Printf("  %s\n", solutions.Message)
		}
		fmt.Println()
	}

	// Example 6: Complex polynomial analysis
	fmt.Println("6. Complex Polynomial Analysis")
	complexPoly, _ := parser.Parse("x^4 - 5*x^2 + 4")

	fmt.Printf("Polynomial: %s\n", complexPoly.String())
	fmt.Printf("LaTeX: %s\n", latex.Format(complexPoly))

	// Expand if needed
	expanded := expand.Expand(complexPoly)
	if expanded.String() != complexPoly.String() {
		fmt.Printf("Expanded: %s\n", expanded.String())
	}

	// Try to solve
	solutions := solve.Solve(complexPoly)
	fmt.Printf("Solving: %s = 0\n", complexPoly.String())
	fmt.Printf("Status: %s\n", solutions.Message)

	if solutions.HasSolutions {
		for i, sol := range solutions.Solutions {
			fmt.Printf("  x_%d = %s\n", i+1, sol.Value.String())
		}
	}
	fmt.Println()
}

// Helper function to verify a solution by substitution
func verifySolution(solExpr, originalExpr ast.Expr, variable string) bool {
	vars := make(map[string]*big.Float)

	// Parse and evaluate the solution
	//solExpr, err := parser.Parse(solution)
	//if err != nil {
	//	return false
	//}

	solVal, err := solExpr.Eval(make(map[string]*big.Float))
	if err != nil {
		return false
	}

	// Substitute back into original expression
	vars[variable] = solVal
	result, err := originalExpr.Eval(vars)
	if err != nil {
		return false
	}

	// Check if result is close to zero
	tolerance := big.NewFloat(1e-10)
	return result.Abs(result).Cmp(tolerance) <= 0
}

// Helper function to verify equation solution (lhs = rhs)
func verifyEquationSolution(solutionExpr, lhs, rhs ast.Expr, variable string) bool {
	vars := make(map[string]*big.Float)

	// Evaluate the solution
	solVal, err := solutionExpr.Eval(make(map[string]*big.Float))
	if err != nil {
		return false
	}

	// Substitute into both sides
	vars[variable] = solVal

	lhsResult, err := lhs.Eval(vars)
	if err != nil {
		return false
	}

	rhsResult, err := rhs.Eval(vars)
	if err != nil {
		return false
	}

	// Check if lhs ≈ rhs
	tolerance := big.NewFloat(1e-10)
	diff := big.NewFloat(0).Sub(lhsResult, rhsResult)
	return diff.Abs(diff).Cmp(tolerance) <= 0
}
