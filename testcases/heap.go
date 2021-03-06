package testcases

func init() {
	TestCases = append(TestCases, TestCase{
		Expect: `12`,
		ReadableSource: `SS STL # push 1
SS STSSL # push 4
SS STSL # push 2  => [2 4 1]
TTS # store 2 to heap[4] => [1], {4: 2}
TLST # print 1
SS STSSL # push 4
TTT #  load heap[4] => [2]
TLST # print 2
`,
	})
}
