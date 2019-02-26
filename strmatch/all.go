package strmatch

// All returns a matcher that returns true when all submatches are true.
func All(submatches ...Matcher) Matcher {
	return func(cmp string) bool {
		for _, submatch := range submatches {
			if !submatch(cmp) {
				return false
			}
		}
		return true
	}
}
