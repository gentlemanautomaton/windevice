package devselect

import (
	"syscall"

	"github.com/gentlemanautomaton/windevice/setupapi"
)

// Any returns a selector that returns true when any selector returns true.
func Any(selectors ...Selector) Selector {
	return func(devices syscall.Handle, device setupapi.DevInfoData) (bool, error) {
		for _, selector := range selectors {
			ok, err := selector(devices, device)
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
