package hostingv4

import (
	"errors"
	"strconv"

	"github.com/PabloPie/go-gandi/hosting"
)

// Private representation of an ip object in v4 of
// Gandi's API, where IDs are integers instead of strings
type iPAddressv4 struct {
	ID       int    `xmlrpc:"id"`
	IP       string `xmlrpc:"ip"`
	RegionID int    `xmlrpc:"datacenter_id"`
	Version  int    `xmlrpc:"version"`
	VM       int    `xmlrpc:"vm_id"`
	State    string `xmlrpc:"state"`
}

// Internally, ips are associated to interfaces, even though
// we abstract those away, we need an internal object for
// API responses
type iface struct {
	IPs      []iPAddressv4 `xmlrpc:"ips"`
	RegionID int           `xmlrpc:"datacenter_id"`
	ID       int           `xmlrpc:"id"`
	VMID     int           `xmlrpc:"vm_id"`
}

// CreateIP creates an ip object that represents a public IP, either v4 or v6
//
// It requires a valid Region object, whose only mandatory field is its ID
// An ipv6 is always created for the interface, even when only an ipv4 is requested
func (h Hostingv4) CreateIP(region hosting.Region, version hosting.IPVersion) (hosting.IPAddress, error) {
	if version != hosting.IPv4 && version != hosting.IPv6 {
		return hosting.IPAddress{}, errors.New("Bad IP version")
	}

	regionID, err := strconv.Atoi(region.ID)
	if err != nil {
		return hosting.IPAddress{}, internalParseError("Region", "ID")
	}

	var ip iPAddressv4
	var response = Operation{}
	err = h.Send("hosting.iface.create", []interface{}{
		map[string]interface{}{
			"datacenter_id": regionID,
			"ip_version":    int(version),
			"bandwidth":     hosting.DefaultBandwidth,
		}}, &response)
	if err != nil {
		return hosting.IPAddress{}, err
	}
	if err = h.waitForOp(response); err != nil {
		return hosting.IPAddress{}, err
	}

	if err = h.Send("hosting.ip.info", []interface{}{response.IPID}, &ip); err != nil {
		return hosting.IPAddress{}, err
	}

	return toIPAddress(ip), nil
}

// CreatePrivateIP creates a private IP within a specified vlan
func (h Hostingv4) CreatePrivateIP(vlan hosting.Vlan, ip string) (hosting.IPAddress, error) {
	var fn = "CreatePrivateIP"
	if vlan.RegionID == "" || vlan.ID == "" {
		return hosting.IPAddress{}, &HostingError{fn, "Vlan", "ID/RegionID", ErrNotProvided}
	}
	regionid, err := strconv.Atoi(vlan.RegionID)
	if err != nil {
		return hosting.IPAddress{}, &HostingError{fn, "Vlan", "RegionID", ErrParse}
	}
	vlanid, err := strconv.Atoi(vlan.ID)
	if err != nil {
		return hosting.IPAddress{}, &HostingError{fn, "Vlan", "ID", ErrParse}
	}

	var ipv4 iPAddressv4
	var response = Operation{}
	err = h.Send("hosting.iface.create", []interface{}{
		map[string]interface{}{
			"datacenter_id": regionid,
			"bandwidth":     hosting.DefaultBandwidth,
			"ip":            ip,
			"vlan":          vlanid,
		}}, &response)
	if err != nil {
		return hosting.IPAddress{}, err
	}
	if err = h.waitForOp(response); err != nil {
		return hosting.IPAddress{}, err
	}

	if err = h.Send("hosting.ip.info", []interface{}{response.IPID}, &ipv4); err != nil {
		return hosting.IPAddress{}, err
	}

	return toIPAddress(ipv4), nil
}

// ListIPs returns a list of ips filtered with the options provided in `diskFilter`
func (h Hostingv4) ListIPs(ipfilter hosting.IPFilter) ([]hosting.IPAddress, error) {
	ipmap, err := ipFilterToMap(ipfilter)
	if err != nil {
		return nil, err
	}

	var response = []iPAddressv4{}
	if err = h.Send("hosting.ip.list", []interface{}{ipmap}, &response); err != nil {
		return nil, err
	}

	var ips []hosting.IPAddress
	for _, iip := range response {
		ips = append(ips, toIPAddress(iip))
	}

	return ips, nil
}

// DeleteIP deletes an IP Address
//
// It will also delete the associated interface, so if it is
// an ipv4, the corresponding ipv6 will also be deleted
func (h Hostingv4) DeleteIP(ip hosting.IPAddress) error {
	ipid, err := strconv.Atoi(ip.ID)
	if err != nil {
		return internalParseError("hosting.IPAddress", "ID")
	}

	var response = Operation{}
	err = h.Send("hosting.ip.info", []interface{}{ipid}, &response)
	if err != nil {
		return err
	}
	err = h.Send("hosting.iface.delete", []interface{}{response.IfaceID}, &response)
	if err != nil {
		return err
	}

	return h.waitForOp(response)
}

// Get the interface associated to a specific IP
func (h Hostingv4) ifaceIDFromIPID(ipid int) (int, error) {
	// An operation already contains a field for iface_id
	// we avoid defining a new struct
	response := Operation{}
	err := h.Send("hosting.ip.info", []interface{}{ipid}, &response)
	if err != nil {
		return 0, err
	}
	return response.IfaceID, nil
}

// Helper function to get an IP object from its v4 ID
func (h Hostingv4) ipFromID(ipid int) (hosting.IPAddress, error) {
	response := iPAddressv4{}
	err := h.Send("hosting.ip.info", []interface{}{ipid}, &response)
	if err != nil {
		return hosting.IPAddress{}, err
	}
	return toIPAddress(response), nil
}

// Internal methods to convert Hosting structures to v4 structures

func ipFilterToMap(ipfilter hosting.IPFilter) (map[string]interface{}, error) {
	ipmap := make(map[string]interface{})
	var err error

	if ipfilter.Version != 0 {
		if ipfilter.Version != hosting.IPv4 && ipfilter.Version != hosting.IPv6 {
			return nil, internalParseError("hosting.IPFilter", "Version")
		}
		ipmap["version"] = int(ipfilter.Version)
	}

	if ipfilter.ID != "" {
		ipmap["id"], err = strconv.Atoi(ipfilter.ID)
		if err != nil {
			return nil, internalParseError("hosting.IPFilter", "ID")
		}
	}
	if ipfilter.RegionID != "" {
		ipmap["datacenter_id"], err = strconv.Atoi(ipfilter.RegionID)
		if err != nil {
			return nil, internalParseError("hosting.IPFilter", "ID")
		}
	}

	if ipfilter.IP != "" {
		ipmap["ip"] = ipfilter.IP
	}

	return ipmap, nil
}

// v4 API IP -> Hosting hosting.IPAddress
func toIPAddress(iip iPAddressv4) (ip hosting.IPAddress) {
	ip.ID = strconv.Itoa(iip.ID)
	ip.IP = iip.IP
	ip.RegionID = strconv.Itoa(iip.RegionID)
	ip.State = iip.State
	ip.VM = strconv.Itoa(iip.VM)

	if iip.Version == 6 {
		ip.Version = hosting.IPv6
	} else {
		ip.Version = hosting.IPv4
	}

	return
}
