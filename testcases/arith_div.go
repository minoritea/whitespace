package testcases

func init() {
	TestCases = append(TestCases, TestCase{
		Expect: `43`,
		ReadableSource: `
SSSTTSTL # push 13
SSSTSSSL # push 8
SSSTSL   # push 2 => [2 8 13]
TSTS    # 8 / 2 =>  [4 13]
SLSTLST # print 4
TSTS    # 13 / 4  => [3]
TLST    # print 3
`,
	})
}
