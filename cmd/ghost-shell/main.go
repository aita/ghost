package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/aita/ghost/shell"
)

const prompt = ">> "

func die(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func main() {
	evaluator := shell.NewEvaluator(bytes.NewReader(nil), os.Stdout, os.Stderr)
	if len(os.Args) > 1 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			die(err)
		}
		prog, err := shell.Parse(file)
		evaluator.Eval(prog)
	} else {
		for {
			fmt.Printf(prompt)
			prog, err := shell.Parse(os.Stdin)
			if err != nil {
				die(err)
			}
			evaluator.Eval(prog)
		}
	}
}
