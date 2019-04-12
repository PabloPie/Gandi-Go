package hosting

// Hosting represents Gandi's api and contains every functionality
// implemented for the IaaS platform
type Hosting interface {
	// VMManager
	DiskManager
	// IPManager
	// SSHKeyManager
	// VlanManager
	RegionManager
	ImageManager
}
