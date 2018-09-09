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
					Start:   Position{Line: 1, Column: 1},
					End:     Position{Line: 1, Column: 5},
				},
				&Token{
					Kind:    STRING,
					Literal: "hello",
					Start:   Position{Line: 1, Column: 6},
					End:     Position{Line: 1, Column: 11},
				},
				&Token{
					Kind:    STRING,
					Literal: "world",
					Start:   Position{Line: 1, Column: 12},
					End:     Position{Line: 1, Column: 17},
				},
				&Token{
					Kind:    NEWLINE,
					Literal: "\r\n",
					Start:   Position{Line: 1, Column: 17},
					End:     Position{Line: 2, Column: 1},
				},
				&Token{
					Kind:    EOF,
					Literal: "",
					Start:   Position{Line: 2, Column: 1},
					End:     Position{Line: 2, Column: 1},
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