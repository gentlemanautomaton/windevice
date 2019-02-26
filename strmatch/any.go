package strmatch

// Any returns a matcher that returns true when any submatch is true.
func Any(submatches ...Matcher) Matcher {
	return func(cmp string) bool {
		for _, submatch := range submatches {
			if submatch(cmp) {
				return true
			}
		}
		return false
	}
}
