package hosting

// ImageManager represents a service capable of getting information about
// Gandi Disk images
type ImageManager interface {

	// ImageByName returns the image with label `name` in the Region
	// provided, if it is a valid one
	ImageByName(name string, region Region) (DiskImage, error)

	// ListImagesInRegion returns every Image that can be found in
	// the Region provided
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
