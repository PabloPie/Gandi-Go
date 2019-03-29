package hosting

import (
	"github.com/PabloPie/Gandi-Go/client"
	"github.com/PabloPie/Gandi-Go/hosting"
)

// Hosting represents Gandi's api and contains every functionality
// implemented for the IaaS platform
type Hosting interface {
	hosting.VMManager
	hosting.DiskManager
	hosting.IPManager
	hosting.SSHKeyManager
	hosting.VlanManager
	hosting.RegionManager
	hosting.ImageManager
}

func Newv4Hosting(client client.V4Caller) Hosting {
	return nil
}

func NewV5Hosting(client client.V5Caller) Hosting {
	return nil
}
