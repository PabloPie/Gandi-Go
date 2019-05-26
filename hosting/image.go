package hosting

// ImageManager represents a service capable of getting information about
// Gandi Disk images
type ImageManager interface {
	ImageByName(name string, region Region) (DiskImage, error)
	ListImagesInRegion(region Region) ([]DiskImage, error)
}

// DiskImage is an image offered by Gandi
// with an OS, used to create system Disks
//
// TODO: Add kernel version and other info
type DiskImage struct {
	ID       string
	DiskID   string
	RegionID string
	Name     string
	Size     int
}
