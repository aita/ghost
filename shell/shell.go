package shell

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

type Shell struct {
	env  *Environment
	eval *Evaluator
}

func (sh *Shell) Init() {
	sh.env = &Environment{
		In:  bytes.NewReader(nil),
		Out: ioutil.Discard,
	}

	commands := map[string]Command{}
	for name, cmd := range builtins {
		commands[name] = cmd
	}
	sh.eval = &Evaluator{
		Commands: commands,
	}
}

func (sh *Shell) Exec(w io.Writer, script string) {
	sh.env.Out = w

	prog, err := Parse(strings.NewReader(script))
	if err != nil {
		fmt.Fprintln(w, "ghost:", err.Error())
		return
	}
	sh.eval.Eval(sh.env, prog)
}
