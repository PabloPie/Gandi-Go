package hosting

type RegionManager interface {
	ListRegions() []Region
	RegionbyCode() Region
}

type Region struct {
	ID    int
	Name  string
	State int
}
