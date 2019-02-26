package devselect

import (
	"syscall"

	"github.com/gentlemanautomaton/windevice/deviceproperty"
	"github.com/gentlemanautomaton/windevice/setupapi"
)

// Description returns a selector that matches device descriptions.
func Description(matcher StringMatcher) Selector {
	return func(devices syscall.Handle, device setupapi.DevInfoData) (bool, error) {
		desc, err := setupapi.GetDeviceRegistryString(devices, device, deviceproperty.Description)
		if err != nil && err != setupapi.ErrInvalidData {
			return false, err
		}
		return matcher.Match(desc), nil
	}
}
