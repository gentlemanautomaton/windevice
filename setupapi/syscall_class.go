package setupapi

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	procSetupDiGetClassDevsExW      = modsetupapi.NewProc("SetupDiGetClassDevsExW")
	procSetupDiClassGuidsFromNameEx = modsetupapi.NewProc("SetupDiClassGuidsFromNameExW")
)

// SetupDiClassGuidsFromNameEx returns the list of GUIDs associated with
// a class name.
//
// https://docs.microsoft.com/en-us/windows/desktop/api/setupapi/nf-setupapi-setupdiclassguidsfromnameexw
func SetupDiClassGuidsFromNameEx(className, machine string) (guids []windows.GUID, err error) {
	cp, err := syscall.UTF16PtrFromString(className)
	if err != nil {
		return nil, err
	}

	var mp *uint16
	if machine != "" {
		mp, err = syscall.UTF16PtrFromString(machine)
		if err != nil {
			return nil, err
		}
	}

	guids = make([]windows.GUID, 1)

	// Make up to 3 attempts to get the class data.
	const rounds = 3
	for i := 0; i < rounds; i++ {
		var length uint32
		length, err = setupDiClassGuidsFromNameEx(cp, mp, guids)
		if err == nil {
			if length == 0 {
				return nil, nil
			}
			return guids, nil
		}
		if err == syscall.ERROR_INSUFFICIENT_BUFFER && i < rounds {
			guids = make([]windows.GUID, length)
		} else {
			return nil, err
		}
	}

	return nil, syscall.ERROR_INSUFFICIENT_BUFFER
}

func setupDiClassGuidsFromNameEx(className, machine *uint16, buffer []windows.GUID) (reqSize uint32, err error) {
	var gp *windows.GUID
	if len(buffer) > 0 {
		gp = &buffer[0]
	}
	r0, _, e := syscall.Syscall6(
		procSetupDiClassGuidsFromNameEx.Addr(),
		6,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(gp)),
		uintptr(len(buffer)),
		uintptr(unsafe.Pointer(&reqSize)),
		uintptr(unsafe.Pointer(machine)),
		0)
	if r0 == 0 {
		if e != 0 {
			err = syscall.Errno(e)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
