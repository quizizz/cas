// Example demonstrating basic CAS operations
package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/quizizz/cas/pkg/latex"
	"github.com/quizizz/cas/pkg/parser"
)

func main() {
	fmt.Println("=== Basic CAS Usage Examples ===\n")

	// Example 1: Parse and evaluate expressions
	fmt.Println("1. Expression Parsing and Evaluation")
	expr, err := parser.Parse("x^2 + 3*x + 2")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Parsed expression: %s\n", expr.String())
	fmt.Printf("LaTeX: %s\n", latex.Format(expr))

	// Evaluate with specific variable values
	vars := map[string]*big.Float{
		"x": big.NewFloat(2.0),
	}
	result, err := expr.Eval(vars)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("When x = 2: %s\n", result.Text('f', 2))
	fmt.Println()

	// Example 2: Working with constants
	fmt.Println("2. Mathematical Constants")
	piExpr, _ := parser.Parse("pi * r^2")
	fmt.Printf("Area formula: %s\n", piExpr.String())
	fmt.Printf("LaTeX: %s\n", latex.Format(piExpr))

	varsWithRadius := map[string]*big.Float{
		"r": big.NewFloat(5.0),
	}
	area, _ := piExpr.Eval(varsWithRadius)
	fmt.Printf("Circle area (r=5): %s\n", area.Text('f', 4))
	fmt.Println()

	// Example 3: Complex expressions
	fmt.Println("3. Complex Mathematical Expressions")
	complexExpr, _ := parser.Parse("sin(x) + cos(x^2) * ln(y)")
	fmt.Printf("Expression: %s\n", complexExpr.String())
	fmt.Printf("LaTeX: %s\n", latex.Format(complexExpr))
	fmt.Printf("Variables: %v\n", complexExpr.Variables())

	// Try to evaluate (will show variables needed)
	_, err = complexExpr.Eval(make(map[string]*big.Float))
	if err != nil {
		fmt.Printf("Error (expected): %v\n", err)
	}
	fmt.Println()

	// Example 4: Simplification
	fmt.Println("4. Expression Simplification")
	redundantExpr, _ := parser.Parse("x + 0 + x*1")
	fmt.Printf("Original: %s\n", redundantExpr.String())

	simplified := redundantExpr.Simplify()
	fmt.Printf("Simplified: %s\n", simplified.String())
	fmt.Printf("LaTeX: %s\n", latex.Format(simplified))
	fmt.Println()

	// Example 5: Function expressions
	fmt.Println("5. Mathematical Functions")
	functions := []string{
		"sqrt(16)",
		"abs(-5)",
		"ln(e)",
		"sin(pi/2)",
		"exp(0)",
	}

	for _, funcStr := range functions {
		funcExpr, _ := parser.Parse(funcStr)
		result, err := funcExpr.Eval(make(map[string]*big.Float))
		if err != nil {
			fmt.Printf("%s = Error: %v\n", funcStr, err)
		} else {
			fmt.Printf("%s = %s\n", funcStr, result.Text('g', 6))
		}
	}
}
