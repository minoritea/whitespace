package testcases

func init() {
	TestCases = append(TestCases, TestCase{
		Expect: `5`,
		ReadableSource: `
SS STST L # push 5
TLST # print
LLL # halt
SS STSS L # push 4
TLST # print

`,
	})
}
