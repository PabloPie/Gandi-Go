package hosting

// for a more readable code
type IPVersion int

const (
	IPv4             IPVersion = 4
	IPv6             IPVersion = 6
	DefaultBandwidth float32   = 102400.0
)

// IPManager represents a service capable of manipulating Gandi IPs
type IPManager interface {
	CreateIP(region Region, version IPVersion) (IPAddress, error)
	DescribeIP(ipfilter IPFilter) ([]IPAddress, error)
	DeleteIP(ip IPAddress) error
}

// IPAddress is the go representation of a Gandi IP
type IPAddress struct {
	ID       string
	IP       string
	RegionID string
	Version  IPVersion
	VM       string
	State    string
}

// IPFilter is given to DescribeIP to filter results
type IPFilter struct {
	ID       string
	RegionID string
	Version  IPVersion
	IP       string
}
