package devselect

import (
	"syscall"

	"github.com/gentlemanautomaton/windevice/deviceproperty"
	"github.com/gentlemanautomaton/windevice/setupapi"
)

// Class returns a selector that matches device class names.
func Class(matcher StringMatcher) Selector {
	return func(devices syscall.Handle, device setupapi.DevInfoData) (bool, error) {
		class, err := setupapi.GetDeviceRegistryString(devices, device, deviceproperty.Class)
		if err != nil && err != setupapi.ErrInvalidData {
			return false, err
		}
		return matcher.Match(class), nil
	}
}
