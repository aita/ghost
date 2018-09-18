package shell

import (
	"fmt"
	"io"

	multierror "github.com/hashicorp/go-multierror"
)

func Parse(r io.Reader) (prog *Program, err error) {
	p := newParser(r)
	defer func() {
		err = p.errors.ErrorOrNil()
	}()

	prog = p.parse()
	return
}

type parser struct {
	scanner *Scanner
	errors  *multierror.Error

	tok *Token // one token look-ahead
}

func newParser(r io.Reader) *parser {
	p := &parser{}
	p.scanner = NewScanner(r, p.error)
	return p
}

func (p *parser) error(pos Position, msg string) {
	p.errors = multierror.Append(p.errors, fmt.Errorf("%d:%d %s", pos.Line, pos.Column, msg))
}

func (p *parser) next() {
	p.tok = p.scanner.Scan()
}

func (p *parser) acceptString(lit string) bool {
	return p.tok.Kind == STRING && p.tok.Literal == lit
}

func (p *parser) expect(kind TokenKind) bool {
	if p.tok.Kind == kind {
		return true
	}
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", kind, p.tok.Kind)
	p.error(p.tok.Pos, msg)
	return false
}

func (p *parser) expectString(lit string) bool {
	if p.tok.Kind == STRING && p.tok.Literal == lit {
		return true
	}
	msg := fmt.Sprintf("expected next token to be STRING(%s), got %s(%s) instead", lit, p.tok.Kind, p.tok.Literal)
	p.error(p.tok.Pos, msg)
	return false
}

func (p *parser) parse() *Program {
	prog := &Program{}
	prog.Body = p.parseBlock()
	return prog
}

func (p *parser) parseBlock() *BlockStmt {
	block := &BlockStmt{}
loopEnd:
	for {
		p.next()
		switch p.tok.Kind {
		case TERMINATOR:
			continue
		case STRING:
			if p.acceptString("end") {
				break loopEnd
			}
			stmt := p.parseStmt()
			block.List = append(block.List, stmt)
		default:
			// unexpected token
			break loopEnd
		}
	}
	return block
}

func (p *parser) parseStmt() Stmt {
	if p.acceptString("if") {
		return p.parseIfStmt()
	}
	return p.parseCommand()
}

func (p *parser) parseIfStmt() *IfStmt {
	p.next()
	ifStmt := &IfStmt{}
	ifStmt.Cond = p.parseCommand()
	if !p.expect(TERMINATOR) {
		return nil
	}
	p.next()
	ifStmt.Body = p.parseBlock()
	if p.acceptString("else") {
		p.next()
		if p.acceptString("if") {
			ifStmt.Else = p.parseIfStmt()
			return ifStmt
		}
		if !p.expect(TERMINATOR) {
			return nil
		}
		ifStmt.Else = p.parseStmt()
	}
	if !p.expectString("end") {
		return nil
	}
	return ifStmt
}

func (p *parser) parseCommand() *Command {
	cmd := &Command{}
	cmd.Command = p.parseWord()
	for {
		p.next()
		switch p.tok.Kind {
		case TERMINATOR, EOF:
			return cmd
		case STRING:
			word := p.parseWord()
			cmd.Args = append(cmd.Args, word)
		default:
		}
	}
}

func (p *parser) parseWord() *Word {
	return &Word{
		Token: p.tok,
		Value: p.tok.Literal,
	}
}
