package utils

import "strings"

func TrimExcessiveSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
