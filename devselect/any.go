package devselect

import "github.com/gentlemanautomaton/windevice"

// Any returns a selector that returns true when any selector returns true.
func Any(selectors ...Selector) Selector {
	return func(device windevice.Device) (bool, error) {
		for _, selector := range selectors {
			ok, err := selector(device)
			if err != nil {
				return false, err
			}
			if ok {
				return true, nil
			}
		}
		return false, nil
	}
}
