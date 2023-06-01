package compiler

import "fmt"

type parser struct {
	lexer     *lexer
	curToken  token
	peekToken token
}

func (p *parser) checkToken(kind tokenKind) bool {
	return kind == p.curToken.kind
}

func (p *parser) checkPeek(kind tokenKind) bool {
	return kind == p.peekToken.kind
}

func (p *parser) match(kind tokenKind) {
	if !p.checkToken(kind) {
		p.abort(fmt.Sprintf("Expected %d, got %d", kind, p.curToken.kind))
	}
	p.nextToken()
}

func (p *parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.getToken()
}

func (p *parser) abort(message string) {
	panic(fmt.Sprintf("parsing error: %s", message))
}

func (p *parser) Parse() {

}

func NewParser(lexer *lexer) *parser {
	newParser := parser{
		lexer: lexer,
	}
	// initialize current and peek tokens
	newParser.nextToken()
	newParser.nextToken()
	return &newParser
}
