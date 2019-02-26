package devselect

import (
	"syscall"

	"github.com/gentlemanautomaton/windevice/deviceproperty"
	"github.com/gentlemanautomaton/windevice/setupapi"
)

// FriendlyName returns a selector that matches device descriptions and
// friendly names.
func FriendlyName(matcher StringMatcher) Selector {
	return func(devices syscall.Handle, device setupapi.DevInfoData) (bool, error) {
		name, err := setupapi.GetDeviceRegistryString(devices, device, deviceproperty.FriendlyName)
		if err != nil && err != setupapi.ErrInvalidData {
			return false, err
		}
		return matcher.Match(name), nil
	}
}
