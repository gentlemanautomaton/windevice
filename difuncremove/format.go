package difuncremove

// Format maps flags to their string representations.
type Format map[Flags]string

// FormatGo maps flags to Go-style constant strings.
var FormatGo = Format{
	Global:         "Global",
	ConfigSpecific: "ConfigSpecific",
}
