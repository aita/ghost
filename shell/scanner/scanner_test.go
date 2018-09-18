package scanner

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aita/ghost/shell/token"
)

func TestScanner(t *testing.T) {
	for _, tt := range []struct {
		input string
		toks  []*token.Token
	}{
		{
			"echo hello;\necho world",
			[]*token.Token{
				&token.Token{
					Kind:    token.STRING,
					Literal: "echo",
					Pos:     token.Position{Offset: 0, Line: 1, Column: 1},
				},
				&token.Token{
					Kind:    token.STRING,
					Literal: "hello",
					Pos:     token.Position{Offset: 5, Line: 1, Column: 6},
				},
				&token.Token{
					Kind:    token.TERMINATOR,
					Literal: ";",
					Pos:     token.Position{Offset: 10, Line: 1, Column: 11},
				},
				&token.Token{
					Kind:    token.STRING,
					Literal: "echo",
					Pos:     token.Position{Offset: 12, Line: 2, Column: 1},
				},
				&token.Token{
					Kind:    token.STRING,
					Literal: "world",
					Pos:     token.Position{Offset: 17, Line: 2, Column: 6},
				},
				&token.Token{
					Kind:    token.TERMINATOR,
					Literal: "",
					Pos:     token.Position{Offset: 22, Line: 2, Column: 11},
				},
				&token.Token{
					Kind:    token.EOF,
					Literal: "",
					Pos:     token.Position{Offset: 22, Line: 2, Column: 11},
				},
			},
		},
		{
			`if test 1; echo 'one'; else; echo "other"; end`,
			[]*token.Token{
				&token.Token{
					Kind:    token.STRING,
					Literal: "if",
					Pos:     token.Position{Offset: 0, Line: 1, Column: 1},
				},
				&token.Token{
					Kind:    token.STRING,
					Literal: "test",
					Pos:     token.Position{Offset: 3, Line: 1, Column: 4},
				},
				&token.Token{
					Kind:    token.STRING,
					Literal: "1",
					Pos:     token.Position{Offset: 8, Line: 1, Column: 9},
				},
				&token.Token{
					Kind:    token.TERMINATOR,
					Literal: ";",
					Pos:     token.Position{Offset: 9, Line: 1, Column: 10},
				},
				&token.Token{
					Kind:    token.STRING,
					Literal: "echo",
					Pos:     token.Position{Offset: 11, Line: 1, Column: 12},
				},
				&token.Token{
					Kind:    token.STRING,
					Literal: "'one'",
					Pos:     token.Position{Offset: 16, Line: 1, Column: 17},
				},
				&token.Token{
					Kind:    token.TERMINATOR,
					Literal: ";",
					Pos:     token.Position{Offset: 21, Line: 1, Column: 22},
				},
				&token.Token{
					Kind:    token.STRING,
					Literal: "else",
					Pos:     token.Position{Offset: 23, Line: 1, Column: 24},
				},
				&token.Token{
					Kind:    token.TERMINATOR,
					Literal: ";",
					Pos:     token.Position{Offset: 27, Line: 1, Column: 28},
				},
				&token.Token{
					Kind:    token.STRING,
					Literal: "echo",
					Pos:     token.Position{Offset: 29, Line: 1, Column: 30},
				},
				&token.Token{
					Kind:    token.STRING,
					Literal: `"other"`,
					Pos:     token.Position{Offset: 34, Line: 1, Column: 35},
				},
				&token.Token{
					Kind:    token.TERMINATOR,
					Literal: ";",
					Pos:     token.Position{Offset: 41, Line: 1, Column: 42},
				},
				&token.Token{
					Kind:    token.STRING,
					Literal: "end",
					Pos:     token.Position{Offset: 43, Line: 1, Column: 44},
				},
				&token.Token{
					Kind:    token.TERMINATOR,
					Literal: "",
					Pos:     token.Position{Offset: 46, Line: 1, Column: 47},
				},
				&token.Token{
					Kind:    token.EOF,
					Literal: "",
					Pos:     token.Position{Offset: 46, Line: 1, Column: 47},
				},
			},
		},
		{
			`echo "hello world"  # comment`,
			[]*token.Token{
				&token.Token{
					Kind:    token.STRING,
					Literal: "echo",
					Pos:     token.Position{Offset: 0, Line: 1, Column: 1},
				},
				&token.Token{
					Kind:    token.STRING,
					Literal: `"hello world"`,
					Pos:     token.Position{Offset: 5, Line: 1, Column: 6},
				},
				&token.Token{
					Kind:    token.TERMINATOR,
					Literal: "",
					Pos:     token.Position{Offset: 29, Line: 1, Column: 30},
				},
				&token.Token{
					Kind:    token.EOF,
					Literal: "",
					Pos:     token.Position{Offset: 29, Line: 1, Column: 30},
				},
			},
		},
		{
			`echo hello \
world`,
			[]*token.Token{
				&token.Token{
					Kind:    token.STRING,
					Literal: "echo",
					Pos:     token.Position{Offset: 0, Line: 1, Column: 1},
				},
				&token.Token{
					Kind:    token.STRING,
					Literal: "hello",
					Pos:     token.Position{Offset: 5, Line: 1, Column: 6},
				},
				&token.Token{
					Kind:    token.STRING,
					Literal: "world",
					Pos:     token.Position{Offset: 13, Line: 2, Column: 1},
				},
				&token.Token{
					Kind:    token.TERMINATOR,
					Literal: "",
					Pos:     token.Position{Offset: 18, Line: 2, Column: 6},
				},
				&token.Token{
					Kind:    token.EOF,
					Literal: "",
					Pos:     token.Position{Offset: 18, Line: 2, Column: 6},
				},
			},
		},
	} {
		r := strings.NewReader(tt.input)
		scanner := NewScanner(r, nil)
		for _, expected := range tt.toks {
			tok := scanner.Scan()
			assert.Equal(t, expected, tok)
		}
	}
}

func TestReadString(t *testing.T) {
	for _, tt := range []struct {
		input    string
		expected string
	}{
		{`hello'`, `hello`},
		{`hello;`, `hello`},
		{"hello\n", `hello`},

		// tests with escape sequence
		{`hello\'`, `hello\'`},
		{`hello\n`, `hello\n`},
		{`hello\ world`, `hello\ world`},
	} {
		r := strings.NewReader(tt.input)
		scanner := NewScanner(r, nil)
		scanner.next()
		lit := scanner.scanString("")
		assert.Equal(t, tt.expected, lit)
	}
}

func TestReadQuotedString(t *testing.T) {
	type Error struct {
		tok token.Position
		msg string
	}

	for _, tt := range []struct {
		input    string
		expected string
		err      *Error
	}{
		{`"hello"`, `"hello"`, nil},
		{`'hello'`, `'hello'`, nil},
		{`"hello world\n"`, `"hello world\n"`, nil},
		{`'hello world\n'`, `'hello world\n'`, nil},
		{`"\"double quote\""`, `""double quote""`, nil},
		{`'It\'s a small world'`, `'It's a small world'`, nil},
		{
			`"hello`,
			`"hello`,
			&Error{
				token.Position{Offset: 6, Line: 1, Column: 7},
				"unexpected end of string",
			},
		},
		{
			`"hello\`,
			`"hello`,
			&Error{
				token.Position{Offset: 7, Line: 1, Column: 8},
				"unexpected end of string",
			},
		},
	} {
		r := strings.NewReader(tt.input)
		var err *Error
		scanner := NewScanner(r, func(pos token.Position, msg string) {
			err = &Error{pos, msg}
		})
		scanner.next()
		lit := scanner.scanQuotedString()
		assert.Equal(t, tt.expected, lit)
		assert.Equal(t, tt.err, err)
	}
}
