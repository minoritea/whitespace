package main

import "fmt"

type command interface {
	exec(*Runtime)
}

type commandWithParam interface {
	command
	setParam(int)
}
type paramHolder struct{ param int }

func (c *paramHolder) setParam(p int) { c.param = p }
func (c *paramHolder) getParam() int  { return c.param }

type labelHolder struct{ label string }

func (c *labelHolder) setLabel(label string) { c.label = label }
func (c *labelHolder) getLabel() string      { return c.label } // rename it

type commandWithLabel interface {
	command
	setLabel(string)
}

type CommandStackPush struct{ paramHolder }
type CommandStackDupN struct{ paramHolder }
type CommandStackDup struct{}
type CommandStackSlide struct{ paramHolder }
type CommandStackSwap struct{}
type CommandStackDiscard struct{}
type CommandArithAdd struct{}
type CommandArithSub struct{}
type CommandArithMul struct{}
type CommandArithDiv struct{}
type CommandArithMod struct{}
type CommandHeapStore struct{}
type CommandHeapLoad struct{}
type CommandLabel struct{}

type CommandFlowJump struct{ labelHolder }
type CommandFlowGoSub struct {
	labelHolder
}
type CommandFlowBEZ struct{ labelHolder }
type CommandFlowBLTZ struct{ labelHolder }
type CommandFlowEndSub struct{}
type CommandFlowHalt struct{}
type CommandIOPutNum struct{}
type CommandIOPutRune struct{}
type CommandIOReadNum struct{}
type CommandIOReadRune struct{}

func (c *CommandStackPush) exec(rt *Runtime) { rt.Push(c.getParam()) }
func (c *CommandStackDupN) exec(rt *Runtime) { rt.CopyToTop(c.getParam()) }
func (c *CommandStackDup) exec(rt *Runtime)  { rt.CopyToTop(0) }
func (c *CommandStackSlide) exec(rt *Runtime) {
	topv := rt.Pop()
	rt.DiscardN(c.getParam() + 1)
	rt.Push(topv)
}
func (c *CommandStackSwap) exec(rt *Runtime)    { rt.Swap() }
func (c *CommandStackDiscard) exec(rt *Runtime) { rt.DiscardN(1) }

func (c *CommandArithAdd) exec(rt *Runtime) { rt.Push(rt.Pop() + rt.Pop()) }
func (c *CommandArithSub) exec(rt *Runtime) { v1, v2 := rt.Pop2(); rt.Push(v2 - v1) }
func (c *CommandArithMul) exec(rt *Runtime) { rt.Push(rt.Pop() * rt.Pop()) }
func (c *CommandArithDiv) exec(rt *Runtime) { v1, v2 := rt.Pop2(); rt.Push(v2 / v1) }
func (c *CommandArithMod) exec(rt *Runtime) { v1, v2 := rt.Pop2(); rt.Push(v2 % v1) }

func (c *CommandHeapStore) exec(rt *Runtime) { v1, v2 := rt.Pop2(); rt.heap[v2] = v1 }
func (c *CommandHeapLoad) exec(rt *Runtime) {
	addr := rt.Pop()
	v, ok := rt.heap[addr]
	if !ok {
		panic(fmt.Sprintf("No such address: %d in the heap", addr))
	}
	rt.Push(v)
}

func (p CommandLabel) exec(rt *Runtime) { /* nothing to do */ }

func (c *CommandFlowJump) exec(rt *Runtime) { rt.JumpToLabel(c.getLabel()) }
func (c *CommandFlowGoSub) exec(rt *Runtime) {
	rt.PushToCallstack(rt.index)
	rt.JumpToLabel(c.getLabel())
}

func (c *CommandFlowBEZ) exec(rt *Runtime) {
	if rt.Pop() == 0 {
		rt.JumpToLabel(c.getLabel())
	}
}

func (c *CommandFlowBLTZ) exec(rt *Runtime) {
	if rt.Pop() < 0 {
		rt.JumpToLabel(c.getLabel())
	}
}
func (c *CommandFlowEndSub) exec(rt *Runtime) { rt.JumpToIndex(rt.PopFromCallstack()) }
func (c CommandFlowHalt) exec(rt *Runtime)    { rt.index = len(rt.commands) }

func (c *CommandIOPutNum) exec(rt *Runtime) { fmt.Fprintf(rt.stdout, "%d", rt.Pop()) }
func (c *CommandIOPutRune) exec(rt *Runtime) {
	fmt.Fprintf(rt.stdout, "%s", string([]rune{rune(int32(rt.Pop()))}))
}

func (c *CommandIOReadNum) exec(rt *Runtime) {
	i := new(int)
	fmt.Fscanf(rt.stdin, "%d", i)
	rt.heap[rt.Pop()] = *i
}

func (c *CommandIOReadRune) exec(rt *Runtime) {
	r, _, err := rt.stdin.ReadRune()
	if err != nil {
		panic(err)
	}
	rt.heap[rt.Pop()] = int(int32(r))
}
