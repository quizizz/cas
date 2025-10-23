// Package calculus implements symbolic differentiation and integration.
package calculus

import (
	"fmt"

	"github.com/quizizz/cas/pkg/ast"
	"github.com/quizizz/cas/pkg/simplify"
)

// Derivative computes the symbolic derivative of an expression with respect to a variable
func Derivative(expr ast.Expr, variable string) (ast.Expr, error) {
	return differentiate(expr, variable)
}

// differentiate implements the differentiation rules
func differentiate(expr ast.Expr, variable string) (ast.Expr, error) {
	switch e := expr.(type) {
	case *ast.Int:
		// d/dx(c) = 0
		return ast.NewInt(0), nil

	case *ast.Float:
		// d/dx(c) = 0
		return ast.NewInt(0), nil

	case *ast.Rational:
		// d/dx(c) = 0
		return ast.NewInt(0), nil

	case *ast.Const:
		// d/dx(π) = 0, d/dx(e) = 0
		return ast.NewInt(0), nil

	case *ast.Var:
		// d/dx(x) = 1, d/dx(y) = 0
		if e.Name() == variable {
			return ast.NewInt(1), nil
		}
		return ast.NewInt(0), nil

	case *ast.Add:
		// d/dx(f + g) = d/dx(f) + d/dx(g)
		return differentiateAdd(e, variable)

	case *ast.Mul:
		// d/dx(f * g) = f' * g + f * g' (product rule)
		return differentiateMul(e, variable)

	case *ast.Pow:
		// d/dx(f^g) = f^g * (g' * ln(f) + g * f'/f) (general power rule)
		// d/dx(f^n) = n * f^(n-1) * f' (power rule for constant exponent)
		return differentiatePow(e, variable)

	case *ast.Func:
		// d/dx(f(g)) = f'(g) * g' (chain rule)
		return differentiateFunc(e, variable)

	default:
		return nil, fmt.Errorf("cannot differentiate expression of type %T", expr)
	}
}

// differentiateAdd handles addition (sum rule)
func differentiateAdd(add *ast.Add, variable string) (ast.Expr, error) {
	terms := add.Terms()
	derivativeTerms := make([]ast.Expr, len(terms))

	for i, term := range terms {
		derivative, err := differentiate(term, variable)
		if err != nil {
			return nil, err
		}
		derivativeTerms[i] = derivative
	}

	result := ast.NewAdd(derivativeTerms...)
	return simplify.Collect(result), nil
}

// differentiateMul handles multiplication (product rule)
func differentiateMul(mul *ast.Mul, variable string) (ast.Expr, error) {
	factors := mul.Terms()

	if len(factors) == 0 {
		return ast.NewInt(0), nil
	}

	if len(factors) == 1 {
		return differentiate(factors[0], variable)
	}

	if len(factors) == 2 {
		// Binary product rule: (fg)' = f'g + fg'
		f := factors[0]
		g := factors[1]

		fPrime, err := differentiate(f, variable)
		if err != nil {
			return nil, err
		}

		gPrime, err := differentiate(g, variable)
		if err != nil {
			return nil, err
		}

		term1 := ast.NewMul(fPrime, g)
		term2 := ast.NewMul(f, gPrime)
		result := ast.NewAdd(term1, term2)

		return simplify.Collect(result), nil
	}

	// Generalized product rule for multiple factors
	// (f₁f₂...fₙ)' = Σᵢ(f₁...fᵢ₋₁·fᵢ'·fᵢ₊₁...fₙ)
	var terms []ast.Expr

	for i, factor := range factors {
		derivative, err := differentiate(factor, variable)
		if err != nil {
			return nil, err
		}

		// Create product of all factors except the i-th, with i-th replaced by its derivative
		var termFactors []ast.Expr
		for j, f := range factors {
			if i == j {
				termFactors = append(termFactors, derivative)
			} else {
				termFactors = append(termFactors, f)
			}
		}

		terms = append(terms, ast.NewMul(termFactors...))
	}

	result := ast.NewAdd(terms...)
	return simplify.Collect(result), nil
}

// differentiatePow handles exponentiation (power rule and general power rule)
func differentiatePow(pow *ast.Pow, variable string) (ast.Expr, error) {
	base := pow.Base()
	exponent := pow.Exponent()

	// Check if exponent is constant
	if !containsVariable(exponent, variable) {
		// Power rule: d/dx(f^n) = n * f^(n-1) * f'
		basePrime, err := differentiate(base, variable)
		if err != nil {
			return nil, err
		}

		// n * f^(n-1)
		newExponent := ast.NewAdd(exponent, ast.NewInt(-1))
		powerTerm := ast.NewPow(base, newExponent)

		result := ast.NewMul(exponent, powerTerm, basePrime)
		return simplify.Collect(result), nil
	}

	// Check if base is constant
	if !containsVariable(base, variable) {
		// d/dx(a^f) = a^f * ln(a) * f'
		exponentPrime, err := differentiate(exponent, variable)
		if err != nil {
			return nil, err
		}

		ln := ast.NewFunc("ln", base)
		result := ast.NewMul(pow, ln, exponentPrime)
		return simplify.Collect(result), nil
	}

	// General case: d/dx(f^g) = f^g * (g' * ln(f) + g * f'/f)
	basePrime, err := differentiate(base, variable)
	if err != nil {
		return nil, err
	}

	exponentPrime, err := differentiate(exponent, variable)
	if err != nil {
		return nil, err
	}

	// g' * ln(f)
	ln := ast.NewFunc("ln", base)
	term1 := ast.NewMul(exponentPrime, ln)

	// g * f'/f
	ratio := ast.NewMul(basePrime, ast.NewPow(base, ast.NewInt(-1)))
	term2 := ast.NewMul(exponent, ratio)

	// f^g * (...)
	innerDerivative := ast.NewAdd(term1, term2)
	result := ast.NewMul(pow, innerDerivative)

	return simplify.Collect(result), nil
}

// differentiateFunc handles function derivatives (chain rule)
func differentiateFunc(fn *ast.Func, variable string) (ast.Expr, error) {
	args := fn.Args()
	if len(args) != 1 {
		return nil, fmt.Errorf("differentiation of multi-argument functions not yet supported")
	}

	arg := args[0]
	argPrime, err := differentiate(arg, variable)
	if err != nil {
		return nil, err
	}

	// Get the derivative of the outer function
	outerDerivative, err := getFunctionDerivative(fn.Name(), arg)
	if err != nil {
		return nil, err
	}

	// Chain rule: (f(g))' = f'(g) * g'
	result := ast.NewMul(outerDerivative, argPrime)
	return simplify.Collect(result), nil
}

// getFunctionDerivative returns the derivative of standard mathematical functions
func getFunctionDerivative(funcName string, arg ast.Expr) (ast.Expr, error) {
	switch funcName {
	case "sin":
		// d/dx(sin(u)) = cos(u)
		return ast.NewFunc("cos", arg), nil

	case "cos":
		// d/dx(cos(u)) = -sin(u)
		sin := ast.NewFunc("sin", arg)
		return ast.NewMul(ast.NewInt(-1), sin), nil

	case "tan":
		// d/dx(tan(u)) = sec²(u) = 1/cos²(u)
		cos := ast.NewFunc("cos", arg)
		cosSquared := ast.NewPow(cos, ast.NewInt(2))
		return ast.NewPow(cosSquared, ast.NewInt(-1)), nil

	case "sec":
		// d/dx(sec(u)) = sec(u)tan(u)
		sec := ast.NewFunc("sec", arg)
		tan := ast.NewFunc("tan", arg)
		return ast.NewMul(sec, tan), nil

	case "csc":
		// d/dx(csc(u)) = -csc(u)cot(u)
		csc := ast.NewFunc("csc", arg)
		cot := ast.NewFunc("cot", arg)
		return ast.NewMul(ast.NewInt(-1), csc, cot), nil

	case "cot":
		// d/dx(cot(u)) = -csc²(u)
		csc := ast.NewFunc("csc", arg)
		cscSquared := ast.NewPow(csc, ast.NewInt(2))
		return ast.NewMul(ast.NewInt(-1), cscSquared), nil

	case "arcsin":
		// d/dx(arcsin(u)) = 1/√(1-u²)
		uSquared := ast.NewPow(arg, ast.NewInt(2))
		oneMinusUSquared := ast.NewAdd(ast.NewInt(1), ast.NewMul(ast.NewInt(-1), uSquared))
		sqrt := ast.NewFunc("sqrt", oneMinusUSquared)
		return ast.NewPow(sqrt, ast.NewInt(-1)), nil

	case "arccos":
		// d/dx(arccos(u)) = -1/√(1-u²)
		uSquared := ast.NewPow(arg, ast.NewInt(2))
		oneMinusUSquared := ast.NewAdd(ast.NewInt(1), ast.NewMul(ast.NewInt(-1), uSquared))
		sqrt := ast.NewFunc("sqrt", oneMinusUSquared)
		return ast.NewMul(ast.NewInt(-1), ast.NewPow(sqrt, ast.NewInt(-1))), nil

	case "arctan":
		// d/dx(arctan(u)) = 1/(1+u²)
		uSquared := ast.NewPow(arg, ast.NewInt(2))
		onePlusUSquared := ast.NewAdd(ast.NewInt(1), uSquared)
		return ast.NewPow(onePlusUSquared, ast.NewInt(-1)), nil

	case "sinh":
		// d/dx(sinh(u)) = cosh(u)
		return ast.NewFunc("cosh", arg), nil

	case "cosh":
		// d/dx(cosh(u)) = sinh(u)
		return ast.NewFunc("sinh", arg), nil

	case "tanh":
		// d/dx(tanh(u)) = sech²(u) = 1/cosh²(u)
		cosh := ast.NewFunc("cosh", arg)
		coshSquared := ast.NewPow(cosh, ast.NewInt(2))
		return ast.NewPow(coshSquared, ast.NewInt(-1)), nil

	case "ln":
		// d/dx(ln(u)) = 1/u
		return ast.NewPow(arg, ast.NewInt(-1)), nil

	case "log":
		// d/dx(log(u)) = 1/(u * ln(10))
		ln10 := ast.NewFunc("ln", ast.NewInt(10))
		denominator := ast.NewMul(arg, ln10)
		return ast.NewPow(denominator, ast.NewInt(-1)), nil

	case "exp":
		// d/dx(e^u) = e^u
		return ast.NewFunc("exp", arg), nil

	case "sqrt":
		// d/dx(√u) = 1/(2√u) = (1/2) * u^(-1/2)
		half := ast.NewRational(1, 2)
		exponent := ast.NewMul(ast.NewInt(-1), half)
		return ast.NewMul(half, ast.NewPow(arg, exponent)), nil

	case "abs":
		// d/dx(|u|) = u/|u| (for u ≠ 0)
		abs := ast.NewFunc("abs", arg)
		return ast.NewMul(arg, ast.NewPow(abs, ast.NewInt(-1))), nil

	default:
		return nil, fmt.Errorf("derivative of function %s not implemented", funcName)
	}
}

// containsVariable checks if an expression contains a specific variable
func containsVariable(expr ast.Expr, variable string) bool {
	variables := expr.Variables()
	for _, v := range variables {
		if v == variable {
			return true
		}
	}
	return false
}

// PartialDerivative computes partial derivatives for multivariable expressions
func PartialDerivative(expr ast.Expr, variable string) (ast.Expr, error) {
	return Derivative(expr, variable)
}

// NthDerivative computes the nth derivative of an expression
func NthDerivative(expr ast.Expr, variable string, n int) (ast.Expr, error) {
	if n < 0 {
		return nil, fmt.Errorf("derivative order must be non-negative")
	}

	if n == 0 {
		return expr, nil
	}

	current := expr
	for i := 0; i < n; i++ {
		derivative, err := Derivative(current, variable)
		if err != nil {
			return nil, err
		}
		current = derivative
	}

	return current, nil
}

// Gradient computes the gradient (vector of partial derivatives) for a multivariable expression
func Gradient(expr ast.Expr, variables []string) (map[string]ast.Expr, error) {
	gradient := make(map[string]ast.Expr)

	for _, variable := range variables {
		partial, err := PartialDerivative(expr, variable)
		if err != nil {
			return nil, fmt.Errorf("error computing partial derivative with respect to %s: %v", variable, err)
		}
		gradient[variable] = partial
	}

	return gradient, nil
}
