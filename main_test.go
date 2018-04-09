package main

import (
	"bytes"
	"github.com/minoritea/whitespace/testcases"
	"io"
	"testing"
)

func toTest(tc testcases.TestCase) func(t *testing.T) {
	return func(t *testing.T) {
		src := bytes.NewBufferString(tc.GetSource())
		output := bytes.NewBuffer(nil)
		var input io.Reader
		if tc.Stdin != "" {
			input = bytes.NewBufferString(tc.Stdin)
		}
		NewParser(src).Parse().Run(
			SetStdin(input),
			SetStdout(output),
		)
		if result := output.String(); result != tc.Expect {
			t.Fatalf(`
The expected result is "%s",
but the acutal result is "%s".
			`, tc.Expect, result)
		}
	}
}

func TestAll(t *testing.T) {
	for _, tc := range testcases.TestCases {
		t.Log(tc.Name)
		t.Run(tc.Name, toTest(tc))
	}
}
