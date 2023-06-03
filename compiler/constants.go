package compiler

const eofChar byte = '\000'

type tokenKind int16

const (
	tokEOF tokenKind = iota - 1
	tokNEWLINE
	tokNUMBER
	tokIDENT
	tokSTRING
)

// Keywords
const (
	tokLABEL tokenKind = iota + 101
	tokGOTO
	tokPRINT
	tokINPUT
	tokLET
	tokIF
	tokTHEN
	tokENDIF
	tokWHILE
	tokREPEAT
	tokENDWHILE
)

var keywordTokenMap = map[string]tokenKind{
	"LABEL":    tokLABEL,
	"GOTO":     tokGOTO,
	"PRINT":    tokPRINT,
	"INPUT":    tokINPUT,
	"LET":      tokLET,
	"IF":       tokIF,
	"THEN":     tokTHEN,
	"ENDIF":    tokENDIF,
	"WHILE":    tokWHILE,
	"REPEAT":   tokREPEAT,
	"ENDWHILE": tokENDWHILE,
}

// Operators
const (
	tokEQ tokenKind = iota + 201
	tokPLUS
	tokMINUS
	tokASTERISK
	tokSLASH
	tokEQEQ
	tokNOTEQ
	tokLT
	tokLTEQ
	tokGT
	tokGTEQ
)

type token struct {
	text string
	kind tokenKind
}
