package devselect

import (
	"github.com/gentlemanautomaton/windevice"
)

// A Selector is capable of matching devices in a device info list.
type Selector func(windevice.Device) (bool, error)

// Select returns true if the selector matches the given device.
func (s Selector) Select(device windevice.Device) (bool, error) {
	if s == nil {
		return true, nil
	}
	return s(device)
}
