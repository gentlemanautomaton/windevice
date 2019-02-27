package strmatch

import "strings"

// Contains returns a matcher for value. The matcher returns true for any
// string that contains value.
func Contains(value string) Matcher {
	return func(cmp string) bool {
		return strings.Contains(cmp, value)
	}
}

// ContainsInsensitive returns a matcher for value. The matcher returns true
// for anystring that contains value when ignoring case.
func ContainsInsensitive(value string) Matcher {
	value = strings.ToLower(value)
	return func(cmp string) bool {
		return strings.Contains(strings.ToLower(cmp), value)
	}
}
