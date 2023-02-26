package main

import (
	"errors"
	"fmt"
	"syscall/js"

	"yy/eval"
	"yy/lexer"
	"yy/object"
	"yy/parser"
)

const version = "v0.0.1"

var debug = false

func main() {
	js.Global().Set("interpret", wasmWrapper())
	<-make(chan bool)
}

func wasmWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return map[string]any{
				"error": fmt.Sprintf("wrong number of args (got %d, want 1)\n", len(args)),
			}
		}

		input := args[0].String()
		if err := interpret(input); err != nil {
			return map[string]any{
				"error": fmt.Sprintf("%s\n", err),
			}
		}
		return nil
	})
}

func interpret(src string) error {
	l := lexer.New(src)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) == 1 {
		return errors.New("parser error: " + p.Errors()[0])
	} else if len(p.Errors()) > 0 {
		errMsg := "parser errors:\n"
		for _, err := range p.Errors() {
			errMsg += fmt.Sprintf("    %q\n", err)
		}
		return errors.New(errMsg)
	}

	env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()

	eval.DefineMacros(program, macroEnv)
	expanded := eval.ExpandMacros(program, macroEnv)

	result := eval.Eval(expanded, env)
	if evalError, ok := result.(*object.Error); ok {
		return fmt.Errorf("yy runtime error: %s\n", evalError.Msg)
	}
	return nil
}
