package testcases

func init() {
	TestCases = append(TestCases, TestCase{
		Name:   "IO_READ_RUNE",
		Expect: `Hello, World🌏`,
		ReadableSource: `
SS    ST L SS   STS L SS   STT L SS  STSS L
SS  STST L SS  STTS L SS  STTT L SS STSSS L
SS STSST L SS STSTS L SS STSTT L SS STTSS L
SS STTST L # [13 12 11 10 9 8 7 6 5 4 3 2 1]

TLTS TLTS TLTS TLTS # read rune 4 times
TLTS TLTS TLTS TLTS # read rune 4 times
TLTS TLTS TLTS TLTS # read rune 4 times
TLTS # total 13 times

SS SS L
SS    ST L SS   STS L SS   STT L SS  STSS L
SS  STST L SS  STTS L SS  STTT L SS STSSS L
SS STSST L SS STSTS L SS STSTT L SS STTSS L
SS STTST L # [13 12 11 10 9 8 7 6 5 4 3 2 1 0]


LSS TTTTTT L
TTT TLSS
SLS LTS SSSSSS L # refer the top of the stack, then jump to "ssssss" if the value is zero
LSL TTTTTT L # else jump to "tttttt"

LSS SSSSSS L
`,
		Stdin: `Hello, World🌏`,
	})
}
