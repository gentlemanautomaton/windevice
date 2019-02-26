package strmatch

import "strings"

// Equal returns a matcher that returns true for anything that
// matches value exactly.
func Equal(value string) Matcher {
	return func(cmp string) bool {
		return cmp == value
	}
}

// EqualFold returns a matcher that returns true for anything that
// case-insenstiviely matches value.
func EqualFold(value string) Matcher {
	return func(cmp string) bool {
		return strings.EqualFold(cmp, value)
	}
}
