package devselect

// StringMatcher is a function capable of matching strings.
type StringMatcher interface {
	Match(string) bool
}
