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
	CreateIP(regionid int, version IPVersion) ([]IPAddress, error)
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
	ID       int
	RegionID int
	Version  IPVersion
	IP       string
}

func (h Hostingv4) CreateIP(regionid int, version IPVersion) ([]IPAddress, error) {
	if version != IPv4 && version != IPv6 {
		return nil, errors.New("Bad IP version")
	}
	
	var err error

	var response = Operation{}
	parameters := map[string]interface{}{"datacenter_id": regionid, "ip_version": int(version), "bandwidth" : defaultBandwidth}
	err = h.Send("hosting.iface.create", []interface{}{parameters}, &response)
	if err != nil {
		return nil, err
	}
	
	var ips = []IPAddress{}
	parameters = map[string]interface{}{"iface_id": response.IfaceID}
	err = h.Send("hosting.ip.list", []interface{}{parameters}, &ips)
	if err != nil {
		return nil, err
	}
	
	// An IPv4 iface creates an IPv4 and an IPv6
	// but IPv6 creation takes time...
	if version == IPv4 && len(ips) != 2 {
		ips = append(ips, IPAddress{
				Version : 6,
				State : "being_created",
				RegionID : regionid,
				ifaceID : response.IfaceID,
			})
	}

	return ips, err
}

func (h Hostingv4) DescribeIP(ipfilter IPFilter) ([]IPAddress, error) {
	ipmap, err := ipFilterToMap(&ipfilter)
	if err != nil {
		return nil, err
	}

	var response = []IPAddress{}
	err = h.Send("hosting.ip.list", []interface{}{ipmap}, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
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

func (h Hostingv4) DeleteIP(ipid int) (err error) {
	var response = Operation{}
	err = h.Send("hosting.ip.info", []interface{}{ipid}, &response)
	if err != nil {
		return
	}

	err = h.Send("hosting.iface.delete", []interface{}{response.IfaceID}, &response)
	return
}
