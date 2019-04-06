package hosting

import (
	"errors"
)

// RegionManager represents a service capable of getting info
// about Gandi Datacenters
type RegionManager interface {
	ListRegions() []Region
	RegionbyCode(code string) Region
}

// Region represents a Gandi datacenter
type Region struct {
	ID      int    `xmlrpc:"id"`
	Name    string `xmlrpc:"dc_code"`
	Country string `xmlrpc:"country"`
}

// ListRegions lists every Gandi datacenter
func (h Hostingv4) ListRegions() ([]Region, error) {
	var res = []Region{}
	request := []interface{}{}
	err := h.Send("hosting.datacenter.list", request, &res)
	if err != nil {
		return []Region{}, err
	}
	return res, nil
}

// RegionbyCode returns the region with code `code` if it exists
func (h Hostingv4) RegionbyCode(code string) (Region, error) {
	var res = []Region{}
	var filter = map[string]string{"dc_code": code}
	request := []interface{}{filter}
	err := h.Send("hosting.datacenter.list", request, &res)
	if err != nil {
		return Region{}, err
	}
	if len(res) != 1 {
		return Region{}, errors.New("Region not found")
	}
	return res[0], nil
}
