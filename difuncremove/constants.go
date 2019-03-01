package difuncremove

// Windows device installation function removal flags.
//
// https://docs.microsoft.com/en-us/windows/desktop/api/setupapi/ns-setupapi-_sp_removedevice_params
const (
	Global         = 0x00000001 // DI_REMOVEDEVICE_GLOBAL
	ConfigSpecific = 0x00000002 // DI_REMOVEDEVICE_CONFIGSPECIFIC
)
