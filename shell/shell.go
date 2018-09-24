package shell

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

type Command interface {
	Run(shell *Shell, env *Environment, args []string) int
}

type Shell struct {
	status   int
	topLevel *Environment
	commands map[string]Command

	In  io.Reader
	Out io.Writer
}

func (sh *Shell) Init() {
	if sh.In == nil {
		sh.In = bytes.NewReader(nil)
	}
	if sh.Out == nil {
		sh.Out = ioutil.Discard
	}
	sh.topLevel = &Environment{}
	sh.commands = map[string]Command{}
	for _, cmd := range builtins {
		sh.AddCommand(cmd.name, cmd)
	}
}

func (sh *Shell) AddCommand(name string, cmd Command) {
	sh.commands[name] = cmd
}

func (sh *Shell) FindCommand(name string) Command {
	return sh.commands[name]
}

func (sh *Shell) Exec(script string) {
	prog, err := Parse(strings.NewReader(script))
	if err != nil {
		fmt.Fprintln(sh.Out, "ghost:", err.Error())
		return
	}
	env := &Environment{
		outer: sh.topLevel,
	}
	sh.Eval(env, prog)
}

func (sh *Shell) error(env *Environment, msg string) {
	fmt.Fprintln(sh.Out, "ghost:", msg)
	sh.status = 127
}

func (sh *Shell) Eval(env *Environment, node Node) {
	switch node := node.(type) {
	case *Program:
		sh.evalProgram(env, node)

	case *IfStmt:
		sh.evalIfStmt(env, node)

	case *BlockStmt:
		sh.evalBlockStmt(env, node)

	case *CommandStmt:
		sh.evalCommandStmt(env, node)

	case *BadStmt:
		sh.error(env, "bad statement")
	}
}

func (sh *Shell) evalProgram(env *Environment, prog *Program) {
	for _, stmt := range prog.Body {
		sh.Eval(env, stmt)
	}
}

func (sh *Shell) evalIfStmt(env *Environment, ifStmt *IfStmt) {
	sh.Eval(env, ifStmt.Cond)
	if sh.status == 0 {
		sh.Eval(env, ifStmt.Body)
	} else if ifStmt.Else != nil {
		sh.Eval(env, ifStmt.Else)
	}
}

func (sh *Shell) evalBlockStmt(env *Environment, blockStmt *BlockStmt) {
	for _, stmt := range blockStmt.List {
		sh.Eval(env, stmt)
	}
}

func (sh *Shell) evalCommandStmt(env *Environment, cmdStmt *CommandStmt) {
	sh.expandWord(env, cmdStmt.Command)
	command := sh.FindCommand(cmdStmt.Command.Value)
	if command == nil {
		sh.error(env, fmt.Sprintf("unknown command %q", cmdStmt.Command.Value))
		return
	}
	args := []string{cmdStmt.Command.Value}
	for _, arg := range cmdStmt.Args {
		sh.expandWord(env, arg)
		args = append(args, arg.Value)
	}
	sh.status = command.Run(sh, env, args)
}

func (sh *Shell) expandWord(env *Environment, word *Word) {
	needsExpand := strings.ContainsAny(word.Value, `$"'\`)
	if !needsExpand {
		return
	}

	switch word.Value[0] {
	case '\'':
		word.Value = word.Value[1 : len(word.Value)-1]
	case '"':
		word.Value = word.Value[1 : len(word.Value)-1]
		word.Value = expandDollar(env, word.Value)
	default:
		word.Value = expandEscape(word.Value)
		word.Value = expandDollar(env, word.Value)
	}
}
