// Example demonstrating advanced LaTeX formatting capabilities
package main

import (
	"fmt"
	"log"

	"github.com/quizizz/cas/pkg/calculus"
	"github.com/quizizz/cas/pkg/latex"
	"github.com/quizizz/cas/pkg/parser"
)

func main() {
	fmt.Println("=== LaTeX Formatting Examples ===\n")

	// Example 1: Basic expression formatting
	fmt.Println("1. Basic Expression Formatting")
	basicExpressions := []string{
		"x^2",
		"x^2 + 2*x + 1",
		"2*x^3 - x^2 + 5*x - 3",
		"sqrt(x^2 + y^2)",
		"1/x",
		"x^(-1)",
		"x^(-2)",
		"(x + 1)^2",
	}

	for _, exprStr := range basicExpressions {
		expr, err := parser.Parse(exprStr)
		if err != nil {
			log.Printf("Parse error for %s: %v", exprStr, err)
			continue
		}

		fmt.Printf("Expression: %s\n", exprStr)
		fmt.Printf("Basic LaTeX: %s\n", expr.LaTeX())
		fmt.Printf("Enhanced LaTeX: %s\n", latex.Format(expr))
		fmt.Println()
	}

	// Example 2: Function formatting
	fmt.Println("2. Function Formatting")
	functionExpressions := []string{
		"sin(x)",
		"cos(x^2)",
		"tan(pi/4)",
		"ln(x + 1)",
		"log(x)",
		"exp(x^2)",
		"abs(x - 1)",
		"sqrt(x^2 + 1)",
		"sinh(x)",
		"cosh(2*x)",
		"tanh(x/2)",
	}

	for _, exprStr := range functionExpressions {
		expr, err := parser.Parse(exprStr)
		if err != nil {
			continue
		}

		fmt.Printf("Expression: %s\n", exprStr)
		fmt.Printf("LaTeX: %s\n", latex.Format(expr))
		fmt.Println()
	}

	// Example 3: Complex expressions with mixed operations
	fmt.Println("3. Complex Mixed Expressions")
	complexExpressions := []string{
		"sin(x^2) * cos(x) + ln(x + 1)",
		"e^(x^2) / sqrt(x + 1)",
		"(sin(x) + cos(x))^2",
		"x^2 * e^(-x^2/2) / sqrt(2*pi)",
		"(a*x + b) / (c*x + d)",
		"sqrt(a^2 + b^2) * sin(theta)",
	}

	for _, exprStr := range complexExpressions {
		expr, err := parser.Parse(exprStr)
		if err != nil {
			continue
		}

		fmt.Printf("Expression: %s\n", exprStr)
		fmt.Printf("LaTeX: %s\n", latex.Format(expr))
		fmt.Println()
	}

	// Example 4: Formatting options
	fmt.Println("4. Custom Formatting Options")
	expr, _ := parser.Parse("1/2 + pi + (x + y)^2")

	// Default options
	defaultLatex := latex.Format(expr)
	fmt.Printf("Default: %s\n", defaultLatex)

	// Custom options - no fractions
	noFractionsOpts := latex.FormatOptions{
		UseFractions:   false,
		UseSymbols:     true,
		UseParentheses: true,
	}
	noFractionsLatex := latex.Format(expr, noFractionsOpts)
	fmt.Printf("No fractions: %s\n", noFractionsLatex)

	// Custom options - no symbols
	noSymbolsOpts := latex.FormatOptions{
		UseFractions:   true,
		UseSymbols:     false,
		UseParentheses: true,
	}
	noSymbolsLatex := latex.Format(expr, noSymbolsOpts)
	fmt.Printf("No symbols: %s\n", noSymbolsLatex)

	// Custom options - no parentheses
	noParensOpts := latex.FormatOptions{
		UseFractions:   true,
		UseSymbols:     true,
		UseParentheses: false,
	}
	noParensLatex := latex.Format(expr, noParensOpts)
	fmt.Printf("No parentheses: %s\n", noParensLatex)
	fmt.Println()

	// Example 5: Special formatting functions
	fmt.Println("5. Special Formatting Functions")

	// Equation formatting
	lhs, _ := parser.Parse("x^2 + 2*x")
	rhs, _ := parser.Parse("8")
	equation := latex.FormatEquation(lhs, rhs)
	fmt.Printf("Equation: %s\n", equation)

	// Derivative formatting
	expr, _ = parser.Parse("x^3 + sin(x)")
	derivativeLatex := latex.FormatDerivative(expr, "x", 1)
	fmt.Printf("First derivative: %s\n", derivativeLatex)

	secondDerivativeLatex := latex.FormatDerivative(expr, "x", 2)
	fmt.Printf("Second derivative: %s\n", secondDerivativeLatex)

	// Integral formatting (indefinite)
	integralLatex := latex.FormatIntegral(expr, "x", false, nil, nil)
	fmt.Printf("Indefinite integral: %s\n", integralLatex)

	// Integral formatting (definite)
	lower, _ := parser.Parse("0")
	upper, _ := parser.Parse("1")
	definiteIntegralLatex := latex.FormatIntegral(expr, "x", true, lower, upper)
	fmt.Printf("Definite integral: %s\n", definiteIntegralLatex)
	fmt.Println()

	// Example 6: Greek letters and variables
	fmt.Println("6. Greek Letters and Special Variables")
	greekExamples := []string{
		"alpha + beta",
		"gamma * delta",
		"pi * r^2",
		"theta + phi",
		"lambda * x",
		"mu + sigma",
		"omega * t",
	}

	for _, exprStr := range greekExamples {
		expr, err := parser.Parse(exprStr)
		if err != nil {
			continue
		}

		fmt.Printf("Expression: %s\n", exprStr)
		fmt.Printf("LaTeX: %s\n", latex.Format(expr))
		fmt.Println()
	}

	// Example 7: Polynomial with subscripts
	fmt.Println("7. Variables with Subscripts")
	subscriptExamples := []string{
		"x_1 + x_2",
		"a_0 + a_1*x + a_2*x^2",
		"y_max - y_min",
	}

	for _, exprStr := range subscriptExamples {
		expr, err := parser.Parse(exprStr)
		if err != nil {
			fmt.Printf("Could not parse %s (subscripts may need manual creation)\n", exprStr)
			continue
		}

		fmt.Printf("Expression: %s\n", exprStr)
		fmt.Printf("LaTeX: %s\n", latex.Format(expr))
		fmt.Println()
	}

	// Example 8: Matrix formatting (if supported)
	fmt.Println("8. Matrix Formatting")

	// Create a simple 2x2 matrix manually for demonstration
	//x11, _ := parser.Parse("1")
	//x12, _ := parser.Parse("2")
	//x21, _ := parser.Parse("3")
	//x22, _ := parser.Parse("4")

	//matrix := [][]interface{}{
	//	{x11, x12},
	//	{x21, x22},
	//}

	// Note: This would require the matrix formatting function to accept interface{}
	fmt.Printf("Matrix example:\n")
	fmt.Printf("LaTeX: \\begin{pmatrix} 1 & 2 \\\\ 3 & 4 \\end{pmatrix}\n")
	fmt.Println()

	// Example 9: Real-world mathematical expressions
	fmt.Println("9. Real-World Mathematical Expressions")
	realWorldExamples := []string{
		"x^2 + y^2",               // Circle equation
		"a*x^2 + b*x + c",         // Quadratic formula
		"e^(-x^2/2) / sqrt(2*pi)", // Normal distribution
		"sin(2*pi*f*t + phi)",     // Sinusoidal wave
		"V * e^(-t/(R*C))",        // RC circuit discharge
	}

	descriptions := []string{
		"Circle equation",
		"General quadratic",
		"Normal distribution (unnormalized)",
		"Sinusoidal wave",
		"RC circuit discharge",
	}

	for i, exprStr := range realWorldExamples {
		expr, err := parser.Parse(exprStr)
		if err != nil {
			continue
		}

		fmt.Printf("%s: %s\n", descriptions[i], exprStr)
		fmt.Printf("LaTeX: %s\n", latex.Format(expr))

		// Try to compute derivative if single variable
		vars := expr.Variables()
		if len(vars) == 1 {
			derivative, err := calculus.Derivative(expr, vars[0])
			if err == nil {
				fmt.Printf("d/d%s: %s\n", vars[0], latex.Format(derivative))
			}
		}
		fmt.Println()
	}
}
