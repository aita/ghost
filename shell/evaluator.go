package shell

import (
	"fmt"
)

type Command interface {
	Run(env *Environment, args []string) int
}

type Evaluator struct {
	Commands map[string]Command
}

func (e *Evaluator) error(env *Environment, msg string) {
	fmt.Fprintln(env.Stdout(), "ghost:", msg)
	env.SetStatus(127)
}

func (e *Evaluator) Eval(env *Environment, node Node) {
	switch node := node.(type) {
	case *Program:
		e.evalProgram(env, node)

	case *IfStmt:
		e.evalIfStmt(env, node)

	case *BlockStmt:
		e.evalBlockStmt(env, node)

	case *CommandStmt:
		e.evalCommandStmt(env, node)

	case *BadStmt:
		e.error(env, "bad statement")
	}
}

func (e *Evaluator) evalProgram(env *Environment, prog *Program) {
	for _, stmt := range prog.Body {
		e.Eval(env, stmt)
	}
}

func (e *Evaluator) evalIfStmt(env *Environment, ifStmt *IfStmt) {
	e.Eval(env, ifStmt.Cond)
	if env.GetStatus() == 0 {
		e.Eval(env, ifStmt.Body)
	} else if ifStmt.Else != nil {
		e.Eval(env, ifStmt.Else)
	}
}

func (e *Evaluator) evalBlockStmt(env *Environment, blockStmt *BlockStmt) {
	for _, stmt := range blockStmt.List {
		e.Eval(env, stmt)
	}
}

func (e *Evaluator) evalCommandStmt(env *Environment, cmdStmt *CommandStmt) {
	name := cmdStmt.Command.Value
	command := e.findCommand(name)
	if command == nil {
		e.error(env, fmt.Sprintf("unknown command %q", name))
		return
	}
	args := []string{}
	for _, arg := range cmdStmt.Args {
		args = append(args, arg.Value)
	}
	status := command.Run(env, args)
	env.SetStatus(status)
}

func (e *Evaluator) findCommand(name string) Command {
	return e.Commands[name]
}
