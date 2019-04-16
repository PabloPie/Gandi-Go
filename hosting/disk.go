package hosting

// DiskManager represents a service capable of manipulating Gandi Disks
type DiskManager interface {
	CreateDisk(disk DiskSpec) (Disk, error)
	CreateDiskFromImage(disk DiskSpec, src DiskImage) (Disk, error)
	ListDisks() ([]Disk, error)
	DiskFromName(name string) Disk
	DescribeDisks(diskFilter DiskFilter) ([]Disk, error)
	DeleteDisk(disk Disk) error
	ExtendDisk(disk Disk, size uint) (Disk, error)
	RenameDisk(disk Disk, name string) (Disk, error)
}

// Disk is a Gandi disk object
type Disk struct {
	ID       string
	Name     string
	Size     uint
	RegionID string
	State    string
	Type     string
	VM       []string
	BootDisk bool
}

// DiskSpec contains the parameters to create a new Disk
type DiskSpec struct {
	RegionID string
	Name     string
	Size     uint
}

//DiskFilter is used to filter the results DescribeDisks returns,
type DiskFilter struct {
	ID       string
	RegionID string
	Name     string
	VMID     string
}
