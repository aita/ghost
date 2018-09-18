package shell

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanner(t *testing.T) {
	for _, tt := range []struct {
		input string
		toks  []*Token
	}{
		{
			"echo hello;\necho world",
			[]*Token{
				&Token{
					Kind:    STRING,
					Literal: "echo",
					Pos:     Position{Offset: 0, Line: 1, Column: 1},
				},
				&Token{
					Kind:    STRING,
					Literal: "hello",
					Pos:     Position{Offset: 5, Line: 1, Column: 6},
				},
				&Token{
					Kind:    TERMINATOR,
					Literal: ";",
					Pos:     Position{Offset: 10, Line: 1, Column: 11},
				},
				&Token{
					Kind:    STRING,
					Literal: "echo",
					Pos:     Position{Offset: 12, Line: 2, Column: 1},
				},
				&Token{
					Kind:    STRING,
					Literal: "world",
					Pos:     Position{Offset: 17, Line: 2, Column: 6},
				},
				&Token{
					Kind:    TERMINATOR,
					Literal: "",
					Pos:     Position{Offset: 22, Line: 2, Column: 11},
				},
				&Token{
					Kind:    EOF,
					Literal: "",
					Pos:     Position{Offset: 22, Line: 2, Column: 11},
				},
			},
		},
		{
			`if test 1; echo 'one'; else; echo "other"; end`,
			[]*Token{
				&Token{
					Kind:    STRING,
					Literal: "if",
					Pos:     Position{Offset: 0, Line: 1, Column: 1},
				},
				&Token{
					Kind:    STRING,
					Literal: "test",
					Pos:     Position{Offset: 3, Line: 1, Column: 4},
				},
				&Token{
					Kind:    STRING,
					Literal: "1",
					Pos:     Position{Offset: 8, Line: 1, Column: 9},
				},
				&Token{
					Kind:    TERMINATOR,
					Literal: ";",
					Pos:     Position{Offset: 9, Line: 1, Column: 10},
				},
				&Token{
					Kind:    STRING,
					Literal: "echo",
					Pos:     Position{Offset: 11, Line: 1, Column: 12},
				},
				&Token{
					Kind:    STRING,
					Literal: "'one'",
					Pos:     Position{Offset: 16, Line: 1, Column: 17},
				},
				&Token{
					Kind:    TERMINATOR,
					Literal: ";",
					Pos:     Position{Offset: 21, Line: 1, Column: 22},
				},
				&Token{
					Kind:    STRING,
					Literal: "else",
					Pos:     Position{Offset: 23, Line: 1, Column: 24},
				},
				&Token{
					Kind:    STRING,
					Literal: "echo",
					Pos:     Position{Offset: 11, Line: 1, Column: 12},
				},
				&Token{
					Kind:    STRING,
					Literal: `"other"`,
					Pos:     Position{Offset: 16, Line: 1, Column: 17},
				},
				&Token{
					Kind:    TERMINATOR,
					Literal: ";",
					Pos:     Position{Offset: 21, Line: 1, Column: 22},
				},
				&Token{
					Kind:    STRING,
					Literal: "end",
					Pos:     Position{Offset: 25, Line: 1, Column: 26},
				},
				&Token{
					Kind:    TERMINATOR,
					Literal: "",
					Pos:     Position{Offset: 28, Line: 1, Column: 29},
				},
				&Token{
					Kind:    EOF,
					Literal: "",
					Pos:     Position{Offset: 28, Line: 1, Column: 29},
				},
			},
		},
		{
			`echo "hello world"  # comment`,
			[]*Token{
				&Token{
					Kind:    STRING,
					Literal: "echo",
					Pos:     Position{Offset: 0, Line: 1, Column: 1},
				},
				&Token{
					Kind:    STRING,
					Literal: `"hello world"`,
					Pos:     Position{Offset: 5, Line: 1, Column: 6},
				},
				&Token{
					Kind:    TERMINATOR,
					Literal: "",
					Pos:     Position{Offset: 29, Line: 1, Column: 30},
				},
				&Token{
					Kind:    EOF,
					Literal: "",
					Pos:     Position{Offset: 29, Line: 1, Column: 30},
				},
			},
		},
		{
			`echo hello \
world`,
			[]*Token{
				&Token{
					Kind:    STRING,
					Literal: "echo",
					Pos:     Position{Offset: 0, Line: 1, Column: 1},
				},
				&Token{
					Kind:    STRING,
					Literal: "hello",
					Pos:     Position{Offset: 5, Line: 1, Column: 6},
				},
				&Token{
					Kind:    STRING,
					Literal: "world",
					Pos:     Position{Offset: 13, Line: 2, Column: 1},
				},
				&Token{
					Kind:    TERMINATOR,
					Literal: "",
					Pos:     Position{Offset: 18, Line: 2, Column: 6},
				},
				&Token{
					Kind:    EOF,
					Literal: "",
					Pos:     Position{Offset: 18, Line: 2, Column: 6},
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
		tok Position
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
				Position{Offset: 6, Line: 1, Column: 7},
				"unexpected end of string",
			},
		},
		{
			`"hello\`,
			`"hello`,
			&Error{
				Position{Offset: 7, Line: 1, Column: 8},
				"unexpected end of string",
			},
		},
	} {
		r := strings.NewReader(tt.input)
		var err *Error
		scanner := NewScanner(r, func(pos Position, msg string) {
			err = &Error{pos, msg}
		})
		scanner.next()
		lit := scanner.scanQuotedString()
		assert.Equal(t, tt.expected, lit)
		assert.Equal(t, tt.err, err)
	}
}
