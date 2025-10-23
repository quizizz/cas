package ast

import (
	"fmt"
	"math/big"
	"strings"
)

// Var represents a variable
type Var struct {
	name string
}

// NewVar creates a new variable expression
func NewVar(name string) *Var {
	return &Var{name: name}
}

func (v *Var) String() string {
	return v.name
}

func (v *Var) LaTeX() string {
	return v.name
}

func (v *Var) Eval(vars map[string]*big.Float) (*big.Float, error) {
	val, ok := vars[v.name]
	if !ok {
		return nil, fmt.Errorf("undefined variable: %s", v.name)
	}
	return new(big.Float).Copy(val), nil
}

func (v *Var) Simplify() Expr {
	return v.Clone()
}

func (v *Var) Equal(other Expr) bool {
	if other.Type() != TypeVar {
		return false
	}
	otherVar := other.(*Var)
	return v.name == otherVar.name
}

func (v *Var) Clone() Expr {
	return &Var{name: v.name}
}

func (v *Var) Variables() []string {
	return []string{v.name}
}

func (v *Var) Type() ExprType {
	return TypeVar
}

func (v *Var) Name() string {
	return v.name
}

// Const represents a mathematical constant (pi, e, etc.)
type Const struct {
	name  string
	value *big.Float
}

// Common mathematical constants
var (
	Pi = &Const{
		name:  "pi",
		value: func() *big.Float {
			pi := new(big.Float)
			pi.SetString("3.1415926535897932384626433832795028841971693993751")
			return pi
		}(),
	}
	E = &Const{
		name:  "e",
		value: func() *big.Float {
			e := new(big.Float)
			e.SetString("2.7182818284590452353602874713526624977572470937000")
			return e
		}(),
	}
)

// NewConst creates a new constant expression
func NewConst(name string, value *big.Float) *Const {
	return &Const{name: name, value: new(big.Float).Copy(value)}
}

func (c *Const) String() string {
	return c.name
}

func (c *Const) LaTeX() string {
	switch c.name {
	case "pi":
		return "\\pi"
	default:
		return c.name
	}
}

func (c *Const) Eval(vars map[string]*big.Float) (*big.Float, error) {
	return new(big.Float).Copy(c.value), nil
}

func (c *Const) Simplify() Expr {
	return c.Clone()
}

func (c *Const) Equal(other Expr) bool {
	if other.Type() != TypeConst {
		return false
	}
	otherConst := other.(*Const)
	return c.name == otherConst.name
}

func (c *Const) Clone() Expr {
	return &Const{
		name:  c.name,
		value: new(big.Float).Copy(c.value),
	}
}

func (c *Const) Variables() []string {
	return []string{}
}

func (c *Const) Type() ExprType {
	return TypeConst
}

func (c *Const) Name() string {
	return c.name
}

func (c *Const) Value() *big.Float {
	return new(big.Float).Copy(c.value)
}

// Func represents a function call
type Func struct {
	name string
	args []Expr
}

// NewFunc creates a new function expression
func NewFunc(name string, args ...Expr) *Func {
	return &Func{name: name, args: args}
}

func (f *Func) String() string {
	if len(f.args) == 0 {
		return f.name + "()"
	}

	argStrs := make([]string, len(f.args))
	for i, arg := range f.args {
		argStrs[i] = arg.String()
	}
	return fmt.Sprintf("%s(%s)", f.name, strings.Join(argStrs, ", "))
}

func (f *Func) LaTeX() string {
	switch f.name {
	case "sqrt":
		if len(f.args) == 1 {
			return fmt.Sprintf("\\sqrt{%s}", f.args[0].LaTeX())
		}
	case "log":
		if len(f.args) == 1 {
			return fmt.Sprintf("\\log{%s}", f.args[0].LaTeX())
		}
	case "ln":
		if len(f.args) == 1 {
			return fmt.Sprintf("\\ln{%s}", f.args[0].LaTeX())
		}
	}

	// Default function representation
	argStrs := make([]string, len(f.args))
	for i, arg := range f.args {
		argStrs[i] = arg.LaTeX()
	}
	return fmt.Sprintf("\\mathrm{%s}(%s)", f.name, strings.Join(argStrs, ", "))
}

func (f *Func) Eval(vars map[string]*big.Float) (*big.Float, error) {
	// Evaluate arguments first
	argVals := make([]*big.Float, len(f.args))
	for i, arg := range f.args {
		val, err := arg.Eval(vars)
		if err != nil {
			return nil, err
		}
		argVals[i] = val
	}

	// Apply function
	switch f.name {
	case "sqrt":
		if len(argVals) != 1 {
			return nil, fmt.Errorf("sqrt expects 1 argument, got %d", len(argVals))
		}
		result := new(big.Float)
		result.Sqrt(argVals[0])
		return result, nil
	case "abs":
		if len(argVals) != 1 {
			return nil, fmt.Errorf("abs expects 1 argument, got %d", len(argVals))
		}
		result := new(big.Float).Copy(argVals[0])
		result.Abs(result)
		return result, nil
	case "ln":
		if len(argVals) != 1 {
			return nil, fmt.Errorf("ln expects 1 argument, got %d", len(argVals))
		}
		return evaluateNaturalLog(argVals[0])
	case "log":
		if len(argVals) == 1 {
			// Base 10 logarithm
			return evaluateLog10(argVals[0])
		} else if len(argVals) == 2 {
			// Logarithm with custom base: log_b(x) = ln(x) / ln(b)
			return evaluateLogBase(argVals[0], argVals[1])
		}
		return nil, fmt.Errorf("log expects 1 or 2 arguments, got %d", len(argVals))
	case "sin":
		if len(argVals) != 1 {
			return nil, fmt.Errorf("sin expects 1 argument, got %d", len(argVals))
		}
		return evaluateSin(argVals[0])
	case "cos":
		if len(argVals) != 1 {
			return nil, fmt.Errorf("cos expects 1 argument, got %d", len(argVals))
		}
		return evaluateCos(argVals[0])
	case "tan":
		if len(argVals) != 1 {
			return nil, fmt.Errorf("tan expects 1 argument, got %d", len(argVals))
		}
		return evaluateTan(argVals[0])
	case "arcsin", "asin":
		if len(argVals) != 1 {
			return nil, fmt.Errorf("arcsin expects 1 argument, got %d", len(argVals))
		}
		return evaluateArcsin(argVals[0])
	case "arccos", "acos":
		if len(argVals) != 1 {
			return nil, fmt.Errorf("arccos expects 1 argument, got %d", len(argVals))
		}
		return evaluateArccos(argVals[0])
	case "arctan", "atan":
		if len(argVals) != 1 {
			return nil, fmt.Errorf("arctan expects 1 argument, got %d", len(argVals))
		}
		return evaluateArctan(argVals[0])
	case "sinh":
		if len(argVals) != 1 {
			return nil, fmt.Errorf("sinh expects 1 argument, got %d", len(argVals))
		}
		return evaluateSinh(argVals[0])
	case "cosh":
		if len(argVals) != 1 {
			return nil, fmt.Errorf("cosh expects 1 argument, got %d", len(argVals))
		}
		return evaluateCosh(argVals[0])
	case "tanh":
		if len(argVals) != 1 {
			return nil, fmt.Errorf("tanh expects 1 argument, got %d", len(argVals))
		}
		return evaluateTanh(argVals[0])
	default:
		return nil, fmt.Errorf("unsupported function: %s", f.name)
	}
}

func (f *Func) Simplify() Expr {
	simplifiedArgs := make([]Expr, len(f.args))
	for i, arg := range f.args {
		simplifiedArgs[i] = arg.Simplify()
	}
	return &Func{name: f.name, args: simplifiedArgs}
}

func (f *Func) Equal(other Expr) bool {
	if other.Type() != TypeFunc {
		return false
	}
	otherFunc := other.(*Func)
	if f.name != otherFunc.name || len(f.args) != len(otherFunc.args) {
		return false
	}
	for i, arg := range f.args {
		if !arg.Equal(otherFunc.args[i]) {
			return false
		}
	}
	return true
}

func (f *Func) Clone() Expr {
	clonedArgs := make([]Expr, len(f.args))
	for i, arg := range f.args {
		clonedArgs[i] = arg.Clone()
	}
	return &Func{name: f.name, args: clonedArgs}
}

func (f *Func) Variables() []string {
	vars := []string{}
	for _, arg := range f.args {
		vars = append(vars, arg.Variables()...)
	}
	return removeDuplicates(vars)
}

func (f *Func) Type() ExprType {
	return TypeFunc
}

func (f *Func) Name() string {
	return f.name
}

func (f *Func) Args() []Expr {
	result := make([]Expr, len(f.args))
	for i, arg := range f.args {
		result[i] = arg.Clone()
	}
	return result
}

// Helper function to remove duplicate strings
func removeDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}