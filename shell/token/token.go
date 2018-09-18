package token

type Position struct {
	Offset int // offset, starting at 0
	Line   int // line number, starting at 1
	Column int // column number, starting at 1 (character count per line)
}

type TokenKind int

// The list of kinds of token
const (
	EOF    = -1
	STRING = iota
	TERMINATOR
)

var tokens = map[TokenKind]string{
	EOF:        "EOF",
	STRING:     "STRING",
	TERMINATOR: "TERMINATOR",
}

func (kind TokenKind) String() string {
	name, ok := tokens[kind]
	if !ok {
		return "UNKNOWN"
	}
	return name
}

type Token struct {
	Kind    TokenKind
	Literal string
	Pos     Position
}
