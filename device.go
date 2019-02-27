package windevice

import (
	"syscall"

	"github.com/gentlemanautomaton/windevice/deviceproperty"
	"github.com/gentlemanautomaton/windevice/setupapi"
)

// Device provides access to Windows device information while executing a query.
type Device struct {
	list  syscall.Handle
	entry setupapi.DevInfoData
}

// Sys returns the low-level device list handle and information data for the
// device.
func (device Device) Sys() (list syscall.Handle, entry setupapi.DevInfoData) {
	return device.list, device.entry
}

// Description returns the description of the device.
func (device Device) Description() (string, error) {
	return setupapi.GetDeviceRegistryString(device.list, device.entry, deviceproperty.Description)
}

// HardwareID returns the set of hardware IDs associated with the device.
func (device Device) HardwareID() ([]string, error) {
	return setupapi.GetDeviceRegistryStrings(device.list, device.entry, deviceproperty.HardwareID)
}

// CompatibleID returns the set of compatible IDs associated with the device.
func (device Device) CompatibleID() ([]string, error) {
	return setupapi.GetDeviceRegistryStrings(device.list, device.entry, deviceproperty.CompatibleID)
}

// Service returns the service for the device.
func (device Device) Service() (string, error) {
	return setupapi.GetDeviceRegistryString(device.list, device.entry, deviceproperty.Service)
}

// Class returns the class name of the device.
func (device Device) Class() (string, error) {
	return setupapi.GetDeviceRegistryString(device.list, device.entry, deviceproperty.Class)
}

// Driver returns the driver for the device.
func (device Device) Driver() (string, error) {
	return setupapi.GetDeviceRegistryString(device.list, device.entry, deviceproperty.Driver)
}

// Manufacturer returns the manufacturer of the device.
func (device Device) Manufacturer() (string, error) {
	return setupapi.GetDeviceRegistryString(device.list, device.entry, deviceproperty.Manufacturer)
}

// FriendlyName returns the friendly name of the device.
func (device Device) FriendlyName() (string, error) {
	return setupapi.GetDeviceRegistryString(device.list, device.entry, deviceproperty.FriendlyName)
}

// LocationInformation returns the location information for the device.
func (device Device) LocationInformation() (string, error) {
	return setupapi.GetDeviceRegistryString(device.list, device.entry, deviceproperty.LocationInformation)
}

// PhysicalDeviceObjectName returns the physical object name of the device.
func (device Device) PhysicalDeviceObjectName() (string, error) {
	return setupapi.GetDeviceRegistryString(device.list, device.entry, deviceproperty.PhysicalDeviceObjectName)
}

// EnumeratorName returns the name of the device's enumerator.
func (device Device) EnumeratorName() (string, error) {
	return setupapi.GetDeviceRegistryString(device.list, device.entry, deviceproperty.EnumeratorName)
}
