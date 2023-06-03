package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"teenytiny/compiler"
)

func main() {
	filePtr := flag.String("f", "", "input file")
	flag.Parse()
	if filePtr == nil || *filePtr == "" || !strings.HasSuffix(*filePtr, ".tt") {
		panic("invalid file")
	}

	source, readErr := os.ReadFile(*filePtr)
	if readErr != nil {
		panic(fmt.Sprintf("error reading file: %s", readErr.Error()))
	}

	lexer := compiler.NewLexer(string(source))
	emitter := compiler.NewEmitter("out.c")
	parser := compiler.NewParser(lexer, emitter)
	parser.Parse()
	emitter.WriteFile()
	fmt.Println("Compiling complete!")
}
