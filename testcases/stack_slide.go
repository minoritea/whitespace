package testcases

func init() {
	TestCases = append(TestCases, TestCase{
		Expect: `12`,
		ReadableSource: `
SS STSL # push 2
SS STL # push 1
SS STL # push 1
SS STL # push 1 => [1 1 1 2]
STL STL # delete the 2nd and 3rd values => [1 2]
TLSTTLST
`,
	})
}
