package setupapi

import (
	"io"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	procSetupDiEnumDeviceInfo        = modsetupapi.NewProc("SetupDiEnumDeviceInfo")
	procSetupDiCreateDeviceInfoList  = modsetupapi.NewProc("SetupDiCreateDeviceInfoList")
	procSetupDiDestroyDeviceInfoList = modsetupapi.NewProc("SetupDiDestroyDeviceInfoList")
)

// SetupDiEnumDeviceInfo returns device information about a member of a device
// information set. It returns io.EOF if there are not more members in the
// set.
//
// https://docs.microsoft.com/en-us/windows/desktop/api/setupapi/nf-setupapi-setupdienumdeviceinfo
func SetupDiEnumDeviceInfo(devices syscall.Handle, index uint32) (info DevInfoData, err error) {
	const errNoMoreItems = 259

	info.Size = uint32(unsafe.Sizeof(info))

	r0, _, e := syscall.Syscall(
		procSetupDiEnumDeviceInfo.Addr(),
		3,
		uintptr(devices),
		uintptr(index),
		uintptr(unsafe.Pointer(&info)))
	if r0 == 0 {
		switch e {
		case 0:
			err = syscall.EINVAL
		case errNoMoreItems:
			err = io.EOF
		default:
			err = syscall.Errno(e)
		}
	}
	return
}

// SetupDiCreateDeviceInfoList creates a device info list.
//
// https://docs.microsoft.com/en-us/windows/desktop/api/setupapi/nf-setupapi-setupdicreatedeviceinfolist
func SetupDiCreateDeviceInfoList(guid *windows.GUID) (handle syscall.Handle, err error) {
	r0, _, e := syscall.Syscall(
		procSetupDiCreateDeviceInfoList.Addr(),
		2,
		uintptr(unsafe.Pointer(guid)),
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

// SetupDiDestroyDeviceInfoList destroys a device info list.
//
// https://docs.microsoft.com/en-us/windows/desktop/api/setupapi/nf-setupapi-setupdidestroydeviceinfolist
func SetupDiDestroyDeviceInfoList(devices syscall.Handle) error {
	r0, _, e := syscall.Syscall(
		procSetupDiDestroyDeviceInfoList.Addr(),
		1,
		uintptr(devices),
		0,
		0)
	if r0 == 0 {
		if e != 0 {
			return syscall.Errno(e)
		}
		return syscall.EINVAL
	}
	return nil
}

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
