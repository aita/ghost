package shell

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

func newToken(kind TokenKind, lit string, pos Position) *Token {
	return &Token{
		Kind:    kind,
		Literal: lit,
		Pos:     pos,
	}
}

type ErrorHandler func(pos Position, msg string)

type Scanner struct {
	ErrorCount       int
	errHandler       ErrorHandler
	src              *bufio.Reader
	ch               rune
	insertTerminator bool
	lastSize         int
	pos              Position
}

func NewScanner(r io.Reader, errHandler ErrorHandler) *Scanner {
	return &Scanner{
		ErrorCount:       0,
		errHandler:       errHandler,
		src:              bufio.NewReader(r),
		ch:               ' ',
		insertTerminator: false,
		pos: Position{
			Line:   1,
			Column: 1,
		},
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
	if s.lastSize > 0 {
		s.pos.Column++
	}
	s.pos.Offset += s.lastSize
	s.lastSize = size
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
		if !unicode.IsSpace(s.ch) {
			break
		}
		if s.ch == '\n' && s.insertTerminator {
			s.insertTerminator = false
			return newToken(TERMINATOR, "\n", pos)
		}
		s.next()
		pos = s.pos
	}

	switch s.ch {
	case EOF:
		if s.insertTerminator {
			s.insertTerminator = false
			return newToken(TERMINATOR, "", pos)
		}
		return newToken(EOF, "", pos)
	case ';':
		s.insertTerminator = false
		s.next()
		return newToken(TERMINATOR, ";", pos)
	case '#':
		s.skipComment()
		goto scanAgain
	case '\'', '"':
		s.insertTerminator = true
		lit := s.scanQuotedString()
		return newToken(STRING, lit, pos)
	default:
		head := ""
		if s.ch == '\\' {
			s.next()
			if s.ch == '\n' {
				s.next()
				goto scanAgain
			}
			head = "\\"
		}

		s.insertTerminator = true
		lit := s.scanString(head)
		return newToken(STRING, lit, pos)
	}
}

func (s *Scanner) skipComment() {
	for {
		if s.ch == EOF || s.ch == '\r' || s.ch == '\n' {
			break
		}
		s.next()
	}
}

func (s *Scanner) scanString(head string) string {
	var sb strings.Builder
	sb.WriteString(head)
scanEnd:
	for {
		switch s.ch {
		case EOF, ';', '\'', '"':
			break scanEnd
		case '\\':
			sb.WriteRune(s.ch)
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
		sb.WriteRune(s.ch)
		s.next()
	}
	return sb.String()
}

func (s *Scanner) scanQuotedString() string {
	quote := s.ch

	var sb strings.Builder
	sb.WriteRune(quote)
	for {
		s.next()
		if s.ch == EOF {
			s.error("unexpected end of string")
			break
		} else if s.ch == '\'' || s.ch == '"' {
			if s.ch == quote {
				sb.WriteRune(s.ch)
				break
			}
		} else if s.ch == '\\' {
			s.next()
			if s.ch == EOF {
				s.error("unexpected end of string")
				break
			}
			if s.ch != quote {
				sb.WriteRune('\\')
			}
		}
		sb.WriteRune(s.ch)
	}
	s.next()
	return sb.String()
}
