package shell

type Node interface {
}

type Stmt interface {
	Node
}

type Program struct {
	Body *BlockStmt
}

type Word struct {
	Token *Token
	Value string
}

type Command struct {
	Command *Word
	Args    []*Word
}

type BlockStmt struct {
	List []Stmt
}

type IfStmt struct {
	Cond *Command
	Body *BlockStmt
	Else Stmt
}
