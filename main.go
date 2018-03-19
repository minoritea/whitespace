package main

import (
	"bufio"
	"fmt"
	//	"github.com/pkg/errors"
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

type Command interface {
	AddBitToParam(bool)
	Exec(*Runtime)
	FinishReadParam()
}

type Runtime struct {
	stdin     io.Reader
	stdout    io.Writer
	stack     Stack
	commands  Program
	heap      map[int]int
	labels    map[int]int
	callstack []int
	index     int
}

func (rt *Runtime) stepIn() bool {
	rt.commands[rt.index].Exec(rt)
	// log.Printf("%d:\t%T\t%+v\n", rt.index, rt.commands[rt.index], rt.commands[rt.index])
	rt.index += 1
	return rt.index < len(rt.commands)
}

func (rt *Runtime) Run() {
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

func (rt *Runtime) JumpToLabel(label int) {
	jumpto, ok := rt.labels[label]
	if !ok {
		panic("No such label was defined")
	}
	rt.index = jumpto
}

func (rt *Runtime) JumpToIndex(index int) { rt.index = index }

type CommandWithParam struct{ Param }

func (p *CommandWithParam) AddBitToParam(b bool) { p.Param = append(p.Param, b) }
func (p *CommandWithParam) FinishReadParam()     {}

type CommandWithNoParam struct{}

func (p CommandWithNoParam) AddBitToParam(_ bool) { panic("DO NOT ADD BIT TO THIS") }
func (p CommandWithNoParam) FinishReadParam()     {}

type CommandStackPush struct{ CommandWithParam }
type CommandStackDupN struct{ CommandWithParam }
type CommandStackDup struct{ CommandWithNoParam }
type CommandStackSlide struct{ CommandWithParam }
type CommandStackSwap struct{ CommandWithNoParam }
type CommandStackDiscard struct{ CommandWithNoParam }
type CommandArithAdd struct{ CommandWithNoParam }
type CommandArithSub struct{ CommandWithNoParam }
type CommandArithMul struct{ CommandWithNoParam }
type CommandArithDiv struct{ CommandWithNoParam }
type CommandArithMod struct{ CommandWithNoParam }
type CommandHeapStore struct{ CommandWithNoParam }
type CommandHeapLoad struct{ CommandWithNoParam }
type CommandLabel struct {
	CommandWithParam
	index  int
	labels map[int]int
}

type CommandFlowJump struct{ CommandWithParam }
type CommandFlowGoSub struct {
	index int
	CommandWithParam
}
type CommandFlowBEZ struct{ CommandWithParam }
type CommandFlowBLTZ struct{ CommandWithParam }
type CommandFlowEndSub struct{ CommandWithNoParam }
type CommandFlowHalt struct{ CommandWithNoParam }
type CommandIOPutNum struct{ CommandWithNoParam }
type CommandIOPutRune struct{ CommandWithNoParam }
type CommandIOReadNum struct{ CommandWithNoParam }
type CommandIOReadRune struct{ CommandWithNoParam }

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

func (p CommandLabel) Exec(rt *Runtime) { /* nothing to do */ }
func (p CommandLabel) FinishReadParam() {
	label := p.ParamToInt()
	if _, exists := p.labels[label]; exists {
		panic("The label is already defined")
	}
	p.labels[label] = p.index
}

func (c *CommandFlowJump) Exec(rt *Runtime) { rt.JumpToLabel(c.ParamToInt()) }
func (c *CommandFlowGoSub) Exec(rt *Runtime) {
	rt.PushToCallstack(c.index)
	rt.JumpToLabel(c.ParamToInt())
}

func (c *CommandFlowBEZ) Exec(rt *Runtime) {
	if rt.Pop() == 0 {
		rt.JumpToLabel(c.ParamToInt())
	}
}

func (c *CommandFlowBLTZ) Exec(rt *Runtime) {
	if rt.Pop() < 0 {
		rt.JumpToLabel(c.ParamToInt())
	}
}
func (c *CommandFlowEndSub) Exec(rt *Runtime) { rt.JumpToIndex(rt.PopFromCallstack()) }
func (c CommandFlowHalt) Exec(rt *Runtime)    { rt.index = len(rt.commands) }

func (c *CommandIOPutNum) Exec(rt *Runtime) { fmt.Fprintf(rt.stdout, "%d", rt.Pop()) }
func (c *CommandIOPutRune) Exec(rt *Runtime) {
	fmt.Fprintf(rt.stdout, "%s", string([]rune{rune(int32(rt.Pop()))}))
}

type Op int
type Param []bool

func (p Param) ParamToInt() int {
	if len(p) < 2 {
		panic("A signed integer must have at least two bits")
	}
	i := 0
	for _, b := range p[1:] {
		i <<= 1

		if b {
			i += 1
		}
	}
	if p[0] {
		i *= -1
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
	index  int
	labels map[int]int
	Program
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		state:  S_START,
		reader: bufio.NewReader(r),
		labels: make(map[int]int),
	}
}

func (p *Parser) Parse() *Parser {
	for p.parse() {
	}
	return p
}

func (p Program) String() string {
	result := ""
	for _, cmd := range p {
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

func (p *Parser) Run(options ...RunOption) {
	rt := &Runtime{
		commands: p.Program,
		heap:     make(map[int]int),
		labels:   p.labels,
	}

	RunOptions(options).apply(rt)

	rt.Run()
}

func SetStdin(r io.Reader) RunOption  { return func(rt *Runtime) { rt.stdin = r } }
func SetStdout(w io.Writer) RunOption { return func(rt *Runtime) { rt.stdout = w } }

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
			// TLSS
			p.AddCommand(&CommandIOPutRune{})
			p.state = S_START
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
			p.AddCommand(&CommandLabel{index: p.index, labels: p.labels})
			p.state = S_READ_PARAM
		case T:
			// LST
			p.AddCommand(&CommandFlowGoSub{index: p.index})
			p.state = S_READ_PARAM
		case L:
			// LSL
			p.AddCommand(&CommandFlowJump{})
			p.state = S_READ_PARAM
		}

	case S_LT:
		switch b {
		case S:
			// LTS
			p.AddCommand(&CommandFlowBEZ{})
			p.state = S_READ_PARAM
		case T:
			// LTT
			p.AddCommand(&CommandFlowBLTZ{})
			p.state = S_READ_PARAM
		case L:
			// LTL
			p.AddCommand(&CommandFlowEndSub{})
			p.state = S_START
		}
	case S_LL:
		switch b {
		case S:
			goto UNKNOWN_STATE
		case T:
			goto UNKNOWN_STATE
		case L:
			// LLL
			p.AddCommand(&CommandFlowHalt{})
			p.state = S_START
		}

	case S_READ_PARAM:
		switch b {
		case S:
			p.LastCommand().AddBitToParam(false)
		case T:
			p.LastCommand().AddBitToParam(true)
		case L:
			p.LastCommand().FinishReadParam()
			p.state = S_START
		}
	}

	return true

UNKNOWN_STATE:
	log.Fatalf("UNKNOWN_STATE: %s, %s\n", s2s[p.state], b2s[b])
	return false // unreachable
}

func (p *Parser) AddCommand(cmd Command) {
	p.Program = append(p.Program, cmd)
	p.index += 1
}

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

var s2s = map[State]string{
	S_START: "S_START ",
	S_S:     "S_S",
	S_ST:    "S_ST",
	S_SL:    "S_SL",

	S_T:   "S_T",
	S_TS:  "S_TS",
	S_TSS: "S_TSS",
	S_TST: "S_TST",
	S_TT:  "S_TT",
	S_TL:  "S_TL",
	S_TLS: "S_TLS",

	S_L:  "S_L",
	S_LS: "S_LS",
	S_LT: "S_LT",
	S_LL: "S_LL",

	S_READ_PARAM: "S_READ_PARAM",
}

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
		NewParser(f).Parse().Run()
	}
}
