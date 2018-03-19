package main

import (
	"bytes"
	"github.com/minoritea/whitespace/testcases"
	"testing"
)

func TestAll(t *testing.T) {
	for _, tc := range testcases.TestCases {
		src := bytes.NewBufferString(tc.GetSource())
		output := bytes.NewBuffer(nil)
		NewParser(src).Parse().Run(
			SetStdin(nil),
			SetStdout(output),
		)
		if result := output.String(); result != tc.Expect {
			t.Errorf(`
The expected result is "%s",
but the acutal result is "%s".
			`, tc.Expect, result)
		}
	}
}
