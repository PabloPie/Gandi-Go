package hostingv4

import (
	"errors"
	"strconv"

	"github.com/PabloPie/go-gandi/hosting"
)

type regionv4 struct {
	ID      int    `xmlrpc:"id"`
	Name    string `xmlrpc:"dc_code"`
	Country string `xmlrpc:"country"`
}

// ListRegions lists every Gandi datacenter
func (h Hostingv4) ListRegions() ([]hosting.Region, error) {
	response := []regionv4{}
	request := []interface{}{}
	err := h.Send("hosting.datacenter.list", request, &response)
	if err != nil {
		return []hosting.Region{}, err
	}

	var regions = []hosting.Region{}
	for _, region := range response {
		regions = append(regions, fromRegionv4(region))
	}
	return regions, nil
}

// RegionbyCode returns the region with code `code` if it exists
func (h Hostingv4) RegionbyCode(code string) (hosting.Region, error) {
	response := []regionv4{}
	filter := map[string]string{"dc_code": code}
	request := []interface{}{filter}
	err := h.Send("hosting.datacenter.list", request, &response)
	if err != nil {
		return hosting.Region{}, err
	}
	if len(response) < 1 {
		return hosting.Region{}, errors.New("hosting.Region not found")
	}

	return fromRegionv4(response[0]), nil
}

// Conversion functions

// regionv4 -> Hosting hosting.Region
func fromRegionv4(region regionv4) hosting.Region {
	id := strconv.Itoa(region.ID)
	return hosting.Region{
		ID:      id,
		Name:    region.Name,
		Country: region.Country,
	}
}
