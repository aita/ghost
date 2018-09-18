package shell

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/aita/ghost/shell/eval"
	"github.com/aita/ghost/shell/parser"
)

type Shell struct {
	env  *eval.Environment
	eval *eval.Evaluator
}

func (sh *Shell) Init() {
	sh.env = &eval.Environment{
		In:  bytes.NewReader(nil),
		Out: ioutil.Discard,
	}
	commands := map[string]eval.Command{}
	for name, cmd := range builtins {
		commands[name] = cmd
	}
	sh.eval = &eval.Evaluator{
		Commands: commands,
	}
}

func (sh *Shell) Exec(w io.Writer, script string) {
	sh.env.Out = w

	prog, err := parser.Parse(strings.NewReader(script))
	if err != nil {
		fmt.Fprintln(w, "ghost:", err.Error())
		return
	}
	sh.eval.Eval(sh.env, prog)
}
