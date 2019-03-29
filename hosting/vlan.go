package hosting

type VlanManager interface {
	CreateVlan(vlan VlanSpec) Vlan
	// another filter
	InfoVlan(vlanid int) Vlan
	ListVlan() []Vlan
	UpdateVlanGW(vlan *Vlan) error
	RenameVlan(vlan *Vlan) error
}

type Vlan struct {
	ID       int
	Name     string
	Gateway  string
	Subnet   string
	RegionID int
}

type VlanSpec struct {
	Name     string
	Gateway  string
	Subnet   string
	RegionID int
}
