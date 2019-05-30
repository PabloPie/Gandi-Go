// Package hosting contains the interfaces and data structures that a user
// will use to interact with the lib
package hosting

// Hosting represents Gandi's API and contains every functionality
// implemented for the IaaS platform
//
// This interface defines a common behaviour for every version of Gandi's API
type Hosting interface {
	// VMManager is an interface containing the operations related to
	// Gandi Virtual Machines
	//
	// - VM Creation / Deletion / Update
	// - VM search
	// - Disk and IP attachment and detachment
	// - VM Stop / Start / Reboot
	VMManager

	// DiskManager is an interface containing the operations related to
	// Gandi Disks
	//
	// - Disk Creation / Deletion / Update
	// - Disk search
	DiskManager

	// IPManager is an interface containing the operations related to
	// Gandi IPs
	//
	// In the case of Hostingv4, it abstracts the Gandi Interfaces
	// - IP Creation / Deletion
	IPManager

	// SSHKeyManager is an interface containing the operations related
	// to SSH Keys in Gandi's platform
	//
	// - Key Creation / Deletion
	// - Key search from Name
	SSHKeyManager

	// XXX: Implement Vlan management
	// VlanManager

	// RegionManager is an interface containing the operations to
	// obtain information about Regions/Datacenters
	//
	// - Region listing / search from DC code
	RegionManager

	// ImageManager is an interface containing the operations to
	// obtain information about Gandi DiskImages
	//
	// - Listing images in a Region
	// - Searching images by name
	ImageManager
}
