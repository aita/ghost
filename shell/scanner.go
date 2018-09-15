package shell

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

type Position struct {
	Offset int // offset, starting at 0
	Line   int // line number, starting at 1
	Column int // column number, starting at 1 (character count per line)
}

// The list of kinds of token
const (
	EOF = -1

	NEWLINE = iota
	COMMENT
	STRING
	SEMICOLON
)

var TokenNames = map[int]string{
	EOF:       "EOF",
	NEWLINE:   "NEWLINE",
	COMMENT:   "COMMENT",
	STRING:    "STRING",
	SEMICOLON: "SEMICOLON",
}

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

type ErrorHandler func(pos Position, msg string)

type Scanner struct {
	src        *bufio.Reader
	ch         rune
	newline    bool
	pos        Position
	sb         strings.Builder
	errHandler ErrorHandler
	ErrorCount int
}

func NewScanner(r io.Reader, errHandler ErrorHandler) *Scanner {
	return &Scanner{
		src:     bufio.NewReader(r),
		ch:      ' ',
		newline: false,
		pos: Position{
			Line:   1,
			Column: 0,
		},
		errHandler: errHandler,
		ErrorCount: 0,
	}
}

func (s *Scanner) error(message string) {
	if s.errHandler != nil {
		s.errHandler(s.pos, message)
	}
	s.ErrorCount++
}

func (s *Scanner) next() {
	ch, size, err := s.src.ReadRune()
	if err != nil {
		if err == io.EOF {
			ch = EOF
		} else {
			s.error(err.Error())
		}
	}
	s.pos.Offset += size
	if size > 0 {
		s.pos.Column++
	}
	if s.ch == '\n' {
		s.pos.Line++
		s.pos.Column = 1
	}
	s.ch = ch
}

func (s *Scanner) Scan() *Token {
scanAgain:
	pos := s.pos
	for {
		if s.ch == EOF {
			return newToken(EOF, "", pos)
		}
		if !unicode.IsSpace(s.ch) {
			break
		}
		if s.ch == '\n' && !s.newline {
			s.newline = true
			return newToken(NEWLINE, "\n", pos)
		}
		s.next()
		pos = s.pos
	}
	s.newline = false

	switch s.ch {
	case ';':
		s.next()
		return newToken(SEMICOLON, ";", pos)
	case '#':
		lit := s.scanComment()
		return newToken(COMMENT, lit, pos)
	case '\'', '"':
		lit := s.scanQuotedString()
		return newToken(STRING, lit, pos)
	default:
		lit := s.scanString()
		if lit == "\n" {
			goto scanAgain
		}
		return newToken(STRING, lit, pos)
	}
}

func (s *Scanner) scanComment() string {
	s.sb.Reset()
	for {
		if s.ch == EOF || s.ch == '\r' || s.ch == '\n' {
			break
		}
		s.sb.WriteRune(s.ch)
		s.next()
	}
	return s.sb.String()
}

func (s *Scanner) scanString() string {
	s.sb.Reset()
scanEnd:
	for {
		switch s.ch {
		case EOF, ';', '\'', '"':
			break scanEnd
		case '\\':
			s.sb.WriteRune(s.ch)
			s.next()
			if s.ch == EOF {
				s.error("unexpected end of string")
				break scanEnd
			}
		default:
			if unicode.IsSpace(s.ch) {
				break scanEnd
			}
		}
		s.sb.WriteRune(s.ch)
		s.next()
	}
	return s.sb.String()
}

func (s *Scanner) scanQuotedString() string {
	quote := s.ch

	s.sb.Reset()
	s.sb.WriteRune(quote)
	for {
		s.next()
		if s.ch == EOF {
			s.error("unexpected end of string")
			break
		} else if s.ch == '\'' || s.ch == '"' {
			if s.ch == quote {
				s.sb.WriteRune(s.ch)
				break
			}
		} else if s.ch == '\\' {
			s.next()
			if s.ch == EOF {
				s.error("unexpected end of string")
				break
			}
			if s.ch != quote {
				s.sb.WriteRune('\\')
			}
		}
		s.sb.WriteRune(s.ch)
	}
	s.next()
	return s.sb.String()
}
