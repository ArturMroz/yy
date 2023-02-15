package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"yy/eval"
	"yy/lexer"
	"yy/object"
	"yy/parser"
)

const version = "v0.0.1"

var debug = flag.Bool("debug", false, "turns on debug mode")

func main() {
	flag.Parse()

	switch len(flag.Args()) {
	case 0:
		repl(os.Stdin, os.Stdout)

	case 1:
		runFile(flag.Args()[0])

	default:
		fmt.Println("usage: yy [--debug] [path_to_script] ")
	}
}

func runFile(f string) {
	src, err := os.ReadFile(f)
	if err != nil {
		fmt.Println("error: couldn't read file: " + f)
		os.Exit(1)
	}

	l := lexer.New(string(src))
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		for _, msg := range p.Errors() {
			fmt.Printf("parser error: %q\n", msg)
		}
		fmt.Println()
		os.Exit(1)
	}

	env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()

	eval.DefineMacros(program, macroEnv)
	expanded := eval.ExpandMacros(program, macroEnv)

	if *debug {
		fmt.Println("after macro expansion:")
		fmt.Println(expanded)
		fmt.Println()
	}

	result := eval.Eval(expanded, env)
	if evalError, ok := result.(*object.Error); ok {
		fmt.Printf("runtime error: %s\n", evalError.Msg)
	}
}

const (
	greet   = "YeetYoink " + version
	prompt  = "yy> "
	padLeft = "    "
)

func repl(in io.Reader, out io.Writer) {
	fmt.Println(greet)

	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()

	for {
		fmt.Fprint(out, prompt)
		if !scanner.Scan() {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if *debug {
			io.WriteString(out, padLeft)
			io.WriteString(out, program.String())
			io.WriteString(out, "\n")
		}

		if len(p.Errors()) != 0 {
			for _, msg := range p.Errors() {
				io.WriteString(out, "\t"+msg+"\n")
			}
			continue
		}

		eval.DefineMacros(program, macroEnv)
		expanded := eval.ExpandMacros(program, macroEnv)

		if *debug {
			io.WriteString(out, expanded.String())
			io.WriteString(out, "\n")
		}

		evaluated := eval.Eval(expanded, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
