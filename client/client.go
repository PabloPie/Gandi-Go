package client

import (
	"errors"

	"github.com/kolo/xmlrpc"
)

const (
	defaultV4URL = "https://rpc.gandi.net/xmlrpc/"
)

// V4Caller defines the methods a client for v4 of Gandi's API needs
type V4Caller interface {
	Send(method string, args []interface{}, reply interface{}) error
}

// Clientv4 represents a wrapper for an xmlrpc client that also
// includes the Gandi APIkey and the API url.
type Clientv4 struct {
	APIKey string
	URL    string

	*xmlrpc.Client
}

// NewClientv4 returns a client to connect to Gandi v4 xmlrpc API.
// If no URL is provided ("") default value is used, an api key is mandatory.
func NewClientv4(URL string, APIKey string) (V4Caller, error) {
	if APIKey == "" {
		return nil, errors.New("Apikey required but not provided")
	}
	if URL == "" {
		URL = defaultV4URL
	}

	client, err := xmlrpc.NewClient(URL, nil)
	if err != nil {
		return nil, err
	}

	return Clientv4{APIKey, URL, client}, nil
}

// Send invokes the named function, waits for it to complete, and returns its error status.
func (c Clientv4) Send(serviceMethod string, args []interface{}, reply interface{}) error {
	params := []interface{}{c.APIKey}
	if len(args) > 0 {
		params = append(params, args...)
	}
	return c.Call(serviceMethod, params, reply)
}
