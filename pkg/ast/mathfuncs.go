package ast

import (
	"fmt"
	"math"
	"math/big"
)

// Mathematical function evaluation helpers using Go's math package
// These convert big.Float to float64 for computation and back to big.Float

// evaluateNaturalLog computes the natural logarithm
func evaluateNaturalLog(x *big.Float) (*big.Float, error) {
	if x.Sign() <= 0 {
		return nil, fmt.Errorf("ln: domain error (argument must be positive)")
	}

	// Special cases
	if x.Cmp(E.value) == 0 {
		return big.NewFloat(1), nil
	}
	if x.Cmp(big.NewFloat(1)) == 0 {
		return big.NewFloat(0), nil
	}

	// Convert to float64, compute, and convert back
	xFloat, _ := x.Float64()
	result := math.Log(xFloat)
	return big.NewFloat(result), nil
}

// evaluateLog10 computes the base-10 logarithm
func evaluateLog10(x *big.Float) (*big.Float, error) {
	if x.Sign() <= 0 {
		return nil, fmt.Errorf("log: domain error (argument must be positive)")
	}

	// Special cases
	if x.Cmp(big.NewFloat(10)) == 0 {
		return big.NewFloat(1), nil
	}
	if x.Cmp(big.NewFloat(1)) == 0 {
		return big.NewFloat(0), nil
	}

	xFloat, _ := x.Float64()
	result := math.Log10(xFloat)
	return big.NewFloat(result), nil
}

// evaluateLogBase computes logarithm with arbitrary base
func evaluateLogBase(x, base *big.Float) (*big.Float, error) {
	if x.Sign() <= 0 {
		return nil, fmt.Errorf("log: domain error (argument must be positive)")
	}
	if base.Sign() <= 0 || base.Cmp(big.NewFloat(1)) == 0 {
		return nil, fmt.Errorf("log: domain error (base must be positive and not equal to 1)")
	}

	// log_b(x) = ln(x) / ln(b)
	lnX, err := evaluateNaturalLog(x)
	if err != nil {
		return nil, err
	}
	lnBase, err := evaluateNaturalLog(base)
	if err != nil {
		return nil, err
	}

	result := new(big.Float).Quo(lnX, lnBase)
	return result, nil
}

// evaluateSin computes the sine function
func evaluateSin(x *big.Float) (*big.Float, error) {
	xFloat, _ := x.Float64()
	result := math.Sin(xFloat)
	return big.NewFloat(result), nil
}

// evaluateCos computes the cosine function
func evaluateCos(x *big.Float) (*big.Float, error) {
	xFloat, _ := x.Float64()
	result := math.Cos(xFloat)
	return big.NewFloat(result), nil
}

// evaluateTan computes the tangent function
func evaluateTan(x *big.Float) (*big.Float, error) {
	xFloat, _ := x.Float64()
	result := math.Tan(xFloat)
	return big.NewFloat(result), nil
}

// evaluateArcsin computes the arcsine function
func evaluateArcsin(x *big.Float) (*big.Float, error) {
	xFloat, _ := x.Float64()
	if xFloat < -1 || xFloat > 1 {
		return nil, fmt.Errorf("arcsin: domain error (argument must be in [-1, 1])")
	}
	result := math.Asin(xFloat)
	return big.NewFloat(result), nil
}

// evaluateArccos computes the arccosine function
func evaluateArccos(x *big.Float) (*big.Float, error) {
	xFloat, _ := x.Float64()
	if xFloat < -1 || xFloat > 1 {
		return nil, fmt.Errorf("arccos: domain error (argument must be in [-1, 1])")
	}
	result := math.Acos(xFloat)
	return big.NewFloat(result), nil
}

// evaluateArctan computes the arctangent function
func evaluateArctan(x *big.Float) (*big.Float, error) {
	xFloat, _ := x.Float64()
	result := math.Atan(xFloat)
	return big.NewFloat(result), nil
}

// evaluateSinh computes the hyperbolic sine function
func evaluateSinh(x *big.Float) (*big.Float, error) {
	xFloat, _ := x.Float64()
	result := math.Sinh(xFloat)
	return big.NewFloat(result), nil
}

// evaluateCosh computes the hyperbolic cosine function
func evaluateCosh(x *big.Float) (*big.Float, error) {
	xFloat, _ := x.Float64()
	result := math.Cosh(xFloat)
	return big.NewFloat(result), nil
}

// evaluateTanh computes the hyperbolic tangent function
func evaluateTanh(x *big.Float) (*big.Float, error) {
	xFloat, _ := x.Float64()
	result := math.Tanh(xFloat)
	return big.NewFloat(result), nil
}