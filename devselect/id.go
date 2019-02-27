package devselect

import "github.com/gentlemanautomaton/windevice"

// ID returns a selector that matches device hardware identifiers.
func ID(matcher StringMatcher) Selector {
	return func(device windevice.Device) (bool, error) {
		ids, err := device.HardwareID()
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
