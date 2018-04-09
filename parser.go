package main

import (
	"bufio"
	"errors"
	"io"
)

func NewParser(r interface {
	Read(p []byte) (n int, err error)
}) *Parser {
	return &Parser{
		reader: bufio.NewReader(r),
		buffer: make([]byte, 10),
		labels: make(map[string]int),
	}
}

type byteReader interface{ ReadByte() (byte, error) }
type bytesReader interface{ ReadBytes(byte) ([]byte, error) }

type Parser struct {
	reader interface {
		byteReader
		bytesReader
	}
	buffer []byte
	bindex int
	index  int
	labels map[string]int
}

func withParam(cmd commandWithParam, r bytesReader) (command, error) {
	p, err := readParam(r)
	if err != nil {
		return nil, err
	}
	cmd.setParam(p)
	return cmd, nil
}

func withLabel(cmd commandWithLabel, r bytesReader) (command, error) {
	l, err := readLabel(r)
	if err != nil {
		return nil, err
	}
	cmd.setLabel(l)
	return cmd, nil
}

func (p *Parser) parseCommand(r interface {
	byteReader
	bytesReader
}) (command, error) {
	p.bindex = 0
	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		switch b {
		case S:
			p.buffer[p.bindex] = 0x53
			p.bindex++
		case T:
			p.buffer[p.bindex] = 0x54
			p.bindex++
		case L:
			p.buffer[p.bindex] = 0x4c
			p.bindex++
		}

		switch string(p.buffer[0:p.bindex]) {
		case "SS":
			return withParam(&CommandStackPush{}, r)
		case "STS":
			return withParam(&CommandStackDupN{}, r)
		case "STL":
			return withParam(&CommandStackSlide{}, r)
		case "SLS":
			return &CommandStackDup{}, nil
		case "SLT":
			return &CommandStackSwap{}, nil
		case "SLL":
			return &CommandStackDiscard{}, nil
		case "TTS":
			return &CommandHeapStore{}, nil
		case "TTT":
			return &CommandHeapLoad{}, nil
		case "TSSS":
			return &CommandArithAdd{}, nil
		case "TSST":
			return &CommandArithSub{}, nil
		case "TSSL":
			return &CommandArithMul{}, nil
		case "TSTS":
			return &CommandArithDiv{}, nil
		case "TSTT":
			return &CommandArithMod{}, nil
		case "TLSS":
			return &CommandIOPutRune{}, nil
		case "TLST":
			return &CommandIOPutNum{}, nil
		case "TLTT":
			return &CommandIOReadNum{}, nil
		case "TLTS":
			return &CommandIOReadRune{}, nil
		case "LSS":
			if label, err := readLabel(r); err != nil {
				return nil, err
			} else {
				p.labels[label] = p.index
				return &CommandLabel{}, nil // empty command
			}
		case "LST":
			return withLabel(&CommandFlowGoSub{}, r)
		case "LSL":
			return withLabel(&CommandFlowJump{}, r)
		case "LTS":
			return withLabel(&CommandFlowBEZ{}, r)
		case "LTT":
			return withLabel(&CommandFlowBLTZ{}, r)
		case "LTL":
			return &CommandFlowEndSub{}, nil
		case "LLL":
			return &CommandFlowHalt{}, nil
		}
	}
}

func readLabel(r bytesReader) (string, error) {
	bs, err := r.ReadBytes(L)
	if err != nil {
		return "", err
	}
	l := len(bs)
	label := make([]byte, 0, l)
	for _, b := range bs[0 : l-1] {
		switch b {
		case S, T:
			label = append(label, b)
		}
	}
	if len(label) == 1 {
		return "", errors.New("A label must have at least one byte")
	}
	return string(label), nil
}

func readParam(r bytesReader) (int, error) {
	bs, err := r.ReadBytes(L)
	if err != nil {
		return 0, err
	}
	i := 0
	first := true
	minus := false
	c := 0
	for _, b := range bs[0 : len(bs)-1] {
		switch b {
		case S:
			if first {
				first = false
			} else {
				i <<= 1
			}
			c++
		case T:
			if first {
				minus = true
				first = false
			} else {
				i <<= 1
				i += 1
			}
			c++
		}
	}
	if c < 2 {
		return 0, errors.New("A signed integer must have at least two bits")
	}
	if minus {
		i *= -1
	}
	return i, nil
}

func (p *Parser) parse() (program, error) {
	var program program
	for {
		if cmd, err := p.parseCommand(p.reader); err == nil {
			p.index++
			program.commands = append(program.commands, cmd)
		} else if err == io.EOF && p.bindex == 0 {
			program.labels = p.labels
			return program, nil
		} else {
			return program, err
		}
	}
}

func (p *Parser) Parse() program {
	program, err := p.parse()
	if err != nil {
		panic(err)
	}
	return program
}

func (p program) Run(options ...RunOption) {
	rt := &Runtime{
		commands: p.commands,
		heap:     make(map[int]int),
		labels:   p.labels,
	}

	RunOptions(options).apply(rt)

	rt.Run()
}
