package hosting

// RegionManager represents a service capable of getting info
// about Gandi Datacenters
type RegionManager interface {

	// Lists every existing Region
	ListRegions() ([]Region, error)

	// Return a Region object given its datacenter code
	RegionbyCode(code string) (Region, error)
}

// Region represents a Gandi datacenter
type Region struct {
	ID      string
	Name    string
	Country string
}
