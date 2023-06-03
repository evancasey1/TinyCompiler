package compiler

import "fmt"

type parser struct {
	lexer     *lexer
	emitter   *emitter
	curToken  token
	peekToken token

	symbols        map[string]bool
	labelsDeclared map[string]bool
	labelsGotoed   map[string]bool
}

// BEGIN CONTROL FUNCTIONS

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
	p.program()
}

// BEGIN GRAMMAR RULES
/*
program

    ::= {statement}
*/
func (p *parser) program() {
	p.emitter.headerLine("#include <stdio.h>") // make printf and scanf available
	p.emitter.headerLine("int main(void){")

	for p.checkToken(tokNEWLINE) {
		p.nextToken()
	}

	for !p.checkToken(tokEOF) {
		p.statement()
	}

	for label := range p.labelsGotoed {
		if _, found := p.labelsDeclared[label]; !found {
			p.abort(fmt.Sprintf("Attempting to GOTO undeclared label: %s", label))
		}
	}

	p.emitter.emitLine("return 0;")
	p.emitter.emitLine("}")
}

/*
statement

	    ::= "PRINT" (expression | string) nl
		| "IF" comparison "THEN" nl {statement} "ENDIF" nl
		| "WHILE" comparison "REPEAT" nl {statement} "ENDWHILE" nl
		| "LABEL" ident nl
		| "GOTO" ident nl
		| "LET" ident "=" expression nl
		| "INPUT" ident nl
*/
func (p *parser) statement() {
	if p.checkToken(tokPRINT) {
		p.nextToken()

		if p.checkToken(tokSTRING) {
			p.emitter.emitLine(fmt.Sprintf(`printf("%s\n");`, p.curToken.text))
			p.nextToken()
		} else {
			p.emitter.emit(`printf("%.2f\n", (float)(`)
			p.expression()
			p.emitter.emitLine("));")
		}
	} else if p.checkToken(tokIF) {
		p.nextToken()
		p.emitter.emit("if(")
		p.comparison()
		p.match(tokTHEN)
		p.nl()
		p.emitter.emitLine("){")

		// zero or more statements are allowed in the body
		for !p.checkToken(tokENDIF) {
			p.statement()
		}
		p.match(tokENDIF)
		p.emitter.emitLine("}")
	} else if p.checkToken(tokWHILE) {
		p.nextToken()
		p.emitter.emit("while(")
		p.comparison()
		p.match(tokREPEAT)
		p.nl()
		p.emitter.emitLine("){")

		// zero or more statements are allowed in the body
		for !p.checkToken(tokENDWHILE) {
			p.statement()
		}
		p.match(tokENDWHILE)
		p.emitter.emitLine("}")
	} else if p.checkToken(tokLABEL) {
		p.nextToken()

		if _, found := p.labelsDeclared[p.curToken.text]; found {
			p.abort(fmt.Sprintf("Label already declared: %s", p.curToken.text))
		}
		p.labelsDeclared[p.curToken.text] = true
		p.emitter.emitLine(p.curToken.text + ":")
		p.match(tokIDENT)
	} else if p.checkToken(tokGOTO) {
		p.nextToken()
		p.labelsGotoed[p.curToken.text] = true
		p.emitter.emitLine(fmt.Sprintf("goto %s;", p.curToken.text))
		p.match(tokIDENT)
	} else if p.checkToken(tokLET) {
		p.nextToken()

		if _, found := p.symbols[p.curToken.text]; !found {
			p.symbols[p.curToken.text] = true
			p.emitter.headerLine(fmt.Sprintf("float %s;", p.curToken.text))
		}
		p.emitter.emit(p.curToken.text + " = ")
		p.match(tokIDENT)
		p.match(tokEQ)
		p.expression()
		p.emitter.emitLine(";")
	} else if p.checkToken(tokINPUT) {
		p.nextToken()
		if _, found := p.symbols[p.curToken.text]; !found {
			p.symbols[p.curToken.text] = true
			p.emitter.headerLine(fmt.Sprintf("float %s;", p.curToken.text))
		}
		p.emitter.emitLine(fmt.Sprintf(`if(0 == scanf("%%f", &%s)) {`, p.curToken.text))
		p.emitter.emitLine(p.curToken.text + " = 0;")
		p.emitter.emit(`scanf("%`)
		p.emitter.emitLine(`*s");`)
		p.emitter.emitLine("}")
		p.match(tokIDENT)
	} else {
		p.abort(fmt.Sprintf("Invalid statement at %s (%d)", p.curToken.text, p.curToken.kind))
	}
	p.nl()
}

/*
comparison

	::= expression (("==" | "!=" | ">" | ">=" | "<" | "<=") expression)+
*/
func (p *parser) comparison() {
	p.expression()
	if p.isComparisonOperator() {
		p.emitter.emit(p.curToken.text)
		p.nextToken()
		p.expression()
	} else {
		p.abort(fmt.Sprintf("Expected comparison operator at %s (%d)", p.curToken.text, p.curToken.kind))
	}

	for p.isComparisonOperator() {
		p.emitter.emit(p.curToken.text)
		p.nextToken()
		p.expression()
	}
}

func (p *parser) isComparisonOperator() bool {
	return p.checkToken(tokGT) ||
		p.checkToken(tokGTEQ) ||
		p.checkToken(tokLT) ||
		p.checkToken(tokLTEQ) ||
		p.checkToken(tokEQEQ) ||
		p.checkToken(tokNOTEQ)
}

/*
expression

	::= term {( "-" | "+" ) term}
*/
func (p *parser) expression() {
	p.term()
	for p.checkToken(tokPLUS) || p.checkToken(tokMINUS) {
		p.emitter.emit(p.curToken.text)
		p.nextToken()
		p.term()
	}
}

/*
term

	::= unary {( "/" | "*" ) unary}
*/
func (p *parser) term() {
	p.unary()
	for p.checkToken(tokSLASH) || p.checkToken(tokASTERISK) {
		p.emitter.emit(p.curToken.text)
		p.nextToken()
		p.unary()
	}
}

/*
unary

	::= ["+" | "-"] primary
*/
func (p *parser) unary() {
	if p.checkToken(tokPLUS) || p.checkToken(tokMINUS) {
		p.emitter.emit(p.curToken.text)
		p.nextToken()
	}
	p.primary()
}

/*
primary

	::= number | ident
*/
func (p *parser) primary() {
	if p.checkToken(tokNUMBER) {
		p.emitter.emit(p.curToken.text)
		p.nextToken()
	} else if p.checkToken(tokIDENT) {
		if _, found := p.symbols[p.curToken.text]; !found {
			p.abort(fmt.Sprintf("Referencing variable before assignment: %s", p.curToken.text))
		}

		p.emitter.emit(p.curToken.text)
		p.nextToken()
	} else {
		p.abort(fmt.Sprintf("Unexpected token at %s (%d)", p.curToken.text, p.curToken.kind))
	}
}

/*
nl

	::= '\n'+
*/
func (p *parser) nl() {
	p.match(tokNEWLINE)

	// allow for extra newline tokens
	for p.checkToken(tokNEWLINE) {
		p.nextToken()
	}
}

func NewParser(lexer *lexer, emitter *emitter) *parser {
	newParser := parser{
		lexer:          lexer,
		emitter:        emitter,
		symbols:        make(map[string]bool),
		labelsDeclared: make(map[string]bool),
		labelsGotoed:   make(map[string]bool),
	}
	// initialize current and peek tokens
	newParser.nextToken()
	newParser.nextToken()
	return &newParser
}
