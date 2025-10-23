// Example demonstrating calculus operations: differentiation and gradients
package main

import (
	"fmt"
	"log"

	"github.com/quizizz/cas/pkg/calculus"
	"github.com/quizizz/cas/pkg/latex"
	"github.com/quizizz/cas/pkg/parser"
)

func main() {
	fmt.Println("=== Calculus Operations Examples ===\n")

	// Example 1: Basic differentiation
	fmt.Println("1. Basic Differentiation")
	examples := []struct {
		expr     string
		variable string
	}{
		{"x^3", "x"},
		{"2*x^2 + 3*x + 1", "x"},
		{"sin(x)", "x"},
		{"e^x", "x"},
		{"ln(x)", "x"},
		{"sqrt(x)", "x"},
	}

	for _, ex := range examples {
		expr, err := parser.Parse(ex.expr)
		if err != nil {
			log.Printf("Parse error for %s: %v", ex.expr, err)
			continue
		}

		derivative, err := calculus.Derivative(expr, ex.variable)
		if err != nil {
			log.Printf("Differentiation error for %s: %v", ex.expr, err)
			continue
		}

		fmt.Printf("d/d%s(%s) = %s\n", ex.variable, ex.expr, derivative.String())
		fmt.Printf("LaTeX: %s\n", latex.FormatDerivative(expr, ex.variable, 1))
		fmt.Printf("Result LaTeX: %s\n", latex.Format(derivative))
		fmt.Println()
	}

	// Example 2: Chain rule examples
	fmt.Println("2. Chain Rule Applications")
	chainRuleExamples := []string{
		"sin(x^2)",
		"ln(x^2 + 1)",
		"e^(x^3)",
		"sqrt(x^2 + 1)",
		"(x^2 + 1)^3",
	}

	for _, exprStr := range chainRuleExamples {
		expr, _ := parser.Parse(exprStr)
		derivative, err := calculus.Derivative(expr, "x")
		if err != nil {
			continue
		}

		fmt.Printf("f(x) = %s\n", exprStr)
		fmt.Printf("f'(x) = %s\n", derivative.String())
		fmt.Printf("LaTeX: %s\n", latex.Format(derivative))
		fmt.Println()
	}

	// Example 3: Product rule
	fmt.Println("3. Product Rule")
	productExamples := []string{
		"x * sin(x)",
		"x^2 * ln(x)",
		"e^x * cos(x)",
		"x * sqrt(x + 1)",
	}

	for _, exprStr := range productExamples {
		expr, _ := parser.Parse(exprStr)
		derivative, err := calculus.Derivative(expr, "x")
		if err != nil {
			continue
		}

		fmt.Printf("d/dx(%s) = %s\n", exprStr, derivative.String())
		fmt.Printf("LaTeX: %s\n", latex.Format(derivative))
		fmt.Println()
	}

	// Example 4: Higher order derivatives
	fmt.Println("4. Higher Order Derivatives")
	expr, _ := parser.Parse("x^4 + 2*x^3 - x^2 + 5")

	fmt.Printf("f(x) = %s\n", expr.String())
	fmt.Printf("LaTeX: %s\n", latex.Format(expr))
	fmt.Println()

	for order := 1; order <= 4; order++ {
		derivative, err := calculus.NthDerivative(expr, "x", order)
		if err != nil {
			fmt.Printf("Error computing %d-th derivative: %v\n", order, err)
			continue
		}

		fmt.Printf("f%s(x) = %s\n", getSuperscript(order), derivative.String())
		fmt.Printf("LaTeX: %s\n", latex.FormatDerivative(expr, "x", order))
		fmt.Printf("Result: %s\n", latex.Format(derivative))
		fmt.Println()
	}

	// Example 5: Multivariable calculus - Gradients
	fmt.Println("5. Gradient (Partial Derivatives)")
	multiVarExamples := []struct {
		expr      string
		variables []string
	}{
		{"x^2 + y^2", []string{"x", "y"}},
		{"x*y + x^2", []string{"x", "y"}},
		{"sin(x) * cos(y)", []string{"x", "y"}},
		{"x^2 + y^2 + z^2", []string{"x", "y", "z"}},
	}

	for _, ex := range multiVarExamples {
		expr, _ := parser.Parse(ex.expr)
		gradient, err := calculus.Gradient(expr, ex.variables)
		if err != nil {
			fmt.Printf("Gradient error for %s: %v\n", ex.expr, err)
			continue
		}

		fmt.Printf("f(%s) = %s\n", joinVariables(ex.variables), ex.expr)
		fmt.Printf("∇f = (")

		for i, variable := range ex.variables {
			if i > 0 {
				fmt.Printf(", ")
			}
			//partial := gradient[variable]
			fmt.Printf("∂f/∂%s", variable)
		}
		fmt.Printf(")\n")

		fmt.Printf("∇f = (")
		for i, variable := range ex.variables {
			if i > 0 {
				fmt.Printf(", ")
			}
			partial := gradient[variable]
			fmt.Printf("%s", partial.String())
		}
		fmt.Printf(")\n")

		fmt.Printf("LaTeX: ∇f = \\left(")
		for i, variable := range ex.variables {
			if i > 0 {
				fmt.Printf(", ")
			}
			partial := gradient[variable]
			fmt.Printf("%s", latex.Format(partial))
		}
		fmt.Printf("\\right)\n")
		fmt.Println()
	}
}

func getSuperscript(n int) string {
	superscripts := map[int]string{
		1: "'",
		2: "''",
		3: "'''",
		4: "''''",
	}
	if s, ok := superscripts[n]; ok {
		return s
	}
	return fmt.Sprintf("^(%d)", n)
}

func joinVariables(vars []string) string {
	if len(vars) == 0 {
		return ""
	}
	if len(vars) == 1 {
		return vars[0]
	}

	result := vars[0]
	for i := 1; i < len(vars)-1; i++ {
		result += ", " + vars[i]
	}
	result += ", " + vars[len(vars)-1]
	return result
}
