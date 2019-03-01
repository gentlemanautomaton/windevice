package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gentlemanautomaton/windevice"
	"github.com/gentlemanautomaton/windevice/deviceclass"
	"github.com/gentlemanautomaton/windevice/devselect"
	"github.com/gentlemanautomaton/windevice/difuncremove"
	"github.com/gentlemanautomaton/windevice/strmatch"
)

func main() {
	var (
		className  string
		enumerator string
		id         string
		name       string
		machine    string
		present    bool
		detail     bool
		remove     bool
	)

	flag.StringVar(&className, "class", "", "include devices from a named device class")
	flag.StringVar(&enumerator, "enum", "", "include devices from a named PnP enumerator")
	flag.StringVar(&id, "id", "", "include devices with a particular hardware identifier")
	flag.StringVar(&name, "name", "", "include devices with a particular description or friendly name")
	flag.StringVar(&machine, "machine", "", "list devices on a remote machine")
	flag.BoolVar(&present, "present", false, "include devices that are present")
	flag.BoolVar(&detail, "detail", false, "print extra detail about each device")
	flag.BoolVar(&remove, "remove", false, "remove a single matched device")

	flag.Parse()

	q := windevice.DeviceQuery{
		Enumerator: enumerator,
		Machine:    machine,
	}

	if present {
		q.Flags = deviceclass.Present
	}

	var selectors []devselect.Selector
	if className != "" {
		class, err := windevice.NewNamedClass(className)
		if err != nil {
			fmt.Printf("Unable to retrieve named class \"%s\": %v\n", className, err)
			os.Exit(1)
		}
		if len(class.Members) == 1 {
			q.Class = class.Members[0]
		} else {
			selectors = append(selectors, devselect.Class(strmatch.EqualFold(className)))
		}
	}
	if id != "" {
		selectors = append(selectors, devselect.ID(strmatch.EqualFold(id)))
	}
	if name != "" {
		matcher := strmatch.ContainsInsensitive(name)
		selectors = append(selectors, devselect.Any(devselect.Description(matcher), devselect.FriendlyName(matcher)))
	}

	if len(selectors) > 0 {
		q.Selector = devselect.All(selectors...)
	}

	count, err := q.Count()
	if err != nil {
		fmt.Printf("Unable to retrieve device list: %v\n", err)
		os.Exit(1)
	}
	if count == 0 {
		fmt.Printf("No devices found.\n")
		return
	}

	if remove && count > 1 {
		fmt.Printf("   More than one device matched. Not removing.\n\n")
	}

	var index int
	q.Each(func(device windevice.Device) {
		if detail {
			printDetail(device, index)
		} else {
			printBasic(device, index)
		}
		if remove && count == 1 {
			removeDevice(device)
		}
		index++
	})
}

func removeDevice(device windevice.Device) {
	fmt.Printf("      --------\n")
	fmt.Printf("      Removing device...\n")

	devID, _ := device.DeviceInstanceID()

	needReboot, err := device.Remove(difuncremove.Global, 0)
	if err != nil {
		fmt.Printf("      Failed: %v\n", err)
	} else {
		fmt.Printf("      Successfully removed %s\n", devID)
		if needReboot {
			fmt.Printf("      A reboot is needed to complete the removal process\n")
		}
	}

	fmt.Printf("      --------\n")
}

func printBasic(device windevice.Device, index int) {
	if desc, err := device.Description(); err != nil {
		fmt.Printf(" %3d: Error: %v\n", index, err)
	} else {
		fmt.Printf(" %3d: %s\n", index, desc)
	}
}

func printDetail(device windevice.Device, index int) {
	desc, err := device.Description()
	if err != nil {
		fmt.Printf(" %3d: Error: %v\n", index, err)
		return
	}
	fmt.Printf(" %3d: Description: %s\n", index, desc)

	if id, _ := device.DeviceInstanceID(); id != "" {
		fmt.Printf("      Device Instance ID: %s\n", id)
	}
	if name, _ := device.FriendlyName(); name != "" {
		fmt.Printf("      Friendly Name: %s\n", name)
	}
	if class, _ := device.Class(); class != "" {
		fmt.Printf("      Class: %s\n", class)
	}
	if guid, _ := device.ClassGUID(); guid != "" {
		fmt.Printf("      Class GUID: %s\n", guid)
	}
	if enum, _ := device.EnumeratorName(); enum != "" {
		fmt.Printf("      Enumerator: %s\n", enum)
	}
	if location, _ := device.LocationInformation(); location != "" {
		fmt.Printf("      Location: %s\n", location)
	}
	if mfg, _ := device.Manufacturer(); mfg != "" {
		fmt.Printf("      Manufacturer: %s\n", mfg)
	}
	if phys, _ := device.PhysicalDeviceObjectName(); phys != "" {
		fmt.Printf("      Physical Device Object: %s\n", phys)
	}
	{
		drivers := device.InstalledDriver()

		count, err := drivers.Count()
		if err != nil {
			fmt.Printf("      Driver Enumeration Failed: %v\n", err)
		}

		switch count {
		case 0:
			fmt.Printf("      Drivers: None\n")
		case 1:
			drivers.Each(func(driver windevice.Driver) {
				fmt.Printf("      Driver: %s, Version: %s, Released: %s\n", driver.Description(), driver.Version(), driver.Date())
			})
		default:
			fmt.Printf("      Drivers:\n")
			driverIndex := 0
			drivers.Each(func(driver windevice.Driver) {
				fmt.Printf("        %2d: Description: %s\n", driverIndex, driver.Description())
				fmt.Printf("            Manufacturer: %s\n", driver.ManufacturerName())
				fmt.Printf("            Provider: %s\n", driver.ProviderName())
				fmt.Printf("            Date: %s\n", driver.Date())
				fmt.Printf("            Version: %s\n", driver.Version())
				driverIndex++
			})
		}
	}
	if service, _ := device.Service(); service != "" {
		fmt.Printf("      Service: %s\n", service)
	}
	if ids, _ := device.HardwareID(); len(ids) > 0 {
		for _, id := range ids {
			fmt.Printf("      Hardware ID: %s\n", id)
		}
	}
	if ids, _ := device.CompatibleID(); len(ids) > 0 {
		for _, id := range ids {
			fmt.Printf("      Compatible ID: %s\n", id)
		}
	}
	if flags, _ := device.ConfigFlags(); flags != 0 {
		fmt.Printf("      Flags: %x\n", flags)
	}
	if devType, err := device.DevType(); err == nil {
		fmt.Printf("      Device Type: %d\n", devType)
	}
	if characteristics, _ := device.Characteristics(); characteristics != 0 {
		fmt.Printf("      Characteristics: %x\n", characteristics)
	}
	if state, err := device.InstallState(); err == nil {
		fmt.Printf("      State: %s\n", state)
	}
}
