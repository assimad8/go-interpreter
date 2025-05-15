package repl



import (
	"bufio"
	"fmt"
	"io"
	"github.com/assimad8/go-interpreter/internal/token"
	"github.com/assimad8/go-interpreter/internal/lexer"
)

const PROMPT = ">> "

func Start(in io.Reader,out io.Writer){
	scanner := bufio.NewScanner(in)

	for{
		fmt.Print(PROMPT)
		scanned := scanner.Scan()

		if(!scanned) {return}

		line := scanner.Text()
		lex := lexer.New(line)

		for tk := lex.NextToken();tk.Type != token.EOF;tk = lex.NextToken() {
			fmt.Printf("%+v\n",tk)
		}
	}
}