package repl

import (
	"bufio"
	"fmt"
	"io"

	// "github.com/assimad8/go-interpreter/internal/token"
	"github.com/assimad8/go-interpreter/internal/evaluator"
	"github.com/assimad8/go-interpreter/internal/lexer"
	"github.com/assimad8/go-interpreter/internal/object"
	"github.com/assimad8/go-interpreter/internal/parser"
)

const MONKEY_FACE = `
		__,__
.--. .-"	 "-.  .--.
/ .. \/ .-. .-.  \/ .. \
| | ' | /   Y   \  |' | |
| | \ \ \ 0 | 0 /  /  / |
\'- ,\.-"""""""-./,-' /
''-'    /_ ^ ^ _\   '-''
	   | \._ _./ |
		\ \ '~' / /
	   '._ '-=-' _.'
	      '-----'

`

const PROMPT = ">> "

func Start(in io.Reader,out io.Writer){
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for{
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "quit" {
			break
		}
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserError(out,p.Errors())
			continue
		}
		evaluated := evaluator.Eval(program,env)
		if evaluated!=nil{
			io.WriteString(out,evaluated.Inspect())
			io.WriteString(out,"\n")
		}
	}
}

func printParserError(out io.Writer,errors []string) {
	io.WriteString(out,MONKEY_FACE)
	io.WriteString(out,"Woops! We ran into some monkey busniss here!\n")
	io.WriteString(out," parser errors:\n")
	for _,msg := range errors {
		io.WriteString(out,"\t"+msg+"\n")
	}
}
