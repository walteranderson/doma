package main

import (
	"bufio"
	"doma/pkg/eval"
	"doma/pkg/lexer"
	"doma/pkg/parser"
	"fmt"
	"io"
	"os"
)

func main() {
	switch len(os.Args) {
	case 1:
		startRepl(os.Stdin, os.Stdout)
	case 2:
		if os.Args[1] == "help" {
			showUsage()
			return
		}
		evalFile(os.Args[1])
	default:
		showUsage()
	}
}

func showUsage() {
	fmt.Println("Usage: ./doma [filename]")
}

func startRepl(in io.Reader, out io.Writer) {
	fmt.Println("Welcome to Doma!")
	scanner := bufio.NewScanner(in)
	env := eval.NewEnv()
	for {
		fmt.Fprint(out, "> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		evalPrint(line, env)
	}
}

func evalFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	contents, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	env := eval.NewEnv()
	evalPrint(string(contents), env)
}

func evalPrint(contents string, env *eval.Env) {
	l := lexer.New(contents)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			fmt.Println(err)
		}
		return
	}
	if len(program.Args) > 0 {
		obj := eval.Eval(program, env)
		if obj != nil {
			fmt.Println(obj.Inspect())
		}
	}
}

func getFileContent(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
