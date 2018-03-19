package testcases

func init() {
	TestCases = append(TestCases, TestCase{
		Expect: `10`,
		ReadableSource: `SS STL # push 1
SS SSL # push 0
SLT # swap
TLSTTLST
`,
	})
}
