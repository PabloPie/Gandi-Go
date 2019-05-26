package hosting

// Hosting represents Gandi's API and contains every functionality
// implemented for the IaaS platform
//
// This interface defines a common behaviour for every version of Gandi's API
type Hosting interface {
	VMManager
	DiskManager
	IPManager
	SSHKeyManager
	// VlanManager
	RegionManager
	ImageManager
}
