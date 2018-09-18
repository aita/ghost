package shell

import (
	"fmt"
	"io"
	"strings"
)

type Command interface {
	Run(stdin io.Reader, stdout, stderr io.Writer, args []string) int
}

var builtins = map[string]Command{
	"echo": CommandFunc(echo),
}

type CommandFunc func(stdin io.Reader, stdout, stderr io.Writer, args []string) int

func (f CommandFunc) Run(stdin io.Reader, stdout, stderr io.Writer, args []string) int {
	return f(stdin, stdout, stderr, args)
}

func echo(stdin io.Reader, stdout, stderr io.Writer, args []string) int {
	s := strings.Join(args, " ")
	fmt.Fprintln(stdout, s)
	return 0
}
