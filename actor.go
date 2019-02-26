package windevice

import (
	"syscall"

	"github.com/gentlemanautomaton/windevice/setupapi"
)

// Actor is a function that can take action on a device.
type Actor func(devices syscall.Handle, device setupapi.DevInfoData)
