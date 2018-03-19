package converter

import "strings"

var (
	fromReadable = strings.NewReplacer(
		" ", "S", "\t", "T", "\n", "L",
		"S", " ", "T", "\t", "L", "\n",
	)

	toReadable = strings.NewReplacer(
		" ", "[S]", "\t", "[T]", "\n", "[L]\n",
	)
)

func FromReadable(src string) string {
	return fromReadable.Replace(src)
}

func ToReadable(src string) string {
	return toReadable.Replace(src)
}
