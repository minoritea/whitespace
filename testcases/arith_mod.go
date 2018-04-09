package testcases

func init() {
	TestCases = append(TestCases, TestCase{
		Name:   "Mod",
		Expect: `1`,
		ReadableSource: `
SSSTTSTL # push 13
SSSTSSL  # push 4 [4 13]
TSTT    # 13 % 4 => 1
TLST    # 1
`,
	})
}
