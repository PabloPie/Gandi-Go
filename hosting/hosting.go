package hosting

import (
	"github.com/PabloPie/Gandi-Go/client"
)

// Hosting represents Gandi's api and contains every functionality
// implemented for the IaaS platform
type Hosting interface {
	// VMManager
	DiskManager
	IPManager
	// SSHKeyManager
	// VlanManager
	// RegionManager
	ImageManager
}

type Hostingv4 struct {
	client.V4Caller
}

func Newv4Hosting(client client.V4Caller) Hosting {
	return Hostingv4{client}
}
