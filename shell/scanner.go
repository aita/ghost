package shell

import (
	"errors"
	"io"
	"unicode"
)

type Position struct {
	Line   int // line number, starting at 1
	Column int // column number, starting at 1 (character count per line)
}

const (
	EOF = iota
	NEWLINE
	COMMENT

	STRING
	SEMICOLON
)

var (
	ErrUnterminatedQuote = errors.New("unexpected end of string")
)

type Token struct {
	Kind    int
	Literal string
	Start   Position
	End     Position
}

func newToken(kind int, lit string, start, end Position) *Token {
	return &Token{
		Kind:    kind,
		Literal: lit,
		Start:   start,
		End:     end,
	}
}

type Scanner struct {
	r   io.RuneScanner
	pos Position
}

func NewScanner(r io.RuneScanner) *Scanner {
	return &Scanner{
		r: r,
		pos: Position{
			Line:   1,
			Column: 1,
		},
	}
}

func (scanner *Scanner) Next() (*Token, error) {
	for {
		var ch rune
		start := scanner.pos
		ch, _, err := scanner.readRune()
		if err == io.EOF {
			return newToken(EOF, "", start, scanner.pos), nil
		} else if err != nil {
			return nil, err
		}
		if unicode.IsSpace(ch) {
			switch ch {
			case '\r':
				scanner.pos.Line++
				scanner.pos.Column = 0
				lit := string(ch)
				ch, _, err := scanner.readRune()
				if err != io.EOF {
					if err != nil {
						return nil, err
					}
					if ch == '\n' {
						lit += string(ch)
					} else {
						err := scanner.unreadRune()
						if err != nil {
							return nil, err
						}
					}
				}
				return newToken(NEWLINE, lit, start, scanner.pos), nil
			case '\n':
				scanner.pos.Line++
				scanner.pos.Column = 0
				return newToken(NEWLINE, "\n", start, scanner.pos), nil
			}
		} else {
			switch ch {
			case ';':
				return newToken(SEMICOLON, ";", start, scanner.pos), nil
			case '\'', '"':
				lit, err := scanner.readQuotedString(ch)
				if err != nil {
					return nil, err
				}
				return newToken(STRING, lit, start, scanner.pos), nil
			default:
				lit, err := scanner.readString(ch)
				if err != nil {
					return nil, err
				}
				return newToken(STRING, lit, start, scanner.pos), nil
			}
		}
	}
}

func (scanner *Scanner) readRune() (r rune, size int, err error) {
	r, size, err = scanner.r.ReadRune()
	if err != nil {
		return
	}
	scanner.pos.Column++
	return
}

func (scanner *Scanner) unreadRune() error {
	err := scanner.r.UnreadRune()
	if err != nil {
		return err
	}
	scanner.pos.Column--
	return nil
}

var escapes = map[rune]rune{
	'a': '\a',
	'b': '\b',
	'f': '\f',
	'n': '\n',
	'r': '\r',
	't': '\t',
	'v': '\v',
}

func (scanner *Scanner) readString(first rune) (string, error) {
	lit := string(first)
	for {
		ch, _, err := scanner.readRune()
		if err == io.EOF {
			return lit, nil
		} else if err != nil {
			return "", err
		}
		if unicode.IsSpace(ch) {
			err = scanner.unreadRune()
			if err != nil {
				return "", err
			}
			return lit, nil
		}
		if ch == '\\' {
			ch, _, err = scanner.readRune()
			if err == io.EOF {
				return lit, nil
			} else if err != nil {
				return "", err
			}
			escape, ok := escapes[ch]
			if ok {
				ch = escape
			}
		}
		lit += string(ch)
	}
}

func (scanner *Scanner) readQuotedString(quote rune) (string, error) {
	lit := string(quote)
	for {
		ch, _, err := scanner.readRune()
		if err == io.EOF {
			return "", ErrUnterminatedQuote
		} else if err != nil {
			return "", err
		}
		switch ch {
		case '\'', '"':
			if ch == quote {
				lit += string(ch)
				return lit, nil
			}
		case '\\':
			ch, _, err = scanner.readRune()
			if err == io.EOF {
				return lit, nil
			} else if err != nil {
				return "", err
			}
			if quote == '"' {
				escape, ok := escapes[ch]
				if ok {
					ch = escape
				}
			} else if quote == '\'' && ch != '\'' {
				lit += "\\"
			}
		}
		lit += string(ch)
	}
}
