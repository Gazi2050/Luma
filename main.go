package main

import (
	"bufio"
	"fmt"
	"luma/evaluator"
	"luma/lexer"
	"luma/parser"
	"os"
)

func runLumaScript(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	l := lexer.NewLexer(string(content))
	p := parser.NewParser(l)
	program := p.ParseProgram()
	env := evaluator.NewEnv()
	evaluator.Eval(program, env)
}

func main() {
	if len(os.Args) > 2 && os.Args[1] == "run" {
		runLumaScript(os.Args[2])
		return
	}

	fmt.Println("Luma REPL v1")
	env := evaluator.NewEnv()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(">> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if line == "" {
			continue
		}
		l := lexer.NewLexer(line)
		p := parser.NewParser(l)
		program := p.ParseProgram()
		evaluator.Eval(program, env)
	}
}
