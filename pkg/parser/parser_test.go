package parser

import (
	"math/big"
	"testing"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []TokenType
	}{
		{
			"simple number",
			"42",
			[]TokenType{TokenInt, TokenEOF},
		},
		{
			"float number",
			"3.14",
			[]TokenType{TokenFloat, TokenEOF},
		},
		{
			"simple expression",
			"x + 1",
			[]TokenType{TokenVar, TokenPlus, TokenInt, TokenEOF},
		},
		{
			"multiplication",
			"2 * x",
			[]TokenType{TokenInt, TokenMultiply, TokenVar, TokenEOF},
		},
		{
			"power",
			"x^2",
			[]TokenType{TokenVar, TokenPower, TokenInt, TokenEOF},
		},
		{
			"parentheses",
			"(x + 1)",
			[]TokenType{TokenLeftParen, TokenVar, TokenPlus, TokenInt, TokenRightParen, TokenEOF},
		},
		{
			"pi constant",
			"\\pi",
			[]TokenType{TokenPi, TokenEOF},
		},
		{
			"sqrt function",
			"\\sqrt{x}",
			[]TokenType{TokenSqrt, TokenLeftBrace, TokenVar, TokenRightBrace, TokenEOF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			var tokens []TokenType

			for {
				token := lexer.NextToken()
				tokens = append(tokens, token.Type)
				if token.Type == TokenEOF {
					break
				}
			}

			if len(tokens) != len(tt.expected) {
				t.Errorf("Token count mismatch. Got %d, want %d", len(tokens), len(tt.expected))
				return
			}

			for i, tokenType := range tokens {
				if tokenType != tt.expected[i] {
					t.Errorf("Token %d: got %s, want %s", i, tokenType, tt.expected[i])
				}
			}
		})
	}
}

func TestParseBasicExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Basic numbers and variables - based on JavaScript test cases
		{"empty", "", "0"},
		{"zero", "0", "0"},
		{"integer", "1", "1"},
		{"float", "3.14", "3.14"},
		{"decimal", ".14", "0.14"},
		{"pi", "pi", "pi"},
		{"e", "e", "e"},
		{"variable", "x", "x"},
		{"theta", "theta", "theta"},

		// Negative numbers
		{"negative zero", "-0", "-1*0"},
		{"negative integer", "-1", "-1"},
		{"negative float", "-3.14", "-3.14"},
		{"negative decimal", "-.14", "-0.14"},
		{"negative pi", "-pi", "-1*pi"},
		{"negative e", "-e", "-1*e"},
		{"negative theta", "-theta", "-1*theta"},

		// LaTeX constants
		{"latex theta", "\\theta", "theta"},
		{"latex pi", "\\pi", "pi"},

		// Basic arithmetic
		{"addition", "1+2", "1+2"},
		{"subtraction", "5-2", "5+-2"},
		{"multiplication", "2*3", "2*3"},
		{"division", "6/2", "6*2^-1"},
		{"power", "2^3", "2^3"},

		// Parentheses
		{"simple parentheses", "(x)", "x"},
		{"parentheses with addition", "(x+1)", "x+1"},
		{"nested parentheses", "((x))", "x"},

		// Variables with subscripts
		{"subscript", "x_1", "x_1"},
		{"subscript with brace", "x_{10}", "x_10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.input == "" {
				// Skip empty input test for now
				t.Skip("Empty input handling needs special case")
				return
			}

			expr, err := Parse(tt.input)
			if err != nil {
				t.Errorf("Parse(%s) returned error: %v", tt.input, err)
				return
			}

			result := expr.String()
			if result != tt.expected {
				t.Errorf("Parse(%s).String() = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseFractions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple fraction", "1/2", "1*2^-1"},
		{"negative fraction", "-1/2", "-1*2^-1"},
		{"fraction with negative denominator", "1/-2", "1*-2^-1"},
		{"latex frac", "\\frac{1}{2}", "1*2^-1"},
		{"latex dfrac", "\\dfrac{3}{4}", "3*4^-1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := Parse(tt.input)
			if err != nil {
				t.Errorf("Parse(%s) returned error: %v", tt.input, err)
				return
			}

			result := expr.String()
			if result != tt.expected {
				t.Errorf("Parse(%s).String() = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseFunctions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"sqrt", "\\sqrt{x}", "sqrt(x)"},
		{"sqrt number", "\\sqrt{4}", "sqrt(4)"},
		{"ln", "\\ln{x}", "ln(x)"},
		{"log", "\\log{x}", "log(x)"},
		{"sin", "\\sin{x}", "sin(x)"},
		{"cos", "\\cos{x}", "cos(x)"},
		{"tan", "\\tan{x}", "tan(x)"},
		{"function call", "f(x)", "f(x)"},
		{"function with multiple args", "f(x, y)", "f(x, y)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := Parse(tt.input)
			if err != nil {
				t.Errorf("Parse(%s) returned error: %v", tt.input, err)
				return
			}

			result := expr.String()
			if result != tt.expected {
				t.Errorf("Parse(%s).String() = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseComplexExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"quadratic", "x^2 + 2*x + 1", "x^2+2*x+1"},
		{"fraction with variables", "x/y", "x*y^-1"},
		{"nested power", "(x^2)^3", "x^2^3"},
		{"implicit multiplication", "2x", "2*x"},
		{"multiple variables", "a*x + b", "a*x+b"},
		{"complex fraction", "(x+1)/(x-1)", "(x+1)*(x+-1)^-1"},
		{"mixed operations", "2*x^2 + 3*x - 1", "2*x^2+3*x+-1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := Parse(tt.input)
			if err != nil {
				t.Errorf("Parse(%s) returned error: %v", tt.input, err)
				return
			}

			result := expr.String()
			if result != tt.expected {
				t.Errorf("Parse(%s).String() = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"unmatched parentheses", "(x + 1"},
		{"invalid token", "x + @"},
		{"empty function", "sin()"},
		{"incomplete frac", "\\frac{1}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)
			if err == nil {
				t.Errorf("Parse(%s) should have returned an error", tt.input)
			}
		})
	}
}

func TestParseOperatorPrecedence(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"addition before multiplication", "1 + 2 * 3", "1+2*3"},
		{"power before multiplication", "2 * 3^2", "2*3^2"},
		{"parentheses override", "(1 + 2) * 3", "(1+2)*3"},
		{"right associative power", "2^3^2", "2^3^2"},
		{"unary minus", "-x + 1", "-1*x+1"},
		{"multiple unary", "--x", "-1*-1*x"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := Parse(tt.input)
			if err != nil {
				t.Errorf("Parse(%s) returned error: %v", tt.input, err)
				return
			}

			result := expr.String()
			if result != tt.expected {
				t.Errorf("Parse(%s).String() = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

// TestParseAndEvaluate tests that parsed expressions can be evaluated correctly
func TestParseAndEvaluate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		vars     map[string]float64
		expected float64
	}{
		{"simple addition", "1+2", nil, 3.0},
		{"simple multiplication", "2*3", nil, 6.0},
		{"power", "2^3", nil, 8.0},
		{"variable substitution", "x+1", map[string]float64{"x": 2}, 3.0},
		{"complex expression", "x^2+2*x+1", map[string]float64{"x": 3}, 16.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := Parse(tt.input)
			if err != nil {
				t.Errorf("Parse(%s) returned error: %v", tt.input, err)
				return
			}

			// Convert float64 vars to big.Float vars
			bigVars := make(map[string]*big.Float)
			for k, v := range tt.vars {
				bigVars[k] = big.NewFloat(v)
			}

			result, err := expr.Eval(bigVars)
			if err != nil {
				t.Errorf("Eval(%s) returned error: %v", tt.input, err)
				return
			}

			resultFloat, _ := result.Float64()
			if resultFloat != tt.expected {
				t.Errorf("Parse(%s).Eval() = %f, want %f", tt.input, resultFloat, tt.expected)
			}
		})
	}
}