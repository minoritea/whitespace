package testcases

func init() {
	TestCases = append(TestCases, TestCase{
		Name:   "Bez",
		Expect: `2`,
		ReadableSource: `
LSL SSSSSS L # jump to main

LSS TTTSSS L # label "tttsss"
SS STS L # push 2
TLST # print
LSL TTTTTT L # jump to the end of the program

# main
LSS SSSSSS L # label "ssssss" is main
SS SS L # push 0
LTS TTTSSS L # jump "tttsss" if the top of the stack is zero

LSS TTTTTT L # the end of the program
`,
	})
}
