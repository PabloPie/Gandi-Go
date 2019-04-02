package mock

import (
	"github.com/PabloPie/Gandi-Go/client"
)

type fn = func([]interface{}, interface{}) error

// we need to map every api function
var funcs = map[string]fn{
	// Disk
	"hosting.disk.update":      hostingDiskUpdate,
	"hosting.disk.list":        hostingDiskList,
	"hosting.disk.create":      hostingDiskCreate,
	"hosting.disk.create_from": hostingDiskCreateFrom,
	"hosting.disk.delete":      hostingDiskDelete,
	"hosting.image.list":       hostingImageList,
	// IP
	"hosting.ip.list":      hostingIPList,
	"hosting.iface.create": hostingIfaceCreate,
	"hosting.iface.delete": hostingIfaceDelete,
	"hosting.iface.info":   hostingIfaceInfo,
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
