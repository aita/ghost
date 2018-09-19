package shell

import (
	"fmt"
	"strings"
)

var builtins = map[string]Command{
	"echo": builtInCommand{
		run: echo,
	},
}

type builtInCommand struct {
	run func(env *Environment, args []string) int
}

func (cmd builtInCommand) Run(env *Environment, args []string) int {
	return cmd.run(env, args)
}

func echo(env *Environment, args []string) int {
	s := strings.Join(args, " ")
	fmt.Fprintln(env.Stdout(), s)
	return 0
}
