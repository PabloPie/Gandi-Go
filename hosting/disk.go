package hosting

// DiskManager represents a service capable of manipulating Gandi Disks
type DiskManager interface {
	CreateDisk(disk DiskSpec) (Disk, error)
	CreateDiskFromImage(disk DiskSpec, src DiskImage) (Disk, error)
	ListAllDisks() ([]Disk, error)
	DiskFromName(name string) Disk
	ListDisks(diskFilter DiskFilter) ([]Disk, error)
	DeleteDisk(disk Disk) error
	ExtendDisk(disk Disk, size uint) (Disk, error)
	RenameDisk(disk Disk, name string) (Disk, error)
}

// Disk is a Gandi disk object
type Disk struct {
	ID       string
	Name     string
	Size     int
	RegionID string
	State    string
	Type     string
	VM       []string
	BootDisk bool
}

// DiskSpec contains the parameters to create a new Disk
//
// The only mandatory field is `RegionID`
type DiskSpec struct {
	RegionID string
	Name     string
	Size     int
}

// DiskFilter is used to search a list of disks
// filtered with the fields defined
//
// TODO: allow a list of elements
type DiskFilter struct {
	ID       string
	RegionID string
	Name     string
	VMID     string
}
