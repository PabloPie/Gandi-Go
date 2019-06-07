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

	// CreateIP creates an IPv4 or an IPv6 in the Region given
	//
	// The IP created can only be public
	CreateIP(region Region, version IPVersion) (IPAddress, error)

	// CreatePrivateIP creates a private IPv4 within a Vlan
	//
	// Region is inferred from the vlan provided, so it is mandatory
	// that it contains a valid RegionID
	CreatePrivateIP(vlan Vlan, ip string) (IPAddress, error)

	// ListIPs return a list of IPs, filtered with the options
	// given in the IPFilter
	//
	// An unset field in `ipfilter` is ignored when making the
	// request
	ListIPs(ipfilter IPFilter) ([]IPAddress, error)

	// DeleteIP deletes the IP given, provided it has
	// a valid ID
	//
	// If the operation was successful a nil error is returned
	DeleteIP(ip IPAddress) error
}

// IPAddress represents a Gandi IP
//
// In the case of Hostingv4, it abstracts away
// the Gandi network interfaces so the user only
// has to manage IPs
type IPAddress struct {
	// ID of the object in the API
	ID string

	// The actual IP address of the object
	IP string

	// ID of the region the IP is in
	RegionID string

	// Version of the IP, either 4 or 6
	Version IPVersion

	// The VM this IP is attached to
	VM string

	// State of the IP: used, free
	State string
}

// IPFilter is used to list IPs, filtered
// with the parameters provided
type IPFilter struct {
	ID       string
	RegionID string
	Version  IPVersion
	IP       string
}
