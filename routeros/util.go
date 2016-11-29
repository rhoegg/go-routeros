package routeros

import (
	"strconv"
)

func stringBool(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func parseBool(s string) bool {
	b, err := strconv.ParseBool(s)
	return (err != nil) && b
}
