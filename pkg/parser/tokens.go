// Package parser implements expression parsing for mathematical expressions.
package parser

import (
	"regexp"
)

// TokenType represents the type of a token
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenInt
	TokenFloat
	TokenVar
	TokenPlus
	TokenMinus
	TokenMultiply
	TokenDivide
	TokenPower
	TokenLeftParen
	TokenRightParen
	TokenLeftBracket
	TokenRightBracket
	TokenLeftBrace
	TokenRightBrace
	TokenPipe
	TokenLeftPipe
	TokenRightPipe
	TokenEquals
	TokenLessEqual
	TokenGreaterEqual
	TokenLess
	TokenGreater
	TokenNotEqual
	TokenSqrt
	TokenFrac
	TokenDfrac
	TokenLeft
	TokenRight
	TokenSubscript
	TokenSuperscript
	TokenLn
	TokenLog
	TokenSin
	TokenCos
	TokenTan
	TokenArcsin
	TokenArccos
	TokenArctan
	TokenCosh
	TokenSinh
	TokenTanh
	TokenSec
	TokenCsc
	TokenCot
	TokenAbs
	TokenPi
	TokenE
	TokenTheta
	TokenPhi
	TokenComma
	TokenExclamation
	TokenError
)

// Token represents a lexical token
type Token struct {
	Type  TokenType
	Value string
	Pos   int
}

// String returns a string representation of the token type
func (t TokenType) String() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenInt:
		return "INT"
	case TokenFloat:
		return "FLOAT"
	case TokenVar:
		return "VAR"
	case TokenPlus:
		return "+"
	case TokenMinus:
		return "-"
	case TokenMultiply:
		return "*"
	case TokenDivide:
		return "/"
	case TokenPower:
		return "^"
	case TokenLeftParen:
		return "("
	case TokenRightParen:
		return ")"
	case TokenSqrt:
		return "sqrt"
	case TokenLn:
		return "ln"
	case TokenLog:
		return "log"
	case TokenSin:
		return "sin"
	case TokenCos:
		return "cos"
	case TokenTan:
		return "tan"
	case TokenPi:
		return "pi"
	case TokenE:
		return "e"
	case TokenComma:
		return ","
	case TokenError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// TokenRule represents a tokenization rule
type TokenRule struct {
	Pattern   *regexp.Regexp
	TokenType TokenType
	Transform func(string) string // Optional transformation function
}

// Lexer tokenizes mathematical expressions
type Lexer struct {
	input string
	pos   int
	rules []TokenRule
}

// NewLexer creates a new lexer instance
func NewLexer(input string) *Lexer {
	lexer := &Lexer{
		input: input,
		pos:   0,
	}
	lexer.initRules()
	return lexer
}

// initRules initializes the tokenization rules
func (l *Lexer) initRules() {
	l.rules = []TokenRule{
		// Skip whitespace and LaTeX spacing
		{regexp.MustCompile(`^\s+`), TokenEOF, nil}, // Will be skipped
		{regexp.MustCompile(`^\\space`), TokenEOF, nil},
		{regexp.MustCompile(`^\\ `), TokenEOF, nil},

		// Numbers - float first to match decimals properly
		{regexp.MustCompile(`^[0-9]+\.[0-9]*`), TokenFloat, nil}, // Handle "1." and "1.23"
		{regexp.MustCompile(`^\.[0-9]+`), TokenFloat, nil},       // Handle ".5"
		{regexp.MustCompile(`^[0-9]+`), TokenInt, nil},

		// Operators
		{regexp.MustCompile(`^\*\*`), TokenPower, nil},
		{regexp.MustCompile(`^\*`), TokenMultiply, nil},
		{regexp.MustCompile(`^\\cdot`), TokenMultiply, nil},
		{regexp.MustCompile(`^\\times`), TokenMultiply, nil},
		{regexp.MustCompile(`^\\ast`), TokenMultiply, nil},
		{regexp.MustCompile(`^/`), TokenDivide, nil},
		{regexp.MustCompile(`^\\div`), TokenDivide, nil},
		{regexp.MustCompile(`^-`), TokenMinus, nil},
		{regexp.MustCompile("^\u2212"), TokenMinus, nil}, // Unicode minus
		{regexp.MustCompile(`^\+`), TokenPlus, nil},
		{regexp.MustCompile(`^\^`), TokenPower, nil},

		// Parentheses and brackets
		{regexp.MustCompile(`^\(`), TokenLeftParen, nil},
		{regexp.MustCompile(`^\)`), TokenRightParen, nil},
		{regexp.MustCompile(`^\\left\(`), TokenLeftParen, nil},
		{regexp.MustCompile(`^\\right\)`), TokenRightParen, nil},
		{regexp.MustCompile(`^\[`), TokenLeftBracket, nil},
		{regexp.MustCompile(`^\]`), TokenRightBracket, nil},
		{regexp.MustCompile(`^\{`), TokenLeftBrace, nil},
		{regexp.MustCompile(`^\}`), TokenRightBrace, nil},
		{regexp.MustCompile(`^\\left\{`), TokenLeftBrace, nil},
		{regexp.MustCompile(`^\\right\}`), TokenRightBrace, nil},

		// Comparison operators
		{regexp.MustCompile(`^<=`), TokenLessEqual, nil},
		{regexp.MustCompile(`^>=`), TokenGreaterEqual, nil},
		{regexp.MustCompile(`^<>`), TokenNotEqual, nil},
		{regexp.MustCompile(`^<`), TokenLess, nil},
		{regexp.MustCompile(`^>`), TokenGreater, nil},
		{regexp.MustCompile(`^=`), TokenEquals, nil},
		{regexp.MustCompile(`^\\le`), TokenLessEqual, func(s string) string { return "<=" }},
		{regexp.MustCompile(`^\\ge`), TokenGreaterEqual, func(s string) string { return ">=" }},
		{regexp.MustCompile(`^\\leq`), TokenLessEqual, func(s string) string { return "<=" }},
		{regexp.MustCompile(`^\\geq`), TokenGreaterEqual, func(s string) string { return ">=" }},
		{regexp.MustCompile(`^=/=`), TokenNotEqual, func(s string) string { return "<>" }},
		{regexp.MustCompile(`^\\ne`), TokenNotEqual, func(s string) string { return "<>" }},

		// Functions and special symbols
		{regexp.MustCompile(`^\\sqrt`), TokenSqrt, nil},
		{regexp.MustCompile(`^\\frac`), TokenFrac, nil},
		{regexp.MustCompile(`^\\dfrac`), TokenDfrac, nil},
		{regexp.MustCompile(`^\\ln`), TokenLn, nil},
		{regexp.MustCompile(`^\\log`), TokenLog, nil},

		// Trigonometric functions
		{regexp.MustCompile(`^\\arcsin`), TokenArcsin, nil},
		{regexp.MustCompile(`^\\arccos`), TokenArccos, nil},
		{regexp.MustCompile(`^\\arctan`), TokenArctan, nil},
		{regexp.MustCompile(`^\\sin`), TokenSin, nil},
		{regexp.MustCompile(`^\\cos`), TokenCos, nil},
		{regexp.MustCompile(`^\\tan`), TokenTan, nil},
		{regexp.MustCompile(`^\\sec`), TokenSec, nil},
		{regexp.MustCompile(`^\\csc`), TokenCsc, nil},
		{regexp.MustCompile(`^\\cot`), TokenCot, nil},

		// Hyperbolic functions
		{regexp.MustCompile(`^\\sinh`), TokenSinh, nil},
		{regexp.MustCompile(`^\\cosh`), TokenCosh, nil},
		{regexp.MustCompile(`^\\tanh`), TokenTanh, nil},

		// Constants (must be before single char variables)
		{regexp.MustCompile(`^pi`), TokenPi, func(s string) string { return "pi" }},
		{regexp.MustCompile(`^\\pi`), TokenPi, func(s string) string { return "pi" }},
		{regexp.MustCompile(`^\\theta`), TokenVar, func(s string) string { return "theta" }},
		{regexp.MustCompile(`^\\phi`), TokenVar, func(s string) string { return "phi" }},

		// Functions (must be before single char variables)
		{regexp.MustCompile(`^sqrt`), TokenSqrt, nil},
		{regexp.MustCompile(`^abs`), TokenAbs, nil},
		{regexp.MustCompile(`^ln`), TokenLn, nil},
		{regexp.MustCompile(`^log`), TokenLog, nil},
		{regexp.MustCompile(`^sin`), TokenSin, nil},
		{regexp.MustCompile(`^cos`), TokenCos, nil},
		{regexp.MustCompile(`^tan`), TokenTan, nil},

		// Known multi-character variables and constants (must be before single char variables)
		{regexp.MustCompile(`^theta`), TokenVar, func(s string) string { return "theta" }},
		{regexp.MustCompile(`^alpha`), TokenVar, func(s string) string { return "alpha" }},
		{regexp.MustCompile(`^beta`), TokenVar, func(s string) string { return "beta" }},
		{regexp.MustCompile(`^gamma`), TokenVar, func(s string) string { return "gamma" }},
		{regexp.MustCompile(`^delta`), TokenVar, func(s string) string { return "delta" }},
		{regexp.MustCompile(`^epsilon`), TokenVar, func(s string) string { return "epsilon" }},
		{regexp.MustCompile(`^phi`), TokenVar, func(s string) string { return "phi" }},
		{regexp.MustCompile(`^psi`), TokenVar, func(s string) string { return "psi" }},
		{regexp.MustCompile(`^omega`), TokenVar, func(s string) string { return "omega" }},

		// Other symbols
		{regexp.MustCompile(`^_`), TokenSubscript, nil},
		{regexp.MustCompile(`^\|`), TokenPipe, nil},
		{regexp.MustCompile(`^\\left\|`), TokenLeftPipe, nil},
		{regexp.MustCompile(`^\\right\|`), TokenRightPipe, nil},
		{regexp.MustCompile(`^,`), TokenComma, nil},
		{regexp.MustCompile(`^!`), TokenExclamation, nil},

		// Single character variables (everything else should be parsed as individual chars for implicit multiplication)
		{regexp.MustCompile(`^[a-zA-Z]`), TokenVar, nil},
	}
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	for l.pos < len(l.input) {
		// Skip whitespace and LaTeX spacing
		if match := regexp.MustCompile(`^\s+`).FindString(l.input[l.pos:]); match != "" {
			l.pos += len(match)
			continue
		}
		if match := regexp.MustCompile(`^\\space`).FindString(l.input[l.pos:]); match != "" {
			l.pos += len(match)
			continue
		}
		if match := regexp.MustCompile(`^\\ `).FindString(l.input[l.pos:]); match != "" {
			l.pos += len(match)
			continue
		}

		// Try to match each rule
		for _, rule := range l.rules {
			if match := rule.Pattern.FindString(l.input[l.pos:]); match != "" {
				token := Token{
					Type:  rule.TokenType,
					Value: match,
					Pos:   l.pos,
				}
				if rule.Transform != nil {
					token.Value = rule.Transform(match)
				}
				l.pos += len(match)

				// Handle special variable cases
				switch token.Value {
				case "pi":
					token.Type = TokenPi
				case "e":
					token.Type = TokenE
				case "ln":
					token.Type = TokenLn
				case "log":
					token.Type = TokenLog
				case "sin":
					token.Type = TokenSin
				case "cos":
					token.Type = TokenCos
				case "tan":
					token.Type = TokenTan
				case "sqrt":
					token.Type = TokenSqrt
				case "abs":
					token.Type = TokenAbs
				}

				return token
			}
		}

		// If no rule matched, it's an invalid character
		char := string(l.input[l.pos])
		return Token{
			Type:  TokenError,
			Value: char,
			Pos:   l.pos,
		}
	}

	return Token{Type: TokenEOF, Pos: l.pos}
}

// Peek returns the next token without advancing the position
func (l *Lexer) Peek() Token {
	savedPos := l.pos
	token := l.NextToken()
	l.pos = savedPos
	return token
}

// Reset resets the lexer position
func (l *Lexer) Reset() {
	l.pos = 0
}

// Position returns current position
func (l *Lexer) Position() int {
	return l.pos
}

// SetPosition sets the lexer position
func (l *Lexer) SetPosition(pos int) {
	l.pos = pos
}