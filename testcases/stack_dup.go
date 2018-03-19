package testcases

func init() {
	TestCases = append(TestCases, TestCase{
		Expect: `445`,
		ReadableSource: `
SS STSTL # push 5
SS STSSL # push 4 => [4 5]
SLS # copy the top value of the stack => [4 4 5]
TLSTTLSTTLST
`,
	})
}
