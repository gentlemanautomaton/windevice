package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gentlemanautomaton/windevice"
	"github.com/gentlemanautomaton/windevice/deviceclass"
	"github.com/gentlemanautomaton/windevice/devselect"
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
	)

	flag.StringVar(&className, "class", "", "include devices from a named device class")
	flag.StringVar(&enumerator, "enum", "", "include devices from a named PnP enumerator")
	flag.StringVar(&id, "id", "", "include devices with a particular hardware identifier")
	flag.StringVar(&name, "name", "", "include devices with a particular description or friendly name")
	flag.StringVar(&machine, "machine", "", "list devices on a remote machine")
	flag.BoolVar(&present, "present", false, "include devices that are present")
	flag.BoolVar(&detail, "detail", false, "print extra detail about each device")

	flag.Parse()

	q := windevice.Query{
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
		matcher := strmatch.Contains(name)
		selectors = append(selectors, devselect.Any(devselect.Description(matcher), devselect.FriendlyName(matcher)))
	}

	if len(selectors) > 0 {
		q.Selector = devselect.All(selectors...)
	}

	var index int
	q.Each(func(device windevice.Device) {
		if detail {
			printDetail(device, index)
		} else {
			printBasic(device, index)
		}
		index++
	})
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

	if name, _ := device.FriendlyName(); name != "" {
		fmt.Printf("      Friendly Name: %s\n", name)
	}
	if class, _ := device.Class(); class != "" {
		fmt.Printf("      Class: %s\n", class)
	}
	if enum, _ := device.EnumeratorName(); enum != "" {
		fmt.Printf("      Enumerator: %s\n", enum)
	}
	if location, _ := device.LocationInformation(); err == nil && location != "" {
		fmt.Printf("      Location: %s\n", location)
	}
	if mfg, _ := device.Manufacturer(); err == nil && mfg != "" {
		fmt.Printf("      Manufacturer: %s\n", mfg)
	}
	if phys, _ := device.PhysicalDeviceObjectName(); err == nil && phys != "" {
		fmt.Printf("      Physical Device Object: %s\n", phys)
	}
	if driver, _ := device.Driver(); err == nil && driver != "" {
		fmt.Printf("      Driver: %s\n", driver)
	}
	if service, _ := device.Service(); err == nil && service != "" {
		fmt.Printf("      Service: %s\n", service)
	}
	if ids, err := device.HardwareID(); err == nil && len(ids) > 0 {
		for _, id := range ids {
			fmt.Printf("      Hardware ID: %s\n", id)
		}
	}
}
