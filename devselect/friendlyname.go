package devselect

import (
	"github.com/gentlemanautomaton/windevice"
	"github.com/gentlemanautomaton/windevice/setupapi"
)

// FriendlyName returns a selector that matches device descriptions and
// friendly names.
func FriendlyName(matcher StringMatcher) Selector {
	return func(device windevice.Device) (bool, error) {
		name, err := device.FriendlyName()
		if err != nil && err != setupapi.ErrInvalidData {
			return false, err
		}
		return matcher.Match(name), nil
	}
}
