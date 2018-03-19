package testcases

func init() {
	TestCases = append(TestCases, TestCase{
		Expect: `2`,
		ReadableSource: `
LSL SSSSSS L # jump to main

LSS TTTSSS L # label "tttsss"
SS STS L # push 2
TLST # print
LSL TTTTTT L # jump to the end of the program

# main
LSS SSSSSS L # "ssssss" is main
SS TT L # push -1
LTT TTTSSS L # jump "tttsss" if the top of the stack is negative

LSS TTTTTT L # the end of the program
`,
	})
}
