package shell

import (
	"fmt"
)

type Environment struct {
	store map[string]string
}

func (e *Environment) Get(name string) (string, bool) {
	s, ok := e.store[name]
	return s, ok
}

func (e *Environment) Set(name, val string) {
	e.store[name] = val
}

func Eval(env *Environment, node Node) {
	switch node := node.(type) {
	case *Program:
		evalProgram(env, node)

	case *BlockStmt:
		evalBlockStmt(env, node)

	case *Command:
		evalCommand(env, node)
	}
}

func evalProgram(env *Environment, prog *Program) {
	Eval(env, prog.Body)
}

func evalBlockStmt(env *Environment, block *BlockStmt) {
	for _, stmt := range block.List {
		Eval(env, stmt)
	}
}

func evalCommand(env *Environment, cmd *Command) {
	fmt.Printf("%#v\n", cmd)
}
