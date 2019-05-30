package hosting

// VlanManager represents a service capable of manipulating
// private networks within Gandi's platform
type VlanManager interface {
	CreateVlan(vlan VlanSpec) Vlan
	InfoVlan(vlan Vlan) Vlan
	ListVlan() []Vlan
	UpdateVlanGW(vlan Vlan) error
	RenameVlan(vlan Vlan) error
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
