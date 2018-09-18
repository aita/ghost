package parser

import (
	"fmt"
	"github.com/aita/ghost/shell/ast"
	"github.com/aita/ghost/shell/token"
	"io"

	"github.com/hashicorp/go-multierror"

	"github.com/aita/ghost/shell/scanner"
)

func Parse(r io.Reader) (prog *ast.Program, err error) {
	p := newParser(r)
	defer func() {
		err = p.errors.ErrorOrNil()
	}()

	prog = p.parse()
	return
}

type parser struct {
	scanner *scanner.Scanner
	errors  *multierror.Error

	tok *token.Token // one token look-ahead
}

func newParser(r io.Reader) *parser {
	p := &parser{}
	p.scanner = scanner.NewScanner(r, p.error)
	return p
}

func (p *parser) error(pos token.Position, msg string) {
	p.errors = multierror.Append(p.errors, fmt.Errorf("%d:%d %s", pos.Line, pos.Column, msg))
}

func (p *parser) next() {
	p.tok = p.scanner.Scan()
}

func (p *parser) accept(kind token.TokenKind) bool {
	return p.tok.Kind == kind
}

func (p *parser) acceptKeyword(keyword string) bool {
	return p.tok.Kind == token.STRING && p.tok.Literal == keyword
}

func (p *parser) expect(kind token.TokenKind) *token.Token {
	if p.tok.Kind != kind {
		msg := fmt.Sprintf("expected next token to be %s, got %s instead", kind, p.tok.Kind)
		p.error(p.tok.Pos, msg)
	}
	tok := p.tok
	p.next() // make progress
	return tok
}

func (p *parser) expectKeyword(keyword string) bool {
	ret := p.tok.Kind == token.STRING && p.tok.Literal == keyword
	if !ret {
		var msg string
		if p.tok.Kind == token.STRING {
			msg = fmt.Sprintf("expected next token to be %q, got %q instead", keyword, p.tok.Literal)
		} else {
			msg = fmt.Sprintf("expected next token to be %q, got %s instead", keyword, p.tok.Kind)
		}
		p.error(p.tok.Pos, msg)
	}
	p.next() // make progress
	return ret
}

func (p *parser) parse() *ast.Program {
	prog := &ast.Program{}
	p.next()
	for {
		if p.accept(token.EOF) {
			break
		}
		stmt := p.parseStmt()
		prog.Body = append(prog.Body, stmt)
	}
	return prog
}

func (p *parser) parseStmt() ast.Stmt {
	if p.acceptKeyword("if") {
		return p.parseIfStmt()
	}
	if p.accept(token.STRING) {
		return p.parseCommand()
	}

	msg := fmt.Sprintf("unexpected token %s(%#v)", p.tok.Kind, p.tok.Literal)
	p.error(p.tok.Pos, msg)
	p.next() // make progress
	return &ast.BadStmt{}
}

func (p *parser) parseIfStmt() *ast.IfStmt {
	p.next()
	ifStmt := &ast.IfStmt{}
	ifStmt.Cond = p.parseCommand()
	ifStmt.Body = p.parseIfBlock()

	expectEnd := true
	if p.acceptKeyword("else") {
		p.next()
		if p.accept(token.TERMINATOR) {
			p.next()
			ifStmt.Else = p.parseIfBlock()
		} else if p.acceptKeyword("if") {
			p.next()
			ifStmt.Else = p.parseIfStmt()
			expectEnd = false
		} else {
			ifStmt.Else = &ast.BadStmt{}
		}
	}
	if expectEnd {
		p.expectKeyword("end")
		p.expect(token.TERMINATOR)
	}
	return ifStmt
}

func (p *parser) parseIfBlock() *ast.BlockStmt {
	block := &ast.BlockStmt{}
	for {
		if p.accept(token.EOF) {
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

func (p *parser) parseCommand() ast.Stmt {
	cmd := &ast.CommandStmt{}
	cmd.Command = p.parseWord()
	for !p.accept(token.TERMINATOR) {
		if p.accept(token.EOF) {
			p.error(p.tok.Pos, "unexpected EOF")
			break
		}
		word := p.parseWord()
		cmd.Args = append(cmd.Args, word)
	}
	p.next()
	return cmd
}

func (p *parser) parseWord() *ast.Word {
	tok := p.expect(token.STRING)
	return &ast.Word{
		Token: tok,
		Value: tok.Literal,
	}
}
