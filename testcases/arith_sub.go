package testcases

func init() {
	TestCases = append(TestCases, TestCase{
		Name:   "Sub",
		Expect: `1`,
		ReadableSource: `
SSSTTL # push 3
SSSTSL # push 2 => [2 3]
TSST # 3 - 2 => [1]
TLST
`,
	})
}
