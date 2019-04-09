package hosting

import (
	"errors"
	"reflect"

	"github.com/PabloPie/Gandi-Go/client"
)

// Hosting represents Gandi's api and contains every functionality
// implemented for the IaaS platform
type Hosting interface {
	// VMManager
	DiskManager
	IPManager
	// SSHKeyManager
	// VlanManager
	RegionManager
	ImageManager
}

type Hostingv4 struct {
	client.V4Caller
}

func Newv4Hosting(client client.V4Caller) Hosting {
	return Hostingv4{client}
}

// structToMap writes every tagged(with key "xmlrpc"), non Zero Value
// from the struct into a map
func structToMap(s interface{}) (map[string]interface{}, error) {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Struct {
		return nil, errors.New("Not a struct")
	}

	out := map[string]interface{}{}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := field.Tag.Get("xmlrpc")
		value := v.Field(i).Interface()
		zeroval := reflect.Zero(field.Type).Interface()
		if reflect.DeepEqual(value, zeroval) || name == "" {
			continue
		}
		out[name] = value
	}
	return out, nil
}
