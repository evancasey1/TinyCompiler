package compiler

type node struct {
	left  *node
	right *node
	value string
}
