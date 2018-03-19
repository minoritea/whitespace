package testcases

func init() {
	TestCases = append(TestCases, TestCase{
		Name:   "IO_READ_NUM",
		Expect: `43`,
		ReadableSource: `
SS STS L # push 2
SS ST L  # push 1 # => [1 2]
TLTT # read a number from stdin and put it(3) to heap[1]
TLTT # read a number from stdin and put it(4) to heap[2]
SS ST L  # push 1
TTT  # push heap[1] to the stack => [3]
SS STS L # push 2 => [2 3]
TTT  # push heap[2] to the stack # [4 3]
TLST # print
TLST # print
`,
		Stdin: "3 4",
	})
}
