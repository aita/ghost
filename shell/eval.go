package shell

import (
	"fmt"
	"io"
	"strconv"
)

type Environment struct {
	store map[string]string
	outer *Environment
}

func (e *Environment) Get(name string) (string, bool) {
	s, ok := e.store[name]
	if !ok && e.outer != nil {
		s, ok = e.outer.Get(name)
	}
	return s, ok
}

func (e *Environment) Set(name, val string) {
	e.store[name] = val
}

type Evaluator struct {
	Stdin    io.Reader
	Stdout   io.Writer
	Stderr   io.Writer
	topLevel *Environment
	Commands map[string]Command
}

func NewEvaluator(stdin io.Reader, stdout, stderr io.Writer) *Evaluator {
	env := &Environment{
		store: map[string]string{
			"?": "0",
		},
		outer: nil,
	}
	return &Evaluator{
		Stdin:    stdin,
		Stdout:   stdout,
		Stderr:   stderr,
		topLevel: env,
		Commands: builtins,
	}
}

func (e *Evaluator) error(msg string) {
	fmt.Fprintln(e.Stderr, "fish:", msg)
}

func (e *Evaluator) Eval(node Node) {
	e.eval(e.topLevel, node)
}

func (e *Evaluator) eval(env *Environment, node Node) {
	switch node := node.(type) {
	case *Program:
		e.evalProgram(env, node)

	case *IfStmt:
		e.evalIfStmt(env, node)

	case *BlockStmt:
		e.evalBlockStmt(env, node)

	case *CommandStmt:
		e.evalCommandStmt(env, node)
	}
}

func (e *Evaluator) evalProgram(env *Environment, prog *Program) {
	e.eval(env, prog.Body)
}

func (e *Evaluator) evalIfStmt(env *Environment, ifStmt *IfStmt) {
	e.eval(env, ifStmt.Cond)
	if e.gettStatus() == 0 {
		e.eval(env, ifStmt.Body)
	} else if ifStmt.Else != nil {
		e.eval(env, ifStmt.Else)
	}
}

func (e *Evaluator) evalBlockStmt(env *Environment, blockStmt *BlockStmt) {
	for _, stmt := range blockStmt.List {
		e.eval(env, stmt)
	}
}

func (e *Evaluator) evalCommandStmt(env *Environment, cmdStmt *CommandStmt) {
	name := cmdStmt.Command.Value
	command := e.findCommand(name)
	if command == nil {
		e.error(fmt.Sprintf("Unknown command %q", name))
		return
	}
	args := []string{}
	for _, arg := range cmdStmt.Args {
		args = append(args, arg.Value)
	}
	status := command.Run(e.Stdin, e.Stdout, e.Stderr, args)
	e.setStatus(status)
}

func (e *Evaluator) gettStatus() int {
	s, _ := e.topLevel.Get("?")
	status, _ := strconv.Atoi(s)
	return status
}

func (e *Evaluator) setStatus(status int) {
	e.topLevel.Set("?", strconv.Itoa(status))
}

func (e *Evaluator) findCommand(name string) Command {
	return e.Commands[name]
}
