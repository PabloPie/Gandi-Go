package hostingv4

import (
	"errors"
	"reflect"

	"github.com/PabloPie/Gandi-Go/client"
	"github.com/PabloPie/Gandi-Go/hosting"
)

// ErrNotProvided indicates that a value was not provided but was needed
var ErrNotProvided = errors.New("Not provided")

// ErrParse indicates that there was an error transforming a value from a struct,
// usually coming from a string to integer conversion
var ErrParse = errors.New("Parsing error")

// ErrMismatch indicates that two values that should be equal are not,
// for example when working when distinct objects that have to be in the
// same datacenter
var ErrMismatch = errors.New("Value mismatch")

// A Hostingv4 contains an xmlrpc client to send requests to
type Hostingv4 struct {
	client.V4Caller
}

// A HostingError records a failed Hosting operation
type HostingError struct {
	Func   string // the failing function
	Struct string // the struct concerned
	Field  string // the field concerned
	Err    error  // the reason the function failed
}

func (e *HostingError) Error() string {
	return "hostingv4." + e.Func + ": field " + e.Field + " in struct " + e.Struct + ": " + e.Err.Error()
}

func internalParseError(s string, f string) error {
	return &HostingError{"_internal_function", s, f, ErrParse}
}

// Newv4Hosting creates a new driver for Gandi v4 API
func Newv4Hosting(client client.V4Caller) hosting.Hosting {
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
		// get the zero value for the type of the field
		zeroval := reflect.Zero(field.Type).Interface()
		// omit the field if it is a zero value or has no name
		if reflect.DeepEqual(value, zeroval) || name == "" {
			continue
		}
		out[name] = value
	}
	return out, nil
}
