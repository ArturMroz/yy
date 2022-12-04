package main

import (
	"fmt"
	"os"

	"ylang/repl"
)

func main() {
	fmt.Println("Hi, welcome to the Y programming language!")
	repl.Start(os.Stdin, os.Stdout)
}
