package ast

import (
	"math/big"
	"testing"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		terms    []Expr
		expected string
	}{
		{"empty", []Expr{}, "0"},
		{"single term", []Expr{NewInt(5)}, "5"},
		{"two positive", []Expr{NewInt(2), NewInt(3)}, "2+3"},
		{"positive and negative", []Expr{NewInt(5), NewInt(-2)}, "5+-2"},
		{"variables", []Expr{NewVar("x"), NewVar("y")}, "x+y"},
		{"mixed", []Expr{NewInt(1), NewVar("x"), NewInt(2)}, "1+x+2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := NewAdd(tt.terms...)
			if expr.String() != tt.expected {
				t.Errorf("Add(%v).String() = %s, want %s", tt.terms, expr.String(), tt.expected)
			}
		})
	}
}

func TestAddEval(t *testing.T) {
	tests := []struct {
		name     string
		terms    []Expr
		vars     map[string]*big.Float
		expected float64
	}{
		{"empty", []Expr{}, make(map[string]*big.Float), 0.0},
		{"constants", []Expr{NewInt(2), NewInt(3), NewInt(4)}, make(map[string]*big.Float), 9.0},
		{"with variable", []Expr{NewInt(1), NewVar("x"), NewInt(2)}, map[string]*big.Float{"x": big.NewFloat(5.0)}, 8.0},
		{"mixed types", []Expr{NewInt(1), NewFloat(2.5), NewRational(1, 2)}, make(map[string]*big.Float), 4.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := NewAdd(tt.terms...)
			result, err := expr.Eval(tt.vars)
			if err != nil {
				t.Errorf("Add(%v).Eval() returned error: %v", tt.terms, err)
				return
			}
			expected := big.NewFloat(tt.expected)
			if result.Cmp(expected) != 0 {
				t.Errorf("Add(%v).Eval() = %s, want %s", tt.terms, result.String(), expected.String())
			}
		})
	}
}

func TestAddSimplify(t *testing.T) {
	tests := []struct {
		name     string
		terms    []Expr
		expected string
	}{
		{"numeric only", []Expr{NewInt(2), NewInt(3)}, "5"},
		{"zero result", []Expr{NewInt(5), NewInt(-5)}, "0"},
		{"mixed", []Expr{NewInt(1), NewVar("x"), NewInt(2)}, "3+x"},
		{"already simplified", []Expr{NewVar("x"), NewVar("y")}, "x+y"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := NewAdd(tt.terms...)
			simplified := expr.Simplify()
			result := simplified.String()
			if result != tt.expected {
				t.Errorf("Add(%v).Simplify().String() = %s, want %s", tt.terms, result, tt.expected)
			}
		})
	}
}

func TestMul(t *testing.T) {
	tests := []struct {
		name     string
		factors  []Expr
		expected string
	}{
		{"empty", []Expr{}, "1"},
		{"single factor", []Expr{NewInt(5)}, "5"},
		{"two factors", []Expr{NewInt(2), NewInt(3)}, "2*3"},
		{"variables", []Expr{NewVar("x"), NewVar("y")}, "x*y"},
		{"mixed", []Expr{NewInt(2), NewVar("x")}, "2*x"},
		{"with addition", []Expr{NewInt(2), NewAdd(NewVar("x"), NewInt(1))}, "2*(x+1)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := NewMul(tt.factors...)
			if expr.String() != tt.expected {
				t.Errorf("Mul(%v).String() = %s, want %s", tt.factors, expr.String(), tt.expected)
			}
		})
	}
}

func TestMulEval(t *testing.T) {
	tests := []struct {
		name     string
		factors  []Expr
		vars     map[string]*big.Float
		expected float64
	}{
		{"empty", []Expr{}, make(map[string]*big.Float), 1.0},
		{"constants", []Expr{NewInt(2), NewInt(3), NewInt(4)}, make(map[string]*big.Float), 24.0},
		{"with variable", []Expr{NewInt(2), NewVar("x")}, map[string]*big.Float{"x": big.NewFloat(5.0)}, 10.0},
		{"with zero", []Expr{NewInt(0), NewVar("x")}, map[string]*big.Float{"x": big.NewFloat(5.0)}, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := NewMul(tt.factors...)
			result, err := expr.Eval(tt.vars)
			if err != nil {
				t.Errorf("Mul(%v).Eval() returned error: %v", tt.factors, err)
				return
			}
			expected := big.NewFloat(tt.expected)
			if result.Cmp(expected) != 0 {
				t.Errorf("Mul(%v).Eval() = %s, want %s", tt.factors, result.String(), expected.String())
			}
		})
	}
}

func TestMulSimplify(t *testing.T) {
	tests := []struct {
		name     string
		factors  []Expr
		expected string
	}{
		{"numeric only", []Expr{NewInt(2), NewInt(3)}, "6"},
		{"with one", []Expr{NewInt(1), NewVar("x")}, "x"},
		{"with zero", []Expr{NewInt(0), NewVar("x")}, "0"},
		{"mixed", []Expr{NewInt(2), NewVar("x"), NewInt(3)}, "6*x"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := NewMul(tt.factors...)
			simplified := expr.Simplify()
			result := simplified.String()
			if result != tt.expected {
				t.Errorf("Mul(%v).Simplify().String() = %s, want %s", tt.factors, result, tt.expected)
			}
		})
	}
}

func TestPow(t *testing.T) {
	tests := []struct {
		name        string
		base        Expr
		exponent    Expr
		expectedStr string
		expectedTeX string
	}{
		{"simple", NewVar("x"), NewInt(2), "x^(2)", "x^{2}"},
		{"complex base", NewAdd(NewVar("x"), NewInt(1)), NewInt(2), "(x+1)^(2)", "(x+1)^{2}"},
		{"complex exponent", NewVar("x"), NewAdd(NewVar("n"), NewInt(1)), "x^(n+1)", "x^{n+1}"},
		{"nested power", NewPow(NewVar("x"), NewInt(2)), NewInt(3), "x^(6)", "x^{6}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := NewPow(tt.base, tt.exponent)
			if expr.String() != tt.expectedStr {
				t.Errorf("Pow(%s, %s).String() = %s, want %s", tt.base.String(), tt.exponent.String(), expr.String(), tt.expectedStr)
			}
			if expr.LaTeX() != tt.expectedTeX {
				t.Errorf("Pow(%s, %s).LaTeX() = %s, want %s", tt.base.String(), tt.exponent.String(), expr.LaTeX(), tt.expectedTeX)
			}
		})
	}
}

func TestPowEval(t *testing.T) {
	tests := []struct {
		name      string
		base      Expr
		exponent  Expr
		vars      map[string]*big.Float
		expected  float64
		shouldErr bool
	}{
		{"2^3", NewInt(2), NewInt(3), make(map[string]*big.Float), 8.0, false},
		{"x^2 with x=3", NewVar("x"), NewInt(2), map[string]*big.Float{"x": big.NewFloat(3.0)}, 9.0, false},
		{"2^0", NewInt(2), NewInt(0), make(map[string]*big.Float), 1.0, false},
		{"5^1", NewInt(5), NewInt(1), make(map[string]*big.Float), 5.0, false},
		{"2^(-2)", NewInt(2), NewInt(-2), make(map[string]*big.Float), 0.25, false},
		{"non-integer exponent", NewInt(2), NewFloat(1.5), make(map[string]*big.Float), 0.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := NewPow(tt.base, tt.exponent)
			result, err := expr.Eval(tt.vars)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("Pow(%s, %s).Eval() should have returned error", tt.base.String(), tt.exponent.String())
				}
				return
			}
			if err != nil {
				t.Errorf("Pow(%s, %s).Eval() returned unexpected error: %v", tt.base.String(), tt.exponent.String(), err)
				return
			}
			expected := big.NewFloat(tt.expected)
			if result.Cmp(expected) != 0 {
				t.Errorf("Pow(%s, %s).Eval() = %s, want %s", tt.base.String(), tt.exponent.String(), result.String(), expected.String())
			}
		})
	}
}

func TestPowSimplify(t *testing.T) {
	tests := []struct {
		name     string
		base     Expr
		exponent Expr
		expected string
	}{
		{"x^0", NewVar("x"), NewInt(0), "1"},
		{"x^1", NewVar("x"), NewInt(1), "x"},
		{"no simplification", NewVar("x"), NewInt(2), "x^(2)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := NewPow(tt.base, tt.exponent)
			simplified := expr.Simplify()
			result := simplified.String()
			if result != tt.expected {
				t.Errorf("Pow(%s, %s).Simplify().String() = %s, want %s", tt.base.String(), tt.exponent.String(), result, tt.expected)
			}
		})
	}
}

func TestExprVariables(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected []string
	}{
		{"constant", NewInt(5), []string{}},
		{"variable", NewVar("x"), []string{"x"}},
		{"addition", NewAdd(NewVar("x"), NewVar("y")), []string{"x", "y"}},
		{"multiplication", NewMul(NewVar("x"), NewVar("x")), []string{"x", "x"}},
		{"power", NewPow(NewVar("x"), NewVar("n")), []string{"x", "n"}},
		{"complex", NewAdd(NewMul(NewVar("a"), NewVar("x")), NewVar("b")), []string{"a", "x", "b"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.expr.Variables()
			// For this basic test, we just check that all expected variables are present
			// Order might vary due to implementation details
			for _, expected := range tt.expected {
				found := false
				for _, actual := range result {
					if actual == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("%s.Variables() = %v, missing expected variable %s", tt.name, result, expected)
				}
			}
		})
	}
}