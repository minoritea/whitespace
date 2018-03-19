package testcases

func init() {
	TestCases = append(TestCases, TestCase{
		Expect: `21`,
		ReadableSource: `
SSSTTTL # push 7
SSSTTL  # push 3
TSSL   # 7 * 3
TLST   # print 21
`,
	})
}
