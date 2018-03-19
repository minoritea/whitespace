package testcases

import "github.com/minoritea/whitespace/converter"

type TestCase struct {
	Name           string
	Expect         string
	ReadableSource string
	Stdin          string
}

func (tc TestCase) GetSource() string {
	return converter.FromReadable(tc.ReadableSource)
}

var TestCases = make([]TestCase, 0)
