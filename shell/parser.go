package shell

import (
	"fmt"
	"io"

	"github.com/hashicorp/go-multierror"
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
		msg := fmt.Sprintf("expected next token to be %s, got %s instead", kind, p.tok.Kind)
		p.error(p.tok.Pos, msg)
	}
	tok := p.tok
	p.next() // make progress
	return tok
}

func (p *parser) expectKeyword(keyword string) bool {
	ret := p.tok.Kind == STRING && p.tok.Literal == keyword
	if !ret {
		var msg string
		if p.tok.Kind == STRING {
			msg = fmt.Sprintf("expected next token to be %q, got %q instead", keyword, p.tok.Literal)
		} else {
			msg = fmt.Sprintf("expected next token to be %q, got %s instead", keyword, p.tok.Kind)
		}
		p.error(p.tok.Pos, msg)
	}
	p.next() // make progress
	return ret
}

func (p *parser) parse() *Program {
	prog := &Program{}
	p.next()
	for {
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

	msg := fmt.Sprintf("unexpected token %s", p.tok.Kind)
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
	for {
		if p.accept(EOF) {
			p.error(p.tok.Pos, "unexpected EOF")
			break
		}
		if p.acceptKeyword("end") || p.acceptKeyword("else") {
			break
		}
		stmt := p.parseStmt()
		block.List = append(block.List, stmt)
	}
	return block
}

func (p *parser) parseCommand() Stmt {
	cmd := &CommandStmt{}
	for !p.accept(TERMINATOR) {
		if p.accept(EOF) {
			p.error(p.tok.Pos, "unexpected EOF")
			break
		}
		word := p.parseWord()
		cmd.List = append(cmd.List, word)
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
