package ghost

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/aita/ghost/shell"
)

type Shell struct {
	env  *shell.Environment
	eval *shell.Evaluator
}

func (sh *Shell) Init() {
	sh.env = &shell.Environment{
		In:  bytes.NewReader(nil),
		Out: ioutil.Discard,
	}
	commands := map[string]shell.Command{}
	for name, cmd := range builtins {
		commands[name] = cmd
	}
	sh.eval = &shell.Evaluator{
		Commands: commands,
	}
}

func (sh *Shell) Exec(w io.Writer, script string) {
	sh.env.Out = w

	prog, err := shell.Parse(strings.NewReader(script))
	if err != nil {
		fmt.Fprintln(w, "ghost:", err.Error())
		return
	}
	sh.eval.Eval(sh.env, prog)
}
