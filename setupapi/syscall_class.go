package setupapi

import (
	"syscall"
	"unsafe"

	"github.com/gentlemanautomaton/windevice/deviceclass"
	"golang.org/x/sys/windows"
)

var (
	procSetupDiGetClassDevsExW      = modsetupapi.NewProc("SetupDiGetClassDevsExW")
	procSetupDiClassGuidsFromNameEx = modsetupapi.NewProc("SetupDiClassGuidsFromNameExW")
)

// SetupDiGetClassDevsEx prepares a device information set.
//
// https://docs.microsoft.com/en-us/windows/desktop/api/setupapi/nf-setupapi-setupdigetclassdevsexw
func SetupDiGetClassDevsEx(guid *windows.GUID, enumerator string, flags uint32, hDevInfoSet syscall.Handle, machineName string) (handle syscall.Handle, err error) {
	var ep *uint16
	if enumerator != "" {
		ep, err = syscall.UTF16PtrFromString(enumerator)
		if err != nil {
			return syscall.InvalidHandle, err
		}
	}

	var mnp *uint16
	if machineName != "" {
		mnp, err = syscall.UTF16PtrFromString(machineName)
		if err != nil {
			return syscall.InvalidHandle, err
		}
	}

	if guid == nil {
		flags |= deviceclass.AllClasses
	}

	r0, _, e := syscall.Syscall9(
		procSetupDiGetClassDevsExW.Addr(),
		7,
		uintptr(unsafe.Pointer(guid)),
		uintptr(unsafe.Pointer(ep)),
		0, // hwndParent
		uintptr(flags),
		uintptr(hDevInfoSet),
		uintptr(unsafe.Pointer(mnp)),
		0,
		0,
		0)
	handle = syscall.Handle(r0)
	if handle == syscall.InvalidHandle {
		if e != 0 {
			err = syscall.Errno(e)
		} else {
			err = syscall.EINVAL
		}
	}
	return
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
