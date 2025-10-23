// Package ast defines the abstract syntax tree nodes for mathematical expressions.
package ast

import (
	"math/big"
)

// Expr represents any mathematical expression node in the AST.
// All expression types implement this interface.
type Expr interface {
	// String returns a string representation of the expression
	String() string

	// LaTeX returns a LaTeX representation of the expression
	LaTeX() string

	// Eval evaluates the expression with given variable values
	Eval(vars map[string]*big.Float) (*big.Float, error)

	// Simplify returns a simplified version of the expression
	Simplify() Expr

	// Equal checks if two expressions are mathematically equivalent
	Equal(other Expr) bool

	// Clone returns a deep copy of the expression
	Clone() Expr

	// Variables returns all variable names used in the expression
	Variables() []string

	// Type returns the expression type for type checking
	Type() ExprType
}

// ExprType represents the type of expression node
type ExprType int

const (
	TypeAdd ExprType = iota
	TypeMul
	TypePow
	TypeVar
	TypeConst
	TypeInt
	TypeFloat
	TypeRational
	TypeFunc
	TypeTrig
	TypeLog
	TypeAbs
	TypeEq
)

// String returns the string representation of the expression type
func (t ExprType) String() string {
	switch t {
	case TypeAdd:
		return "Add"
	case TypeMul:
		return "Mul"
	case TypePow:
		return "Pow"
	case TypeVar:
		return "Var"
	case TypeConst:
		return "Const"
	case TypeInt:
		return "Int"
	case TypeFloat:
		return "Float"
	case TypeRational:
		return "Rational"
	case TypeFunc:
		return "Func"
	case TypeTrig:
		return "Trig"
	case TypeLog:
		return "Log"
	case TypeAbs:
		return "Abs"
	case TypeEq:
		return "Eq"
	default:
		return "Unknown"
	}
}

// Seq represents a sequence of expressions (base for Add and Mul)
type Seq interface {
	Expr
	// Terms returns the terms in the sequence
	Terms() []Expr
	// AddTerm adds a new term to the sequence
	AddTerm(term Expr) Seq
}

// Binary represents a binary operation
type Binary interface {
	Expr
	// Left returns the left operand
	Left() Expr
	// Right returns the right operand
	Right() Expr
}

// Unary represents a unary operation
type Unary interface {
	Expr
	// Operand returns the operand
	Operand() Expr
}

// Symbol represents symbolic expressions (variables, constants, functions)
type Symbol interface {
	Expr
	// Name returns the symbol name
	Name() string
}

// Numeric represents numeric expressions
type Numeric interface {
	Expr
	// Value returns the numeric value
	Value() *big.Float
}