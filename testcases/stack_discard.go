package testcases

func init() {
	TestCases = append(TestCases, TestCase{
		Expect: `2`,
		ReadableSource: `
SS STSL # push 2
SS STTL # push 3 => [3 2]
SLL   # discard
TLST  # print 2
`,
	})
}
