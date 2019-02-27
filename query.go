package windevice

import (
	"io"

	"github.com/gentlemanautomaton/windevice/setupapi"
	"golang.org/x/sys/windows"
)

// Query holds device query information. Its zero value is a valid query for
// all devices.
type Query struct {
	Class      windows.GUID
	Enumerator string
	Flags      uint32
	Machine    string // TODO: Consider removing this if it's not well supported
	Selector   Selector
}

// Count returns the number of devices matching the query.
func (q Query) Count() (int, error) {
	var total int
	err := q.Each(func(Device) {
		total++
	})
	return total, err
}

// Each performs an action on each device that matches the query.
func (q Query) Each(action Actor) error {
	var classPtr *windows.GUID
	if q.Class != zeroGUID {
		classPtr = &q.Class
	}

	devices, err := setupapi.GetClassDevsEx(classPtr, q.Enumerator, q.Flags, 0, q.Machine)
	if err != nil {
		return err
	}
	defer setupapi.DestroyDeviceInfoList(devices)

	i := uint32(0)
	for {
		device, err := setupapi.EnumDeviceInfo(devices, i)
		switch err {
		case nil:
		case io.EOF:
			return nil
		default:
			return err
		}

		i++

		d := Device{
			list:  devices,
			entry: device,
		}

		if q.Selector != nil {
			matched, err := q.Selector.Select(d)
			if err != nil {
				return err
			}
			if !matched {
				continue
			}
		}

		action(d)
	}
}
