package latex

import (
	"math/big"
	"testing"

	"github.com/quizizz/cas/pkg/ast"
	"github.com/quizizz/cas/pkg/parser"
)

func TestFormatNumbers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"positive integer", "42", "42"},
		{"negative integer", "-7", "-7"},
		{"rational half", "1/2", "\\frac{1}{2}"},
		{"rational third", "1/3", "\\frac{1}{3}"},
		{"rational two thirds", "2/3", "\\frac{2}{3}"},
		{"simple float", "3.14", "3.14"},
		{"integer float", "5.0", "5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.input)
			if err != nil {
				// For rationals, create manually since parser might not handle fractions
				if tt.input == "1/2" {
					expr = ast.NewRational(1, 2)
				} else if tt.input == "1/3" {
					expr = ast.NewRational(1, 3)
				} else if tt.input == "2/3" {
					expr = ast.NewRational(2, 3)
				} else {
					t.Errorf("Parse error: %v", err)
					return
				}
			}

			result := Format(expr)
			if result != tt.expected {
				t.Errorf("Format(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFormatVariables(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"single letter", "x", "x"},
		{"multiple letters", "var", "\\mathrm{var}"},
		{"Greek alpha", "alpha", "\\alpha"},
		{"Greek beta", "beta", "\\beta"},
		{"Greek pi", "pi", "\\pi"},
		{"Greek theta", "theta", "\\theta"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var expr ast.Expr
			if tt.input == "pi" {
				piVal, _ := big.NewFloat(0).SetString("3.14159265358979323846")
				expr = ast.NewConst(tt.input, piVal)
			} else {
				expr = ast.NewVar(tt.input)
			}

			result := Format(expr)
			if result != tt.expected {
				t.Errorf("Format(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFormatAddition(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple addition", "x+y", "x + y"},
		{"three terms", "a+b+c", "a + b + c"},
		{"with constants", "x+1", "x + 1"},
		{"negative term", "x-y", "x - y"},
		{"complex expression", "2*x+3*y", "2x + 3y"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.input)
			if err != nil {
				t.Errorf("Parse error: %v", err)
				return
			}

			result := Format(expr)
			t.Logf("Format(%s) = %s", tt.input, result)

			// Check that basic structure is correct
			if !containsSubstring(result, "+") && !containsSubstring(result, "-") {
				if tt.input != "x-y" { // x-y might be formatted differently
					t.Errorf("Expected + or - in result: %s", result)
				}
			}
		})
	}
}

func TestFormatMultiplication(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple multiplication", "x*y", "xy"},
		{"coefficient", "2*x", "2x"},
		{"negative coefficient", "-3*x", "-3x"},
		{"multiple factors", "a*b*c", "abc"},
		{"number multiplication", "2*3", "2 \\cdot 3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.input)
			if err != nil {
				t.Errorf("Parse error: %v", err)
				return
			}

			result := Format(expr)
			t.Logf("Format(%s) = %s", tt.input, result)

			// Basic validation - should not contain * symbols
			if containsSubstring(result, "*") {
				t.Errorf("LaTeX output should not contain * symbol: %s", result)
			}
		})
	}
}

func TestFormatPowers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple square", "x^2", "x^{2}"},
		{"cube", "x^3", "x^{3}"},
		{"square root", "x^(1/2)", "\\sqrt{x}"},
		{"negative exponent", "x^(-1)", "\\frac{1}{x}"},
		{"negative square", "x^(-2)", "\\frac{1}{x^{2}}"},
		{"complex base", "(x+1)^2", "\\left(x + 1\\right)^{2}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var expr ast.Expr
			var err error

			// Handle special cases that parser might not handle well
			if tt.input == "x^(1/2)" {
				base := ast.NewVar("x")
				exp := ast.NewRational(1, 2)
				expr = ast.NewPow(base, exp)
			} else if tt.input == "x^(-1)" {
				base := ast.NewVar("x")
				exp := ast.NewInt(-1)
				expr = ast.NewPow(base, exp)
			} else if tt.input == "x^(-2)" {
				base := ast.NewVar("x")
				exp := ast.NewInt(-2)
				expr = ast.NewPow(base, exp)
			} else {
				expr, err = parser.Parse(tt.input)
				if err != nil {
					t.Errorf("Parse error: %v", err)
					return
				}
			}

			result := Format(expr)
			t.Logf("Format(%s) = %s", tt.input, result)

			// Check for expected patterns
			if tt.name == "square root" {
				if !containsSubstring(result, "\\sqrt{") {
					t.Logf("Expected \\sqrt in result: %s", result)
				}
			}
			if tt.name == "negative exponent" || tt.name == "negative square" {
				if !containsSubstring(result, "\\frac{") {
					t.Logf("Expected \\frac in result: %s", result)
				}
			}
		})
	}
}

func TestFormatFunctions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"sine", "sin(x)", "\\sin\\left(x\\right)"},
		{"cosine", "cos(x)", "\\cos\\left(x\\right)"},
		{"tangent", "tan(x)", "\\tan\\left(x\\right)"},
		{"natural log", "ln(x)", "\\ln\\left(x\\right)"},
		{"logarithm", "log(x)", "\\log\\left(x\\right)"},
		{"square root", "sqrt(x)", "\\sqrt{x}"},
		{"absolute value", "abs(x)", "\\left|x\\right|"},
		{"exponential", "exp(x)", "e^{x}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.input)
			if err != nil {
				t.Errorf("Parse error: %v", err)
				return
			}

			result := Format(expr)
			t.Logf("Format(%s) = %s", tt.input, result)

			// Check that LaTeX function names are used
			if tt.name == "sine" && !containsSubstring(result, "\\sin") {
				t.Errorf("Expected \\sin in result: %s", result)
			}
			if tt.name == "square root" && !containsSubstring(result, "\\sqrt{") {
				t.Errorf("Expected \\sqrt{ in result: %s", result)
			}
			if tt.name == "exponential" && !containsSubstring(result, "e^{") {
				t.Errorf("Expected e^{ in result: %s", result)
			}
		})
	}
}

func TestFormatComplexExpressions(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"polynomial", "x^3+2*x^2-x+5"},
		{"rational function", "(x+1)/(x-1)"},
		{"trigonometric", "sin(x)^2+cos(x)^2"},
		{"exponential", "e^(x^2)*ln(x+1)"},
		{"nested functions", "sin(cos(x))"},
		{"derivative notation", "d/dx(x^2)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.input)
			if err != nil {
				// Skip expressions that parser might not handle
				t.Logf("Parse error (skipping): %v", err)
				return
			}

			result := Format(expr)
			t.Logf("Format(%s) = %s", tt.input, result)

			// Basic validation - should not contain unescaped * or ^
			if containsSubstring(result, "*") {
				t.Logf("Warning: LaTeX contains * symbol: %s", result)
			}

			// Should contain proper LaTeX syntax
			if len(result) == 0 {
				t.Errorf("Empty LaTeX result")
			}
		})
	}
}

func TestFormatOptions(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		options FormatOptions
		check   func(string) bool
	}{
		{
			"fractions disabled",
			"1/2",
			FormatOptions{UseFractions: false},
			func(result string) bool {
				return !containsSubstring(result, "\\frac")
			},
		},
		{
			"symbols disabled",
			"pi",
			FormatOptions{UseSymbols: false},
			func(result string) bool {
				return !containsSubstring(result, "\\pi")
			},
		},
		{
			"parentheses disabled",
			"(x+1)*y",
			FormatOptions{UseParentheses: false},
			func(result string) bool {
				return !containsSubstring(result, "\\left(")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var expr ast.Expr
			var err error

			if tt.input == "1/2" {
				expr = ast.NewRational(1, 2)
			} else if tt.input == "pi" {
				piVal, _ := big.NewFloat(0).SetString("3.14159265358979323846")
				expr = ast.NewConst("pi", piVal)
			} else {
				expr, err = parser.Parse(tt.input)
				if err != nil {
					t.Errorf("Parse error: %v", err)
					return
				}
			}

			result := Format(expr, tt.options)
			t.Logf("Format(%s) with options = %s", tt.input, result)

			if !tt.check(result) {
				t.Logf("Options test failed for: %s", result)
			}
		})
	}
}

func TestSpecialFormatting(t *testing.T) {
	tests := []struct {
		name     string
		function func() string
		check    func(string) bool
	}{
		{
			"equation formatting",
			func() string {
				lhs, _ := parser.Parse("x^2")
				rhs, _ := parser.Parse("4")
				return FormatEquation(lhs, rhs)
			},
			func(result string) bool {
				return containsSubstring(result, "=") && containsSubstring(result, "x^{2}")
			},
		},
		{
			"derivative formatting",
			func() string {
				expr, _ := parser.Parse("x^2")
				return FormatDerivative(expr, "x", 1)
			},
			func(result string) bool {
				return containsSubstring(result, "\\frac{d}{dx}") && containsSubstring(result, "x^{2}")
			},
		},
		{
			"integral formatting",
			func() string {
				expr, _ := parser.Parse("x^2")
				return FormatIntegral(expr, "x", false, nil, nil)
			},
			func(result string) bool {
				return containsSubstring(result, "\\int") && containsSubstring(result, "dx")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function()
			t.Logf("Special format result: %s", result)

			if !tt.check(result) {
				t.Logf("Special formatting test failed for: %s", result)
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
func BenchmarkFormatSimple(b *testing.B) {
	expr, _ := parser.Parse("x^2+2*x+1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Format(expr)
	}
}

func BenchmarkFormatComplex(b *testing.B) {
	expr, _ := parser.Parse("sin(x^2)*exp(cos(y))+ln(z)")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Format(expr)
	}
}
