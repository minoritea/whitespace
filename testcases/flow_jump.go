package testcases

func init() {
	TestCases = append(TestCases, TestCase{
		Expect: `5`,
		ReadableSource: `SSSTSL  # push 2
SSSTTL  # push 3 => [3 2]
LSLTTTTSSSSL # jump to "ttttssss"
SSSTSSL # push 4, it must be skipped
LSSTTTTSSSSL # label: "ttttssss"
TSSS   # sum the top 2 values of the stack
TLST   # print
`,
	})
}
