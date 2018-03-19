package testcases

import "strings"

type TestCase struct {
	Expect         string
	ReadableSource string
	Stdin          string
}

func (tc TestCase) GetSource() string {
	return strings.NewReplacer(
		" ", "S", "\t", "T", "\n", "L",
		"S", " ", "T", "\t", "L", "\n",
	).Replace(tc.ReadableSource)
}

var TestCases = make([]TestCase, 0)
