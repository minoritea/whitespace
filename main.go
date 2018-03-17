package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

type Program []Command

func (p Program) LastCommand() Command {
	if len(p) == 0 {
		panic("No command found.")
	}
	return p[len(p)-1]
}

func (p *Program) AddBitToLastCommandParam(b bool) { p.LastCommand().AddBitToParam(b) }

func (p Program) Run() {
	rt := &Runtime{
		commands: p,
		heap:     make(map[int]int),
		labels:   make(map[int]int),
	}

	rt.Run()
}

type Command interface {
	AddBitToParam(bool)
	Exec(*Runtime)
}

type Runtime struct {
	stack     Stack
	commands  Program
	heap      map[int]int
	labels    map[int]int
	callstack []int
	index     int
}

func (rt *Runtime) stepIn() bool {
	rt.commands[rt.index].Exec(rt)
	rt.index += 1
	return rt.index < len(rt.commands)
}

func (rt *Runtime) Run() {
	for rt.stepIn() {
	}
}

type Stack []int

func (rt *Runtime) Push(i int) { rt.stack = append(rt.stack, i) }
func (rt *Runtime) Pop() int {
	li := len(rt.stack) - 1
	if li < 0 {
		panic("No value exists on the stack.")
	}

	topv, rest := rt.stack[li], rt.stack[0:li]
	rt.stack = rest
	return topv
}

func (rt *Runtime) Pop2() (int, int) {
	li := len(rt.stack) - 2
	if li < 0 {
		panic("No enough values exist on the stack.")
	}
	v1, v2, rest := rt.stack[li+1], rt.stack[li], rt.stack[0:li]
	rt.stack = rest
	return v1, v2
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

func (rt *Runtime) JumpTo(label int) {
	jumpto, ok := rt.labels[c.ParamToInt()]
	if !ok {
		panic("No such label was defined")
	}
	rt.index = jumpto
}

type CommandWithParam struct{ Param }

func (p *CommandWithParam) AddBitToParam(b bool) { p.Param = append(p.Param, b) }

type CommandWithNoParam struct{}

func (p CommandWithNoParam) AddBitToParam(_ bool) { panic("DO NOT ADD BIT TO THIS") }

type CommandStackPush struct{ CommandWithParam }
type CommandStackDupN struct{ CommandWithParam }
type CommandStackDup struct{ CommandWithNoParam }
type CommandStackSlide struct{ CommandWithParam }
type CommandStackSwap struct{ CommandWithNoParam }
type CommandStackDiscard struct{ CommandWithNoParam }
type CommandIOPutNum struct{ CommandWithNoParam }
type CommandArithAdd struct{ CommandWithNoParam }
type CommandArithSub struct{ CommandWithNoParam }
type CommandArithMul struct{ CommandWithNoParam }
type CommandArithDiv struct{ CommandWithNoParam }
type CommandArithMod struct{ CommandWithNoParam }
type CommandHeapStore struct{ CommandWithNoParam }
type CommandHeapLoad struct{ CommandWithNoParam }
type CommandLabel struct{ CommandWithParam }
type CommandFlowGoSub struct{ CommandWithParam }
type CommandFlowJump struct{ CommandWithParam }
type CommandFlowBEZ struct{ CommandWithParam }
type CommandFlowBLTZ struct{ CommandWithParam }
type CommandFlowEndSub struct{ CommandWithNoParam }
type CommandFlowHalt struct{ CommandWithNoParam }

func (c *CommandStackPush) Exec(rt *Runtime) { rt.Push(c.ParamToInt()) }
func (c *CommandStackDupN) Exec(rt *Runtime) { rt.CopyToTop(c.ParamToInt()) }
func (c *CommandStackDup) Exec(rt *Runtime)  { rt.CopyToTop(0) }
func (c *CommandStackSlide) Exec(rt *Runtime) {
	topv := rt.Pop()
	rt.DiscardN(c.ParamToInt() + 1)
	rt.Push(topv)
}
func (c *CommandStackSwap) Exec(rt *Runtime)    { rt.Swap() }
func (c *CommandStackDiscard) Exec(rt *Runtime) { rt.DiscardN(1) }

func (c *CommandArithAdd) Exec(rt *Runtime) { rt.Push(rt.Pop() + rt.Pop()) }
func (c *CommandArithSub) Exec(rt *Runtime) { v1, v2 := rt.Pop2(); rt.Push(v2 - v1) }
func (c *CommandArithMul) Exec(rt *Runtime) { rt.Push(rt.Pop() * rt.Pop()) }
func (c *CommandArithDiv) Exec(rt *Runtime) { v1, v2 := rt.Pop2(); rt.Push(v2 / v1) }
func (c *CommandArithMod) Exec(rt *Runtime) { v1, v2 := rt.Pop2(); rt.Push(v2 % v1) }

func (c *CommandHeapStore) Exec(rt *Runtime) { v1, v2 := rt.Pop2(); rt.heap[v2] = v1 }
func (c *CommandHeapLoad) Exec(rt *Runtime) {
	v, ok := rt.heap[rt.Pop()]
	if !ok {
		panic("No such address in the heap")
	}
	rt.Push(v)
}

func (c *CommandLabel) Exec(rt *Runtime)     { rt.labels[c.ParamToInt()] = rt.index }
func (c *CommandFlowJump) Exec(rt *Runtime)  { rt.JumpTo(c.ParamToInt()) }
func (c *CommandFlowGoSub) Exec(rt *Runtime) { rt.JumpTo(c.ParamToInt()) }

func (c *CommandIOPutNum) Exec(rt *Runtime) { fmt.Printf("%d", rt.Pop()) }

type Op int
type Param []bool

func (p Param) ParamToInt() int {
	i := 0
	for _, b := range p {
		i *= 2

		if b {
			i += 1
		}
	}
	return i
}

type B = byte

var (
	S B = 0x20
	T B = 0x09
	L B = 0x0A
)

var b2s = map[B]string{S: "S", T: "T", L: "L"}

type Parser struct {
	state  State
	reader *bufio.Reader
	Program
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		state:  S_START,
		reader: bufio.NewReader(r),
	}
}

func (p *Parser) Parse() Program {
	for p.parse() {
	}
	return p.Program
}

func (p *Parser) parse() bool {
	b, err := p.reader.ReadByte()
	if err != nil {
		return false
	}

	switch p.state {
	case S_START:
		switch b {
		case S:
			p.state = S_S
		case T:
			p.state = S_T
		case L:
			p.state = S_L
		}

	case S_S:
		switch b {
		case S:
			// SS
			p.AddCommand(&CommandStackPush{})
			p.state = S_READ_PARAM
		case T:
			p.state = S_ST
		case L:
			p.state = S_SL
		}
	case S_ST:
		switch b {
		case S:
			// STS
			p.AddCommand(&CommandStackDupN{})
			p.state = S_READ_PARAM
		case T:
			goto UNKNOWN_STATE
		case L:
			// STL
			p.AddCommand(&CommandStackSlide{})
			p.state = S_READ_PARAM
		}
	case S_SL:
		switch b {
		case S:
			// SLS
			p.AddCommand(&CommandStackDup{})
			p.state = S_START
		case T:
			// SLT
			p.AddCommand(&CommandStackSwap{})
			p.state = S_START
		case L:
			// SLL
			p.AddCommand(&CommandStackDiscard{})
			p.state = S_START
		}

	case S_T:
		switch b {
		case S:
			p.state = S_TS
		case T:
			p.state = S_TT
		case L:
			p.state = S_TL
		}

	case S_TS:
		switch b {
		case S:
			p.state = S_TSS
		case T:
			p.state = S_TST
		case L:
			goto UNKNOWN_STATE
		}

	case S_TT:
		switch b {
		case S:
			// TTS
			p.AddCommand(&CommandHeapStore{})
			p.state = S_START
		case T:
			// TTT
			p.AddCommand(&CommandHeapLoad{})
			p.state = S_START
		case L:
			goto UNKNOWN_STATE
		}

	case S_TSS:
		switch b {
		case S:
			// TSSS
			p.AddCommand(&CommandArithAdd{})
			p.state = S_START
		case T:
			// TSST
			p.AddCommand(&CommandArithSub{})
			p.state = S_START
		case L:
			// TSSL
			p.AddCommand(&CommandArithMul{})
			p.state = S_START
		}

	case S_TST:
		switch b {
		case S:
			// TSTS
			p.AddCommand(&CommandArithDiv{})
			p.state = S_START
		case T:
			// TSTT
			p.AddCommand(&CommandArithMod{})
			p.state = S_START
		case L:
			goto UNKNOWN_STATE
		}

	case S_TL:
		switch b {
		case S:
			p.state = S_TLS
		case T:
			goto UNKNOWN_STATE
		case L:
			goto UNKNOWN_STATE
		}

	case S_TLS:
		switch b {
		case S:
			goto UNKNOWN_STATE
		case T:
			// TLST
			p.AddCommand(&CommandIOPutNum{})
			p.state = S_START
		case L:
			goto UNKNOWN_STATE
		}

	case S_L:
		switch b {
		case S:
			p.state = S_LS
		case T:
			p.state = S_LT
		case L:
			p.state = S_LL
		}

	case S_LS:
		switch b {
		case S:
			// LSS
			p.state = S_LS
		case T:
			p.state = S_LT
		case L:
			p.state = S_LL
		}

	case S_READ_PARAM:
		switch b {
		case S:
			p.LastCommand().AddBitToParam(false)
		case T:
			p.LastCommand().AddBitToParam(true)
		case L:
			p.state = S_START
		}
	}

	return true

UNKNOWN_STATE:
	log.Fatalf("UNKNOWN_STATE: %d, %s\n", p.state, b2s[b])
	return false // unreachable
}

func (p *Parser) AddCommand(cmd Command) { p.Program = append(p.Program, cmd) }

type State int

const (
	S_START State = iota
	S_S
	S_ST
	S_SL

	S_T
	S_TS
	S_TSS
	S_TST
	S_TT
	S_TL
	S_TLS

	S_L
	S_LS
	S_LT
	S_LL

	S_READ_PARAM
)

func main() {
	NewParser(os.Stdin).Parse().Run()
}
