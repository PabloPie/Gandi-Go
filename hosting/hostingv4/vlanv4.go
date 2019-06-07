package hostingv4

import (
	"fmt"
	"log"
	"strconv"

	"github.com/PabloPie/go-gandi/hosting"
)

type vlanv4 struct {
	ID       int    `xmlrpc:"id"`
	Name     string `xmlrpc:"name"`
	Gateway  string `xmlrpc:"gateway"`
	Subnet   string `xmlrpc:"subnet"`
	RegionID int    `xmlrpc:"datacenter_id"`
}

type vlanSpecv4 struct {
	Name     string `xmlrpc:"name"`
	Gateway  string `xmlrpc:"gateway"`
	Subnet   string `xmlrpc:"subnet"`
	RegionID int    `xmlrpc:"datacenter_id"`
}

type vlanFilterv4 struct {
	ID       []int  `xmlrpc:"id"`
	RegionID []int  `xmlrpc:"datacenter_id"`
	Name     string `xmlrpc:"name"`
}

// CreateVlan creates a new vlan from the spec given
func (h Hostingv4) CreateVlan(newVlan hosting.VlanSpec) (hosting.Vlan, error) {
	var fn = "CreateVlan"
	if newVlan.RegionID == "" {
		return hosting.Vlan{}, &HostingError{fn, "VlanSpec", "RegionID", ErrNotProvided}
	}
	if newVlan.Name == "" {
		return hosting.Vlan{}, &HostingError{fn, "VlanSpec", "Name", ErrNotProvided}
	}

	vlanv4, err := toVlanSpecv4(newVlan)
	if err != nil {
		return hosting.Vlan{}, err
	}
	vlan, _ := structToMap(vlanv4)

	response := Operation{}
	params := []interface{}{vlan}
	log.Printf("[INFO] Creating Vlan %s...", newVlan.Name)
	err = h.Send("hosting.vlan.create", params, &response)
	if err != nil {
		return hosting.Vlan{}, err
	}
	if err := h.waitForOp(response); err != nil {
		return hosting.Vlan{}, err
	}
	log.Printf("[INFO] Vlan %s(ID: %d) created!", newVlan.Name, response.DiskID)

	// operations don't contain a vlan's id
	// we need to use its name to get the Vlan
	return h.VlanFromName(newVlan.Name)
}

// VlanFromName is a helper function to get a Vlan given its name
func (h Hostingv4) VlanFromName(name string) (hosting.Vlan, error) {
	if name == "" {
		return hosting.Vlan{}, &HostingError{"VlanFromName", "-", "name", ErrNotProvided}
	}
	vlans, err := h.ListVlans(hosting.VlanFilter{Name: name})
	if err != nil {
		return hosting.Vlan{}, err
	}
	if len(vlans) < 1 {
		return hosting.Vlan{}, fmt.Errorf("Vlan '%s' does not exist", name)
	}

	return vlans[0], nil
}

// ListVlans returns a list of vlans filtered with the options provided in `vlanFilter`
func (h Hostingv4) ListVlans(vlanfilter hosting.VlanFilter) ([]hosting.Vlan, error) {
	filterv4, err := toVlanFilterv4(vlanfilter)
	if err != nil {
		return nil, err
	}
	filter, _ := structToMap(filterv4)
	response := []vlanv4{}
	params := []interface{}{}
	if len(filter) > 0 {
		params = append(params, filter)
	}
	err = h.Send("hosting.vlan.list", params, &response)
	if err != nil {
		return nil, err
	}

	var vlans []hosting.Vlan
	for _, vlan := range response {
		vlans = append(vlans, fromVlanv4(vlan))
	}

	return vlans, nil
}

// UpdateVlanGW updates the gateway of the vlan
func (h Hostingv4) UpdateVlanGW(vlan hosting.Vlan, newGW string) (hosting.Vlan, error) {
	var fn = "UpdateVlanGW"
	if vlan.ID == "" {
		return hosting.Vlan{}, &HostingError{fn, "Vlan", "ID", ErrNotProvided}
	}
	vlanupdate := map[string]string{"gateway": newGW}
	vlanid, err := strconv.Atoi(vlan.ID)
	if err != nil {
		return hosting.Vlan{}, &HostingError{fn, "Vlan", "ID", ErrParse}
	}

	response := Operation{}
	request := []interface{}{vlanid, vlanupdate}
	err = h.Send("hosting.vlan.update", request, &response)
	if err != nil {
		return hosting.Vlan{}, err
	}
	err = h.waitForOp(response)
	if err != nil {
		return hosting.Vlan{}, err
	}

	// Use ListVlans to get the Vlan object instead of implementing vlanFromID
	// specially given that vlan.list does return more info than vlan.info
	ids := []string{vlan.ID}
	vlans, err := h.ListVlans(hosting.VlanFilter{ID: ids})
	if err != nil {
		return hosting.Vlan{}, err
	}
	if len(vlans) < 1 {
		return hosting.Vlan{}, fmt.Errorf("Vlan '%s' does not exist", vlan.ID)
	}
	return vlans[0], nil
}

// RenameVlan renames a private network
func (h Hostingv4) RenameVlan(vlan hosting.Vlan, newName string) (hosting.Vlan, error) {
	var fn = "RenameVlan"
	if vlan.ID == "" {
		return hosting.Vlan{}, &HostingError{fn, "Vlan", "ID", ErrNotProvided}
	}
	vlanupdate := map[string]string{"name": newName}
	vlanid, err := strconv.Atoi(vlan.ID)
	if err != nil {
		return hosting.Vlan{}, &HostingError{fn, "Vlan", "ID", ErrParse}
	}

	response := Operation{}
	request := []interface{}{vlanid, vlanupdate}
	err = h.Send("hosting.vlan.update", request, &response)
	if err != nil {
		return hosting.Vlan{}, err
	}
	err = h.waitForOp(response)
	if err != nil {
		return hosting.Vlan{}, err
	}

	return h.VlanFromName(newName)
}

// DeleteVlan deletes a Vlan
//
// A Vlan won't be deleted if there is any existing private ip
// linked to it
func (h Hostingv4) DeleteVlan(vlan hosting.Vlan) error {
	var fn = "DeleteVlan"
	if vlan.ID == "" {
		return &HostingError{fn, "Vlan", "ID", ErrNotProvided}
	}

	vlanid, err := strconv.Atoi(vlan.ID)
	if err != nil {
		return &HostingError{fn, "Vlan", "ID", ErrParse}
	}

	response := Operation{}
	params := []interface{}{vlanid}
	err = h.Send("hosting.vlan.delete", params, &response)
	if err != nil {
		return err
	}
	err = h.waitForOp(response)
	return err
}

// Conversion functions

// Hosting VlanSpec -> v4 VlanSpec
func toVlanSpecv4(vlan hosting.VlanSpec) (vlanSpecv4, error) {
	region, err := strconv.Atoi(vlan.RegionID)
	if err != nil {
		return vlanSpecv4{}, internalParseError("VlanSpec", "RegionID")
	}
	return vlanSpecv4{
		RegionID: region,
		Name:     vlan.Name,
		Gateway:  vlan.Gateway,
		Subnet:   vlan.Subnet,
	}, nil
}

// Hosting VlanFilter -> v4 VlanFilter
func toVlanFilterv4(vlan hosting.VlanFilter) (vlanFilterv4, error) {
	var regions []int
	for _, r := range vlan.RegionID {
		region := toInt(r)
		if region == -1 {
			return vlanFilterv4{}, internalParseError("VlanFilter", "RegionID")
		}
		regions = append(regions, region)
	}

	var ids []int
	for _, vlanid := range vlan.ID {
		id := toInt(vlanid)
		if id == -1 {
			return vlanFilterv4{}, internalParseError("VlanFilter", "ID")
		}
		ids = append(ids, id)
	}

	return vlanFilterv4{
		RegionID: regions,
		ID:       ids,
		Name:     vlan.Name,
	}, nil
}

// v4 Vlan -> Hosting Vlan
func fromVlanv4(vlan vlanv4) hosting.Vlan {
	id := strconv.Itoa(vlan.ID)
	region := strconv.Itoa(vlan.RegionID)
	return hosting.Vlan{
		ID:       id,
		RegionID: region,
		Gateway:  vlan.Gateway,
		Subnet:   vlan.Subnet,
		Name:     vlan.Name,
	}
}
