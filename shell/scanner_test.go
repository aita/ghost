package shell

import (
	"fmt"
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
		{
			`echo "'hello'\n"; echo '\'ghost\''`,
			[]*Token{
				&Token{
					Kind:    STRING,
					Literal: "echo",
					Start:   Position{Line: 1, Column: 1},
					End:     Position{Line: 1, Column: 5},
				},
				&Token{
					Kind:    STRING,
					Literal: "\"'hello'\n\"",
					Start:   Position{Line: 1, Column: 6},
					End:     Position{Line: 1, Column: 17},
				},
				&Token{
					Kind:    SEMICOLON,
					Literal: ";",
					Start:   Position{Line: 1, Column: 17},
					End:     Position{Line: 1, Column: 18},
				},
				&Token{
					Kind:    STRING,
					Literal: "echo",
					Start:   Position{Line: 1, Column: 19},
					End:     Position{Line: 1, Column: 23},
				},
				&Token{
					Kind:    STRING,
					Literal: "''ghost''",
					Start:   Position{Line: 1, Column: 24},
					End:     Position{Line: 1, Column: 35},
				},
				&Token{
					Kind:    EOF,
					Literal: "",
					Start:   Position{Line: 1, Column: 35},
					End:     Position{Line: 1, Column: 35},
				},
			},
		},
		{
			"echo hello  # put comment here",
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
					Kind:    COMMENT,
					Literal: "# put comment here",
					Start:   Position{Line: 1, Column: 13},
					End:     Position{Line: 1, Column: 31},
				},
				&Token{
					Kind:    EOF,
					Literal: "",
					Start:   Position{Line: 1, Column: 31},
					End:     Position{Line: 1, Column: 31},
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
		// tests with trailing baclash and nothing
		{"hello\\", "hello", ErrUnterminatedString},
	} {
		r := strings.NewReader(tt.input)
		scanner := NewScanner(r)
		ch, _, _ := r.ReadRune()
		err := scanner.readString(ch)
		assert.Equal(t, tt.expected, scanner.b.String())
		assert.Equal(t, tt.err, err)
	}
}

func TestReadQuotedString(t *testing.T) {
	for _, tt := range []struct {
		quote    rune
		s        string
		expected string
		err      error
	}{
		// single-quote string
		// tests with escape sequence
		{'\'', "hello\\'", "'hello''", nil},
		{'\'', "hello\\n", "'hello\n'", nil},
		{'\'', "hello\\ world", "hello world", nil},
		// tests with trailing backslash and newline
		{'\'', "hello\\\r", "hello", nil},
		{'\'', "hello\\\r\n", "hello", nil},
		{'\'', "hello\\\n", "hello", nil},
		// tests with trailing backslash and string
		{'\'', "hello\\\nworld", "helloworld", nil},
		{'\'', "hello\\\rworld", "helloworld", nil},
		{'\'', "hello\\\r\nworld", "helloworld", nil},
		// tests with trailing baclash and nothing
		{'\'', "hello\\", "hello", ErrUnterminatedString},

		// double-quote string
		// TODO
	} {
		r := strings.NewReader(fmt.Sprintf("%s%c", tt.s, tt.quote))
		scanner := NewScanner(r)
		err := scanner.readQuotedString(tt.quote)
		assert.Equal(t, tt.expected, scanner.b.String())
		assert.Equal(t, tt.err, err)
	}
}
