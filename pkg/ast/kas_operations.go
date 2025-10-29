package ast

import (
	"math/big"
)

// Collect implements the Node.js KAS collect() function
// Groups like terms and combines their coefficients
func (a *Add) Collect() Expr {
	// [Expr expr, Numeric coefficient] pairs
	type termPair struct {
		expr  Expr
		coeff Numeric
	}

	var pairs []termPair

	// Process each term to extract coefficient and expression parts
	for _, term := range a.terms {
		if mul, isMul := term.(*Mul); isMul {
			// Partition multiplication into numeric and non-numeric parts
			numPart, exprPart := mul.partition()
			pairs = append(pairs, termPair{expr: exprPart, coeff: numPart})
		} else if num, isNum := term.(Numeric); isNum {
			// Pure number term
			pairs = append(pairs, termPair{expr: NewInt(1), coeff: num})
		} else {
			// Pure expression term (coefficient = 1)
			pairs = append(pairs, termPair{expr: term, coeff: NewInt(1)})
		}
	}

	// Group by normalized expression string
	grouped := make(map[string][]termPair)
	for _, pair := range pairs {
		key := pair.expr.String() // Use string representation as key
		grouped[key] = append(grouped[key], pair)
	}

	// Combine coefficients for each unique expression
	var collectedTerms []Expr
	for _, group := range grouped {
		if len(group) == 0 {
			continue
		}

		expr := group[0].expr

		// Sum all coefficients
		var coeffSum Expr = group[0].coeff
		for i := 1; i < len(group); i++ {
			coeffSum = NewAdd(coeffSum, group[i].coeff)
		}

		// Simplify the coefficient sum
		if add, isAdd := coeffSum.(*Add); isAdd && len(add.terms) == 1 {
			coeffSum = add.terms[0]
		}

		// Create the collected term
		if isOne(coeffSum) {
			collectedTerms = append(collectedTerms, expr)
		} else if isZero(coeffSum) {
			// Skip zero terms
			continue
		} else if isOne(expr) {
			collectedTerms = append(collectedTerms, coeffSum)
		} else {
			collectedTerms = append(collectedTerms, NewMul(coeffSum, expr))
		}
	}

	if len(collectedTerms) == 0 {
		return NewInt(0)
	} else if len(collectedTerms) == 1 {
		return collectedTerms[0]
	}

	return NewAdd(collectedTerms...)
}

// partition separates a multiplication into numeric and non-numeric parts
func (m *Mul) partition() (Numeric, Expr) {
	var numericTerms []Numeric
	var exprTerms []Expr

	for _, term := range m.Terms() {
		if num, isNum := term.(Numeric); isNum {
			numericTerms = append(numericTerms, num)
		} else {
			exprTerms = append(exprTerms, term)
		}
	}

	// Combine numeric terms
	var numPart Numeric = NewInt(1)
	for _, num := range numericTerms {
		if intNum, isInt := numPart.(*Int); isInt {
			if intTerm, isIntTerm := num.(*Int); isIntTerm {
				// Multiply two integers
				result := new(big.Int).Mul(intNum.IntValue(), intTerm.IntValue())
				numPart = &Int{value: result}
			} else {
				// Convert to float and multiply
				val1, _ := intNum.Value().Float64()
				val2, _ := num.Value().Float64()
				numPart = NewFloat(val1 * val2)
			}
		} else {
			// Handle float multiplication
			val1, _ := numPart.Value().Float64()
			val2, _ := num.Value().Float64()
			numPart = NewFloat(val1 * val2)
		}
	}

	// Combine expression terms
	var exprPart Expr
	if len(exprTerms) == 0 {
		exprPart = NewInt(1)
	} else if len(exprTerms) == 1 {
		exprPart = exprTerms[0]
	} else {
		exprPart = NewMul(exprTerms...)
	}

	return numPart, exprPart
}

// Factor implements basic factoring for equations
// Extracts common factors from all terms in an addition
func (a *Add) Factor() Expr {
	if len(a.terms) == 0 {
		return NewInt(0)
	}

	if len(a.terms) == 1 {
		return a.terms[0]
	}

	// Find GCD of all coefficients
	var gcd *big.Int

	for _, term := range a.terms {
		var coeff *big.Int

		if mul, isMul := term.(*Mul); isMul {
			numPart, _ := mul.partition()
			if intNum, isInt := numPart.(*Int); isInt {
				coeff = new(big.Int).Set(intNum.IntValue())
			} else {
				// For non-integer coefficients, can't factor easily
				return a
			}
		} else if intTerm, isInt := term.(*Int); isInt {
			coeff = new(big.Int).Set(intTerm.IntValue())
		} else {
			// Non-numeric term, coefficient is 1
			coeff = big.NewInt(1)
		}

		if gcd == nil {
			gcd = new(big.Int).Abs(coeff)
		} else {
			gcd = new(big.Int).GCD(nil, nil, gcd, new(big.Int).Abs(coeff))
		}

		if gcd.Cmp(big.NewInt(1)) == 0 {
			// GCD is 1, no common factor
			return a
		}
	}

	// If GCD > 1, factor it out
	if gcd.Cmp(big.NewInt(1)) > 0 {
		var factoredTerms []Expr
		gcdExpr := &Int{value: gcd}

		for _, term := range a.terms {
			if mul, isMul := term.(*Mul); isMul {
				numPart, exprPart := mul.partition()
				if intNum, isInt := numPart.(*Int); isInt {
					// Divide coefficient by GCD
					quotient := new(big.Int).Div(intNum.IntValue(), gcd)
					newCoeff := &Int{value: quotient}

					if isOne(newCoeff) {
						factoredTerms = append(factoredTerms, exprPart)
					} else {
						factoredTerms = append(factoredTerms, NewMul(newCoeff, exprPart))
					}
				}
			} else if intTerm, isInt := term.(*Int); isInt {
				quotient := new(big.Int).Div(intTerm.IntValue(), gcd)
				factoredTerms = append(factoredTerms, &Int{value: quotient})
			} else {
				// This shouldn't happen if GCD calculation was correct
				factoredTerms = append(factoredTerms, term)
			}
		}

		var factoredSum Expr
		if len(factoredTerms) == 1 {
			factoredSum = factoredTerms[0]
		} else {
			factoredSum = NewAdd(factoredTerms...)
		}

		return NewMul(gcdExpr, factoredSum)
	}

	return a
}

// DivideThrough implements the Node.js KAS divideThrough() function
func (eq *Eq) DivideThrough(expr Expr) Expr {
	// First try to factor the expression
	var factored Expr = expr
	if add, isAdd := expr.(*Add); isAdd {
		factored = add.Factor()
	}

	// If not a multiplication after factoring, return original
	mul, isMul := factored.(*Mul)
	if !isMul {
		return expr
	}

	// Separate terms into additive expressions and other factors
	var addTerms []Expr
	var otherTerms []Expr

	for _, term := range mul.Terms() {
		if _, isAdd := term.(*Add); isAdd {
			addTerms = append(addTerms, term)
		} else {
			otherTerms = append(otherTerms, term)
		}
	}

	// For equalities, prefer keeping only Add terms if they exist
	if eq.eqType == EqEqual && len(addTerms) > 0 {
		if len(addTerms) == 1 {
			return addTerms[0]
		}
		return NewMul(addTerms...)
	}

	// Remove numeric factors (divide them out)
	var remainingTerms []Expr
	for _, term := range otherTerms {
		if _, isNum := term.(Numeric); !isNum {
			// Keep non-numeric terms
			remainingTerms = append(remainingTerms, term)
		}
	}

	// Combine remaining terms with add terms
	allRemainingTerms := append(addTerms, remainingTerms...)

	if len(allRemainingTerms) == 0 {
		return NewInt(1) // All factors were divided out
	} else if len(allRemainingTerms) == 1 {
		return allRemainingTerms[0]
	}

	return NewMul(allRemainingTerms...)
}
