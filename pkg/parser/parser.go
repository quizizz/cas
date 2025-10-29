package parser

import (
	"fmt"
	"strconv"

	"github.com/quizizz/cas/pkg/ast"
)

// Parser implements a recursive descent parser for mathematical expressions
type Parser struct {
	lexer   *Lexer
	current Token
}

// New creates a new parser instance
func New(input string) *Parser {
	lexer := NewLexer(input)
	parser := &Parser{
		lexer: lexer,
	}
	parser.advance()
	return parser
}

// Parse parses the input expression and returns an AST node
func Parse(input string) (ast.Expr, error) {
	parser := New(input)
	// Check for lexical errors first
	if parser.current.Type == TokenError {
		return nil, fmt.Errorf("invalid character '%s' at position %d", parser.current.Value, parser.current.Pos)
	}
	return parser.parseExpression()
}

// advance moves to the next token
func (p *Parser) advance() {
	p.current = p.lexer.NextToken()
	// Check for lexical errors during parsing
	if p.current.Type == TokenError {
		return
	}
}

// peek returns the next token without advancing
func (p *Parser) peek() Token {
	return p.lexer.Peek()
}

// expect consumes a token of the given type or returns an error
func (p *Parser) expect(tokenType TokenType) error {
	if p.current.Type == TokenError {
		return fmt.Errorf("invalid character '%s' at position %d", p.current.Value, p.current.Pos)
	}
	if p.current.Type != tokenType {
		return fmt.Errorf("expected %s, got %s at position %d", tokenType, p.current.Type, p.current.Pos)
	}
	p.advance()
	return nil
}

// parseExpression parses a full expression (handles equations)
func (p *Parser) parseExpression() (ast.Expr, error) {
	// Check for lexical errors
	if p.current.Type == TokenError {
		return nil, fmt.Errorf("invalid character '%s' at position %d", p.current.Value, p.current.Pos)
	}

	left, err := p.parseArithmeticExpression()
	if err != nil {
		return nil, err
	}

	// Check for comparison operators (equations/inequalities)
	switch p.current.Type {
	case TokenEquals, TokenLess, TokenGreater, TokenLessEqual, TokenGreaterEqual, TokenNotEqual:
		operator := p.current.Type
		p.advance()
		right, err := p.parseArithmeticExpression()
		if err != nil {
			return nil, err
		}

		// Map token types to equation types
		var eqType ast.EqType
		switch operator {
		case TokenEquals:
			eqType = ast.EqEqual
		case TokenLess:
			eqType = ast.EqLess
		case TokenGreater:
			eqType = ast.EqGreater
		case TokenLessEqual:
			eqType = ast.EqLessEqual
		case TokenGreaterEqual:
			eqType = ast.EqGreaterEqual
		case TokenNotEqual:
			eqType = ast.EqNotEqual
		}

		return ast.NewEq(left, right, eqType), nil
	}

	return left, nil
}

// parseArithmeticExpression parses addition and subtraction (lowest precedence)
func (p *Parser) parseArithmeticExpression() (ast.Expr, error) {
	left, err := p.parseMultiplicativeExpression()
	if err != nil {
		return nil, err
	}

	for p.current.Type == TokenPlus || p.current.Type == TokenMinus {
		op := p.current.Type
		p.advance()
		right, err := p.parseMultiplicativeExpression()
		if err != nil {
			return nil, err
		}

		if op == TokenPlus {
			left = ast.NewAdd(left, right)
		} else {
			// Handle subtraction as addition of negative
			negatedRight := ast.NewMul(ast.NewInt(-1), right)
			left = ast.NewAdd(left, negatedRight)
		}
	}

	return left, nil
}

// parseMultiplicativeExpression parses multiplication and division
func (p *Parser) parseMultiplicativeExpression() (ast.Expr, error) {
	left, err := p.parseUnaryExpression()
	if err != nil {
		return nil, err
	}

	for p.current.Type == TokenMultiply || p.current.Type == TokenDivide || p.isImplicitMultiplication() {
		var op TokenType
		if p.current.Type == TokenMultiply || p.current.Type == TokenDivide {
			op = p.current.Type
			p.advance()
		} else {
			// Implicit multiplication
			op = TokenMultiply
		}

		right, err := p.parseUnaryExpression()
		if err != nil {
			return nil, err
		}

		if op == TokenMultiply {
			left = ast.NewMul(left, right)
		} else {
			// Handle division - convert to multiplication by reciprocal for KAS compatibility
			reciprocal := ast.NewPow(right, ast.NewInt(-1))
			left = ast.NewMul(left, reciprocal)
		}
	}

	return left, nil
}

// isImplicitMultiplication checks if the current position indicates implicit multiplication
func (p *Parser) isImplicitMultiplication() bool {
	switch p.current.Type {
	case TokenVar, TokenLeftParen, TokenLeftBrace, TokenSqrt, TokenFrac, TokenDfrac, TokenLn, TokenLog, TokenSin, TokenCos, TokenTan, TokenAbs, TokenPi, TokenE:
		return true
	default:
		return false
	}
}

// parseExponentialExpression parses exponentiation (right-associative)
func (p *Parser) parseExponentialExpression() (ast.Expr, error) {
	left, err := p.parsePrimaryExpression()
	if err != nil {
		return nil, err
	}

	if p.current.Type == TokenPower {
		p.advance()
		// Right-associative: a^b^c = a^(b^c)
		// For unary minus in exponent, we need to parse unary expressions
		right, err := p.parseUnaryExpression()
		if err != nil {
			return nil, err
		}
		return ast.NewPow(left, right), nil
	}

	return left, nil
}

// parseUnaryExpression parses unary operators (negation and unary plus)
func (p *Parser) parseUnaryExpression() (ast.Expr, error) {
	if p.current.Type == TokenPlus {
		// For KAS compatibility, reject unary plus - KAS doesn't parse "+49"
		return nil, fmt.Errorf("unexpected token '+' at position %d", p.current.Pos)
	}
	if p.current.Type == TokenMinus {
		p.advance()
		// Recursively handle multiple minuses like "--x"
		operand, err := p.parseUnaryExpression()
		if err != nil {
			return nil, err
		}

		// Special case: if operand is a positive integer, create negative integer directly
		// Exception: for zero, KAS expects -1*0 representation
		if operand.Type() == ast.TypeInt {
			intVal := operand.(*ast.Int)
			floatVal := intVal.Value()
			if floatVal.IsInt() {
				intValue, _ := floatVal.Int64()
				if intValue == 0 {
					// For KAS compatibility: -0 should be -1*0
					return ast.NewMul(ast.NewInt(-1), operand), nil
				}
				if intVal.Value().Sign() > 0 {
					return ast.NewInt(-intValue), nil
				}
			}
		}

		// Special case: if operand is a positive float, create negative float directly
		if operand.Type() == ast.TypeFloat {
			if floatVal := operand.(*ast.Float); floatVal.Value().Sign() > 0 {
				positiveValue, _ := floatVal.Value().Float64()
				return ast.NewFloat(-positiveValue), nil
			}
		}

		return ast.NewMul(ast.NewInt(-1), operand), nil
	}

	if p.current.Type == TokenPlus {
		p.advance()
		return p.parseExponentialExpression()
	}

	return p.parseExponentialExpression()
}

// parsePrimaryExpression parses primary expressions (atoms)
func (p *Parser) parsePrimaryExpression() (ast.Expr, error) {
	switch p.current.Type {
	case TokenInt:
		return p.parseInteger()
	case TokenFloat:
		return p.parseFloat()
	case TokenVar:
		return p.parseVariable()
	case TokenPi:
		return p.parseConstant()
	case TokenE:
		return p.parseConstant()
	case TokenLeftParen:
		return p.parseParentheses()
	case TokenLeftBrace:
		return p.parseBraces()
	case TokenSqrt:
		return p.parseSqrt()
	case TokenFrac, TokenDfrac:
		return p.parseFrac()
	case TokenLn, TokenLog:
		return p.parseLogFunction()
	case TokenSin, TokenCos, TokenTan, TokenArcsin, TokenArccos, TokenArctan:
		return p.parseTrigFunction()
	case TokenSinh, TokenCosh, TokenTanh:
		return p.parseHyperbolicFunction()
	case TokenAbs, TokenLeftPipe:
		return p.parseAbsoluteValue()
	default:
		return nil, fmt.Errorf("unexpected token %s at position %d", p.current.Type, p.current.Pos)
	}
}

// parseInteger parses integer literals
func (p *Parser) parseInteger() (ast.Expr, error) {
	value := p.current.Value
	p.advance()

	// Handle subscripts
	if p.current.Type == TokenSubscript {
		return p.parseSubscriptedVariable(value)
	}

	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid integer: %s", value)
	}

	return ast.NewInt(intValue), nil
}

// parseFloat parses floating-point literals
func (p *Parser) parseFloat() (ast.Expr, error) {
	value := p.current.Value
	p.advance()

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid float: %s", value)
	}

	return ast.NewFloat(floatValue), nil
}

// parseVariable parses variables and function calls
func (p *Parser) parseVariable() (ast.Expr, error) {
	name := p.current.Value
	p.advance()

	// Handle subscripts
	if p.current.Type == TokenSubscript {
		return p.parseSubscriptedVariable(name)
	}

	// Handle function calls - any variable followed by parentheses is a function call
	if p.current.Type == TokenLeftParen {
		return p.parseFunctionCall(name)
	}

	return ast.NewVar(name), nil
}

// parseSubscriptedVariable parses variables with subscripts (e.g., x_1, x_n)
func (p *Parser) parseSubscriptedVariable(baseName string) (ast.Expr, error) {
	if err := p.expect(TokenSubscript); err != nil {
		return nil, err
	}

	subscript := ""
	if p.current.Type == TokenLeftBrace {
		p.advance()
		for p.current.Type != TokenRightBrace && p.current.Type != TokenEOF {
			subscript += p.current.Value
			p.advance()
		}
		if err := p.expect(TokenRightBrace); err != nil {
			return nil, err
		}
	} else {
		subscript = p.current.Value
		p.advance()
	}

	varName := baseName + "_" + subscript
	return ast.NewVar(varName), nil
}

// parseConstant parses mathematical constants
func (p *Parser) parseConstant() (ast.Expr, error) {
	switch p.current.Type {
	case TokenPi:
		p.advance()
		return ast.Pi, nil
	case TokenE:
		p.advance()
		return ast.E, nil
	default:
		return nil, fmt.Errorf("unknown constant: %s", p.current.Value)
	}
}

// parseParentheses parses parenthesized expressions
func (p *Parser) parseParentheses() (ast.Expr, error) {
	if err := p.expect(TokenLeftParen); err != nil {
		return nil, err
	}

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if err := p.expect(TokenRightParen); err != nil {
		return nil, err
	}

	return expr, nil
}

// parseBraces parses braced expressions (similar to parentheses in LaTeX)
func (p *Parser) parseBraces() (ast.Expr, error) {
	if err := p.expect(TokenLeftBrace); err != nil {
		return nil, err
	}

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if err := p.expect(TokenRightBrace); err != nil {
		return nil, err
	}

	return expr, nil
}

// parseSqrt parses square root expressions, including \sqrt[n]{x} syntax
func (p *Parser) parseSqrt() (ast.Expr, error) {
	if err := p.expect(TokenSqrt); err != nil {
		return nil, err
	}

	// Handle optional root index [n] in \sqrt[n]{x}
	var rootIndex ast.Expr
	if p.current.Type == TokenLeftBracket {
		p.advance() // consume [
		var err error
		rootIndex, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
		if err := p.expect(TokenRightBracket); err != nil {
			return nil, err
		}
	}

	var operand ast.Expr
	var err error

	if p.current.Type == TokenLeftBrace {
		p.advance()
		operand, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
		if err := p.expect(TokenRightBrace); err != nil {
			return nil, err
		}
	} else {
		operand, err = p.parsePrimaryExpression()
		if err != nil {
			return nil, err
		}
	}

	// If root index is specified, create a power expression: x^(1/n)
	if rootIndex != nil {
		// \sqrt[n]{x} = x^(1/n)
		reciprocalIndex := ast.NewPow(rootIndex, ast.NewInt(-1))
		return ast.NewPow(operand, reciprocalIndex), nil
	}

	// Create sqrt function call for standard square root
	return ast.NewFunc("sqrt", operand), nil
}

// parseFrac parses fraction expressions
func (p *Parser) parseFrac() (ast.Expr, error) {
	p.advance() // consume \frac or \dfrac

	if err := p.expect(TokenLeftBrace); err != nil {
		return nil, err
	}

	numerator, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if err := p.expect(TokenRightBrace); err != nil {
		return nil, err
	}

	if err := p.expect(TokenLeftBrace); err != nil {
		return nil, err
	}

	denominator, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if err := p.expect(TokenRightBrace); err != nil {
		return nil, err
	}

	// Create division as multiplication by reciprocal
	reciprocal := ast.NewPow(denominator, ast.NewInt(-1))
	return ast.NewMul(numerator, reciprocal), nil
}

// parseLogFunction parses logarithm functions
func (p *Parser) parseLogFunction() (ast.Expr, error) {
	var funcName string
	switch p.current.Type {
	case TokenLn:
		funcName = "ln"
	case TokenLog:
		funcName = "log"
	}
	p.advance()

	var operand ast.Expr
	var err error

	if p.current.Type == TokenLeftBrace {
		p.advance()
		operand, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
		if err := p.expect(TokenRightBrace); err != nil {
			return nil, err
		}
	} else if p.current.Type == TokenLeftParen {
		return p.parseFunctionCall(funcName)
	} else {
		operand, err = p.parsePrimaryExpression()
		if err != nil {
			return nil, err
		}
	}

	return ast.NewFunc(funcName, operand), nil
}

// parseTrigFunction parses trigonometric functions
func (p *Parser) parseTrigFunction() (ast.Expr, error) {
	var funcName string
	switch p.current.Type {
	case TokenSin:
		funcName = "sin"
	case TokenCos:
		funcName = "cos"
	case TokenTan:
		funcName = "tan"
	case TokenArcsin:
		funcName = "arcsin"
	case TokenArccos:
		funcName = "arccos"
	case TokenArctan:
		funcName = "arctan"
	}
	p.advance()

	var operand ast.Expr
	var err error

	if p.current.Type == TokenLeftBrace {
		p.advance()
		operand, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
		if err := p.expect(TokenRightBrace); err != nil {
			return nil, err
		}
	} else if p.current.Type == TokenLeftParen {
		return p.parseFunctionCall(funcName)
	} else {
		operand, err = p.parsePrimaryExpression()
		if err != nil {
			return nil, err
		}
	}

	return ast.NewFunc(funcName, operand), nil
}

// parseHyperbolicFunction parses hyperbolic functions
func (p *Parser) parseHyperbolicFunction() (ast.Expr, error) {
	var funcName string
	switch p.current.Type {
	case TokenSinh:
		funcName = "sinh"
	case TokenCosh:
		funcName = "cosh"
	case TokenTanh:
		funcName = "tanh"
	}
	p.advance()

	var operand ast.Expr
	var err error

	if p.current.Type == TokenLeftBrace {
		p.advance()
		operand, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
		if err := p.expect(TokenRightBrace); err != nil {
			return nil, err
		}
	} else if p.current.Type == TokenLeftParen {
		return p.parseFunctionCall(funcName)
	} else {
		operand, err = p.parsePrimaryExpression()
		if err != nil {
			return nil, err
		}
	}

	return ast.NewFunc(funcName, operand), nil
}

// parseAbsoluteValue parses absolute value expressions
func (p *Parser) parseAbsoluteValue() (ast.Expr, error) {
	if p.current.Type == TokenAbs {
		return p.parseTrigFunction() // Reuse trig function parsing logic
	}

	// Handle |expression| syntax
	if err := p.expect(TokenLeftPipe); err != nil {
		return nil, err
	}

	operand, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if err := p.expect(TokenRightPipe); err != nil {
		return nil, err
	}

	return ast.NewFunc("abs", operand), nil
}

// parseFunctionCall parses function calls with parentheses
func (p *Parser) parseFunctionCall(funcName string) (ast.Expr, error) {
	if err := p.expect(TokenLeftParen); err != nil {
		return nil, err
	}

	var args []ast.Expr

	// Empty function calls are not allowed
	if p.current.Type == TokenRightParen {
		return nil, fmt.Errorf("empty function call not allowed: %s()", funcName)
	}

	// Parse first argument
	arg, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	args = append(args, arg)

	// Parse additional arguments
	for p.current.Type == TokenComma {
		p.advance()
		arg, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	if err := p.expect(TokenRightParen); err != nil {
		return nil, err
	}

	return ast.NewFunc(funcName, args...), nil
}
