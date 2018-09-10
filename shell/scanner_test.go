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
					Literal: "\"'hello'\\n\"",
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
		r := strings.NewReader(tt.input)
		scanner := NewScanner(r)
		for _, expected := range tt.toks {
			tok, err := scanner.Next()
			assert.Nil(t, err)
			assert.Equal(t, expected, tok)
		}
	}
}

func TestReadString(t *testing.T) {
	for _, tt := range []struct {
		input    string
		expected string
		err      error
	}{
		// tests with escape sequence
		{"hello\\'", "hello'", nil},
		{"hello\\n", "hello\n", nil},
		{"hello\\ world", "hello world", nil},
		// tests with trailing backslash and newline
		{"hello\\\r", "hello", nil},
		{"hello\\\r\n", "hello", nil},
		{"hello\\\n", "hello", nil},
		// tests with trailing backslash and string
		{"hello\\\nworld", "helloworld", nil},
		{"hello\\\rworld", "helloworld", nil},
		{"hello\\\r\nworld", "helloworld", nil},
		// tests with trailing backslash and nothing
		{"hello\\", "hello", ErrUnterminatedString},
	} {
		r := strings.NewReader(tt.input)
		scanner := NewScanner(r)
		ch, _ := scanner.readRune()
		err := scanner.readString(ch)
		assert.Equal(t, tt.expected, scanner.b.String())
		assert.Equal(t, tt.err, err)
	}
}

func TestReadQuotedString(t *testing.T) {
	for _, tt := range []struct {
		input    string
		expected string
		err      error
	}{
		{`"hello"`, `"hello"`, nil},
		{`'hello'`, `'hello'`, nil},
		{`"hello world\n"`, `"hello world\n"`, nil},
		{`'hello world\n'`, `'hello world\n'`, nil},
		{`"\"double quote\""`, `""double quote""`, nil},
		{`'It\'s a small world'`, `'It's a small world'`, nil},
		{`"hello`, `"hello`, ErrUnterminatedString},
		{`"hello\`, `"hello`, ErrUnterminatedString},
	} {
		r := strings.NewReader(tt.input)
		scanner := NewScanner(r)
		q, _ := scanner.readRune()
		err := scanner.readQuotedString(q)
		assert.Equal(t, tt.expected, scanner.b.String())
		assert.Equal(t, tt.err, err)
	}
}
