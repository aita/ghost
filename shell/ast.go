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
	List []*Word
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
