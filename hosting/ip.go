package hosting

type IPManager interface {
	CreateIP(version int) []IPAddress
	DescribeIP(ipfilter IPFilter) []IPAddress
	DeleteIP(ipid int) error
}

type IPAddress struct {
	ID       int    `xmlrpc:"id"`
	IP       string `xmlrpc:"ip"`
	RegionID int    `xmlrpc:"datacenter_id"`
	Version  string `xmlrpc:"version"`
	VM       int    `xmlrpc:"vm_id"`
	State    string `xmlrpc:"state"`
	ifaceID  int    `xmlrpc:"iface_id"`
}

type IPFilter struct {
	ID       int
	RegionID int
	Version  string
	VMID     int
	IP       string
}

type IPAddressSpec struct {
	RegionID int    `xmlrpc:"datacenter_id"`
	Version  string `xmlrpc:"ip_version"`
}

type Iface struct {
	ID  int         `xmlrpc:"id"`
	IPS []IPAddress `xmlrpc:"ips"`
}
