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

func (p *parser) expect(kind TokenKind) *Token {
	if p.tok.Kind != kind {
		p.expectError(kind)
	}
	tok := p.tok
	p.next() // make progress
	return tok
}

func (p *parser) expectKeyword(keyword string) bool {
	ret := p.tok.Kind == STRING && p.tok.Literal == keyword
	if !ret {
		msg := fmt.Sprintf("expected next token to be STRING(%#v), got %s(%#v) instead", keyword, p.tok.Kind, p.tok.Literal)
		p.error(p.tok.Pos, msg)
	}
	p.next() // make progress
	return ret
}

func (p *parser) expectError(kind TokenKind) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", kind, p.tok.Kind)
	p.error(p.tok.Pos, msg)
}

func (p *parser) parse() *Program {
	prog := &Program{}
	for {
		p.next()
		if p.accept(EOF) {
			break
		}
		stmt := p.parseStmt()
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

	msg := fmt.Sprintf("unexpected token %s(%#v)", p.tok.Kind, p.tok.Literal)
	p.error(p.tok.Pos, msg)
	p.next() // make progress
	return &BadStmt{}
}

func (p *parser) parseIfStmt() *IfStmt {
	p.next()
	ifStmt := &IfStmt{}
	ifStmt.Cond = p.parseCommand()
	ifStmt.Body = p.parseIfBlock()

	expectEnd := true
	if p.acceptKeyword("else") {
		p.next()
		if p.accept(TERMINATOR) {
			p.next()
			ifStmt.Else = p.parseIfBlock()
		} else if p.acceptKeyword("if") {
			p.next()
			ifStmt.Else = p.parseIfStmt()
			expectEnd = false
		} else {
			ifStmt.Else = &BadStmt{}
		}
	}
	if expectEnd {
		p.expectKeyword("end")
		p.expect(TERMINATOR)
	}
	return ifStmt
}

func (p *parser) parseIfBlock() *BlockStmt {
	block := &BlockStmt{}
	for !p.acceptKeyword("end") && !p.acceptKeyword("else") {
		stmt := p.parseStmt()
		block.List = append(block.List, stmt)
	}
	return block
}

func (p *parser) parseCommand() Stmt {
	cmd := &CommandStmt{}
	cmd.Command = p.parseWord()
	for !p.accept(TERMINATOR) {
		word := p.parseWord()
		cmd.Args = append(cmd.Args, word)
	}
	p.next()
	return cmd
}

func (p *parser) parseWord() *Word {
	tok := p.expect(STRING)
	return &Word{
		Token: tok,
		Value: tok.Literal,
	}
}
