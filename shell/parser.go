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

func (p *parser) accept(kind TokenKind) bool {
	return p.tok.Kind == kind
}

func (p *parser) acceptKeyword(keyword string) bool {
	return p.tok.Kind == STRING && p.tok.Literal == keyword
}

func (p *parser) expect(kind TokenKind) bool {
	if p.tok.Kind == kind {
		return true
	}
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", kind, p.tok.Kind)
	p.error(p.tok.Pos, msg)
	return false
}

func (p *parser) expectKeyword(keyword string) bool {
	if p.tok.Kind == STRING && p.tok.Literal == keyword {
		return true
	}
	msg := fmt.Sprintf("expected next token to be STRING(%v), got %s(%v) instead", keyword, p.tok.Kind, p.tok.Literal)
	p.error(p.tok.Pos, msg)
	return false
}

func (p *parser) parse() *Program {
	prog := &Program{}
	for {
		p.next()
		if p.accept(EOF) {
			break
		}
		stmt := p.parseStmt()
		if stmt == nil {
			// maybe something wrong
			break
		}
		prog.Body = append(prog.Body, stmt)
	}
	return prog
}

func (p *parser) parseStmt() Stmt {
	if p.acceptKeyword("if") {
		return p.parseIfStmt()
	}
	if p.accept(STRING) {
		return p.parseCommand()
	}

	msg := fmt.Sprintf("unexpected token %s(%v)", p.tok.Kind, p.tok.Literal)
	p.error(p.tok.Pos, msg)
	return nil
}

func (p *parser) parseIfStmt() *IfStmt {
	p.next()
	ifStmt := &IfStmt{}
	ifStmt.Cond = p.parseCommand()
	ifStmt.Body = p.parseIfBlock()

	expectEnd := true
	if p.acceptKeyword("else") {
		p.next()
		if !p.expect(TERMINATOR) {
			return nil
		}
		p.next()
		if p.acceptKeyword("if") {
			p.next()
			ifStmt.Else = p.parseIfStmt()
			expectEnd = false
		} else {
			ifStmt.Else = p.parseIfBlock()
		}
	}
	if expectEnd {
		if !p.expectKeyword("end") {
			return nil
		}
		p.next()
		if !p.expect(TERMINATOR) {
			return nil
		}
	}
	return ifStmt
}

func (p *parser) parseIfBlock() *BlockStmt {
	block := &BlockStmt{}
	for !p.acceptKeyword("end") && !p.acceptKeyword("else") {
		stmt := p.parseStmt()
		if stmt == nil {
			// something maybe wrong
			break
		}
		block.List = append(block.List, stmt)
	}
	return block
}

func (p *parser) parseCommand() *CommandStmt {
	cmd := &CommandStmt{}
	if !p.expect(STRING) {
		return nil
	}
	cmd.Command = p.parseWord()

	p.next()
	for !p.accept(TERMINATOR) {
		if p.accept(STRING) {
			word := p.parseWord()
			cmd.Args = append(cmd.Args, word)
		} else {
			msg := fmt.Sprintf("unexpected token %s(%v)", p.tok.Kind, p.tok.Literal)
			p.error(p.tok.Pos, msg)
			return nil
		}
		p.next()
	}

	p.next()
	return cmd
}

func (p *parser) parseWord() *Word {
	return &Word{
		Token: p.tok,
		Value: p.tok.Literal,
	}
}
