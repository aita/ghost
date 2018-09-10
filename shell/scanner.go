package shell

import (
	"errors"
	"io"
	"strings"
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
	ErrUnterminatedString = errors.New("unexpected end of string")
)

type Token struct {
	Kind    int
	Literal string
	Pos     Position
}

func newToken(kind int, lit string, pos Position) *Token {
	return &Token{
		Kind:    kind,
		Literal: lit,
		Pos:     pos,
	}
}

type Scanner struct {
	r   io.RuneScanner
	pos Position
	b   strings.Builder
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
		start := scanner.pos
		ch, _, err := scanner.readRune()
		if err == io.EOF {
			return newToken(EOF, "", start), nil
		} else if err != nil {
			return nil, err
		}
		if unicode.IsSpace(ch) {
			switch ch {
			case '\r':
				scanner.pos.Line++
				scanner.pos.Column = 0
				lit := "\r"
				ch, _, err := scanner.readRune()
				if err != io.EOF {
					if err != nil {
						return nil, err
					}
					if ch == '\n' {
						lit = "\r\n"
					} else {
						err := scanner.unreadRune()
						if err != nil {
							return nil, err
						}
					}
				}
				return newToken(NEWLINE, lit, start), nil
			case '\n':
				scanner.pos.Line++
				scanner.pos.Column = 0
				return newToken(NEWLINE, "\n", start), nil
			}
		} else {
			switch ch {
			case ';':
				return newToken(SEMICOLON, ";", start), nil
			case '#':
				scanner.b.Reset()
				scanner.b.WriteRune(ch)
				for {
					ch, _, err = scanner.readRune()
					if err == io.EOF {
						break
					} else if err != nil {
						return nil, err
					}
					if ch == '\r' || ch == '\n' {
						err := scanner.unreadRune()
						if err != nil {
							return nil, err
						}
					}
					scanner.b.WriteRune(ch)
				}
				return newToken(COMMENT, scanner.b.String(), start), nil
			case '\'', '"':
				err := scanner.readQuotedString(ch)
				if err != nil {
					return nil, err
				}
				return newToken(STRING, scanner.b.String(), start), nil
			default:
				err := scanner.readString(ch)
				if err != nil {
					return nil, err
				}
				return newToken(STRING, scanner.b.String(), start), nil
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

func (scanner *Scanner) readString(first rune) error {
	scanner.b.Reset()
	scanner.b.WriteRune(first)
	for {
		ch, _, err := scanner.readRune()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		if unicode.IsSpace(ch) {
			return scanner.unreadRune()
		}
		if ch == '\\' {
			ch, _, err = scanner.readRune()
			if err == io.EOF {
				return ErrUnterminatedString
			} else if err != nil {
				return err
			}
			// ignore "\r\n" and "\n"
			if ch == '\r' {
				ch, _, err = scanner.readRune()
				if err == io.EOF {
					return nil
				} else if err != nil {
					return err
				}
				if ch == '\n' {
					continue
				}
			} else if ch == '\n' {
				continue
			} else {
				escape, ok := escapes[ch]
				if ok {
					ch = escape
				}
			}
		}
		scanner.b.WriteRune(ch)
	}
}

func (scanner *Scanner) readQuotedString(quote rune) error {
	scanner.b.Reset()
	scanner.b.WriteRune(quote)
	for {
		ch, _, err := scanner.readRune()
		if err == io.EOF {
			return ErrUnterminatedString
		} else if err != nil {
			return err
		}
		switch ch {
		case '\'', '"':
			if ch == quote {
				scanner.b.WriteRune(ch)
				return nil
			}
		case '\\':
			ch, _, err = scanner.readRune()
			if err == io.EOF {
				return ErrUnterminatedString
			} else if err != nil {
				return err
			}
			if quote == '"' {
				escape, ok := escapes[ch]
				if ok {
					ch = escape
				}
			} else if quote == '\'' && ch != '\'' {
				scanner.b.WriteRune('\\')
			}
		}
		scanner.b.WriteRune(ch)
	}
}
