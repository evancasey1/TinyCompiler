package compiler

import (
	"fmt"
	"unicode"
)

type lexer struct {
	source  string
	CurChar byte
	curPos  int
}

func (l *lexer) nextChar() {
	l.curPos += 1
	if l.curPos >= len(l.source) {
		l.CurChar = eofChar
	} else {
		l.CurChar = l.source[l.curPos]
	}
}

func (l *lexer) peek() byte {
	if l.curPos+1 >= len(l.source) {
		return eofChar
	}
	return l.source[l.curPos+1]
}

func (l *lexer) abort(message string) {
	panic(fmt.Sprintf("lexing error: %s", message))
}

func (l *lexer) skipWhitespace() {
	for l.CurChar == ' ' || l.CurChar == '\t' || l.CurChar == '\r' {
		l.nextChar()
	}
}

func (l *lexer) skipComment() {
	if l.CurChar == '#' {
		for l.CurChar != '\n' {
			l.nextChar()
		}
	}
}
func (l *lexer) getToken() token {
	l.skipWhitespace()
	l.skipComment()

	var currToken token
	switch curChar := l.CurChar; {
	case curChar == '+':
		currToken = token{text: string(l.CurChar), kind: tokPLUS}
	case curChar == '-':
		currToken = token{text: string(l.CurChar), kind: tokMINUS}
	case curChar == '*':
		currToken = token{text: string(l.CurChar), kind: tokASTERISK}
	case curChar == '/':
		currToken = token{text: string(l.CurChar), kind: tokSLASH}
	case curChar == '=':
		if l.peek() == '=' {
			lastChar := l.CurChar
			l.nextChar()
			currToken = token{text: string(lastChar + l.CurChar), kind: tokEQEQ}
		} else {
			currToken = token{text: string(l.CurChar), kind: tokEQ}
		}
	case curChar == '>':
		if l.peek() == '=' {
			lastChar := l.CurChar
			l.nextChar()
			currToken = token{text: string(lastChar + l.CurChar), kind: tokGTEQ}
		} else {
			currToken = token{text: string(l.CurChar), kind: tokGT}
		}
	case curChar == '<':
		if l.peek() == '=' {
			lastChar := l.CurChar
			l.nextChar()
			currToken = token{text: string(lastChar + l.CurChar), kind: tokLTEQ}
		} else {
			currToken = token{text: string(l.CurChar), kind: tokLT}
		}
	case curChar == '!':
		if l.peek() == '=' {
			lastChar := l.CurChar
			l.nextChar()
			currToken = token{text: string(lastChar + l.CurChar), kind: tokNOTEQ}
		} else {
			l.abort(fmt.Sprintf("expected !=, got !%c", l.peek()))
		}
	case curChar == '"':
		l.nextChar()
		startPos := l.curPos

		for l.CurChar != '"' {
			// Don't allow special characters in the string. No escape characters, newlines, tabs, or %.
			// We will be using C's printf on this string.
			if l.CurChar == '\r' || l.CurChar == '\n' || l.CurChar == '\t' || l.CurChar == '\\' || l.CurChar == '%' {
				l.abort(fmt.Sprintf("illegal character in string: %c", l.CurChar))
			}
			l.nextChar()
		}
		tokText := l.source[startPos : l.curPos+1]
		currToken = token{text: tokText, kind: tokSTRING}
	case unicode.IsDigit(rune(curChar)):
		// Leading character is a digit, so this must be a number.
		// Get all consecutive digits and decimal if there is one.
		startPos := l.curPos
		for unicode.IsDigit(rune(l.peek())) {
			l.nextChar()
		}
		if l.peek() == '.' {
			l.nextChar()
			if !unicode.IsDigit(rune(l.peek())) {
				l.abort(fmt.Sprintf("illegal character in number: %c", l.peek()))
			}
			for unicode.IsDigit(rune(l.peek())) {
				l.nextChar()
			}
		}
		tokText := l.source[startPos : l.curPos+1]
		currToken = token{text: tokText, kind: tokNUMBER}
	case unicode.IsLetter(rune(curChar)):
		// Leading character is a letter, so this must be an identifier or a keyword.
		// Get all consecutive alpha numeric characters.
		/*
			startPos = self.curPos
			while self.peek().isalnum():
				self.nextChar()

			# Check if the token is in the list of keywords.
			tokText = self.source[startPos : self.curPos + 1] # Get the substring.
			keyword = Token.checkIfKeyword(tokText)
			if keyword == None: # Identifier
				token = Token(tokText, TokenType.IDENT)
			else:   # Keyword
				token = Token(tokText, keyword)
		*/
		startPos := l.curPos
		for unicode.IsLetter(rune(l.peek())) || unicode.IsDigit(rune(l.peek())) {
			l.nextChar()
		}
		tokText := l.source[startPos : l.curPos+1]
		keyword, ok := keywordTokenMap[tokText]
		if ok {
			currToken = token{text: tokText, kind: keyword}
		} else {
			currToken = token{text: tokText, kind: tokIDENT}
		}

	case curChar == '\n':
		currToken = token{text: string(l.CurChar), kind: tokNEWLINE}
	case curChar == eofChar:
		currToken = token{kind: tokEOF}
	default:
		l.abort(fmt.Sprintf("unknown token: %c", l.CurChar))
	}
	l.nextChar()
	return currToken
}

func (l *lexer) Lex() {
	tok := l.getToken()
	for tok.kind != tokEOF {
		fmt.Printf("%d ", tok.kind)
		tok = l.getToken()
	}
}

func NewLexer(source string) *lexer {
	newLexer := lexer{
		source: source + "\n", // append a newline character to make the last line easier to parse
		curPos: -1,
	}
	newLexer.nextChar()
	return &newLexer
}
