package testcases

func init() {
	TestCases = append(TestCases, TestCase{
		Name:   "IO_PUT_RUNE",
		Expect: `A`,
		ReadableSource: `
SS STSSSSST L # push 65
TLSS # print rune 'A'
`,
	})
}
