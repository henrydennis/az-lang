package main

import (
	"az-lang/interpreter"
	"az-lang/lexer"
	"az-lang/object"
	"az-lang/parser"
	"bufio"
	"fmt"
	"os"
	"strings"
)

const VERSION = "0.1.0"

func main() {
	if len(os.Args) > 1 {
		// File mode
		filename := os.Args[1]
		if !strings.HasSuffix(filename, ".abc") {
			fmt.Println("Error: ABC files must have .abc extension")
			os.Exit(1)
		}
		runFile(filename)
	} else {
		// REPL mode
		runREPL()
	}
}

func runFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}

	env := object.NewEnvironment()
	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		printParserErrors(p.Errors())
		os.Exit(1)
	}

	result := interpreter.Eval(program, env)
	if result != nil {
		if errObj, ok := result.(*object.Error); ok {
			fmt.Println(errObj.Inspect())
			os.Exit(1)
		}
	}
}

func runREPL() {
	fmt.Printf("ABC Language v%s\n", VERSION)
	fmt.Println("An English-like programming language")
	fmt.Println("Type your code below. Press Ctrl+C to exit.\n")

	scanner := bufio.NewScanner(os.Stdin)
	env := object.NewEnvironment()

	for {
		fmt.Print("abc> ")
		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		if line == "" {
			continue
		}

		// Handle multi-line input for blocks
		if needsMoreInput(line) {
			line = readMultiLine(scanner, line)
		}

		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			printParserErrors(p.Errors())
			continue
		}

		result := interpreter.Eval(program, env)
		if result != nil {
			if result.Type() != object.NULL_OBJ {
				fmt.Println(result.Inspect())
			}
		}
	}
}

// needsMoreInput checks if the line has unclosed blocks
func needsMoreInput(line string) bool {
	beginCount := strings.Count(strings.ToLower(line), "begin")
	endCount := strings.Count(strings.ToLower(line), "end")
	return beginCount > endCount
}

// readMultiLine reads additional lines until blocks are balanced
func readMultiLine(scanner *bufio.Scanner, firstLine string) string {
	var builder strings.Builder
	builder.WriteString(firstLine)
	builder.WriteString("\n")

	beginCount := strings.Count(strings.ToLower(firstLine), "begin")
	endCount := strings.Count(strings.ToLower(firstLine), "end")

	for beginCount > endCount {
		fmt.Print("...> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		builder.WriteString(line)
		builder.WriteString("\n")

		beginCount += strings.Count(strings.ToLower(line), "begin")
		endCount += strings.Count(strings.ToLower(line), "end")
	}

	return builder.String()
}

func printParserErrors(errors []string) {
	fmt.Println("Parser errors:")
	for _, msg := range errors {
		fmt.Printf("  %s\n", msg)
	}
}
