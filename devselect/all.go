package devselect

import (
	"syscall"

	"github.com/gentlemanautomaton/windevice/setupapi"
)

// All returns a selector that returns true when all selectors return true.
func All(selectors ...Selector) Selector {
	return func(devices syscall.Handle, device setupapi.DevInfoData) (bool, error) {
		for _, selector := range selectors {
			ok, err := selector(devices, device)
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
