package ast

import (
	"math"
	"math/big"
	"testing"
)

func TestTrigonometricFunctions(t *testing.T) {
	tests := []struct {
		name      string
		function  string
		input     float64
		expected  float64
		tolerance float64
	}{
		// Basic trigonometric values
		{"sin(0)", "sin", 0, 0, 1e-15},
		{"sin(π/2)", "sin", math.Pi / 2, 1, 1e-15},
		{"sin(π)", "sin", math.Pi, 0, 1e-15},
		{"cos(0)", "cos", 0, 1, 1e-15},
		{"cos(π/2)", "cos", math.Pi / 2, 0, 1e-15},
		{"cos(π)", "cos", math.Pi, -1, 1e-15},
		{"tan(0)", "tan", 0, 0, 1e-15},
		{"tan(π/4)", "tan", math.Pi / 4, 1, 1e-15},

		// Inverse trigonometric functions
		{"arcsin(0)", "arcsin", 0, 0, 1e-15},
		{"arcsin(1)", "arcsin", 1, math.Pi / 2, 1e-15},
		{"arccos(1)", "arccos", 1, 0, 1e-15},
		{"arccos(0)", "arccos", 0, math.Pi / 2, 1e-15},
		{"arctan(0)", "arctan", 0, 0, 1e-15},
		{"arctan(1)", "arctan", 1, math.Pi / 4, 1e-15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := big.NewFloat(tt.input)
			var result *big.Float
			var err error

			switch tt.function {
			case "sin":
				result, err = evaluateSin(input)
			case "cos":
				result, err = evaluateCos(input)
			case "tan":
				result, err = evaluateTan(input)
			case "arcsin":
				result, err = evaluateArcsin(input)
			case "arccos":
				result, err = evaluateArccos(input)
			case "arctan":
				result, err = evaluateArctan(input)
			default:
				t.Fatalf("Unknown function: %s", tt.function)
			}

			if err != nil {
				t.Errorf("%s(%f) returned error: %v", tt.function, tt.input, err)
				return
			}

			resultFloat, _ := result.Float64()
			diff := math.Abs(resultFloat - tt.expected)
			if diff > tt.tolerance {
				t.Errorf("%s(%f) = %f, want %f (diff: %e)", tt.function, tt.input, resultFloat, tt.expected, diff)
			}
		})
	}
}

func TestHyperbolicFunctions(t *testing.T) {
	tests := []struct {
		name      string
		function  string
		input     float64
		expected  float64
		tolerance float64
	}{
		{"sinh(0)", "sinh", 0, 0, 1e-15},
		{"sinh(1)", "sinh", 1, math.Sinh(1), 1e-15},
		{"cosh(0)", "cosh", 0, 1, 1e-15},
		{"cosh(1)", "cosh", 1, math.Cosh(1), 1e-15},
		{"tanh(0)", "tanh", 0, 0, 1e-15},
		{"tanh(1)", "tanh", 1, math.Tanh(1), 1e-15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := big.NewFloat(tt.input)
			var result *big.Float
			var err error

			switch tt.function {
			case "sinh":
				result, err = evaluateSinh(input)
			case "cosh":
				result, err = evaluateCosh(input)
			case "tanh":
				result, err = evaluateTanh(input)
			default:
				t.Fatalf("Unknown function: %s", tt.function)
			}

			if err != nil {
				t.Errorf("%s(%f) returned error: %v", tt.function, tt.input, err)
				return
			}

			resultFloat, _ := result.Float64()
			diff := math.Abs(resultFloat - tt.expected)
			if diff > tt.tolerance {
				t.Errorf("%s(%f) = %f, want %f (diff: %e)", tt.function, tt.input, resultFloat, tt.expected, diff)
			}
		})
	}
}

func TestLogarithmicFunctions(t *testing.T) {
	tests := []struct {
		name      string
		function  string
		input     float64
		base      float64 // for log with base
		expected  float64
		tolerance float64
		shouldErr bool
	}{
		{"ln(1)", "ln", 1, 0, 0, 1e-15, false},
		{"ln(e)", "ln", math.E, 0, 1, 1e-15, false},
		{"ln(10)", "ln", 10, 0, math.Log(10), 1e-15, false},
		{"ln(0)", "ln", 0, 0, 0, 1e-15, true},
		{"ln(-1)", "ln", -1, 0, 0, 1e-15, true},

		{"log(1)", "log10", 1, 0, 0, 1e-15, false},
		{"log(10)", "log10", 10, 0, 1, 1e-15, false},
		{"log(100)", "log10", 100, 0, 2, 1e-15, false},
		{"log(0)", "log10", 0, 0, 0, 1e-15, true},

		{"log_2(8)", "log_base", 8, 2, 3, 1e-15, false},
		{"log_3(27)", "log_base", 27, 3, 3, 1e-15, false},
		{"log_10(1000)", "log_base", 1000, 10, 3, 1e-15, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := big.NewFloat(tt.input)
			var result *big.Float
			var err error

			switch tt.function {
			case "ln":
				result, err = evaluateNaturalLog(input)
			case "log10":
				result, err = evaluateLog10(input)
			case "log_base":
				base := big.NewFloat(tt.base)
				result, err = evaluateLogBase(input, base)
			default:
				t.Fatalf("Unknown function: %s", tt.function)
			}

			if tt.shouldErr {
				if err == nil {
					t.Errorf("%s(%f) should have returned an error", tt.function, tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("%s(%f) returned error: %v", tt.function, tt.input, err)
				return
			}

			resultFloat, _ := result.Float64()
			diff := math.Abs(resultFloat - tt.expected)
			if diff > tt.tolerance {
				t.Errorf("%s(%f) = %f, want %f (diff: %e)", tt.function, tt.input, resultFloat, tt.expected, diff)
			}
		})
	}
}

func TestFunctionEvaluation(t *testing.T) {
	tests := []struct {
		name      string
		funcName  string
		args      []Expr
		vars      map[string]*big.Float
		expected  float64
		tolerance float64
		shouldErr bool
	}{
		{
			"sin(π/2)",
			"sin",
			[]Expr{NewFloat(math.Pi / 2)},
			make(map[string]*big.Float),
			1.0,
			1e-15,
			false,
		},
		{
			"cos(0)",
			"cos",
			[]Expr{NewInt(0)},
			make(map[string]*big.Float),
			1.0,
			1e-15,
			false,
		},
		{
			"ln(e)",
			"ln",
			[]Expr{E},
			make(map[string]*big.Float),
			1.0,
			1e-15,
			false,
		},
		{
			"log(10)",
			"log",
			[]Expr{NewInt(10)},
			make(map[string]*big.Float),
			1.0,
			1e-15,
			false,
		},
		{
			"sinh(0)",
			"sinh",
			[]Expr{NewInt(0)},
			make(map[string]*big.Float),
			0.0,
			1e-15,
			false,
		},
		{
			"sin(x) with x=π/4",
			"sin",
			[]Expr{NewVar("x")},
			map[string]*big.Float{"x": big.NewFloat(math.Pi / 4)},
			math.Sin(math.Pi / 4),
			1e-15,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			funcExpr := NewFunc(tt.funcName, tt.args...)
			result, err := funcExpr.Eval(tt.vars)

			if tt.shouldErr {
				if err == nil {
					t.Errorf("Function %s should have returned an error", tt.name)
				}
				return
			}

			if err != nil {
				t.Errorf("Function %s returned error: %v", tt.name, err)
				return
			}

			resultFloat, _ := result.Float64()
			diff := math.Abs(resultFloat - tt.expected)
			if diff > tt.tolerance {
				t.Errorf("Function %s = %f, want %f (diff: %e)", tt.name, resultFloat, tt.expected, diff)
			}
		})
	}
}

func TestDomainErrors(t *testing.T) {
	tests := []struct {
		name string
		fn   func() (*big.Float, error)
	}{
		{"ln(-1)", func() (*big.Float, error) { return evaluateNaturalLog(big.NewFloat(-1)) }},
		{"ln(0)", func() (*big.Float, error) { return evaluateNaturalLog(big.NewFloat(0)) }},
		{"log(-5)", func() (*big.Float, error) { return evaluateLog10(big.NewFloat(-5)) }},
		{"arcsin(2)", func() (*big.Float, error) { return evaluateArcsin(big.NewFloat(2)) }},
		{"arcsin(-2)", func() (*big.Float, error) { return evaluateArcsin(big.NewFloat(-2)) }},
		{"arccos(1.5)", func() (*big.Float, error) { return evaluateArccos(big.NewFloat(1.5)) }},
		{"arccos(-1.5)", func() (*big.Float, error) { return evaluateArccos(big.NewFloat(-1.5)) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.fn()
			if err == nil {
				t.Errorf("%s should have returned a domain error", tt.name)
			}
		})
	}
}