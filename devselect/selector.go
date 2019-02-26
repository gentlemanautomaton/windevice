package devselect

import (
	"syscall"

	"github.com/gentlemanautomaton/windevice/setupapi"
)

// A Selector is capable of matching devices in a device info list.
type Selector func(devices syscall.Handle, device setupapi.DevInfoData) (bool, error)

// Select returns true if the selector matches the given device.
func (s Selector) Select(devices syscall.Handle, device setupapi.DevInfoData) (bool, error) {
	if s == nil {
		return true, nil
	}
	return s(devices, device)
}
