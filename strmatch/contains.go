package strmatch

import "strings"

// Contains returns a matcher that returns true for anything that
// contains value.
func Contains(value string) Matcher {
	return func(cmp string) bool {
		return strings.Contains(cmp, value)
	}
}
