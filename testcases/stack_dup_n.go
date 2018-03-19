package testcases

func init() {
	TestCases = append(TestCases, TestCase{
		Expect: `2232`,
		ReadableSource: `
SS STSL # push 2
SS STTL # push 3 => [3 2]
STS SSTL # copy the 2nd value of the stack => [2 3 2]
STS STSL # copy the 3nd value of the stack => [2 2 3 2]
TLST # print 2
TLST # print 2
TLST # print 3
TLST # print 2
`,
	})
}
