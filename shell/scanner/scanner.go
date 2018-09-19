package scanner

import (
	"bufio"
	"io"
	"strings"
	"unicode"

	"github.com/aita/ghost/shell/token"
)

func newToken(kind token.TokenKind, lit string, pos token.Position) *token.Token {
	return &token.Token{
		Kind:    kind,
		Literal: lit,
		Pos:     pos,
	}
}

type ErrorHandler func(pos token.Position, msg string)

type Scanner struct {
	ErrorCount       int
	errHandler       ErrorHandler
	src              *bufio.Reader
	ch               rune
	insertTerminator bool
	sb               strings.Builder
	lastSize         int
	pos              token.Position
}

func NewScanner(r io.Reader, errHandler ErrorHandler) *Scanner {
	return &Scanner{
		ErrorCount:       0,
		errHandler:       errHandler,
		src:              bufio.NewReader(r),
		ch:               ' ',
		insertTerminator: false,
		pos: token.Position{
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
			ch = token.EOF
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

func (s *Scanner) Scan() *token.Token {
scanAgain:
	pos := s.pos
	for {
		if !unicode.IsSpace(s.ch) {
			break
		}
		if s.ch == '\n' && s.insertTerminator {
			s.insertTerminator = false
			return newToken(token.TERMINATOR, "\n", pos)
		}
		s.next()
		pos = s.pos
	}

	switch s.ch {
	case token.EOF:
		if s.insertTerminator {
			s.insertTerminator = false
			return newToken(token.TERMINATOR, "", pos)
		}
		return newToken(token.EOF, "", pos)
	case ';':
		s.insertTerminator = false
		s.next()
		return newToken(token.TERMINATOR, ";", pos)
	case '#':
		s.skipComment()
		goto scanAgain
	case '\'', '"':
		s.insertTerminator = true
		lit := s.scanQuotedString()
		return newToken(token.STRING, lit, pos)
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
		return newToken(token.STRING, lit, pos)
	}
}

func (s *Scanner) skipComment() {
	for {
		if s.ch == token.EOF || s.ch == '\r' || s.ch == '\n' {
			break
		}
		s.sb.WriteRune(s.ch)
		s.next()
	}
}

func (s *Scanner) scanString(head string) string {
	s.sb.Reset()
	s.sb.WriteString(head)
scanEnd:
	for {
		switch s.ch {
		case token.EOF, ';', '\'', '"':
			break scanEnd
		case '\\':
			s.sb.WriteRune(s.ch)
			s.next()
			if s.ch == token.EOF {
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
		if s.ch == token.EOF {
			s.error("unexpected end of string")
			break
		} else if s.ch == '\'' || s.ch == '"' {
			if s.ch == quote {
				s.sb.WriteRune(s.ch)
				break
			}
		} else if s.ch == '\\' {
			s.next()
			if s.ch == token.EOF {
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
