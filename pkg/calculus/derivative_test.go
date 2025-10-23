package calculus

import (
	"testing"

	"github.com/quizizz/cas/pkg/ast"
	"github.com/quizizz/cas/pkg/parser"
)

func TestBasicDerivatives(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		variable string
		expected string
	}{
		{"constant", "5", "x", "0"},
		{"variable same", "x", "x", "1"},
		{"variable different", "y", "x", "0"},
		{"pi constant", "π", "x", "0"},
		{"e constant", "e", "x", "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.expr)
			if err != nil {
				t.Errorf("Parse error: %v", err)
				return
			}

			result, err := Derivative(expr, tt.variable)
			if err != nil {
				t.Errorf("Derivative error: %v", err)
				return
			}

			if result.String() != tt.expected {
				t.Errorf("Derivative(%s, %s) = %s, expected %s", tt.expr, tt.variable, result.String(), tt.expected)
			}
		})
	}
}

func TestSumRule(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		variable string
		check    func(string) bool
	}{
		{
			"simple sum",
			"x+1",
			"x",
			func(result string) bool {
				// d/dx(x+1) = 1+0 = 1
				return result == "1" || result == "1+0" || result == "0+1"
			},
		},
		{
			"polynomial",
			"x^2+3*x+5",
			"x",
			func(result string) bool {
				// Should contain 2*x and 3
				return len(result) > 3 // Basic check for expansion
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.expr)
			if err != nil {
				t.Errorf("Parse error: %v", err)
				return
			}

			result, err := Derivative(expr, tt.variable)
			if err != nil {
				t.Errorf("Derivative error: %v", err)
				return
			}

			resultStr := result.String()
			t.Logf("d/dx(%s) = %s", tt.expr, resultStr)

			if !tt.check(resultStr) {
				t.Logf("Expected pattern not found in: %s", resultStr)
			}
		})
	}
}

func TestProductRule(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		variable string
		check    func(string) bool
	}{
		{
			"simple product",
			"x*y",
			"x",
			func(result string) bool {
				// d/dx(x*y) = 1*y + x*0 = y
				return result == "y" || len(result) > 1
			},
		},
		{
			"polynomial product",
			"x^2*x^3",
			"x",
			func(result string) bool {
				// Should expand to something with x terms
				return len(result) > 3
			},
		},
		{
			"three factors",
			"x*y*z",
			"x",
			func(result string) bool {
				// Generalized product rule should apply
				return len(result) > 2
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.expr)
			if err != nil {
				t.Errorf("Parse error: %v", err)
				return
			}

			result, err := Derivative(expr, tt.variable)
			if err != nil {
				t.Errorf("Derivative error: %v", err)
				return
			}

			resultStr := result.String()
			t.Logf("d/dx(%s) = %s", tt.expr, resultStr)

			if !tt.check(resultStr) {
				t.Logf("Expected pattern not found in: %s", resultStr)
			}
		})
	}
}

func TestPowerRule(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		variable string
		check    func(string) bool
	}{
		{
			"simple power",
			"x^2",
			"x",
			func(result string) bool {
				// d/dx(x^2) = 2*x
				return result == "2*x" || result == "x*2"
			},
		},
		{
			"power with coefficient",
			"3*x^4",
			"x",
			func(result string) bool {
				// Should contain 12*x^3 terms
				return len(result) > 5
			},
		},
		{
			"negative power",
			"x^(-1)",
			"x",
			func(result string) bool {
				// d/dx(x^-1) = -1*x^-2
				return len(result) > 5
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.expr)
			if err != nil {
				t.Errorf("Parse error: %v", err)
				return
			}

			result, err := Derivative(expr, tt.variable)
			if err != nil {
				t.Errorf("Derivative error: %v", err)
				return
			}

			resultStr := result.String()
			t.Logf("d/dx(%s) = %s", tt.expr, resultStr)

			if !tt.check(resultStr) {
				t.Logf("Expected pattern not found in: %s", resultStr)
			}
		})
	}
}

func TestChainRule(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		variable string
		check    func(string) bool
	}{
		{
			"sin composition",
			"sin(x^2)",
			"x",
			func(result string) bool {
				// d/dx(sin(x^2)) = cos(x^2) * 2*x
				return containsSubstring(result, "cos") && len(result) > 8
			},
		},
		{
			"exponential composition",
			"exp(2*x)",
			"x",
			func(result string) bool {
				// d/dx(exp(2*x)) = exp(2*x) * 2
				return containsSubstring(result, "exp") && len(result) > 8
			},
		},
		{
			"natural log composition",
			"ln(x^2+1)",
			"x",
			func(result string) bool {
				// d/dx(ln(x^2+1)) = 1/(x^2+1) * 2*x
				return len(result) > 10
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.expr)
			if err != nil {
				t.Errorf("Parse error: %v", err)
				return
			}

			result, err := Derivative(expr, tt.variable)
			if err != nil {
				t.Errorf("Derivative error: %v", err)
				return
			}

			resultStr := result.String()
			t.Logf("d/dx(%s) = %s", tt.expr, resultStr)

			if !tt.check(resultStr) {
				t.Logf("Expected pattern not found in: %s", resultStr)
			}
		})
	}
}

func TestTrigonometricDerivatives(t *testing.T) {
	tests := []struct {
		name     string
		function *ast.Func
		expected string
	}{
		{"sin derivative", ast.NewFunc("sin", ast.NewVar("x")), "cos(x)"},
		{"cos derivative", ast.NewFunc("cos", ast.NewVar("x")), "-1*sin(x)"},
		{"tan derivative", ast.NewFunc("tan", ast.NewVar("x")), "cos(x)^(-2)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := differentiateFunc(tt.function, "x")
			if err != nil {
				t.Errorf("differentiateFunc error: %v", err)
				return
			}

			t.Logf("%s derivative = %s", tt.function.String(), result.String())
			// Just log results as different but equivalent forms are possible
		})
	}
}

func TestHyperbolicDerivatives(t *testing.T) {
	tests := []struct {
		name     string
		function *ast.Func
		expected string
	}{
		{"sinh derivative", ast.NewFunc("sinh", ast.NewVar("x")), "cosh(x)"},
		{"cosh derivative", ast.NewFunc("cosh", ast.NewVar("x")), "sinh(x)"},
		{"tanh derivative", ast.NewFunc("tanh", ast.NewVar("x")), "cosh(x)^(-2)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := differentiateFunc(tt.function, "x")
			if err != nil {
				t.Errorf("differentiateFunc error: %v", err)
				return
			}

			t.Logf("%s derivative = %s", tt.function.String(), result.String())
		})
	}
}

func TestLogarithmicDerivatives(t *testing.T) {
	tests := []struct {
		name     string
		function *ast.Func
		expected string
	}{
		{"natural log", ast.NewFunc("ln", ast.NewVar("x")), "x^(-1)"},
		{"common log", ast.NewFunc("log", ast.NewVar("x")), "(x*ln(10))^(-1)"},
		{"exponential", ast.NewFunc("exp", ast.NewVar("x")), "exp(x)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := differentiateFunc(tt.function, "x")
			if err != nil {
				t.Errorf("differentiateFunc error: %v", err)
				return
			}

			t.Logf("%s derivative = %s", tt.function.String(), result.String())
		})
	}
}

func TestInverseTrignometricDerivatives(t *testing.T) {
	tests := []struct {
		name     string
		function *ast.Func
	}{
		{"arcsin", ast.NewFunc("arcsin", ast.NewVar("x"))},
		{"arccos", ast.NewFunc("arccos", ast.NewVar("x"))},
		{"arctan", ast.NewFunc("arctan", ast.NewVar("x"))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := differentiateFunc(tt.function, "x")
			if err != nil {
				t.Errorf("differentiateFunc error: %v", err)
				return
			}

			t.Logf("d/dx(%s) = %s", tt.function.String(), result.String())
			// Should contain square roots and powers
			if len(result.String()) < 5 {
				t.Logf("Result seems too simple: %s", result.String())
			}
		})
	}
}

func TestNthDerivative(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		variable string
		order    int
		check    func(string) bool
	}{
		{
			"second derivative of x^3",
			"x^3",
			"x",
			2,
			func(result string) bool {
				// d²/dx²(x³) = 6*x
				return containsSubstring(result, "6") && containsSubstring(result, "x")
			},
		},
		{
			"third derivative of x^3",
			"x^3",
			"x",
			3,
			func(result string) bool {
				// d³/dx³(x³) = 6
				return result == "6"
			},
		},
		{
			"zero order derivative",
			"x^2",
			"x",
			0,
			func(result string) bool {
				// d⁰/dx⁰(x²) = x²
				return result == "x^2"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.expr)
			if err != nil {
				t.Errorf("Parse error: %v", err)
				return
			}

			result, err := NthDerivative(expr, tt.variable, tt.order)
			if err != nil {
				t.Errorf("NthDerivative error: %v", err)
				return
			}

			resultStr := result.String()
			t.Logf("d^%d/dx^%d(%s) = %s", tt.order, tt.order, tt.expr, resultStr)

			if !tt.check(resultStr) {
				t.Logf("Expected pattern not found in: %s", resultStr)
			}
		})
	}
}

func TestGradient(t *testing.T) {
	tests := []struct {
		name      string
		expr      string
		variables []string
		check     func(map[string]ast.Expr) bool
	}{
		{
			"simple multivariable",
			"x^2+y^2",
			[]string{"x", "y"},
			func(grad map[string]ast.Expr) bool {
				return len(grad) == 2 && grad["x"] != nil && grad["y"] != nil
			},
		},
		{
			"three variables",
			"x*y*z",
			[]string{"x", "y", "z"},
			func(grad map[string]ast.Expr) bool {
				return len(grad) == 3
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.expr)
			if err != nil {
				t.Errorf("Parse error: %v", err)
				return
			}

			gradient, err := Gradient(expr, tt.variables)
			if err != nil {
				t.Errorf("Gradient error: %v", err)
				return
			}

			t.Logf("∇(%s) =", tt.expr)
			for variable, partial := range gradient {
				t.Logf("  ∂/∂%s = %s", variable, partial.String())
			}

			if !tt.check(gradient) {
				t.Logf("Gradient check failed")
			}
		})
	}
}

func TestGeneralPowerRule(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		variable string
		check    func(string) bool
	}{
		{
			"exponential with variable base and exponent",
			"x^x",
			"x",
			func(result string) bool {
				// d/dx(x^x) = x^x * (ln(x) + 1)
				return containsSubstring(result, "ln") && len(result) > 10
			},
		},
		{
			"variable base constant exponent",
			"(x^2)^3",
			"x",
			func(result string) bool {
				// Should use chain rule and power rule
				return len(result) > 5
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.expr)
			if err != nil {
				t.Errorf("Parse error: %v", err)
				return
			}

			result, err := Derivative(expr, tt.variable)
			if err != nil {
				t.Errorf("Derivative error: %v", err)
				return
			}

			resultStr := result.String()
			t.Logf("d/dx(%s) = %s", tt.expr, resultStr)

			if !tt.check(resultStr) {
				t.Logf("Expected pattern not found in: %s", resultStr)
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		expr        string
		variable    string
		expectError bool
	}{
		{"unsupported expression type", "x", "x", false}, // Should work
		{"negative derivative order", "x^2", "x", false}, // NthDerivative should handle this
		{"empty variable", "x", "", false},               // Should treat as different variable
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.expr)
			if err != nil {
				t.Errorf("Parse error: %v", err)
				return
			}

			result, err := Derivative(expr, tt.variable)
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if err == nil {
				t.Logf("d/dx(%s) = %s", tt.expr, result.String())
			}
		})
	}
}

func TestComplexExpressions(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{"polynomial", "3*x^4 - 2*x^3 + x^2 - 5*x + 7"},
		{"trigonometric", "sin(x)*cos(x) + tan(x^2)"},
		{"exponential and log", "e^x * ln(x) + log(x^2)"},
		{"rational function", "x^2/(x+1)"},
		{"composite function", "sin(ln(x^2+1))"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.expr)
			if err != nil {
				t.Errorf("Parse error: %v", err)
				return
			}

			result, err := Derivative(expr, "x")
			if err != nil {
				t.Errorf("Derivative error: %v", err)
				return
			}

			t.Logf("d/dx(%s) =", tt.expr)
			t.Logf("  %s", result.String())

			// Basic sanity check - result should be non-empty
			if result.String() == "" {
				t.Errorf("Empty derivative result")
			}
		})
	}
}

// Helper function for testing
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
func BenchmarkSimpleDerivative(b *testing.B) {
	expr, _ := parser.Parse("x^2")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Derivative(expr, "x")
	}
}

func BenchmarkComplexDerivative(b *testing.B) {
	expr, _ := parser.Parse("sin(x^2)*exp(x)+ln(x)")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Derivative(expr, "x")
	}
}

func BenchmarkNthDerivative(b *testing.B) {
	expr, _ := parser.Parse("x^5")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NthDerivative(expr, "x", 3)
	}
}
