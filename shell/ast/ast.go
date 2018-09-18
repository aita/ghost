package ast

import "github.com/aita/ghost/shell/token"

type Node interface {
}

type Stmt interface {
	Node
}

type Program struct {
	Body []Stmt
}

type Word struct {
	Token *token.Token
	Value string
}

type CommandStmt struct {
	Command *Word
	Args    []*Word
}

type BlockStmt struct {
	List []Stmt
}

type IfStmt struct {
	Cond Stmt
	Body *BlockStmt
	Else Stmt
}

type BadStmt struct {
}
