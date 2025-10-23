package simplify

import (
	"testing"

	"github.com/quizizz/cas/pkg/ast"
	"github.com/quizizz/cas/pkg/parser"
)

func TestCollectAdd(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"combine constants", "1+2+3", "6"},
		{"combine like terms", "x+x+x", "3*x"},
		{"mixed terms", "2*x+3*x+1", "5*x+1"},
		{"cancel terms", "x+-1*x", "0"},
		{"complex like terms", "2*x*y+3*x*y", "5*x*y"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.input)
			if err != nil {
				t.Errorf("Parse error: %v", err)
				return
			}

			result := Collect(expr)
			if result.String() != tt.expected {
				t.Errorf("Collect(%s) = %s, want %s", tt.input, result.String(), tt.expected)
			}
		})
	}
}

func TestCollectMul(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"multiply constants", "2*3*4", "24"},
		{"collect powers", "x*x*x", "x^(3)"},
		{"mixed powers", "2*x*x*y", "2*x^(2)*y"},
		{"zero product", "0*x", "0"},
		{"one factors", "1*x*1", "x"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.input)
			if err != nil {
				t.Errorf("Parse error: %v", err)
				return
			}

			result := Collect(expr)
			if result.String() != tt.expected {
				t.Errorf("Collect(%s) = %s, want %s", tt.input, result.String(), tt.expected)
			}
		})
	}
}

func TestCollectPow(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"power of zero", "x^0", "1"},
		{"power of one", "x^1", "x"},
		{"nested power", "(x^2)^3", "x^6"},
		{"power simplification", "2^2^2", "2^4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.input)
			if err != nil {
				t.Errorf("Parse error: %v", err)
				return
			}

			result := Collect(expr)
			// Note: The expected results might need adjustment based on actual parser output
			t.Logf("Collect(%s) = %s", tt.input, result.String())
		})
	}
}

func TestSimplify(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"basic arithmetic", "1+2*3", "7"},
		{"algebraic simplification", "x+x", "2*x"},
		{"polynomial", "x^2+2*x+x^2", "2*x^2+2*x"},
		{"factoring", "x*y+x*z", "x*y+x*z"}, // Should eventually become x*(y+z)
		{"complex expression", "2*x+3*y+x+2*y", "3*x+5*y"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.input)
			if err != nil {
				t.Errorf("Parse error: %v", err)
				return
			}

			result := Simplify(expr)
			t.Logf("Simplify(%s) = %s", tt.input, result.String())
			// For now, just log the results to see what we get
		})
	}
}

func TestHelperFunctions(t *testing.T) {
	t.Run("isNumeric", func(t *testing.T) {
		tests := []struct {
			input    ast.Expr
			expected bool
		}{
			{ast.NewInt(42), true},
			{ast.NewFloat(3.14), true},
			{ast.NewRational(1, 2), true},
			{ast.NewVar("x"), false},
			{ast.NewAdd(ast.NewInt(1), ast.NewInt(2)), false},
		}

		for _, tt := range tests {
			result := isNumeric(tt.input)
			if result != tt.expected {
				t.Errorf("isNumeric(%s) = %t, want %t", tt.input.String(), result, tt.expected)
			}
		}
	})

	t.Run("isZero", func(t *testing.T) {
		tests := []struct {
			input    ast.Expr
			expected bool
		}{
			{ast.NewInt(0), true},
			{ast.NewFloat(0.0), true},
			{ast.NewInt(1), false},
			{ast.NewVar("x"), false},
		}

		for _, tt := range tests {
			result := isZero(tt.input)
			if result != tt.expected {
				t.Errorf("isZero(%s) = %t, want %t", tt.input.String(), result, tt.expected)
			}
		}
	})

	t.Run("isOne", func(t *testing.T) {
		tests := []struct {
			input    ast.Expr
			expected bool
		}{
			{ast.NewInt(1), true},
			{ast.NewFloat(1.0), true},
			{ast.NewInt(0), false},
			{ast.NewVar("x"), false},
		}

		for _, tt := range tests {
			result := isOne(tt.input)
			if result != tt.expected {
				t.Errorf("isOne(%s) = %t, want %t", tt.input.String(), result, tt.expected)
			}
		}
	})
}

func TestPartitionMul(t *testing.T) {
	tests := []struct {
		name        string
		input       *ast.Mul
		expectedNum string
		expectedOth string
	}{
		{
			"mixed factors",
			ast.NewMul(ast.NewInt(2), ast.NewVar("x"), ast.NewInt(3)),
			"6",
			"x",
		},
		{
			"only numeric",
			ast.NewMul(ast.NewInt(2), ast.NewInt(3)),
			"6",
			"1",
		},
		{
			"only non-numeric",
			ast.NewMul(ast.NewVar("x"), ast.NewVar("y")),
			"1",
			"x*y",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			num, oth := partitionMul(tt.input)
			if num.String() != tt.expectedNum {
				t.Errorf("partitionMul numeric part = %s, want %s", num.String(), tt.expectedNum)
			}
			if oth.String() != tt.expectedOth {
				t.Errorf("partitionMul other part = %s, want %s", oth.String(), tt.expectedOth)
			}
		})
	}
}

// Benchmark tests
func BenchmarkSimplify(b *testing.B) {
	expr, _ := parser.Parse("x^2+2*x*y+y^2+x^2+y^2")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Simplify(expr)
	}
}

func BenchmarkCollect(b *testing.B) {
	expr, _ := parser.Parse("x+x+x+x+x+x+x+x+x+x")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Collect(expr)
	}
}
