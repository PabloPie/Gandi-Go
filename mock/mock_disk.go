package mock

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/PabloPie/Gandi-Go/hosting"
)

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

var disks = []hosting.Disk{
	hosting.Disk{
		ID:       1,
		Name:     "sys_disk1",
		Size:     10240,
		RegionID: 4,
		State:    "created",
		Type:     "data",
		VM:       []int{1},
		BootDisk: true,
	}, hosting.Disk{
		ID:       4,
		Name:     "sys_disk3",
		Size:     10240,
		RegionID: 4,
		State:    "created",
		Type:     "data",
		VM:       []int{3},
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
	},
	hosting.Disk{
		ID:       3,
		Name:     "disk3",
		Size:     1024,
		RegionID: 3,
		State:    "created",
		Type:     "data",
		VM:       []int{2},
		BootDisk: false,
	},
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
		setValue(disks, reply)
		return nil
	}
	filter := reflect.ValueOf(args[0])
	if filter.Kind() != reflect.Map {
		return errors.New("Invalid method parameter: first agument must be a struct")
	}
	if filter.Len() == 0 {
		setValue(disks, reply)
		return nil
	}

	// vm_id == 2
	// datacenter_id == 4
	// name == disk3
	// id == 1
	keys := filter.MapKeys()
	if len(keys) == 1 {
		key := keys[0].Interface().(string)
		switch key {
		case "vm_id":
			setValue(disks[2:], reply)
		case "datacenter_id":
			setValue(disks[:2], reply)
		case "name":
			setValue([]hosting.Disk{disks[3]}, reply)
		case "id":
			setValue([]hosting.Disk{disks[1]}, reply)
		default:
			return errors.New("Unknown key provided in filter")
		}
	} else {
		setValue(disks, reply)
	}
	return nil
}

func hostingDiskCreate(args []interface{}, reply interface{}) error {
	if len(args) < 1 {
		return errors.New("disk.create() takes 1 argument")
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
		return errors.New("disk.create_from() takes 2 argument")
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
		return errors.New("disk.delete() takes 1 argument")
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

// DiskImage functions
func hostingImageList(args []interface{}, reply interface{}) error {
	if len(args) > 1 {
		return errors.New("image.list_images() takes 1 optional parameter")
	}
	// We always receive a struct, no need to verify args size
	imagefilter := reflect.ValueOf(args[0])
	if imagefilter.Kind() != reflect.Map {
		return errors.New("Invalid method parameter: second argument must be a struct")
	}

	if imagefilter.Len() == 0 {
		setValue(diskimages, reply)
		return nil
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
