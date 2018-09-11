package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/aita/ghost/shell"
)

const prompt = ">> "

func die(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func printTokens(toks []*shell.Token) {
	for _, tok := range toks {
		lit := tok.Literal
		if tok.Kind == shell.NEWLINE {
			lit = ""
		}
		fmt.Printf("%d:%d %s %s\n", tok.Pos.Line, tok.Pos.Column, shell.TokenNames[tok.Kind], lit)
	}
}

func main() {
	r := bufio.NewReader(os.Stdin)
	scanner := shell.NewScanner(r)

loopEnd:
	for {
		fmt.Printf(prompt)
		toks := []*shell.Token{}
		for {
			tok, err := scanner.Scan()
			if err != nil {
				die(err)
			}
			toks = append(toks, tok)
			if tok.Kind == shell.NEWLINE {
				printTokens(toks)
				break
			}
			if tok.Kind == shell.EOF {
				fmt.Println()
				printTokens(toks)
				break loopEnd
			}
		}
	}
}
