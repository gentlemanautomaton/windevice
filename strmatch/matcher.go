package strmatch

// Matcher is a function capable of matching strings.
type Matcher func(string) bool

// Match returns true if the matcher matches the value.
func (m Matcher) Match(value string) bool {
	if m == nil {
		return true
	}
	return m(value)
}
