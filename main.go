package main

import (
	"flag"
	"fmt"
	"os"

	"yy/eval"
	"yy/lexer"
	"yy/object"
	"yy/parser"
	"yy/repl"
)

var debug = flag.Bool("debug", false, "turns on debug mode")

func main() {
	flag.Parse()

	switch len(flag.Args()) {
	case 0:
		fmt.Println("YeetYoink v0.0.1")
		repl.Start(os.Stdin, os.Stdout, *debug)

	case 1:
		f := flag.Args()[0]
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
		fmt.Println("usage: yy [script] [--debug=true|false]")
	}
}
