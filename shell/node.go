package shell

type Node interface {
}

type Stmt interface {
	Node
}

type Program struct {
	Body []Stmt
}

type Word struct {
	Token *Token
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
	Cond *CommandStmt
	Body *BlockStmt
	Else Stmt
}
