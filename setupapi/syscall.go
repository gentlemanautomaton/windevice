package setupapi

import (
	"errors"
	"fmt"
	"io"
	"syscall"
	"unsafe"

	"github.com/gentlemanautomaton/windevice/deviceclass"
	"golang.org/x/sys/windows"
)

var (
	// ErrEmptyBuffer is returned when a nil or zero-sized buffer is provided
	// to a system call.
	ErrEmptyBuffer = errors.New("nil or empty buffer provided")

	// ErrInvalidRegistry is returned when an unexpected registry value type
	// is encountered.
	//ErrInvalidRegistry = errors.New("invalid registry type")

	// ErrInvalidData is returned when a property isn't present or isn't valid.
	ErrInvalidData = syscall.Errno(13)
)

var (
	modsetupapi = windows.NewLazySystemDLL("setupapi.dll")

	procSetupDiGetClassDevsExW           = modsetupapi.NewProc("SetupDiGetClassDevsExW")
	procSetupDiEnumDeviceInfo            = modsetupapi.NewProc("SetupDiEnumDeviceInfo")
	procSetupDiCreateDeviceInfoList      = modsetupapi.NewProc("SetupDiCreateDeviceInfoList")
	procSetupDiDestroyDeviceInfoList     = modsetupapi.NewProc("SetupDiDestroyDeviceInfoList")
	procSetupDiGetDeviceRegistryProperty = modsetupapi.NewProc("SetupDiGetDeviceRegistryPropertyW")
	procSetupDiClassGuidsFromNameEx      = modsetupapi.NewProc("SetupDiClassGuidsFromNameExW")
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

// GetDeviceRegistryString retrieves the requested property from the registry
// and returns it as a string.
func GetDeviceRegistryString(devices syscall.Handle, device DevInfoData, property uint32) (value string, err error) {
	var (
		b        [1024 * 2]byte
		buffer   = b[:]
		dataType uint32
	)

	// Make up to 3 attempts to get the property data.
	const rounds = 3
	for i := 0; i < rounds; i++ {
		var length uint32
		length, dataType, err = SetupDiGetDeviceRegistryProperty(devices, device, property, buffer)
		if err == nil {
			buffer = buffer[:length]
			break
		}
		if err == syscall.ERROR_INSUFFICIENT_BUFFER && i < rounds {
			buffer = make([]byte, length)
		} else {
			return "", err
		}
	}

	if err != nil {
		return "", err
	}

	switch dataType {
	case syscall.REG_SZ, syscall.REG_EXPAND_SZ:
		return utf16BytesToString(buffer), nil
	default:
		return "", fmt.Errorf("expected REG_SZ registry type but received type %d", dataType)
	}
}

// GetDeviceRegistryStrings retrieves the requested property from the registry
// and returns it as a slice of strings.
func GetDeviceRegistryStrings(devices syscall.Handle, device DevInfoData, property uint32) (values []string, err error) {
	var (
		b        [1024 * 2]byte
		buffer   = b[:]
		dataType uint32
	)

	// Make up to 3 attempts to get the property data.
	const rounds = 3
	for i := 0; i < rounds; i++ {
		var length uint32
		length, dataType, err = SetupDiGetDeviceRegistryProperty(devices, device, property, buffer)
		if err == nil {
			buffer = buffer[:length]
			break
		}
		if err == syscall.ERROR_INSUFFICIENT_BUFFER && i < rounds {
			buffer = make([]byte, length)
		} else {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	switch dataType {
	case syscall.REG_SZ, syscall.REG_EXPAND_SZ:
		return []string{utf16BytesToString(buffer)}, nil
	case syscall.REG_MULTI_SZ:
		return utf16BytesToSplitString(buffer), nil
	default:
		return nil, fmt.Errorf("expected REG_MULTI_SZ registry type but received type %d", dataType)
	}
}

// SetupDiGetDeviceRegistryProperty retrieves a property from a member of a
// device information set.
//
// https://docs.microsoft.com/en-us/windows/desktop/api/setupapi/nf-setupapi-setupdigetdeviceregistrypropertyw
func SetupDiGetDeviceRegistryProperty(devices syscall.Handle, device DevInfoData, property uint32, buffer []byte) (reqSize uint32, registryDataType uint32, err error) {
	if len(buffer) == 0 {
		return 0, 0, ErrEmptyBuffer
	}

	device.Size = uint32(unsafe.Sizeof(device))

	r0, _, e := syscall.Syscall9(
		procSetupDiGetDeviceRegistryProperty.Addr(),
		7,
		uintptr(devices),
		uintptr(unsafe.Pointer(&device)),
		uintptr(property),
		uintptr(unsafe.Pointer(&registryDataType)),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(len(buffer)),
		uintptr(unsafe.Pointer(&reqSize)),
		0,
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
