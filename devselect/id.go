package devselect

import (
	"syscall"

	"github.com/gentlemanautomaton/windevice/deviceproperty"
	"github.com/gentlemanautomaton/windevice/setupapi"
)

// ID returns a selector that matches device hardware identifiers.
func ID(matcher StringMatcher) Selector {
	return func(devices syscall.Handle, device setupapi.DevInfoData) (bool, error) {
		ids, err := setupapi.GetDeviceRegistryStrings(devices, device, deviceproperty.HardwareID)
		if err != nil {
			return false, err
		}
		for _, id := range ids {
			if matcher.Match(id) {
				return true, nil
			}
		}
		return false, nil
	}
}
