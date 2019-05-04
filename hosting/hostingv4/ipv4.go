package hostingv4

import (
	"errors"
	"strconv"

	"github.com/PabloPie/Gandi-Go/hosting"
)

type (
	IPAddress = hosting.IPAddress
	IPFilter  = hosting.IPFilter
)

type iPAddressv4 struct {
	ID       int    `xmlrpc:"id"`
	IP       string `xmlrpc:"ip"`
	RegionID int    `xmlrpc:"datacenter_id"`
	Version  int    `xmlrpc:"version"`
	VM       int    `xmlrpc:"vm_id"`
	State    string `xmlrpc:"state"`
}

type iface struct {
	IPs      []iPAddressv4 `xmlrpc:"ips"`
	RegionID int           `xmlrpc:"datacenter_id"`
	ID       int           `xmlrpc:"id"`
	VMID     int           `xmlrpc:"vm_id"`
}

func (h Hostingv4) CreateIP(region Region, version hosting.IPVersion) (IPAddress, error) {
	if version != hosting.IPv4 && version != hosting.IPv6 {
		return IPAddress{}, errors.New("Bad IP version")
	}

	var err error
	var iip iPAddressv4
	var region_id_int int
	var response = Operation{}

	region_id_int, err = strconv.Atoi(region.ID)
	if err != nil {
		return IPAddress{}, internalParseError("Region", "ID")
	}

	err = h.Send("hosting.iface.create", []interface{}{
		map[string]interface{}{
			"datacenter_id": region_id_int,
			"ip_version":    int(version),
			"bandwidth":     hosting.DefaultBandwidth,
		}}, &response)
	if err != nil {
		return IPAddress{}, err
	}
	if err = h.waitForOp(response); err != nil {
		return IPAddress{}, err
	}

	if err = h.Send("hosting.ip.info", []interface{}{response.IPID}, &iip); err != nil {
		return IPAddress{}, err
	}

	return toIPAddress(iip), nil
}

func (h Hostingv4) DescribeIP(ipfilter IPFilter) ([]IPAddress, error) {
	ipmap, err := ipFilterToMap(ipfilter)
	if err != nil {
		return nil, err
	}

	var response = []iPAddressv4{}
	if err = h.Send("hosting.ip.list", []interface{}{ipmap}, &response); err != nil {
		return nil, err
	}

	var ips []IPAddress
	for _, iip := range response {
		ips = append(ips, toIPAddress(iip))
	}

	return ips, nil
}

func (h Hostingv4) DeleteIP(ip IPAddress) error {
	ipid, err := strconv.Atoi(ip.ID)
	if err != nil {
		return internalParseError("IPAddress", "ID")
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

func (h Hostingv4) ipFromID(ipid int) (IPAddress, error) {
	response := iPAddressv4{}
	err := h.Send("hosting.ip.info", []interface{}{ipid}, &response)
	if err != nil {
		return IPAddress{}, err
	}
	return toIPAddress(response), nil
}

// Internal methods to convert Hosting structures to v4 structures

func ipFilterToMap(ipfilter IPFilter) (map[string]interface{}, error) {
	ipmap := make(map[string]interface{})
	var err error

	if ipfilter.Version != 0 {
		if ipfilter.Version != hosting.IPv4 && ipfilter.Version != hosting.IPv6 {
			return nil, internalParseError("IPFilter", "Version")
		}
		ipmap["version"] = int(ipfilter.Version)
	}

	if ipfilter.ID != "" {
		ipmap["id"], err = strconv.Atoi(ipfilter.ID)
		if err != nil {
			return nil, internalParseError("IPFilter", "ID")
		}
	}
	if ipfilter.RegionID != "" {
		ipmap["datacenter_id"], err = strconv.Atoi(ipfilter.RegionID)
		if err != nil {
			return nil, internalParseError("IPFilter", "ID")
		}
	}

	if ipfilter.IP != "" {
		ipmap["ip"] = ipfilter.IP
	}

	return ipmap, nil
}

func toIPAddress(iip iPAddressv4) (ip IPAddress) {
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
