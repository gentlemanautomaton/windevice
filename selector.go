package windevice

import (
	"syscall"

	"github.com/gentlemanautomaton/windevice/setupapi"
)

// Selector is an interface capable of selecting devices in a device list.
type Selector interface {
	Select(devices syscall.Handle, device setupapi.DevInfoData) (bool, error)
}
