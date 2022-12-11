package repl

import (
	"bufio"
	"fmt"
	"io"

	"ylang/eval"
	"ylang/lexer"
	"ylang/object"
	"ylang/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			for _, msg := range p.Errors() {
				io.WriteString(out, "\t"+msg+"\n")
			}
			continue
		}

		io.WriteString(out, "   ")
		io.WriteString(out, program.String())
		io.WriteString(out, "\n")

		evaluated := eval.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, "=> ")
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
