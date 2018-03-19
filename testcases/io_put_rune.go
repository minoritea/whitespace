package testcases

func init() {
	TestCases = append(TestCases, TestCase{
		Expect: `A`,
		ReadableSource: `
SS STSSSSST L # push 65
TLSS # print rune 'A'
`,
	})
}
