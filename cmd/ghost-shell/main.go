package main

import (
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
	env := &shell.Environment{}
	if len(os.Args) > 1 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			die(err)
		}
		prog, err := shell.Parse(file)
		shell.Eval(env, prog)
	} else {
		for {
			fmt.Printf(prompt)
			prog, err := shell.Parse(os.Stdin)
			if err != nil {
				die(err)
			}
			shell.Eval(env, prog)
		}
	}
}
