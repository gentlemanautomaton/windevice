package devselect

import (
	"github.com/gentlemanautomaton/windevice"
	"github.com/gentlemanautomaton/windevice/setupapi"
)

// Class returns a selector that matches device class names.
func Class(matcher StringMatcher) Selector {
	return func(device windevice.Device) (bool, error) {
		class, err := device.Class()
		if err != nil && err != setupapi.ErrInvalidData {
			return false, err
		}
		return matcher.Match(class), nil
	}
}
