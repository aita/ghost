package shell

import (
	"fmt"
	"strings"
)

var builtins []builtinCommand

type builtinCommand struct {
	name string
	desc string
	run  func(sh *Shell, env *Environment, args []string) int
}

func (cmd builtinCommand) Run(sh *Shell, env *Environment, args []string) int {
	return cmd.run(sh, env, args)
}

func init() {
	builtins = []builtinCommand{
		{
			name: "help",
			desc: "show help",
			run:  help,
		},
		{
			name: "echo",
			desc: "display a line of text",
			run:  echo,
		},
		{
			name: "set",
			desc: "change shell variables",
			run:  set,
		},
	}
}

func help(sh *Shell, env *Environment, args []string) int {
	for _, cmd := range builtins {
		fmt.Fprintf(sh.Out, "%s\t\t%s\n", cmd.name, cmd.desc)
	}
	return 0
}

func echo(sh *Shell, env *Environment, args []string) int {
	s := strings.Join(args[1:], " ")
	fmt.Fprintln(sh.Out, s)
	return 0
}

func set(sh *Shell, env *Environment, args []string) int {
	if len(args) != 3 {
		fmt.Fprintln(sh.Out, "usage: set VARIABLE_NAME VALUE")
		return 1
	}
	env.Set(args[1], args[2])
	return 0
}
