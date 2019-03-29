package hosting

type DiskManager interface {
	CreateDisk(disk DiskSpec) (Disk, error)
	// Maybe use a pointer receiver also, and populate the struct
	// InfoDisk(disk *Disk) error
	// `disk` needs to have an id
	// OR use this function as a filter to get a list of disks
	// We would need to define a DiskFilter
	InfoDisk(diskid int) (Disk, error)
	ListDisks() ([]Disk, error)
	DeleteDisk(disk *Disk) error

	// Extends `diskid` size by `size` (original size + `size`)
	// Disks cannot shrink in size
	ExtendDisk(disk *Disk, size uint) error

	RenameDisk(disk *Disk, name string) error
}

type ImageManager interface {
	ListImages(region Region) []DiskImage
	ImageByName(name string, region Region) int
}

type Disk struct {
	ID       int
	Name     string
	Size     int
	RegionID int
	State    string
	Type     string
	VM       []int
	BootDisk bool
}

type DiskSpec struct {
	RegionID int
	Name     string
	Size     int
}

type DiskImage struct {
	ID       int
	DiskID   int
	RegionID int
	Name     string
	Os       string
	Size     int
	State    int
}
