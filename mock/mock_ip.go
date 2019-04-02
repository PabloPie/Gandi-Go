package mock

import (
	"errors"
	"reflect"

	"github.com/PabloPie/Gandi-Go/hosting"
)

var ips = []hosting.IPAddress{
	hosting.IPAddress{
		ID:       1,
		IP:       "2001:4b98:dc0:41:216:3eff:fe9b:1c39",
		RegionID: 3,
		Version:  "6",
		VM:       1,
		State:    "created",
	},
	hosting.IPAddress{
		ID:       2,
		IP:       "46.226.108.29",
		RegionID: 3,
		Version:  "4",
		VM:       1,
		State:    "created",
	},
	hosting.IPAddress{
		ID:       3,
		IP:       "2001:4b98:dc2:41:216:3eff:fea8:c071",
		RegionID: 4,
		Version:  "6",
		VM:       2,
		State:    "created",
	},
}

func hostingIfaceCreate(args []interface{}, reply interface{}) error {
	if len(args) < 1 {
		return errors.New("iface.create() takes 1 argument")
	}
	ifacespec := reflect.ValueOf(args[0])
	if ifacespec.Kind() != reflect.Struct {
		return errors.New("Invalid method parameter: first argument must be a struct")
	}
	_, ok := (ifacespec.Interface()).(hosting.IPAddressSpec)
	if !ok {
		return errors.New("Struct provided is not a IPAddressSpec")
	}
	op := hosting.Operation{
		IfaceID: 1,
		Step:    "WAIT",
		Type:    "iface_create",
	}
	setValue(op, reply)
	return nil
}

func hostingIfaceDelete(args []interface{}, reply interface{}) error {
	if len(args) != 1 {
		return errors.New("iface.delete() takes 1 argument")
	}
	ifaceid := reflect.ValueOf(args[0])
	if ifaceid.Kind() != reflect.Int {
		return errors.New("Invalid method parameter: first argument must be an integer")
	}
	op := hosting.Operation{
		IfaceID: ifaceid.Interface().(int),
		Step:    "WAIT",
		Type:    "disk_delete",
	}
	setValue(op, reply)
	return nil
}

// Asking for iface_id 1 returns an ipv4 and an ipv6,
// asking for iface_id 2 returns an ipv6 only
func hostingIfaceInfo(args []interface{}, reply interface{}) error {
	if len(args) != 1 {
		return errors.New("iface.info() takes 1 argument")
	}
	ifaceid := reflect.ValueOf(args[0])
	if ifaceid.Kind() != reflect.Int {
		return errors.New("Invalid method parameter: first argument must be an integer")
	}
	ifaceidInt := ifaceid.Interface().(int)

	iface := hosting.Iface{
		ID:  ifaceidInt,
		IPS: ips[:2],
	}
	if ifaceidInt == 2 {
		iface.IPS = ips[2:3]
	}

	setValue(iface, reply)
	return nil
}

func hostingIPList(args []interface{}, reply interface{}) error {
	if len(args) > 1 {
		return errors.New("ip.list() takes 1 optional argument")
	}
	if len(args) == 0 {
		setValue(ips, reply)
	}
	filter := reflect.ValueOf(args[0])
	if filter.Kind() != reflect.Map {
		return errors.New("Invalid method parameter: first agument must be a struct")
	}

	if filter.Len() == 0 {
		setValue(ips, reply)
	}

	// version == 6
	// vm_id == 1
	// ip == 46.226.108.29
	// datacenter_id == 4
	keys := filter.MapKeys()
	if len(keys) == 1 {
		key := keys[0].Interface().(string)
		switch key {
		case "version":
			setValue([]hosting.IPAddress{ips[0], ips[2]}, reply)
		case "datacenter_id":
			setValue([]hosting.IPAddress{ips[2]}, reply)
		case "ip":
			setValue([]hosting.IPAddress{ips[1]}, reply)
		case "vm_id":
			setValue(ips[:2], reply)
		case "id":
		default:
			return errors.New("Unknown key provided in filter")
		}
	} else {
		setValue(ips, reply)
	}
	return nil
}
