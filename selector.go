package windevice

// Selector is an interface capable of selecting devices in a device list.
type Selector interface {
	Select(Device) (bool, error)
}
