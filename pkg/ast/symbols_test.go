package ast

import (
	"math/big"
	"reflect"
	"testing"
)

func TestVar(t *testing.T) {
	tests := []struct {
		name     string
		varName  string
		expected string
	}{
		{"simple variable", "x", "x"},
		{"multi-character", "theta", "theta"},
		{"subscript", "x_1", "x_1"},
		{"uppercase", "X", "X"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := NewVar(tt.varName)
			if expr.String() != tt.expected {
				t.Errorf("Var(%s).String() = %s, want %s", tt.varName, expr.String(), tt.expected)
			}
			if expr.LaTeX() != tt.expected {
				t.Errorf("Var(%s).LaTeX() = %s, want %s", tt.varName, expr.LaTeX(), tt.expected)
			}
		})
	}
}

func TestVarEval(t *testing.T) {
	tests := []struct {
		name      string
		varName   string
		vars      map[string]*big.Float
		expected  float64
		shouldErr bool
	}{
		{"defined variable", "x", map[string]*big.Float{"x": big.NewFloat(3.0)}, 3.0, false},
		{"undefined variable", "y", map[string]*big.Float{"x": big.NewFloat(3.0)}, 0.0, true},
		{"zero value", "z", map[string]*big.Float{"z": big.NewFloat(0.0)}, 0.0, false},
		{"negative value", "a", map[string]*big.Float{"a": big.NewFloat(-2.5)}, -2.5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := NewVar(tt.varName)
			result, err := expr.Eval(tt.vars)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("Var(%s).Eval() should have returned error", tt.varName)
				}
				return
			}
			if err != nil {
				t.Errorf("Var(%s).Eval() returned unexpected error: %v", tt.varName, err)
				return
			}
			expected := big.NewFloat(tt.expected)
			if result.Cmp(expected) != 0 {
				t.Errorf("Var(%s).Eval() = %s, want %s", tt.varName, result.String(), expected.String())
			}
		})
	}
}

func TestVarVariables(t *testing.T) {
	tests := []struct {
		name     string
		varName  string
		expected []string
	}{
		{"single var", "x", []string{"x"}},
		{"complex name", "theta", []string{"theta"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := NewVar(tt.varName)
			result := expr.Variables()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Var(%s).Variables() = %v, want %v", tt.varName, result, tt.expected)
			}
		})
	}
}

func TestConst(t *testing.T) {
	tests := []struct {
		name        string
		constant    *Const
		expectedStr string
		expectedTeX string
	}{
		{"pi", Pi, "pi", "\\pi"},
		{"e", E, "e", "e"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant.String() != tt.expectedStr {
				t.Errorf("%s.String() = %s, want %s", tt.name, tt.constant.String(), tt.expectedStr)
			}
			if tt.constant.LaTeX() != tt.expectedTeX {
				t.Errorf("%s.LaTeX() = %s, want %s", tt.name, tt.constant.LaTeX(), tt.expectedTeX)
			}
		})
	}
}

func TestConstEval(t *testing.T) {
	tests := []struct {
		name      string
		constant  *Const
		minVal    float64
		maxVal    float64
	}{
		{"pi", Pi, 3.14, 3.15},
		{"e", E, 2.71, 2.72},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.constant.Eval(make(map[string]*big.Float))
			if err != nil {
				t.Errorf("%s.Eval() returned error: %v", tt.name, err)
				return
			}
			resultFloat, _ := result.Float64()
			if resultFloat < tt.minVal || resultFloat > tt.maxVal {
				t.Errorf("%s.Eval() = %f, want between %f and %f", tt.name, resultFloat, tt.minVal, tt.maxVal)
			}
		})
	}
}

func TestFunc(t *testing.T) {
	tests := []struct {
		name        string
		funcName    string
		args        []Expr
		expectedStr string
	}{
		{"no args", "f", []Expr{}, "f()"},
		{"one arg", "sqrt", []Expr{NewVar("x")}, "sqrt(x)"},
		{"two args", "log", []Expr{NewVar("x"), NewVar("y")}, "log(x, y)"},
		{"complex arg", "sin", []Expr{NewAdd(NewVar("x"), NewInt(1))}, "sin(x+1)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := NewFunc(tt.funcName, tt.args...)
			if expr.String() != tt.expectedStr {
				t.Errorf("Func(%s, %v).String() = %s, want %s", tt.funcName, tt.args, expr.String(), tt.expectedStr)
			}
		})
	}
}

func TestFuncLaTeX(t *testing.T) {
	tests := []struct {
		name        string
		funcName    string
		args        []Expr
		expectedTeX string
	}{
		{"sqrt", "sqrt", []Expr{NewVar("x")}, "\\sqrt{x}"},
		{"log", "log", []Expr{NewVar("x")}, "\\log{x}"},
		{"ln", "ln", []Expr{NewVar("x")}, "\\ln{x}"},
		{"generic", "sin", []Expr{NewVar("x")}, "\\mathrm{sin}(x)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := NewFunc(tt.funcName, tt.args...)
			if expr.LaTeX() != tt.expectedTeX {
				t.Errorf("Func(%s, %v).LaTeX() = %s, want %s", tt.funcName, tt.args, expr.LaTeX(), tt.expectedTeX)
			}
		})
	}
}

func TestFuncEval(t *testing.T) {
	tests := []struct {
		name      string
		funcName  string
		args      []Expr
		vars      map[string]*big.Float
		expected  float64
		tolerance float64
		shouldErr bool
	}{
		{"sqrt(4)", "sqrt", []Expr{NewInt(4)}, make(map[string]*big.Float), 2.0, 1e-10, false},
		{"sqrt(x) with x=9", "sqrt", []Expr{NewVar("x")}, map[string]*big.Float{"x": big.NewFloat(9.0)}, 3.0, 1e-10, false},
		{"abs(-5)", "abs", []Expr{NewInt(-5)}, make(map[string]*big.Float), 5.0, 1e-10, false},
		{"abs(x) with x=-3", "abs", []Expr{NewVar("x")}, map[string]*big.Float{"x": big.NewFloat(-3.0)}, 3.0, 1e-10, false},
		{"unsupported function", "unsupported", []Expr{NewInt(1)}, make(map[string]*big.Float), 0.0, 0.0, true},
		{"sqrt wrong args", "sqrt", []Expr{NewInt(1), NewInt(2)}, make(map[string]*big.Float), 0.0, 0.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := NewFunc(tt.funcName, tt.args...)
			result, err := expr.Eval(tt.vars)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("Func(%s).Eval() should have returned error", tt.name)
				}
				return
			}
			if err != nil {
				t.Errorf("Func(%s).Eval() returned unexpected error: %v", tt.name, err)
				return
			}
			expected := big.NewFloat(tt.expected)
			diff := new(big.Float).Sub(result, expected)
			diff.Abs(diff)
			tolerance := big.NewFloat(tt.tolerance)
			if diff.Cmp(tolerance) > 0 {
				t.Errorf("Func(%s).Eval() = %s, want %s (tolerance: %g)", tt.name, result.String(), expected.String(), tt.tolerance)
			}
		})
	}
}

func TestFuncVariables(t *testing.T) {
	tests := []struct {
		name     string
		funcName string
		args     []Expr
		expected []string
	}{
		{"no variables", "sqrt", []Expr{NewInt(4)}, []string{}},
		{"one variable", "sin", []Expr{NewVar("x")}, []string{"x"}},
		{"multiple variables", "pow", []Expr{NewVar("x"), NewVar("y")}, []string{"x", "y"}},
		{"duplicate variables", "add", []Expr{NewVar("x"), NewVar("x")}, []string{"x"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := NewFunc(tt.funcName, tt.args...)
			result := expr.Variables()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Func(%s).Variables() = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}