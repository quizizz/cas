// Package expand implements polynomial and expression expansion algorithms.
package expand

import (
	"math/big"

	"github.com/quizizz/cas/pkg/ast"
)

// Options controls expansion behavior
type Options struct {
	// MaxDegree limits the maximum polynomial degree to expand
	MaxDegree int
	// ExpandLogs enables logarithm expansion (ln(xy) = ln(x) + ln(y))
	ExpandLogs bool
	// ExpandTrig enables trigonometric expansion
	ExpandTrig bool
}

// DefaultOptions returns the default expansion options
func DefaultOptions() Options {
	return Options{
		MaxDegree:  10,
		ExpandLogs: false,
		ExpandTrig: false,
	}
}

// Expand performs polynomial and algebraic expansion on expressions
func Expand(expr ast.Expr, opts ...Options) ast.Expr {
	options := DefaultOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	return expandExpr(expr, options)
}

// expandExpr recursively expands expressions
func expandExpr(expr ast.Expr, opts Options) ast.Expr {
	switch e := expr.(type) {
	case *ast.Add:
		return expandAdd(e, opts)
	case *ast.Mul:
		return expandMul(e, opts)
	case *ast.Pow:
		return expandPow(e, opts)
	case *ast.Func:
		return expandFunc(e, opts)
	default:
		return expr
	}
}

// expandAdd expands addition expressions
func expandAdd(add *ast.Add, opts Options) ast.Expr {
	terms := add.Terms()
	expandedTerms := make([]ast.Expr, len(terms))

	for i, term := range terms {
		expandedTerms[i] = expandExpr(term, opts)
	}

	return ast.NewAdd(expandedTerms...)
}

// expandMul expands multiplication expressions (distributive property)
func expandMul(mul *ast.Mul, opts Options) ast.Expr {
	factors := mul.Terms()

	// First, expand each factor recursively
	expandedFactors := make([]ast.Expr, len(factors))
	for i, factor := range factors {
		expandedFactors[i] = expandExpr(factor, opts)
	}

	// Find additive factors to distribute
	var addFactors []*ast.Add
	var otherFactors []ast.Expr

	for _, factor := range expandedFactors {
		if add, ok := factor.(*ast.Add); ok {
			addFactors = append(addFactors, add)
		} else {
			otherFactors = append(otherFactors, factor)
		}
	}

	// If no additive factors, return the multiplication
	if len(addFactors) == 0 {
		if len(otherFactors) == 1 {
			return otherFactors[0]
		}
		return ast.NewMul(otherFactors...)
	}

	// Distribute multiplication over addition
	return distributeMultiplication(addFactors, otherFactors, opts)
}

// expandPow expands power expressions
func expandPow(pow *ast.Pow, opts Options) ast.Expr {
	base := expandExpr(pow.Base(), opts)
	exp := expandExpr(pow.Exponent(), opts)

	// Handle (ab)^c = a^c * b^c for multiplication base
	if mul, ok := base.(*ast.Mul); ok {
		factors := mul.Terms()
		expandedFactors := make([]ast.Expr, len(factors))
		for i, factor := range factors {
			expandedFactors[i] = ast.NewPow(factor, exp)
		}
		return expandExpr(ast.NewMul(expandedFactors...), opts)
	}

	// Handle (a+b)^n expansion for positive integer exponents
	if add, ok := base.(*ast.Add); ok {
		if intExp, ok := exp.(*ast.Int); ok {
			return expandPolynomial(add, intExp, opts)
		}
	}

	return ast.NewPow(base, exp)
}

// expandFunc expands function expressions
func expandFunc(fn *ast.Func, opts Options) ast.Expr {
	args := fn.Args()
	expandedArgs := make([]ast.Expr, len(args))

	for i, arg := range args {
		expandedArgs[i] = expandExpr(arg, opts)
	}

	expandedFunc := ast.NewFunc(fn.Name(), expandedArgs...)

	// Apply specific function expansion rules
	switch fn.Name() {
	case "ln", "log":
		if opts.ExpandLogs {
			return expandLogarithm(expandedFunc)
		}
	case "tan":
		if opts.ExpandTrig {
			return expandTangent(expandedFunc)
		}
	case "sec":
		if opts.ExpandTrig {
			return expandSecant(expandedFunc)
		}
	case "csc":
		if opts.ExpandTrig {
			return expandCosecant(expandedFunc)
		}
	case "cot":
		if opts.ExpandTrig {
			return expandCotangent(expandedFunc)
		}
	}

	return expandedFunc
}

// distributeMultiplication distributes multiplication over addition using FOIL and distributive property
func distributeMultiplication(addFactors []*ast.Add, otherFactors []ast.Expr, opts Options) ast.Expr {
	if len(addFactors) == 0 {
		return ast.NewMul(otherFactors...)
	}

	// Start with the first additive factor
	result := addFactors[0].Terms()

	// Multiply with each subsequent additive factor
	for i := 1; i < len(addFactors); i++ {
		result = multiplyTermSets(result, addFactors[i].Terms(), opts)
	}

	// Multiply with other factors
	if len(otherFactors) > 0 {
		otherFactor := ast.NewMul(otherFactors...)
		for j, term := range result {
			result[j] = expandExpr(ast.NewMul(term, otherFactor), opts)
		}
	}

	return ast.NewAdd(result...)
}

// multiplyTermSets multiplies two sets of terms (distributive property)
func multiplyTermSets(set1, set2 []ast.Expr, opts Options) []ast.Expr {
	var result []ast.Expr

	for _, term1 := range set1 {
		for _, term2 := range set2 {
			product := expandExpr(ast.NewMul(term1, term2), opts)
			result = append(result, product)
		}
	}

	return result
}

// expandPolynomial expands (a+b+...)^n using binomial/multinomial theorem
func expandPolynomial(add *ast.Add, exp *ast.Int, opts Options) ast.Expr {
	// Get the exponent value
	expVal, err := exp.Eval(make(map[string]*big.Float))
	if err != nil {
		return ast.NewPow(add, exp)
	}

	expInt, _ := expVal.Int64()

	// Limit expansion to prevent excessive computation
	if expInt < 0 || expInt > int64(opts.MaxDegree) {
		return ast.NewPow(add, exp)
	}

	if expInt == 0 {
		return ast.NewInt(1)
	}

	if expInt == 1 {
		return add
	}

	// For small exponents, use direct multiplication
	if expInt <= 4 {
		return expandByRepeatedMultiplication(add, int(expInt), opts)
	}

	// For larger exponents, we could implement binomial expansion
	// For now, fall back to the original expression
	return ast.NewPow(add, exp)
}

// expandByRepeatedMultiplication expands (a+b)^n by repeated multiplication
func expandByRepeatedMultiplication(add *ast.Add, n int, opts Options) ast.Expr {
	result := add

	for i := 1; i < n; i++ {
		// Multiply result by add
		mul := ast.NewMul(result, add)
		result = expandExpr(mul, opts).(*ast.Add)
	}

	return result
}

// Logarithm expansion functions
func expandLogarithm(fn *ast.Func) ast.Expr {
	args := fn.Args()
	if len(args) != 1 {
		return fn
	}

	arg := args[0]

	// ln(xy) = ln(x) + ln(y)
	if mul, ok := arg.(*ast.Mul); ok {
		factors := mul.Terms()
		logTerms := make([]ast.Expr, len(factors))
		for i, factor := range factors {
			logTerms[i] = ast.NewFunc(fn.Name(), factor)
		}
		return ast.NewAdd(logTerms...)
	}

	// ln(x^y) = y * ln(x)
	if pow, ok := arg.(*ast.Pow); ok {
		base := pow.Base()
		exp := pow.Exponent()
		return ast.NewMul(exp, ast.NewFunc(fn.Name(), base))
	}

	return fn
}

// Trigonometric expansion functions
func expandTangent(fn *ast.Func) ast.Expr {
	args := fn.Args()
	if len(args) != 1 {
		return fn
	}

	// tan(x) = sin(x) / cos(x)
	arg := args[0]
	sin := ast.NewFunc("sin", arg)
	cos := ast.NewFunc("cos", arg)

	// Create division as multiplication by reciprocal
	reciprocal := ast.NewPow(cos, ast.NewInt(-1))
	return ast.NewMul(sin, reciprocal)
}

func expandSecant(fn *ast.Func) ast.Expr {
	args := fn.Args()
	if len(args) != 1 {
		return fn
	}

	// sec(x) = 1 / cos(x)
	arg := args[0]
	cos := ast.NewFunc("cos", arg)
	reciprocal := ast.NewPow(cos, ast.NewInt(-1))
	return reciprocal
}

func expandCosecant(fn *ast.Func) ast.Expr {
	args := fn.Args()
	if len(args) != 1 {
		return fn
	}

	// csc(x) = 1 / sin(x)
	arg := args[0]
	sin := ast.NewFunc("sin", arg)
	reciprocal := ast.NewPow(sin, ast.NewInt(-1))
	return reciprocal
}

func expandCotangent(fn *ast.Func) ast.Expr {
	args := fn.Args()
	if len(args) != 1 {
		return fn
	}

	// cot(x) = cos(x) / sin(x)
	arg := args[0]
	cos := ast.NewFunc("cos", arg)
	sin := ast.NewFunc("sin", arg)

	// Create division as multiplication by reciprocal
	reciprocal := ast.NewPow(sin, ast.NewInt(-1))
	return ast.NewMul(cos, reciprocal)
}

// ExpandFully performs complete expansion including nested expressions
func ExpandFully(expr ast.Expr, opts ...Options) ast.Expr {
	options := DefaultOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	current := expr
	maxIterations := 5 // Prevent infinite loops

	for i := 0; i < maxIterations; i++ {
		expanded := Expand(current, options)
		if expanded.String() == current.String() {
			break // No more changes
		}
		current = expanded
	}

	return current
}
