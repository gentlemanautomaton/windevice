package hwprofile

// Format maps scopes to their string representations.
type Format map[Scope]string

// FormatGo maps flags to Go-style constant strings.
var FormatGo = Format{
	Global:         "Global",
	ConfigSpecific: "ConfigSpecific",
	ConfigGeneral:  "ConfigGeneral",
}
