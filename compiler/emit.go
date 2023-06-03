package compiler

import (
	"fmt"
	"os"
)

type emitter struct {
	filePath string
	header   string
	code     string
}

func (e *emitter) emit(code string) {
	e.code += code
}

func (e *emitter) emitLine(code string) {
	e.code += code + "\n"
}

func (e *emitter) headerLine(code string) {
	e.header += code + "\n"
}

func (e *emitter) WriteFile() {
	if err := os.WriteFile(e.filePath, []byte(e.header+e.code), 0644); err != nil {
		panic(fmt.Sprintf("Error writing to file: %s", err.Error()))
	}
}

func NewEmitter(filePath string) *emitter {
	return &emitter{filePath: filePath}
}
