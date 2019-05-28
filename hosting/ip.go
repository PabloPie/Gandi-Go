package hosting

// IPVersion represents the possible versions of an ip
//
// limits possible input parameters
type IPVersion int

const (
	// IPv4 is just the int value 4
	IPv4 IPVersion = 4
	// IPv6 is just the int value 6
	IPv6 IPVersion = 6
	// DefaultBandwidth is the bandwidth an interface
	// is created with, so it doesn't default to 1Mbps
	DefaultBandwidth float32 = 102400.0
)

// IPManager represents a service capable of manipulating Gandi IPs
type IPManager interface {
	CreateIP(region Region, version IPVersion) (IPAddress, error)
	ListIPs(ipfilter IPFilter) ([]IPAddress, error)
	DeleteIP(ip IPAddress) error
}

// IPAddress represents a Gandi IP
//
// It abstracts away the Gandi network interfaces
// so the user only has to worry about ips
type IPAddress struct {
	ID       string
	IP       string
	RegionID string
	Version  IPVersion
	VM       string
	State    string
}

// IPFilter is used to list IPs, filtered
// with the parameters provided
type IPFilter struct {
	ID       string
	RegionID string
	Version  IPVersion
	IP       string
}
