package repl

import (
	"bufio"
	"fmt"
	"io"
	"yap/evaluator"
	"yap/lexer"
	"yap/parser"
)

const PROMPT = "> "

func Start(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	for {
		fmt.Fprint(out, PROMPT)
		scan := sc.Scan()

		if !scan {
			return
		}

		line := sc.Text()

		if len(line) == 0 {
			return
		}

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParserProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
