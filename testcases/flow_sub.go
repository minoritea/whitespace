package testcases

func init() {
	TestCases = append(TestCases, TestCase{
		Expect: `5`,
		ReadableSource: `LSLSSSSSSL # jump to main: "ssssss"

LSSTSSSSTL # label "tsssst"
SSSTSTL     # push 5
LTL # end of the subroutine: "tsssst"


# main
LSSSSSSSSL # label "ssssss" is main
SSSTSSL     # push 4
LSTTSSSSTL # call "tsssst" => [5 4]
TLST       # print
`,
	})
}
