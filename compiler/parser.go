package compiler

import "fmt"

type parser struct {
	lexer     *lexer
	curToken  token
	peekToken token

	symbols        []string
	labelsDeclared []string
	labelsGotoed   []string
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
	fmt.Println("PROGRAM")

	for p.checkToken(tokNEWLINE) {
		p.nextToken()
	}

	for !p.checkToken(tokEOF) {
		p.statement()
	}

	for _, label := range p.labelsGotoed {
		if !elementInSet(p.labelsDeclared, label) {
			p.abort(fmt.Sprintf("Attempting to GOTO undeclared label: %s", label))
		}
	}
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
		fmt.Println("STATEMENT-PRINT")
		p.nextToken()

		if p.checkToken(tokSTRING) {
			p.nextToken()
		} else {
			p.expression()
		}
	} else if p.checkToken(tokIF) {
		fmt.Println("STATEMENT-IF")
		p.nextToken()
		p.comparison()
		p.match(tokTHEN)
		p.nl()

		// zero or more statements are allowed in the body
		for !p.checkToken(tokENDIF) {
			p.statement()
		}
		p.match(tokENDIF)
	} else if p.checkToken(tokWHILE) {
		fmt.Println("STATEMENT-WHILE")
		p.nextToken()
		p.comparison()
		p.match(tokREPEAT)
		p.nl()

		// zero or more statements are allowed in the body
		for !p.checkToken(tokENDWHILE) {
			p.statement()
		}
		p.match(tokENDWHILE)
	} else if p.checkToken(tokLABEL) {
		fmt.Println("STATEMENT-LABEL")
		p.nextToken()

		if elementInSet(p.labelsDeclared, p.curToken.text) {
			p.abort(fmt.Sprintf("Label already declared: %s", p.curToken.text))
		}
		p.labelsDeclared = setAdd(p.labelsDeclared, p.curToken.text)

		p.match(tokIDENT)
	} else if p.checkToken(tokGOTO) {
		fmt.Println("STATEMENT-GOTO")
		p.nextToken()
		p.labelsGotoed = setAdd(p.labelsGotoed, p.curToken.text)
		p.match(tokIDENT)
	} else if p.checkToken(tokLET) {
		fmt.Println("STATEMENT-LET")
		p.nextToken()
		p.symbols = setAdd(p.symbols, p.curToken.text)
		p.match(tokIDENT)
		p.match(tokEQ)
		p.expression()
	} else if p.checkToken(tokINPUT) {
		fmt.Println("STATEMENT-INPUT")
		p.nextToken()
		p.symbols = setAdd(p.symbols, p.curToken.text)
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
	fmt.Println("COMPARISON")
	p.expression()
	if p.isComparisonOperator() {
		p.nextToken()
		p.expression()
	} else {
		p.abort(fmt.Sprintf("Expected comparison operator at %s (%d)", p.curToken.text, p.curToken.kind))
	}

	for p.isComparisonOperator() {
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
	fmt.Println("EXPRESSION")
	p.term()
	for p.checkToken(tokPLUS) || p.checkToken(tokMINUS) {
		p.nextToken()
		p.term()
	}
}

/*
term

	::= unary {( "/" | "*" ) unary}
*/
func (p *parser) term() {
	fmt.Println("TERM")
	p.unary()
	for p.checkToken(tokSLASH) || p.checkToken(tokASTERISK) {
		p.nextToken()
		p.unary()
	}
}

/*
unary

	::= ["+" | "-"] primary
*/
func (p *parser) unary() {
	fmt.Println("UNARY")
	if p.checkToken(tokPLUS) || p.checkToken(tokMINUS) {
		p.nextToken()
	}
	p.primary()
}

/*
primary

	::= number | ident
*/
func (p *parser) primary() {
	fmt.Printf("PRIMARY (%s)\n", p.curToken.text)

	if p.checkToken(tokNUMBER) {
		p.nextToken()
	} else if p.checkToken(tokIDENT) {
		if !elementInSet(p.symbols, p.curToken.text) {
			p.abort(fmt.Sprintf("Referencing variable before assignment: %s", p.curToken.text))
		}
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
	fmt.Println("NEWLINE")
	p.match(tokNEWLINE)

	// allow for extra newline tokens
	for p.checkToken(tokNEWLINE) {
		p.nextToken()
	}
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
