package hosting

// RegionManager represents a service capable of getting info
// about Gandi Datacenters
type RegionManager interface {
	ListRegions() ([]Region, error)
	RegionbyCode(code string) (Region, error)
}

// Region represents a Gandi datacenter
type Region struct {
	ID      string
	Name    string
	Country string
}
