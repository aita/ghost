package shell

type Node interface {
}

type Program struct {
	Body []Node
}

type WordNode struct {
	Token *Token
	Value string
}

type CommandNode struct {
	List []*WordNode
}

type BlockNode struct {
	List []Node
}

type IfNode struct {
	Cond Node
	Body *BlockNode
	Else Node
}

type BadNode struct {
}
