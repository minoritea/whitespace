package testcases

func init() {
	TestCases = append(TestCases, TestCase{
		Expect: `5`,
		ReadableSource: `SS STSSL # push 4
SS STSTL # push 5 => [4 5]
SLSTLST # print 5
`,
	})
}
