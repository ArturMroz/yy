package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"yy/eval"
	"yy/lexer"
	"yy/object"
	"yy/parser"
	"yy/yikes"
)

const version = "v0.0.1"

var debug = false

func main() {
	switch len(os.Args) {
	case 1:
		repl()

	case 2:
		runFile(os.Args[1])

	default:
		fmt.Println("usage: yy [path_to_script] ")
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
		for _, errMsg := range p.Errors() {
			fmt.Println(yikes.PrettyError(src, errMsg.Offset, errMsg.Msg))
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()

	eval.DefineMacros(program, macroEnv)
	expanded := eval.ExpandMacros(program, macroEnv)

	if debug {
		fmt.Println("after macro expansion:")
		fmt.Println(expanded)
		fmt.Println()
	}

	result := eval.Eval(expanded, env)
	if evalError, ok := result.(*object.Error); ok {
		fmt.Println(yikes.PrettyError(src, evalError.Pos, evalError.Msg))
	}
}

const (
	greet   = "YeetYoink " + version
	prompt  = "yy> "
	padLeft = "    "
)

func repl() {
	in := os.Stdin
	out := os.Stdout
	scanner := bufio.NewScanner(in)

	env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()

	fmt.Println(greet)

	for {
		fmt.Fprint(out, prompt)
		if !scanner.Scan() {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if debug {
			io.WriteString(out, padLeft)
			io.WriteString(out, program.String())
			io.WriteString(out, "\n")
		}

		if len(p.Errors()) > 0 {
			for _, msg := range p.Errors() {
				io.WriteString(out, msg.Error()+"\n")
			}
			continue
		}

		eval.DefineMacros(program, macroEnv)
		expanded := eval.ExpandMacros(program, macroEnv)

		if debug {
			io.WriteString(out, padLeft)
			io.WriteString(out, expanded.String())
			io.WriteString(out, "\n")
		}

		evaluated := eval.Eval(expanded, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.String())
			io.WriteString(out, "\n")
		}
	}
}
