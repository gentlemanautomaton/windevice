package diflagex

import "strings"

// Value describes a set of extended flags for device installation and
// enumeration.
type Value uint32

// Match reports whether v contains all of the extended flags specified by c.
func (v Value) Match(c Value) bool {
	return v&c == c
}

// String returns a string representation of the extended flags using a
// default separator and format.
func (v Value) String() string {
	return v.Join("|", FormatGo)
}

// Join returns a string representation of the extended flags using the given
// separator and format.
func (v Value) Join(sep string, format Format) string {
	if s, ok := format[v]; ok {
		return s
	}

	var matched []string
	for i := 0; i < 32; i++ {
		flag := Value(1 << uint32(i))
		if v.Match(flag) {
			if s, ok := format[flag]; ok {
				matched = append(matched, s)
			}
		}
	}

	return strings.Join(matched, sep)
}
