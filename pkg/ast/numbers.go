package ast

import (
	"fmt"
	"math/big"
)

// Int represents an integer constant
type Int struct {
	value *big.Int
}

// NewInt creates a new integer expression
func NewInt(value int64) *Int {
	return &Int{value: big.NewInt(value)}
}

// NewIntFromString creates an integer from string representation
func NewIntFromString(s string) (*Int, error) {
	val := new(big.Int)
	if _, ok := val.SetString(s, 10); !ok {
		return nil, fmt.Errorf("invalid integer: %s", s)
	}
	return &Int{value: val}, nil
}

func (i *Int) String() string {
	return i.value.String()
}

func (i *Int) LaTeX() string {
	return i.value.String()
}

func (i *Int) Eval(vars map[string]*big.Float) (*big.Float, error) {
	result := new(big.Float)
	result.SetInt(i.value)
	return result, nil
}

func (i *Int) Simplify() Expr {
	return i.Clone()
}

func (i *Int) Equal(other Expr) bool {
	if other.Type() != TypeInt {
		return false
	}
	otherInt := other.(*Int)
	return i.value.Cmp(otherInt.value) == 0
}

func (i *Int) Clone() Expr {
	newVal := new(big.Int).Set(i.value)
	return &Int{value: newVal}
}

func (i *Int) Variables() []string {
	return []string{}
}

func (i *Int) Type() ExprType {
	return TypeInt
}

func (i *Int) Value() *big.Float {
	result := new(big.Float)
	result.SetInt(i.value)
	return result
}

// IntValue returns the underlying *big.Int value
func (i *Int) IntValue() *big.Int {
	return new(big.Int).Set(i.value)
}

// Float represents a floating-point constant
type Float struct {
	value *big.Float
}

// NewFloat creates a new float expression
func NewFloat(value float64) *Float {
	return &Float{value: big.NewFloat(value)}
}

// NewFloatFromString creates a float from string representation
func NewFloatFromString(s string) (*Float, error) {
	val, _, err := big.ParseFloat(s, 10, 256, big.ToNearestEven)
	if err != nil {
		return nil, fmt.Errorf("invalid float: %s", s)
	}
	return &Float{value: val}, nil
}

func (f *Float) String() string {
	return f.value.Text('g', -1)
}

func (f *Float) LaTeX() string {
	return f.value.Text('g', -1)
}

func (f *Float) Eval(vars map[string]*big.Float) (*big.Float, error) {
	return new(big.Float).Copy(f.value), nil
}

func (f *Float) Simplify() Expr {
	return f.Clone()
}

func (f *Float) Equal(other Expr) bool {
	if other.Type() != TypeFloat {
		return false
	}
	otherFloat := other.(*Float)
	return f.value.Cmp(otherFloat.value) == 0
}

func (f *Float) Clone() Expr {
	return &Float{value: new(big.Float).Copy(f.value)}
}

func (f *Float) Variables() []string {
	return []string{}
}

func (f *Float) Type() ExprType {
	return TypeFloat
}

func (f *Float) Value() *big.Float {
	return new(big.Float).Copy(f.value)
}

// Rational represents a rational number (fraction)
type Rational struct {
	numerator   *big.Int
	denominator *big.Int
}

// NewRational creates a new rational expression
func NewRational(num, den int64) *Rational {
	r := &Rational{
		numerator:   big.NewInt(num),
		denominator: big.NewInt(den),
	}
	r.reduce()
	return r
}

// NewRationalPreserved creates a new rational expression without automatic reduction
// Used for KAS compatibility to preserve original fraction representation
func NewRationalPreserved(num, den int64) *Rational {
	return &Rational{
		numerator:   big.NewInt(num),
		denominator: big.NewInt(den),
	}
}

// NewRationalFromInts creates a rational from big integers
func NewRationalFromInts(num, den *big.Int) *Rational {
	r := &Rational{
		numerator:   new(big.Int).Set(num),
		denominator: new(big.Int).Set(den),
	}
	r.reduce()
	return r
}

// reduce simplifies the rational to lowest terms
func (r *Rational) reduce() {
	gcd := new(big.Int).GCD(nil, nil, r.numerator, r.denominator)
	r.numerator.Div(r.numerator, gcd)
	r.denominator.Div(r.denominator, gcd)

	// Ensure denominator is positive
	if r.denominator.Sign() < 0 {
		r.numerator.Neg(r.numerator)
		r.denominator.Neg(r.denominator)
	}
}

func (r *Rational) String() string {
	// For KAS compatibility, always show fractions as num/den even when den=1
	return fmt.Sprintf("%s/%s", r.numerator.String(), r.denominator.String())
}

func (r *Rational) LaTeX() string {
	if r.denominator.Cmp(big.NewInt(1)) == 0 {
		return r.numerator.String()
	}
	return fmt.Sprintf("\\frac{%s}{%s}", r.numerator.String(), r.denominator.String())
}

func (r *Rational) Eval(vars map[string]*big.Float) (*big.Float, error) {
	num := new(big.Float).SetInt(r.numerator)
	den := new(big.Float).SetInt(r.denominator)
	result := new(big.Float).Quo(num, den)
	return result, nil
}

func (r *Rational) Simplify() Expr {
	return r.Clone()
}

func (r *Rational) Equal(other Expr) bool {
	if other.Type() != TypeRational {
		return false
	}
	otherRat := other.(*Rational)
	return r.numerator.Cmp(otherRat.numerator) == 0 && r.denominator.Cmp(otherRat.denominator) == 0
}

func (r *Rational) Clone() Expr {
	return &Rational{
		numerator:   new(big.Int).Set(r.numerator),
		denominator: new(big.Int).Set(r.denominator),
	}
}

func (r *Rational) Variables() []string {
	return []string{}
}

func (r *Rational) Type() ExprType {
	return TypeRational
}

func (r *Rational) Value() *big.Float {
	num := new(big.Float).SetInt(r.numerator)
	den := new(big.Float).SetInt(r.denominator)
	return new(big.Float).Quo(num, den)
}

// Numerator returns the numerator
func (r *Rational) Numerator() *big.Int {
	return new(big.Int).Set(r.numerator)
}

// Denominator returns the denominator
func (r *Rational) Denominator() *big.Int {
	return new(big.Int).Set(r.denominator)
}
