package main

import (
	"fmt"
	"os"

	"ylang/eval"
	"ylang/lexer"
	"ylang/object"
	"ylang/parser"
	"ylang/repl"
)

func main() {
	switch len(os.Args) {
	case 1:
		fmt.Println("Welcome to the Y programming language REPL!")
		repl.Start(os.Stdin, os.Stdout)

	case 2:
		f := os.Args[1]
		src, err := os.ReadFile(f)
		if err != nil {
			fmt.Println("error: couldn't read file: " + f)
			return
		}

		l := lexer.New(string(src))
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) > 0 {
			for _, msg := range p.Errors() {
				fmt.Printf("parser error: %q\n", msg)
			}
			fmt.Println()
			return
		}

		result := eval.Eval(program, object.NewEnvironment())
		if evalError, ok := result.(*object.Error); ok {
			fmt.Printf("runtime error: %s\n", evalError.Msg)
		}

	default:
		fmt.Println("usage: ylang [script]")
	}
}
