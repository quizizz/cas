// Package latex provides advanced LaTeX formatting for mathematical expressions.
package latex

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/quizizz/cas/pkg/ast"
)

// FormatOptions controls LaTeX formatting behavior
type FormatOptions struct {
	// UseDisplayMode enables display-style formatting for large expressions
	UseDisplayMode bool
	// UseFractions converts divisions to proper fractions
	UseFractions bool
	// UseParentheses controls when to show parentheses
	UseParentheses bool
	// SimplifyRoots converts √(...) to simpler forms when possible
	SimplifyRoots bool
	// UseSymbols enables special mathematical symbols (π, e, etc.)
	UseSymbols bool
	// MaxDecimalPlaces limits decimal precision in output
	MaxDecimalPlaces int
}

// DefaultFormatOptions returns the default LaTeX formatting options
func DefaultFormatOptions() FormatOptions {
	return FormatOptions{
		UseDisplayMode:   false,
		UseFractions:     true,
		UseParentheses:   true,
		SimplifyRoots:    true,
		UseSymbols:       true,
		MaxDecimalPlaces: 6,
	}
}

// Format converts an expression to properly formatted LaTeX
func Format(expr ast.Expr, opts ...FormatOptions) string {
	options := DefaultFormatOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	return formatExpression(expr, options, 0)
}

// formatExpression formats an expression with proper precedence handling
func formatExpression(expr ast.Expr, opts FormatOptions, parentPrec int) string {
	switch e := expr.(type) {
	case *ast.Int:
		return formatInteger(e, opts)
	case *ast.Float:
		return formatFloat(e, opts)
	case *ast.Rational:
		return formatRational(e, opts)
	case *ast.Var:
		return formatVariable(e, opts)
	case *ast.Const:
		return formatConstant(e, opts)
	case *ast.Add:
		return formatAddition(e, opts, parentPrec)
	case *ast.Mul:
		return formatMultiplication(e, opts, parentPrec)
	case *ast.Pow:
		return formatPower(e, opts, parentPrec)
	case *ast.Func:
		return formatFunction(e, opts)
	default:
		return expr.String()
	}
}

func formatInteger(i *ast.Int, opts FormatOptions) string {
	val, _ := i.Eval(make(map[string]*big.Float))
	intVal, _ := val.Int64()

	if intVal < 0 {
		return fmt.Sprintf("-%d", -intVal)
	}
	return fmt.Sprintf("%d", intVal)
}

func formatFloat(f *ast.Float, opts FormatOptions) string {
	val, _ := f.Eval(make(map[string]*big.Float))

	// Check if it's effectively an integer
	if val.IsInt() {
		intVal, _ := val.Int64()
		return fmt.Sprintf("%d", intVal)
	}

	// Format with limited precision
	text := val.Text('f', opts.MaxDecimalPlaces)
	text = strings.TrimRight(text, "0")
	text = strings.TrimRight(text, ".")
	return text
}

func formatRational(r *ast.Rational, opts FormatOptions) string {
	num := r.Numerator()
	den := r.Denominator()

	if den.Cmp(big.NewInt(1)) == 0 {
		return num.String()
	}

	if opts.UseFractions {
		// Check for special cases
		if num.Cmp(big.NewInt(1)) == 0 && den.Cmp(big.NewInt(2)) == 0 {
			return "\\frac{1}{2}"
		}
		if num.Cmp(big.NewInt(1)) == 0 && den.Cmp(big.NewInt(3)) == 0 {
			return "\\frac{1}{3}"
		}
		if num.Cmp(big.NewInt(2)) == 0 && den.Cmp(big.NewInt(3)) == 0 {
			return "\\frac{2}{3}"
		}

		return fmt.Sprintf("\\frac{%s}{%s}", num.String(), den.String())
	}

	// Fallback to decimal
	rat := new(big.Rat).SetFrac(num, den)
	float, _ := rat.Float64()
	return fmt.Sprintf("%.6g", float)
}

func formatVariable(v *ast.Var, opts FormatOptions) string {
	name := v.Name()

	// Handle subscripts
	if strings.Contains(name, "_") {
		parts := strings.Split(name, "_")
		if len(parts) == 2 {
			return fmt.Sprintf("%s_{%s}", parts[0], parts[1])
		}
	}

	// Handle Greek letters and special symbols
	if opts.UseSymbols {
		switch name {
		case "alpha":
			return "\\alpha"
		case "beta":
			return "\\beta"
		case "gamma":
			return "\\gamma"
		case "delta":
			return "\\delta"
		case "epsilon":
			return "\\epsilon"
		case "theta":
			return "\\theta"
		case "lambda":
			return "\\lambda"
		case "mu":
			return "\\mu"
		case "pi":
			return "\\pi"
		case "sigma":
			return "\\sigma"
		case "phi":
			return "\\phi"
		case "omega":
			return "\\omega"
		}
	}

	// Multi-character variables get proper formatting
	if len(name) > 1 {
		return fmt.Sprintf("\\mathrm{%s}", name)
	}

	return name
}

func formatConstant(c *ast.Const, opts FormatOptions) string {
	if !opts.UseSymbols {
		return c.String()
	}

	switch c.Name() {
	case "pi", "π":
		return "\\pi"
	case "e":
		return "e"
	default:
		return c.Name()
	}
}

func formatAddition(add *ast.Add, opts FormatOptions, parentPrec int) string {
	terms := add.Terms()
	if len(terms) == 0 {
		return "0"
	}

	var parts []string
	for i, term := range terms {
		formatted := formatExpression(term, opts, 1)

		// Handle negative terms
		if i > 0 {
			if strings.HasPrefix(formatted, "-") {
				parts = append(parts, " - "+formatted[1:])
			} else {
				parts = append(parts, " + "+formatted)
			}
		} else {
			parts = append(parts, formatted)
		}
	}

	result := strings.Join(parts, "")

	// Add parentheses if needed
	if parentPrec > 1 && opts.UseParentheses {
		return fmt.Sprintf("\\left(%s\\right)", result)
	}

	return result
}

func formatMultiplication(mul *ast.Mul, opts FormatOptions, parentPrec int) string {
	factors := mul.Terms()
	if len(factors) == 0 {
		return "1"
	}

	var parts []string
	var hasNegative bool

	for i, factor := range factors {
		formatted := formatExpression(factor, opts, 3)

		// Check for negative coefficient
		if i == 0 && strings.HasPrefix(formatted, "-1") && len(factors) > 1 {
			hasNegative = true
			if formatted == "-1" {
				continue // Skip the -1 coefficient
			} else {
				formatted = "-" + formatted[2:] // Remove -1* and keep just -
			}
		}

		// Special formatting for common patterns
		if i > 0 && needsMultiplicationSpace(factors[i-1], factor) {
			parts = append(parts, " \\cdot "+formatted)
		} else if i > 0 {
			parts = append(parts, formatted)
		} else {
			parts = append(parts, formatted)
		}
	}

	result := strings.Join(parts, "")
	if hasNegative && !strings.HasPrefix(result, "-") {
		result = "-" + result
	}

	// Add parentheses if needed
	if parentPrec > 2 && opts.UseParentheses && len(factors) > 1 {
		return fmt.Sprintf("\\left(%s\\right)", result)
	}

	return result
}

func needsMultiplicationSpace(left, right ast.Expr) bool {
	// Add space between numbers
	if _, ok := left.(ast.Numeric); ok {
		if _, ok := right.(ast.Numeric); ok {
			return true
		}
	}

	// Add space between number and function
	if _, ok := left.(ast.Numeric); ok {
		if _, ok := right.(*ast.Func); ok {
			return true
		}
	}

	return false
}

func formatPower(pow *ast.Pow, opts FormatOptions, parentPrec int) string {
	base := formatExpression(pow.Base(), opts, 4)
	exp := formatExpression(pow.Exponent(), opts, 0)

	// Special handling for square roots
	if opts.SimplifyRoots {
		if exp == "\\frac{1}{2}" || exp == "0.5" {
			return fmt.Sprintf("\\sqrt{%s}", base)
		}

		// Cube roots and higher
		if strings.HasPrefix(exp, "\\frac{1}{") && strings.HasSuffix(exp, "}") {
			root := exp[8 : len(exp)-1] // Extract the denominator
			if root == "3" {
				return fmt.Sprintf("\\sqrt[3]{%s}", base)
			}
			if len(root) == 1 {
				return fmt.Sprintf("\\sqrt[%s]{%s}", root, base)
			}
		}
	}

	// Handle negative exponents as fractions
	if opts.UseFractions && strings.HasPrefix(exp, "-") {
		posExp := exp[1:]
		if posExp == "1" {
			return fmt.Sprintf("\\frac{1}{%s}", base)
		}
		return fmt.Sprintf("\\frac{1}{%s^{%s}}", base, posExp)
	}

	// Special case for exponent = 1
	if exp == "1" {
		return base
	}

	// Special case for exponent = 0
	if exp == "0" {
		return "1"
	}

	// General case
	if needsBaseBraces(pow.Base()) {
		return fmt.Sprintf("\\left(%s\\right)^{%s}", base, exp)
	}

	return fmt.Sprintf("%s^{%s}", base, exp)
}

func needsBaseBraces(base ast.Expr) bool {
	switch base.(type) {
	case *ast.Add:
		return true
	case *ast.Mul:
		// Only if it has more than one factor
		if mul, ok := base.(*ast.Mul); ok {
			return len(mul.Terms()) > 1
		}
	}
	return false
}

func formatFunction(fn *ast.Func, opts FormatOptions) string {
	name := fn.Name()
	args := fn.Args()

	if len(args) == 0 {
		return name
	}

	// Format arguments
	var argStrs []string
	for _, arg := range args {
		argStrs = append(argStrs, formatExpression(arg, opts, 0))
	}

	// Special LaTeX functions
	switch name {
	case "sqrt":
		if len(args) == 1 {
			return fmt.Sprintf("\\sqrt{%s}", argStrs[0])
		}
	case "sin", "cos", "tan", "sec", "csc", "cot":
		return fmt.Sprintf("\\%s\\left(%s\\right)", name, strings.Join(argStrs, ", "))
	case "arcsin", "arccos", "arctan":
		latexName := strings.Replace(name, "arc", "\\arcsin", 1)
		latexName = strings.Replace(latexName, "arccos", "\\arccos", 1)
		latexName = strings.Replace(latexName, "arctan", "\\arctan", 1)
		return fmt.Sprintf("%s\\left(%s\\right)", latexName, strings.Join(argStrs, ", "))
	case "sinh", "cosh", "tanh":
		return fmt.Sprintf("\\%s\\left(%s\\right)", name, strings.Join(argStrs, ", "))
	case "ln":
		return fmt.Sprintf("\\ln\\left(%s\\right)", strings.Join(argStrs, ", "))
	case "log":
		if len(args) == 1 {
			return fmt.Sprintf("\\log\\left(%s\\right)", argStrs[0])
		} else if len(args) == 2 {
			return fmt.Sprintf("\\log_{%s}\\left(%s\\right)", argStrs[1], argStrs[0])
		}
	case "abs":
		return fmt.Sprintf("\\left|%s\\right|", strings.Join(argStrs, ", "))
	case "exp":
		// Use e^x notation for exponential
		return fmt.Sprintf("e^{%s}", strings.Join(argStrs, ", "))
	}

	// Generic function formatting
	return fmt.Sprintf("\\mathrm{%s}\\left(%s\\right)", name, strings.Join(argStrs, ", "))
}

// FormatEquation formats an equation with proper LaTeX styling
func FormatEquation(lhs, rhs ast.Expr, opts ...FormatOptions) string {
	options := DefaultFormatOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	leftSide := Format(lhs, options)
	rightSide := Format(rhs, options)

	return fmt.Sprintf("%s = %s", leftSide, rightSide)
}

// FormatDerivative formats a derivative expression with proper notation
func FormatDerivative(expr ast.Expr, variable string, order int, opts ...FormatOptions) string {
	options := DefaultFormatOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	exprFormatted := Format(expr, options)

	if order == 1 {
		return fmt.Sprintf("\\frac{d}{d%s}\\left(%s\\right)", variable, exprFormatted)
	} else if order > 1 {
		return fmt.Sprintf("\\frac{d^{%d}}{d%s^{%d}}\\left(%s\\right)", order, variable, order, exprFormatted)
	}

	return exprFormatted
}

// FormatIntegral formats an integral expression
func FormatIntegral(expr ast.Expr, variable string, definite bool, lower, upper ast.Expr, opts ...FormatOptions) string {
	options := DefaultFormatOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	exprFormatted := Format(expr, options)

	if definite && lower != nil && upper != nil {
		lowerFormatted := Format(lower, options)
		upperFormatted := Format(upper, options)
		return fmt.Sprintf("\\int_{%s}^{%s} %s \\, d%s", lowerFormatted, upperFormatted, exprFormatted, variable)
	}

	return fmt.Sprintf("\\int %s \\, d%s", exprFormatted, variable)
}

// FormatMatrix formats matrix expressions (placeholder for future implementation)
func FormatMatrix(matrix [][]ast.Expr, opts ...FormatOptions) string {
	options := DefaultFormatOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	var rows []string
	for _, row := range matrix {
		var cols []string
		for _, col := range row {
			cols = append(cols, Format(col, options))
		}
		rows = append(rows, strings.Join(cols, " & "))
	}

	return fmt.Sprintf("\\begin{pmatrix}\n%s\n\\end{pmatrix}", strings.Join(rows, " \\\\\n"))
}
