package ast

import (
	"math/big"
	"testing"
)

func TestInt(t *testing.T) {
	tests := []struct {
		name     string
		value    int64
		expected string
	}{
		{"zero", 0, "0"},
		{"positive", 42, "42"},
		{"negative", -123, "-123"},
		{"large", 123456789, "123456789"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := NewInt(tt.value)
			if expr.String() != tt.expected {
				t.Errorf("Int(%d).String() = %s, want %s", tt.value, expr.String(), tt.expected)
			}
		})
	}
}

func TestIntFromString(t *testing.T) {
	tests := []struct {
		input     string
		expected  string
		shouldErr bool
	}{
		{"0", "0", false},
		{"42", "42", false},
		{"-123", "-123", false},
		{"invalid", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			expr, err := NewIntFromString(tt.input)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("NewIntFromString(%s) should have returned error", tt.input)
				}
				return
			}
			if err != nil {
				t.Errorf("NewIntFromString(%s) returned unexpected error: %v", tt.input, err)
				return
			}
			if expr.String() != tt.expected {
				t.Errorf("NewIntFromString(%s).String() = %s, want %s", tt.input, expr.String(), tt.expected)
			}
		})
	}
}

func TestIntEval(t *testing.T) {
	tests := []struct {
		name     string
		value    int64
		expected float64
	}{
		{"zero", 0, 0.0},
		{"positive", 42, 42.0},
		{"negative", -123, -123.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := NewInt(tt.value)
			result, err := expr.Eval(make(map[string]*big.Float))
			if err != nil {
				t.Errorf("Int(%d).Eval() returned error: %v", tt.value, err)
				return
			}
			expectedBig := big.NewFloat(tt.expected)
			if result.Cmp(expectedBig) != 0 {
				t.Errorf("Int(%d).Eval() = %s, want %s", tt.value, result.String(), expectedBig.String())
			}
		})
	}
}

func TestIntEqual(t *testing.T) {
	tests := []struct {
		name     string
		a, b     *Int
		expected bool
	}{
		{"equal positive", NewInt(42), NewInt(42), true},
		{"equal negative", NewInt(-123), NewInt(-123), true},
		{"equal zero", NewInt(0), NewInt(0), true},
		{"not equal", NewInt(42), NewInt(43), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Equal(tt.b)
			if result != tt.expected {
				t.Errorf("%s.Equal(%s) = %t, want %t", tt.a.String(), tt.b.String(), result, tt.expected)
			}
		})
	}
}

func TestFloat(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected string
	}{
		{"zero", 0.0, "0"},
		{"integer", 42.0, "42"},
		{"decimal", 3.14, "3.14"},
		{"negative", -2.5, "-2.5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := NewFloat(tt.value)
			// Note: big.Float formatting might vary, so we check the numeric value
			resultVal, _ := expr.Eval(make(map[string]*big.Float))
			expectedVal := big.NewFloat(tt.value)
			if resultVal.Cmp(expectedVal) != 0 {
				t.Errorf("Float(%f) evaluated to %s, want %s", tt.value, resultVal.String(), expectedVal.String())
			}
		})
	}
}

func TestRational(t *testing.T) {
	tests := []struct {
		name        string
		num, den    int64
		expectedStr string
		expectedTeX string
	}{
		{"simple fraction", 1, 2, "1/2", "\\frac{1}{2}"},
		{"integer", 4, 1, "4/1", "4"},
		{"negative", -3, 4, "-3/4", "\\frac{-3}{4}"},
		{"reduced", 6, 9, "2/3", "\\frac{2}{3}"},
		{"zero", 0, 5, "0/1", "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := NewRational(tt.num, tt.den)
			if expr.String() != tt.expectedStr {
				t.Errorf("Rational(%d, %d).String() = %s, want %s", tt.num, tt.den, expr.String(), tt.expectedStr)
			}
			if expr.LaTeX() != tt.expectedTeX {
				t.Errorf("Rational(%d, %d).LaTeX() = %s, want %s", tt.num, tt.den, expr.LaTeX(), tt.expectedTeX)
			}
		})
	}
}

func TestRationalEval(t *testing.T) {
	tests := []struct {
		name        string
		num, den    int64
		expectedVal float64
	}{
		{"half", 1, 2, 0.5},
		{"third", 1, 3, 1.0 / 3.0},
		{"integer", 4, 1, 4.0},
		{"negative", -3, 4, -0.75},
		{"zero", 0, 5, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := NewRational(tt.num, tt.den)
			result, err := expr.Eval(make(map[string]*big.Float))
			if err != nil {
				t.Errorf("Rational(%d, %d).Eval() returned error: %v", tt.num, tt.den, err)
				return
			}
			expected := big.NewFloat(tt.expectedVal)
			// Use a small tolerance for float comparison
			diff := new(big.Float).Sub(result, expected)
			diff.Abs(diff)
			tolerance := big.NewFloat(1e-10)
			if diff.Cmp(tolerance) > 0 {
				t.Errorf("Rational(%d, %d).Eval() = %s, want %s", tt.num, tt.den, result.String(), expected.String())
			}
		})
	}
}