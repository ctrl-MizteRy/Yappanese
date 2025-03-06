package main

import (
	"bufio"
	"fmt"
	"os"
	//"os/user"
	"strings"
	"yap/evaluator"
	"yap/lexer"
	"yap/object"
	"yap/parser"
	//"yap/repl"
)

/*
import (
	"fmt"
	"os"
	"os/user"
	"yap/repl"
)*/

func main() {
	/*
		user, err := user.Current()
		if err != nil {
			panic(err)
		} else {
			fmt.Printf("Hello %s! Welcome to yappanese!\n", user.Username)
			fmt.Println("Try to write some code into the command-line")
			repl.Start(os.Stdin, os.Stdout)
		}*/

	args := os.Args
	if len(args) != 2 {
		fmt.Println("Please enter a file you want to interpret")
		os.Exit(1)
	}

	fileName := args[1]
	fileType := strings.Split(fileName, ".")

	if len(fileType) != 2 {
		fmt.Printf("Please provide a file and not %s\n", fileName)
		os.Exit(1)
	}
	if fileType[1] != "txt" {
		fmt.Printf("Please provide a '.txt' file, not .%s", fileType)
		os.Exit(1)
	}

	lines := []string{}
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Could not open %s file", fileName)
		os.Exit(2)
	}

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	l := lexer.New(strings.Join(lines, " "))
	p := parser.New(l)
	program := p.ParserProgram()
	if len(p.Errors()) != 0 {
		for _, error := range p.Errors() {
			fmt.Printf("\t" + error + "\n")
		}
		os.Exit(3)
	}
	env := object.NewEnviroment()
	eval := evaluator.Eval(program, env)
	if eval != nil {
		fmt.Println(eval.Inspect())
	}

}
