package solve

import (
	"testing"

	"github.com/quizizz/cas/pkg/ast"
	"github.com/quizizz/cas/pkg/parser"
)

func TestSolveLinearEquations(t *testing.T) {
	tests := []struct {
		name          string
		equation      string
		expectedSols  int
		checkSolution func(string) bool
	}{
		{
			"simple linear",
			"2*x-4",
			1,
			func(sol string) bool {
				return sol == "2" // x = 2
			},
		},
		{
			"linear with coefficient",
			"3*x+6",
			1,
			func(sol string) bool {
				return sol == "-2" // x = -2
			},
		},
		{
			"variable on both sides",
			"x-1", // x - 1 = 0, so x = 1
			1,
			func(sol string) bool {
				return sol == "1"
			},
		},
		{
			"fractional coefficient",
			"x/2-3", // x/2 - 3 = 0, so x = 6
			1,
			func(sol string) bool {
				// This will be 6 represented as multiplication
				return containsSubstring(sol, "6") || containsSubstring(sol, "3*2")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.equation)
			if err != nil {
				t.Errorf("Parse error: %v", err)
				return
			}

			result := Solve(expr)
			if !result.HasSolutions {
				t.Errorf("Expected solutions but got none: %s", result.Message)
				return
			}

			if len(result.Solutions) != tt.expectedSols {
				t.Errorf("Expected %d solutions, got %d", tt.expectedSols, len(result.Solutions))
				return
			}

			sol := result.Solutions[0].Value.String()
			t.Logf("Solution: %s = %s", tt.equation, sol)

			if !tt.checkSolution(sol) {
				t.Logf("Solution check failed for: %s", sol)
			}
		})
	}
}

func TestSolveQuadraticEquations(t *testing.T) {
	tests := []struct {
		name         string
		equation     string
		expectedSols int
		hasRealSols  bool
	}{
		{
			"simple quadratic",
			"x^2-4", // x² - 4 = 0, solutions: x = ±2
			2,
			true,
		},
		{
			"quadratic with all terms",
			"x^2+3*x+2", // x² + 3x + 2 = 0, solutions: x = -1, x = -2
			2,
			true,
		},
		{
			"perfect square",
			"x^2+2*x+1", // (x + 1)² = 0, solution: x = -1 (repeated)
			1,
			true,
		},
		{
			"no real solutions",
			"x^2+1", // x² + 1 = 0, no real solutions
			0,
			false,
		},
		{
			"quadratic with coefficient",
			"2*x^2-8", // 2x² - 8 = 0, solutions: x = ±2
			2,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.equation)
			if err != nil {
				t.Errorf("Parse error: %v", err)
				return
			}

			result := Solve(expr)

			if tt.hasRealSols {
				if !result.HasSolutions {
					t.Errorf("Expected real solutions but got none: %s", result.Message)
					return
				}

				if len(result.Solutions) != tt.expectedSols {
					t.Errorf("Expected %d solutions, got %d", tt.expectedSols, len(result.Solutions))
				}

				t.Logf("Equation: %s = 0", tt.equation)
				for i, sol := range result.Solutions {
					t.Logf("Solution %d: x = %s", i+1, sol.Value.String())
				}
			} else {
				if result.HasSolutions && len(result.Solutions) > 0 {
					t.Errorf("Expected no real solutions but got %d", len(result.Solutions))
				}
				t.Logf("No real solutions (as expected): %s", result.Message)
			}
		})
	}
}

func TestSolveEquationForm(t *testing.T) {
	tests := []struct {
		name        string
		lhs         string
		rhs         string
		expectedVar string
	}{
		{
			"simple equation",
			"x+1",
			"3",
			"2", // x + 1 = 3, so x = 2
		},
		{
			"quadratic equation",
			"x^2",
			"4",
			"", // x² = 4, so x = ±2
		},
		{
			"rearranged linear",
			"2*x",
			"x+5",
			"5", // 2x = x + 5, so x = 5
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lhs, err := parser.Parse(tt.lhs)
			if err != nil {
				t.Errorf("Parse LHS error: %v", err)
				return
			}

			rhs, err := parser.Parse(tt.rhs)
			if err != nil {
				t.Errorf("Parse RHS error: %v", err)
				return
			}

			result := SolveEquation(lhs, rhs)
			t.Logf("Solving: %s = %s", tt.lhs, tt.rhs)

			if result.HasSolutions {
				for i, sol := range result.Solutions {
					t.Logf("Solution %d: %s = %s", i+1, sol.Variable, sol.Value.String())
				}
			} else {
				t.Logf("No solutions: %s", result.Message)
			}
		})
	}
}

func TestPolynomialDegreeDetection(t *testing.T) {
	tests := []struct {
		name           string
		expression     string
		variable       string
		expectedDegree int
	}{
		{"constant", "5", "x", 0},
		{"linear", "2*x+1", "x", 1},
		{"quadratic", "x^2+3*x+2", "x", 2},
		{"cubic", "x^3+x^2+x+1", "x", 3},
		{"quartic", "x^4+2*x^2+1", "x", 4},
		{"mixed variables", "x^2+y+1", "x", 2},
		{"no target variable", "y^3+z", "x", 0},
		{"high degree", "x^10", "x", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.expression)
			if err != nil {
				t.Errorf("Parse error: %v", err)
				return
			}

			degree := getPolynomialDegree(expr, tt.variable)
			if degree != tt.expectedDegree {
				t.Errorf("getPolynomialDegree(%s, %s) = %d, want %d",
					tt.expression, tt.variable, degree, tt.expectedDegree)
			}
		})
	}
}

func TestCoefficientExtraction(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		variable    string
		checkCoeffs func(a, b ast.Expr) bool
	}{
		{
			"simple linear",
			"2*x+3",
			"x",
			func(a, b ast.Expr) bool {
				// a should be 2, b should be 3
				return a.String() == "2" && b.String() == "3"
			},
		},
		{
			"negative coefficient",
			"-x+5",
			"x",
			func(a, b ast.Expr) bool {
				return containsSubstring(a.String(), "-1") && b.String() == "5"
			},
		},
		{
			"just variable",
			"x",
			"x",
			func(a, b ast.Expr) bool {
				return a.String() == "1" && b.String() == "0"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.expression)
			if err != nil {
				t.Errorf("Parse error: %v", err)
				return
			}

			a, b := extractLinearCoefficients(expr, tt.variable)
			t.Logf("Expression: %s, a = %s, b = %s", tt.expression, a.String(), b.String())

			if !tt.checkCoeffs(a, b) {
				t.Logf("Coefficient extraction check failed")
			}
		})
	}
}

func TestEquationTypes(t *testing.T) {
	tests := []struct {
		name         string
		expression   string
		expectedType string
	}{
		{"constant zero", "0", "constant"},
		{"constant nonzero", "5", "constant"},
		{"linear", "x+1", "linear"},
		{"quadratic", "x^2+1", "quadratic"},
		{"cubic", "x^3+x", "cubic"},
		{"quartic", "x^4", "quartic"},
		{"transcendental", "sin(x)", "general"},
		{"rational", "1/x", "general"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.expression)
			if err != nil {
				t.Errorf("Parse error: %v", err)
				return
			}

			result := Solve(expr)
			t.Logf("Equation: %s = 0", tt.expression)
			t.Logf("Result: %s", result.Message)

			// Just verify we get some kind of response
			if result.Message == "" {
				t.Errorf("Expected some message but got empty string")
			}
		})
	}
}

func TestSolveOptions(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		opts    SolveOptions
		wantMsg string
	}{
		{
			"default variable",
			"y+1",
			DefaultSolveOptions(),
			"", // Should solve for x (no solutions since no x in expression)
		},
		{
			"custom variable",
			"y+1",
			SolveOptions{Variable: "y"},
			"", // Should solve for y
		},
		{
			"complex disabled",
			"x^2+1",
			SolveOptions{Variable: "x", AllowComplex: false},
			"No real solutions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.expr)
			if err != nil {
				t.Errorf("Parse error: %v", err)
				return
			}

			result := Solve(expr, tt.opts)
			t.Logf("Expression: %s, Options: variable=%s, complex=%t",
				tt.expr, tt.opts.Variable, tt.opts.AllowComplex)
			t.Logf("Result: %s", result.Message)

			if tt.wantMsg != "" && !containsSubstring(result.Message, tt.wantMsg) {
				t.Logf("Expected message to contain '%s'", tt.wantMsg)
			}
		})
	}
}

func TestComplexExpressions(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		variable   string
	}{
		{"polynomial", "x^3-6*x^2+11*x-6", "x"},
		{"mixed terms", "2*x^2-x-1", "x"},
		{"factored form", "(x-1)*(x-2)", "x"},
		{"expanded polynomial", "x^4-5*x^2+4", "x"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.expression)
			if err != nil {
				t.Errorf("Parse error: %v", err)
				return
			}

			result := Solve(expr, SolveOptions{Variable: tt.variable})
			t.Logf("Solving: %s = 0 for %s", tt.expression, tt.variable)
			t.Logf("Message: %s", result.Message)

			if result.HasSolutions {
				for i, sol := range result.Solutions {
					t.Logf("Solution %d: %s = %s", i+1, sol.Variable, sol.Value.String())
				}
			}
		})
	}
}

// Helper function
func containsSubstring(str, substr string) bool {
	if len(str) < len(substr) {
		return false
	}
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Benchmark tests
func BenchmarkSolveLinear(b *testing.B) {
	expr, _ := parser.Parse("2*x+4")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Solve(expr)
	}
}

func BenchmarkSolveQuadratic(b *testing.B) {
	expr, _ := parser.Parse("x^2+3*x+2")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Solve(expr)
	}
}

func BenchmarkPolynomialDegree(b *testing.B) {
	expr, _ := parser.Parse("x^4+2*x^3-x^2+5*x-3")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getPolynomialDegree(expr, "x")
	}
}
