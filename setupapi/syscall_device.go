package setupapi

import (
	"io"
	"syscall"
	"unsafe"

	"github.com/gentlemanautomaton/windevice/deviceclass"
	"golang.org/x/sys/windows"
)

var (
	procSetupDiEnumDeviceInfo        = modsetupapi.NewProc("SetupDiEnumDeviceInfo")
	procSetupDiCreateDeviceInfoList  = modsetupapi.NewProc("SetupDiCreateDeviceInfoList")
	procSetupDiDestroyDeviceInfoList = modsetupapi.NewProc("SetupDiDestroyDeviceInfoList")
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
