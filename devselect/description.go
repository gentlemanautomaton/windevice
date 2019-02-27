package devselect

import (
	"github.com/gentlemanautomaton/windevice"
	"github.com/gentlemanautomaton/windevice/setupapi"
)

// Description returns a selector that matches device descriptions.
func Description(matcher StringMatcher) Selector {
	return func(device windevice.Device) (bool, error) {
		desc, err := device.Description()
		if err != nil && err != setupapi.ErrInvalidData {
			return false, err
		}
		return matcher.Match(desc), nil
	}
}
