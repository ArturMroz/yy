package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"syscall/js"

	"yy/eval"
	"yy/lexer"
	"yy/object"
	"yy/parser"
)

const version = "v0.0.1"

var debug = flag.Bool("debug", false, "turns on debug mode")

func main() {
	// fmt.Println("args", os.Args)
	flag.Parse()
	if len(os.Args) == 1 && os.Args[0] == "js" {
		fmt.Println("web ass baby, interpret")
		js.Global().Set("interpret", interpretWrapper())
		<-make(chan bool)
		return
	}

	switch len(flag.Args()) {
	case 0:
		repl()

	case 1:
		runFile(flag.Args()[0])

	default:
		fmt.Println("usage: yy [--debug] [path_to_script] ")
		// TODO work around this limitation for better UX (this is how Go stdlib flag parsing works)
		fmt.Println("note that the --debug flag must come before the path to the script")
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

func interpretWrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return "no arguments passed"
		}

		input := args[0].String()
		fmt.Printf("input %s\n", input)

		output, err := interpret(input)
		if err != nil {
			return map[string]any{
				"error": fmt.Sprintf("error: %s\n", err),
			}
		}
		return output
	})

	return jsonFunc
}

func interpret(src string) (string, error) {
	l := lexer.New(src)
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

	result := eval.Eval(expanded, env)
	if evalError, ok := result.(*object.Error); ok {
		return "", errors.New(fmt.Sprintf("runtime error: %s\n", evalError.Msg))
	}
	return result.String(), nil
}
