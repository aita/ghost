package shell

import (
	"fmt"
	"strings"

	"github.com/aita/ghost/shell/eval"
)

var builtins = map[string]eval.Command{
	"echo": builtInCommand{
		run: echo,
	},
}

type builtInCommand struct {
	run func(env *eval.Environment, args []string) int
}

func (cmd builtInCommand) Run(env *eval.Environment, args []string) int {
	return cmd.run(env, args)
}

func echo(env *eval.Environment, args []string) int {
	s := strings.Join(args, " ")
	fmt.Fprintln(env.Stdout(), s)
	return 0
}
