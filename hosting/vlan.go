package hosting

// VlanManager represents a service capable of manipulating
// private networks within Gandi's platform
type VlanManager interface {
	CreateVlan(vlan VlanSpec) (Vlan, error)

	// List return a list of Vlans, filtered with the options
	// given in the VlanFilter
	//
	// An unset field in `vlanfilter` is ignored when making the
	// request
	ListVlans(vlanfilter VlanFilter) ([]Vlan, error)
	UpdateVlanGW(vlan Vlan) (Vlan, error)
	RenameVlan(vlan Vlan) (Vlan, error)
	DeleteVlan(vlan Vlan) error
}

// Vlan represents a private network
type Vlan struct {
	ID       string
	Name     string
	Gateway  string
	Subnet   string
	RegionID string
}

// VlanSpec contains the information needed
// to create a private network
type VlanSpec struct {
	Name     string
	Gateway  string
	Subnet   string
	RegionID string
}

// VlanFilter is a struct to define filtering criteria
// when listing Vlan objects
type VlanFilter struct {
	ID       []string
	RegionID []string
	Name     string
}
