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
	run func(sh *Shell, args []string) int
}

func (cmd builtInCommand) Run(sh *Shell, args []string) int {
	return cmd.run(sh, args)
}

func echo(sh *Shell, args []string) int {
	s := strings.Join(args[1:], " ")
	fmt.Fprintln(sh.Out, s)
	return 0
}
