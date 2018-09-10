package shell

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanner(t *testing.T) {
	for _, x := range []struct {
		input string
		toks  []*Token
	}{
		{
			"echo hello world\r\n",
			[]*Token{
				&Token{
					Kind:    STRING,
					Literal: "echo",
					Pos:     Position{Line: 1, Column: 1},
				},
				&Token{
					Kind:    STRING,
					Literal: "hello",
					Pos:     Position{Line: 1, Column: 6},
				},
				&Token{
					Kind:    STRING,
					Literal: "world",
					Pos:     Position{Line: 1, Column: 12},
				},
				&Token{
					Kind:    NEWLINE,
					Literal: "\r\n",
					Pos:     Position{Line: 1, Column: 17},
				},
				&Token{
					Kind:    EOF,
					Literal: "",
					Pos:     Position{Line: 2, Column: 1},
				},
			},
		},
		{
			`echo "'hello'\n"; echo '\'ghost\''`,
			[]*Token{
				&Token{
					Kind:    STRING,
					Literal: "echo",
					Pos:     Position{Line: 1, Column: 1},
				},
				&Token{
					Kind:    STRING,
					Literal: "\"'hello'\n\"",
					Pos:     Position{Line: 1, Column: 6},
				},
				&Token{
					Kind:    SEMICOLON,
					Literal: ";",
					Pos:     Position{Line: 1, Column: 17},
				},
				&Token{
					Kind:    STRING,
					Literal: "echo",
					Pos:     Position{Line: 1, Column: 19},
				},
				&Token{
					Kind:    STRING,
					Literal: "''ghost''",
					Pos:     Position{Line: 1, Column: 24},
				},
				&Token{
					Kind:    EOF,
					Literal: "",
					Pos:     Position{Line: 1, Column: 35},
				},
			},
		},
		{
			"echo hello  # put comment here",
			[]*Token{
				&Token{
					Kind:    STRING,
					Literal: "echo",
					Pos:     Position{Line: 1, Column: 1},
				},
				&Token{
					Kind:    STRING,
					Literal: "hello",
					Pos:     Position{Line: 1, Column: 6},
				},
				&Token{
					Kind:    COMMENT,
					Literal: "# put comment here",
					Pos:     Position{Line: 1, Column: 13},
				},
				&Token{
					Kind:    EOF,
					Literal: "",
					Pos:     Position{Line: 1, Column: 31},
				},
			},
		},
	} {
		r := strings.NewReader(x.input)
		scanner := NewScanner(r)
		for _, expected := range x.toks {
			tok, err := scanner.Next()
			assert.Nil(t, err)
			assert.Equal(t, expected, tok)
		}
	}
}
