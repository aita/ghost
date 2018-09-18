package ghost

import (
	"fmt"
	"strings"

	"github.com/aita/ghost/shell"
)

var builtins = map[string]shell.Command{
	"echo": builtInCommand{
		run: echo,
	},
}

type builtInCommand struct {
	run func(env *shell.Environment, args []string) int
}

func (cmd builtInCommand) Run(env *shell.Environment, args []string) int {
	return cmd.run(env, args)
}

func echo(env *shell.Environment, args []string) int {
	s := strings.Join(args, " ")
	fmt.Fprintln(env.Stdout(), s)
	return 0
}
