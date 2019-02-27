package devselect

import "github.com/gentlemanautomaton/windevice"

// All returns a selector that returns true when all selectors return true.
func All(selectors ...Selector) Selector {
	return func(device windevice.Device) (bool, error) {
		for _, selector := range selectors {
			ok, err := selector(device)
			if err != nil {
				return false, err
			}
			if !ok {
				return false, nil
			}
		}
		return true, nil
	}
}
