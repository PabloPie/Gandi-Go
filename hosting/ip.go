package hosting

import (
	"errors"
)

// for a more readable code
type IPVersion int

const (
	IPv4 IPVersion = 4
	IPv6 IPVersion = 6
	defaultBandwidth float32 = 102400.0
)

type IPManager interface {
	CreateIP(region Region, version IPVersion) (IPAddress, error)
	DescribeIP(ipfilter IPFilter) ([]IPAddress, error)
	DeleteIP(ipid int) error
}

type IPAddress struct {
	ID       int    `xmlrpc:"id"`
	IP       string `xmlrpc:"ip"`
	RegionID int    `xmlrpc:"datacenter_id"`
	Version  int `xmlrpc:"version"`
	VM       int    `xmlrpc:"vm_id"`
	State    string `xmlrpc:"state"`
	ifaceID  int    `xmlrpc:"iface_id"`
}

type IPFilter struct {
	ID       int       `xmlrpc:"id"`
	RegionID int       `xmlrpc:"datacenter_id"`
	Version  IPVersion `xmlrpc:"version"`
	IP       string    `xmlrpc:"ip"`
}

type Iface struct {
	ID int
	IPS []IPAddress
}

func (h Hostingv4) CreateIP(region Region, version IPVersion) (IPAddress, error) {
	var ip IPAddress
	
	if version != IPv4 && version != IPv6 {
		return ip, errors.New("Bad IP version")
	}

	var err error

	var response = Operation{}
	err = h.Send("hosting.iface.create", []interface{} {
					map[string]interface{} {
						"datacenter_id": region.ID,
						"ip_version": int(version),
						"bandwidth" : defaultBandwidth,
						}}, &response)
	if err != nil{
		return ip, err
	}
	if err = h.waitForOp(response) ; err != nil {
		return ip, err
	}

	err = h.Send("hosting.ip.info", []interface{}{response.IPID}, &ip)

	return ip, err
}

func (h Hostingv4) DescribeIP(ipfilter IPFilter) ([]IPAddress, error) {
	ipmap, err := ipFilterToMap(&ipfilter)
	if err != nil {
		return nil, err
	}

	var response = []IPAddress{}
	err = h.Send("hosting.ip.list", []interface{}{ipmap}, &response)

	return response, err
}

func ipFilterToMap(ipfilter* IPFilter) (map[string]interface{}, error) {
	ipmap := map[string]interface{}{}
	if(ipfilter.Version != 0) {
		if(ipfilter.Version != IPv4 && ipfilter.Version != IPv6) {
			return nil, errors.New("Bad IP version")
		}
		ipmap["version"] = int(ipfilter.Version)
	}
	if(ipfilter.ID != 0) {
		ipmap["id"] = ipfilter.ID
	}
	if(ipfilter.IP != "") {
		ipmap["ip"] = ipfilter.IP
	}
	if(ipfilter.RegionID != 0) {
		ipmap["datacenter_id"] = ipfilter.RegionID
	}
	return ipmap, nil
}

func (h Hostingv4) DeleteIP(ipid int) error {
	var response = Operation{}
	err := h.Send("hosting.ip.info", []interface{}{ipid}, &response)
	if err != nil {
		return err
	}

	err = h.Send("hosting.iface.delete", []interface{}{response.IfaceID}, &response)
	
	return h.waitForOp(response)
}
