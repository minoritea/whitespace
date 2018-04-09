package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

type program struct {
	commands []command
	labels   map[string]int
}

type Runtime struct {
	stdin     *bufio.Reader
	stdout    io.Writer
	stack     Stack
	commands  []command
	labels    map[string]int
	heap      map[int]int
	callstack []int
	index     int
}

func (rt *Runtime) stepIn() bool {
	rt.commands[rt.index].exec(rt)
	// log.Printf("%d:\t%T\t%+v\n", rt.index, rt.commands[rt.index], rt.commands[rt.index])
	rt.index += 1
	return rt.index < len(rt.commands)
}

func (rt *Runtime) Run() {
	if len(rt.commands) < 1 {
		return
	}

	for rt.stepIn() {
	}
}

type Stack []int

func (rt *Runtime) Push(i int) { rt.stack = push(rt.stack, i) }
func (rt *Runtime) Pop() int {
	tv, rest := popN(rt.stack, 1)
	rt.stack = rest
	return tv[0]
}

func push(stack []int, i int) []int { return append(stack, i) }
func popN(stack []int, i int) (topN []int, rest []int) {
	l := len(stack)
	li := l - i
	if li < 0 {
		panic("No enough values exist on the stack.")
	}
	return stack[li:l], stack[0:li]
}

func (rt *Runtime) Pop2() (int, int) {
	tv, rest := popN(rt.stack, 2)
	rt.stack = rest
	return tv[1], tv[0]
}

func (rt *Runtime) PushToCallstack(i int) { rt.callstack = push(rt.callstack, i) }
func (rt *Runtime) PopFromCallstack() int {
	tv, rest := popN(rt.callstack, 1)
	rt.callstack = rest
	return tv[0]
}

func (rt *Runtime) Swap() {
	li := len(rt.stack) - 2
	if li < 0 {
		panic("No values on the stack.")
	}
	rt.stack[li+1], rt.stack[li] = rt.stack[li], rt.stack[li+1]
}

func (rt *Runtime) CopyToTop(i int) {
	i = len(rt.stack) - i - 1
	if i < 0 {
		panic("Invalid index number")
	}
	rt.Push(rt.stack[i])
}

func (rt *Runtime) DiscardN(i int) {
	i = len(rt.stack) - i
	if i < 0 {
		panic("Invalid index number")
	}
	rt.stack = rt.stack[0:i]
}

func (rt *Runtime) JumpToLabel(label string) {
	jumpto, ok := rt.labels[label]
	if !ok {
		panic("No such label was defined")
	}
	rt.index = jumpto
}

func (rt *Runtime) JumpToIndex(index int) { rt.index = index }

type B = byte

var (
	S B = 0x20
	T B = 0x09
	L B = 0x0A
)

func (p program) String() string {
	result := ""
	for _, cmd := range p.commands {
		result += fmt.Sprintf("%T\n", cmd)
	}
	return result
}

type RunOption func(*Runtime)
type RunOptions []RunOption

func (options RunOptions) apply(rt *Runtime) {
	for _, f := range options {
		f(rt)
	}
}

func SetStdin(r io.Reader) RunOption  { return func(rt *Runtime) { rt.stdin = bufio.NewReader(r) } }
func SetStdout(w io.Writer) RunOption { return func(rt *Runtime) { rt.stdout = w } }

type State int

func openSourceFile() (*os.File, error) {
	if len(os.Args) < 2 {
		log.Fatal(`Usage: whitespace [file]`)
	}
	return os.Open(os.Args[1])
}

func main() {
	if f, err := openSourceFile(); err != nil {
		log.Fatal(err)
	} else {
		NewParser(f).Parse().Run(SetStdout(os.Stdout), SetStdin(os.Stdin))
	}
}
