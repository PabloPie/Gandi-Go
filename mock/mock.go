package mock

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/PabloPie/Gandi-Go/client"
	"github.com/PabloPie/Gandi-Go/hosting"
)

type fn = func([]interface{}, interface{}) error

// we need to map every api function
var funcs = map[string]fn{
	"hosting.disk.update":      hostingDiskUpdate,
	"hosting.disk.list":        hostingDiskList,
	"hosting.disk.create":      hostingDiskCreate,
	"hosting.disk.create_from": hostingDiskCreateFrom,
	"hosting.disk.delete":      hostingDiskDelete,
	"hosting.image.list":       hostingImageList,
}

// Clientv4 is a mock client for Gandi's v4 API
type Clientv4 struct{}

// NewMockClientv4 creates a mock client for v4
func NewMockClientv4() client.V4Caller {
	return Clientv4{}
}

// Send invokes the correct method to treat the rpc call based on funcs. that translates
// an xmlrpc method to a go function
func (m Clientv4) Send(serviceMethod string, args []interface{}, reply interface{}) error {
	err := funcs[serviceMethod](args, reply)
	return err
}

// Disk functions
func hostingDiskUpdate(args []interface{}, reply interface{}) error {
	if len(args) != 2 {
		return errors.New("update() takes 2 arguments")
	}

	id, ok := args[0].(int)
	if !ok {
		return errors.New("Invalid method parameter: first argument must be an int")
	}
	diskupdate := reflect.ValueOf(args[1])
	if diskupdate.Kind() != reflect.Map {
		return errors.New("Invalid method parameter: second argument must be a struct")
	}

	if diskupdate.Len() == 0 {
		return errors.New("No valid parameter given in struct")
	}
	for _, e := range diskupdate.MapKeys() {
		key := e.Interface().(string)
		if key != "size" && key != "name" {
			return fmt.Errorf("Unknown %s key in struct", key)
		}
	}

	op := hosting.Operation{
		DiskID: id,
		Step:   "WAIT",
		Type:   "disk_update",
	}
	setValue(op, reply)
	return nil
}

func hostingDiskList(args []interface{}, reply interface{}) error {
	if len(args) > 1 {
		return errors.New("list() takes 1 optional argument")
	}
	if len(args) == 0 {
		// Do we keep a state or do we send random data?
		disks := []hosting.Disk{
			hosting.Disk{
				ID:       1,
				Name:     "sys_disk1",
				Size:     10240,
				RegionID: 3,
				State:    "created",
				Type:     "data",
				VM:       []int{1},
				BootDisk: true,
			},
			hosting.Disk{
				ID:       2,
				Name:     "sys_disk2",
				Size:     20480,
				RegionID: 3,
				State:    "created",
				Type:     "data",
				VM:       []int{2},
				BootDisk: true,
			}}
		setValue(disks, reply)
		return nil
	}
	// else we send the disks filtered by args[0]
	return nil
}

func hostingDiskCreate(args []interface{}, reply interface{}) error {
	if len(args) < 1 {
		return errors.New("create() takes 1 argument")
	}
	diskspec := reflect.ValueOf(args[0])
	if diskspec.Kind() != reflect.Struct {
		return errors.New("Invalid method parameter: first argument must be a struct")
	}
	_, ok := (diskspec.Interface()).(hosting.DiskSpec)
	if !ok {
		return errors.New("Struct provided is not a DiskSpec")
	}
	op := hosting.Operation{
		DiskID: 1,
		Step:   "WAIT",
		Type:   "disk_create",
	}
	setValue(op, reply)
	return nil
}

func hostingDiskCreateFrom(args []interface{}, reply interface{}) error {
	if len(args) < 2 {
		return errors.New("create_from() takes 2 argument")
	}
	srcdisk := reflect.ValueOf(args[1])
	if srcdisk.Kind() != reflect.Int {
		return errors.New("Invalid method parameter: second argument must be an integer")
	}
	args = []interface{}{args[0]}
	return hostingDiskCreate(args, reply)
}

func hostingDiskDelete(args []interface{}, reply interface{}) error {
	if len(args) != 1 {
		return errors.New("delete() takes 1 argument")
	}
	diskid := reflect.ValueOf(args[0])
	if diskid.Kind() != reflect.Int {
		return errors.New("Invalid method parameter: first argument must be an integer")
	}
	op := hosting.Operation{
		DiskID: diskid.Interface().(int),
		Step:   "WAIT",
		Type:   "disk_delete",
	}
	setValue(op, reply)
	return nil

}

// Pretty unsafe, add type checks at least
func setValue(str interface{}, reply interface{}) {
	replyPtrValue := reflect.ValueOf(reply)
	replyValue := reflect.Indirect(replyPtrValue)
	replyValue.Set(reflect.ValueOf(str))
}

// DiskImage functions
func hostingImageList(args []interface{}, reply interface{}) error {
	if len(args) > 1 {
		return errors.New("list_images() takes 1 optional parameter")
	}
	imagefilter := reflect.ValueOf(args[0])
	if imagefilter.Kind() != reflect.Map {
		return errors.New("Invalid method parameter: second argument must be a struct")
	}

	if imagefilter.Len() == 0 {
		setValue(diskimages, reply)
	}

	for _, e := range imagefilter.MapKeys() {
		key := e.Interface().(string)
		//assume he asked for Debian 9
		if key == "label" {
			setValue([]hosting.DiskImage{diskimages[1]}, reply)
			return nil
		}
		if key == "system" {
			setValue(diskimages[:2], reply)
			return nil
		}
	}
	return errors.New("Unknown key provided in filter")
}

var diskimages = []hosting.DiskImage{
	hosting.DiskImage{
		ID:       1,
		DiskID:   100,
		RegionID: 3,
		Name:     "Debian 8",
		Os:       "Debian",
		Version:  "8",
		Size:     10240,
	},
	hosting.DiskImage{
		ID:       2,
		DiskID:   101,
		RegionID: 3,
		Name:     "Debian 9",
		Os:       "Debian",
		Version:  "9",
		Size:     10240,
	},
	hosting.DiskImage{
		ID:       1,
		DiskID:   102,
		RegionID: 3,
		Name:     "FreeBSD 11.2 UFS",
		Os:       "FreeBSD",
		Version:  "11.2",
		Size:     10240,
	},
}
