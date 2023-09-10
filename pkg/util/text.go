package util

import (
	"regexp"
)

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var stripRegex = regexp.MustCompile(ansi)

// StripANSI removes ANSI escape codes from a string.
func StripANSI(str string) string {
	return stripRegex.ReplaceAllString(str, "")
}
