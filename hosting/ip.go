package hosting

type IPManager interface {
	CreateIP(version int) IPAddress
	// Same as disk and vm, make it a filter function
	InfoIP(ipid int) IPAddress
	ListIP() []IPAddress
	DeleteIP(ipid int) error
}

type IPAddress struct {
	ID       int
	IP       string
	RegionID int
	Version  int
	VM       int
	State    string
	Vlan     int // if private
}

type IPAddressSpec struct {
	RegionID int
	Version  int
	Vlan     int
	IP       int // if vlan defined
}
