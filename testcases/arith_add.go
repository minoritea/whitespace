package testcases

func init() {
	TestCases = append(TestCases, TestCase{
		Name:   "Add",
		Expect: `52`,
		ReadableSource: `
SSSTSL # push 2
SSSTTL # push 3 => [3 2]
TSSS # 3 + 2 => [5]
TLST  # print 5
SSSTTL  # push 3
SSTTL  # push -1 => [-1 3]
TSSS # -1 + 3 => 2
TLST  # print 2
`,
	})
}
