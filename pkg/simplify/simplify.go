// Package simplify implements advanced algebraic simplification algorithms.
package simplify

import (
	"math/big"

	"github.com/quizizz/cas/pkg/ast"
)

// Options controls simplification behavior
type Options struct {
	// Once indicates whether to perform only one simplification pass
	Once bool
	// KeepNegative prevents factoring out common negative factors
	KeepNegative bool
	// MaxIterations limits the number of simplification iterations
	MaxIterations int
}

// DefaultOptions returns the default simplification options
func DefaultOptions() Options {
	return Options{
		Once:          false,
		KeepNegative:  false,
		MaxIterations: 10,
	}
}

// Simplify performs advanced simplification on an expression
func Simplify(expr ast.Expr, opts ...Options) ast.Expr {
	options := DefaultOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	current := expr
	iteration := 0

	for iteration < options.MaxIterations {
		// Factor and collect
		step1 := Factor(current, options)
		step2 := Collect(step1, options)

		// Rollback if collect didn't do anything
		if Equal(step1, step2) {
			step2 = step1
		}

		// Expand if we're stuck
		if Equal(current, step2) {
			step3 := Expand(step2, options)
			if !Equal(step2, step3) {
				step2 = Collect(step3, options)
			}
		}

		// Stop if no change or if doing only one iteration
		if Equal(current, step2) || options.Once {
			return step2
		}

		current = step2
		iteration++
	}

	return current
}

// Collect combines like terms and simplifies expressions
func Collect(expr ast.Expr, opts ...Options) ast.Expr {
	options := DefaultOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	switch e := expr.(type) {
	case *ast.Add:
		return collectAdd(e, options)
	case *ast.Mul:
		return collectMul(e, options)
	case *ast.Pow:
		return collectPow(e, options)
	case ast.Numeric:
		return collectNumeric(e)
	default:
		return expr
	}
}

// Factor extracts common factors from expressions
func Factor(expr ast.Expr, opts ...Options) ast.Expr {
	options := DefaultOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	switch e := expr.(type) {
	case *ast.Add:
		return factorAdd(e, options)
	case *ast.Mul:
		return factorMul(e, options)
	default:
		return expr
	}
}

// Expand expands products and powers
func Expand(expr ast.Expr, opts ...Options) ast.Expr {
	options := DefaultOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	switch e := expr.(type) {
	case *ast.Mul:
		return expandMul(e, options)
	case *ast.Pow:
		return expandPow(e, options)
	default:
		return expr
	}
}

// Equal checks if two expressions are structurally equal after normalization
func Equal(a, b ast.Expr) bool {
	normA := Normalize(a)
	normB := Normalize(b)
	return normA.String() == normB.String()
}

// Normalize puts an expression in canonical form
func Normalize(expr ast.Expr) ast.Expr {
	switch e := expr.(type) {
	case *ast.Add:
		return normalizeAdd(e)
	case *ast.Mul:
		return normalizeMul(e)
	case *ast.Pow:
		return normalizePow(e)
	default:
		return expr
	}
}

// collectAdd combines like terms in an addition
func collectAdd(add *ast.Add, opts Options) ast.Expr {
	// First, flatten nested additions recursively
	addExpr := ast.Expr(add)
	allTerms := flattenAdd(&addExpr)

	// Recursively collect each term
	for i, term := range allTerms {
		allTerms[i] = Collect(term, opts)
	}

	// Group terms by their "like" structure
	termGroups := make(map[string][]ast.Expr)
	coefficients := make(map[string]*big.Float)

	for _, term := range allTerms {
		baseForm, coeff := extractCoefficientAndBase(term)
		key := baseForm.String()

		termGroups[key] = append(termGroups[key], term)
		if coefficients[key] == nil {
			coefficients[key] = new(big.Float)
		}
		coefficients[key].Add(coefficients[key], coeff)
	}

	// Reconstruct the result
	var result []ast.Expr

	for key, termGroup := range termGroups {
		coeff := coefficients[key]
		baseForm := termGroup[0]
		if len(termGroup) > 1 || coeff.Cmp(big.NewFloat(1)) != 0 {
			// Extract base form from first term
			baseForm, _ = extractCoefficientAndBase(termGroup[0])
		}

		// Skip zero terms
		if coeff.Sign() == 0 {
			continue
		}

		// Create the simplified term
		if coeff.Cmp(big.NewFloat(1)) == 0 {
			// Coefficient is 1, just use the base
			if baseForm.String() == "1" {
				result = append(result, ast.NewInt(1))
			} else {
				result = append(result, baseForm)
			}
		} else if coeff.Cmp(big.NewFloat(-1)) == 0 {
			// Coefficient is -1
			if baseForm.String() == "1" {
				result = append(result, ast.NewInt(-1))
			} else {
				result = append(result, ast.NewMul(ast.NewInt(-1), baseForm))
			}
		} else {
			// General coefficient
			if baseForm.String() == "1" {
				// Pure numeric
				if coeff.IsInt() {
					intVal, _ := coeff.Int64()
					result = append(result, ast.NewInt(intVal))
				} else {
					floatExpr, _ := ast.NewFloatFromString(coeff.Text('g', -1))
					result = append(result, floatExpr)
				}
			} else {
				// Coefficient times base
				if coeff.IsInt() {
					intVal, _ := coeff.Int64()
					result = append(result, ast.NewMul(ast.NewInt(intVal), baseForm))
				} else {
					floatExpr, _ := ast.NewFloatFromString(coeff.Text('g', -1))
					result = append(result, ast.NewMul(floatExpr, baseForm))
				}
			}
		}
	}

	if len(result) == 0 {
		return ast.NewInt(0)
	} else if len(result) == 1 {
		return result[0]
	} else {
		return ast.NewAdd(result...)
	}
}

// extractCoefficientAndBase separates a term into coefficient and base parts
func extractCoefficientAndBase(term ast.Expr) (base ast.Expr, coeff *big.Float) {
	coeff = big.NewFloat(1)
	base = term

	switch t := term.(type) {
	case *ast.Int:
		val, _ := t.Eval(make(map[string]*big.Float))
		coeff = val
		base = ast.NewInt(1)
	case *ast.Float:
		val, _ := t.Eval(make(map[string]*big.Float))
		coeff = val
		base = ast.NewInt(1)
	case *ast.Mul:
		// Flatten multiplication and separate numeric/non-numeric parts
		factors := flattenMul(term)
		var nonNumericFactors []ast.Expr

		for _, factor := range factors {
			if isNumeric(factor) {
				val, err := factor.Eval(make(map[string]*big.Float))
				if err == nil {
					coeff.Mul(coeff, val)
				}
			} else {
				nonNumericFactors = append(nonNumericFactors, factor)
			}
		}

		if len(nonNumericFactors) == 0 {
			base = ast.NewInt(1)
		} else if len(nonNumericFactors) == 1 {
			base = nonNumericFactors[0]
		} else {
			// Sort factors for consistent representation
			base = normalizeMul(ast.NewMul(nonNumericFactors...))
		}
	default:
		// Variables, functions, powers, etc.
		base = term
		coeff = big.NewFloat(1)
	}

	return base, coeff
}

// flattenMul recursively flattens nested multiplication expressions
func flattenMul(mulExpr ast.Expr) []ast.Expr {
	switch e := mulExpr.(type) {
	case *ast.Mul:
		terms := e.Terms()
		var allTerms []ast.Expr
		for _, term := range terms {
			subTerms := flattenMul(term)
			allTerms = append(allTerms, subTerms...)
		}
		return allTerms
	default:
		return []ast.Expr{mulExpr}
	}
}

// collectMul simplifies multiplication expressions
func collectMul(mul *ast.Mul, opts Options) ast.Expr {
	// First, flatten nested multiplications
	allFactors := flattenMul(mul)

	// Collect each factor recursively
	for i, factor := range allFactors {
		allFactors[i] = Collect(factor, opts)
	}

	// Separate numeric and non-numeric parts
	numPart, exprPart := partitionMul(ast.NewMul(allFactors...))

	// Check for zero
	if isZero(numPart) {
		return ast.NewInt(0)
	}

	// Collect powers of same base
	powers := collectPowers(exprPart)

	// Reconstruct
	var factors []ast.Expr
	if !isOne(numPart) {
		factors = append(factors, numPart)
	}

	for _, powerInfo := range powers {
		exp := powerInfo.exponent
		base := powerInfo.base
		if isOne(exp) {
			factors = append(factors, base)
		} else if isZero(exp) {
			// a^0 = 1, so skip this factor
			continue
		} else {
			factors = append(factors, ast.NewPow(base, exp))
		}
	}

	if len(factors) == 0 {
		return ast.NewInt(1)
	} else if len(factors) == 1 {
		return factors[0]
	} else {
		return ast.NewMul(factors...)
	}
}

// collectPow simplifies power expressions
func collectPow(pow *ast.Pow, opts Options) ast.Expr {
	base := Collect(pow.Base(), opts)
	exp := Collect(pow.Exponent(), opts)

	// Handle (a^b)^c = a^(bc)
	if basePow, ok := base.(*ast.Pow); ok {
		newExp := Collect(ast.NewMul(basePow.Exponent(), exp), opts)
		return ast.NewPow(basePow.Base(), newExp)
	}

	// Handle a^0 = 1
	if isZero(exp) {
		return ast.NewInt(1)
	}

	// Handle a^1 = a
	if isOne(exp) {
		return base
	}

	return ast.NewPow(base, exp)
}

// collectNumeric simplifies numeric expressions
func collectNumeric(num ast.Numeric) ast.Expr {
	switch n := num.(type) {
	case *ast.Rational:
		// Reduce to lowest terms
		gcd := new(big.Int).GCD(nil, nil, n.Numerator(), n.Denominator())
		newNum := new(big.Int).Div(n.Numerator(), gcd)
		newDen := new(big.Int).Div(n.Denominator(), gcd)

		if newDen.Cmp(big.NewInt(1)) == 0 {
			return &ast.Int{}
		}
		return ast.NewRationalFromInts(newNum, newDen)
	default:
		return num.(ast.Expr)
	}
}

// Helper functions

func partitionMul(mul *ast.Mul) (numeric ast.Expr, others ast.Expr) {
	terms := mul.Terms()
	var numTerms, otherTerms []ast.Expr

	for _, term := range terms {
		if isNumeric(term) {
			numTerms = append(numTerms, term)
		} else {
			otherTerms = append(otherTerms, term)
		}
	}

	// Multiply all numeric terms
	if len(numTerms) == 0 {
		numeric = ast.NewInt(1)
	} else if len(numTerms) == 1 {
		numeric = numTerms[0]
	} else {
		result := numTerms[0]
		vars := make(map[string]*big.Float)
		for i := 1; i < len(numTerms); i++ {
			val1, _ := result.Eval(vars)
			val2, _ := numTerms[i].Eval(vars)
			product := new(big.Float).Mul(val1, val2)
			result, _ = ast.NewFloatFromString(product.Text('g', -1))
		}
		numeric = result
	}

	// Combine other terms
	if len(otherTerms) == 0 {
		others = ast.NewInt(1)
	} else if len(otherTerms) == 1 {
		others = otherTerms[0]
	} else {
		others = ast.NewMul(otherTerms...)
	}

	return
}

func collectPowers(expr ast.Expr) map[string]powerInfo {
	powers := make(map[string]powerInfo)

	switch e := expr.(type) {
	case *ast.Mul:
		for _, term := range e.Terms() {
			if pow, ok := term.(*ast.Pow); ok {
				base := pow.Base()
				exp := pow.Exponent()
				key := base.String()

				if existing, ok := powers[key]; ok {
					// Add exponents: a^m * a^n = a^(m+n)
					newExp := Collect(ast.NewAdd(existing.exponent, exp))
					powers[key] = powerInfo{base: base, exponent: newExp}
				} else {
					powers[key] = powerInfo{base: base, exponent: exp}
				}
			} else {
				// Term without explicit exponent has exponent 1
				key := term.String()
				if existing, ok := powers[key]; ok {
					// Add 1 to existing exponent
					newExp := Collect(ast.NewAdd(existing.exponent, ast.NewInt(1)))
					powers[key] = powerInfo{base: term, exponent: newExp}
				} else {
					powers[key] = powerInfo{base: term, exponent: ast.NewInt(1)}
				}
			}
		}
	case *ast.Pow:
		key := e.Base().String()
		powers[key] = powerInfo{base: e.Base(), exponent: e.Exponent()}
	default:
		if !isOne(expr) {
			key := expr.String()
			powers[key] = powerInfo{base: expr, exponent: ast.NewInt(1)}
		}
	}

	return powers
}

type powerInfo struct {
	base     ast.Expr
	exponent ast.Expr
}

func isNumeric(expr ast.Expr) bool {
	_, ok := expr.(ast.Numeric)
	return ok
}

func isZero(expr ast.Expr) bool {
	if num, ok := expr.(ast.Numeric); ok {
		val, err := num.Eval(make(map[string]*big.Float))
		if err != nil {
			return false
		}
		return val.Sign() == 0
	}
	return false
}

func isOne(expr ast.Expr) bool {
	if num, ok := expr.(ast.Numeric); ok {
		val, err := num.Eval(make(map[string]*big.Float))
		if err != nil {
			return false
		}
		return val.Cmp(big.NewFloat(1)) == 0
	}
	return false
}

// factorAdd extracts common factors from addition terms
func factorAdd(add *ast.Add, opts Options) ast.Expr {
	terms := add.Terms()
	if len(terms) < 2 {
		return add
	}

	// Look for common multiplicative factors
	commonFactor := findCommonFactor(terms)
	if commonFactor != nil && !isOne(commonFactor) {
		// Factor out the common factor
		var factorizedTerms []ast.Expr
		for _, term := range terms {
			if mul, ok := term.(*ast.Mul); ok {
				// Remove the common factor from this term
				remaining := removeFactor(mul, commonFactor)
				factorizedTerms = append(factorizedTerms, remaining)
			} else if termsEqual(term, commonFactor) {
				factorizedTerms = append(factorizedTerms, ast.NewInt(1))
			} else {
				// Check if term contains the factor
				if containsFactor(term, commonFactor) {
					remaining := ast.NewMul(term, ast.NewPow(commonFactor, ast.NewInt(-1)))
					factorizedTerms = append(factorizedTerms, Collect(remaining, opts))
				} else {
					factorizedTerms = append(factorizedTerms, term)
				}
			}
		}

		factorizedSum := ast.NewAdd(factorizedTerms...)
		return ast.NewMul(commonFactor, factorizedSum)
	}

	return add
}

func factorMul(mul *ast.Mul, opts Options) ast.Expr {
	// For multiplication, we can factor by collecting like terms
	return Collect(mul, opts)
}

func expandMul(mul *ast.Mul, opts Options) ast.Expr {
	terms := mul.Terms()
	result := terms[0]

	// Multiply each subsequent term using distributive property
	for i := 1; i < len(terms); i++ {
		result = distributiveMultiply(result, terms[i], opts)
	}

	return result
}

func expandPow(pow *ast.Pow, opts Options) ast.Expr {
	base := pow.Base()
	exp := pow.Exponent()

	// Handle (a+b)^n expansion for small integer exponents
	if add, ok := base.(*ast.Add); ok {
		if intExp, ok := exp.(*ast.Int); ok {
			val, err := intExp.Eval(make(map[string]*big.Float))
			if err == nil {
				expInt, _ := val.Int64()
				if expInt >= 0 && expInt <= 4 {
					return expandBinomial(add, int(expInt), opts)
				}
			}
		}
	}

	return pow
}

func normalizeAdd(add *ast.Add) ast.Expr {
	terms := add.Terms()

	// Sort terms by their string representation for consistent ordering
	sortedTerms := make([]ast.Expr, len(terms))
	copy(sortedTerms, terms)

	// Simple bubble sort by string representation
	for i := 0; i < len(sortedTerms)-1; i++ {
		for j := 0; j < len(sortedTerms)-i-1; j++ {
			if sortedTerms[j].String() > sortedTerms[j+1].String() {
				sortedTerms[j], sortedTerms[j+1] = sortedTerms[j+1], sortedTerms[j]
			}
		}
	}

	return ast.NewAdd(sortedTerms...)
}

func normalizeMul(mul *ast.Mul) ast.Expr {
	terms := mul.Terms()

	// Sort factors by their string representation
	sortedTerms := make([]ast.Expr, len(terms))
	copy(sortedTerms, terms)

	// Simple bubble sort by string representation
	for i := 0; i < len(sortedTerms)-1; i++ {
		for j := 0; j < len(sortedTerms)-i-1; j++ {
			if sortedTerms[j].String() > sortedTerms[j+1].String() {
				sortedTerms[j], sortedTerms[j+1] = sortedTerms[j+1], sortedTerms[j]
			}
		}
	}

	return ast.NewMul(sortedTerms...)
}

func normalizePow(pow *ast.Pow) ast.Expr {
	// Power expressions are already in canonical form
	return pow
}

// Helper functions for advanced simplification

// flattenAdd recursively flattens nested addition expressions
func flattenAdd(add *ast.Expr) []ast.Expr {
	switch e := (*add).(type) {
	case *ast.Add:
		terms := e.Terms()
		var allTerms []ast.Expr
		for _, term := range terms {
			subTerms := flattenAdd(&term)
			allTerms = append(allTerms, subTerms...)
		}
		return allTerms
	default:
		return []ast.Expr{*add}
	}
}

func findCommonFactor(terms []ast.Expr) ast.Expr {
	if len(terms) == 0 {
		return nil
	}

	// Start with factors of the first term
	firstFactors := extractFactors(terms[0])
	if len(firstFactors) == 0 {
		return nil
	}

	// Find factors that appear in all terms
	var commonFactors []ast.Expr
	for _, factor := range firstFactors {
		isCommon := true
		for i := 1; i < len(terms); i++ {
			if !containsFactor(terms[i], factor) {
				isCommon = false
				break
			}
		}
		if isCommon {
			commonFactors = append(commonFactors, factor)
		}
	}

	if len(commonFactors) == 0 {
		return nil
	} else if len(commonFactors) == 1 {
		return commonFactors[0]
	} else {
		return ast.NewMul(commonFactors...)
	}
}

func extractFactors(expr ast.Expr) []ast.Expr {
	switch e := expr.(type) {
	case *ast.Mul:
		return e.Terms()
	case ast.Numeric:
		return []ast.Expr{expr}
	case *ast.Var:
		return []ast.Expr{expr}
	case *ast.Pow:
		return []ast.Expr{expr}
	default:
		return []ast.Expr{expr}
	}
}

func containsFactor(expr, factor ast.Expr) bool {
	factors := extractFactors(expr)
	for _, f := range factors {
		if termsEqual(f, factor) {
			return true
		}
	}
	return false
}

func removeFactor(mul *ast.Mul, factor ast.Expr) ast.Expr {
	terms := mul.Terms()
	var remaining []ast.Expr
	factorRemoved := false

	for _, term := range terms {
		if !factorRemoved && termsEqual(term, factor) {
			factorRemoved = true
			continue
		}
		remaining = append(remaining, term)
	}

	if len(remaining) == 0 {
		return ast.NewInt(1)
	} else if len(remaining) == 1 {
		return remaining[0]
	} else {
		return ast.NewMul(remaining...)
	}
}

func termsEqual(a, b ast.Expr) bool {
	return a.String() == b.String()
}

func distributiveMultiply(a, b ast.Expr, opts Options) ast.Expr {
	// Handle multiplication with addition using distributive property
	if addA, ok := a.(*ast.Add); ok {
		// (a1 + a2 + ...) * b = a1*b + a2*b + ...
		terms := addA.Terms()
		var products []ast.Expr
		for _, term := range terms {
			product := ast.NewMul(term, b)
			products = append(products, Expand(product, opts))
		}
		return ast.NewAdd(products...)
	}

	if addB, ok := b.(*ast.Add); ok {
		// a * (b1 + b2 + ...) = a*b1 + a*b2 + ...
		terms := addB.Terms()
		var products []ast.Expr
		for _, term := range terms {
			product := ast.NewMul(a, term)
			products = append(products, Expand(product, opts))
		}
		return ast.NewAdd(products...)
	}

	return ast.NewMul(a, b)
}

func expandBinomial(add *ast.Add, n int, opts Options) ast.Expr {
	if n == 0 {
		return ast.NewInt(1)
	}
	if n == 1 {
		return add
	}

	terms := add.Terms()
	if len(terms) == 2 {
		// Use binomial theorem for (a+b)^n
		a, b := terms[0], terms[1]
		var expandedTerms []ast.Expr

		for k := 0; k <= n; k++ {
			// C(n,k) * a^(n-k) * b^k
			coeff := binomialCoeff(n, k)
			aPower := ast.NewPow(a, ast.NewInt(int64(n-k)))
			bPower := ast.NewPow(b, ast.NewInt(int64(k)))

			term := ast.NewMul(ast.NewInt(coeff), aPower, bPower)
			expandedTerms = append(expandedTerms, term)
		}

		return ast.NewAdd(expandedTerms...)
	}

	// For more than 2 terms, use repeated multiplication
	result := add
	for i := 1; i < n; i++ {
		result = distributiveMultiply(result, add, opts).(*ast.Add)
	}
	return result
}

func binomialCoeff(n, k int) int64 {
	if k > n || k < 0 {
		return 0
	}
	if k == 0 || k == n {
		return 1
	}

	// Use the identity C(n,k) = C(n,n-k) to minimize computation
	if k > n-k {
		k = n - k
	}

	result := int64(1)
	for i := 0; i < k; i++ {
		result = result * int64(n-i) / int64(i+1)
	}
	return result
}
