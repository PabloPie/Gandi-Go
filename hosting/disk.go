package hosting

import (
	"errors"
)

const (
	defaultSize   int = 10240
	defaultRegion int = 3
)

// DiskManager represents a service capable of manipulating Gandi Disks
type DiskManager interface {
	CreateDisk(disk DiskSpec, src *DiskImage) (Disk, error)
	DescribeDisks(diskFilter DiskFilter) ([]Disk, error)
	DeleteDisk(disk *Disk) error
	ExtendDisk(disk *Disk, size uint) error
	RenameDisk(disk *Disk, name string) error
}

// ImageManager represents a service capable of getting information about
// Gandi Disk images
type ImageManager interface {
	ImageByName(name string, region Region) (DiskImage, error)
	ImageByNameVersion(os string, version string, region Region) (DiskImage, error)
}

// Disk is a Gandi disk object
type Disk struct {
	ID       int    `xmlrpc:"id"`
	Name     string `xmlrpc:"name"`
	Size     int    `xmlrpc:"size"`
	RegionID int    `xmlrpc:"datacenter_id"`
	State    string `xmlrpc:"state"`
	Type     string `xmlrpc:"type"`
	VM       []int  `xmlrpc:"vms_id"`
	BootDisk bool   `xmlrpc:"is_boot_disk"`
}

// DiskSpec contains the parameters to create a new Disk
type DiskSpec struct {
	RegionID int    `xmlrpc:"datacenter_id"`
	Name     string `xmlrpc:"name"`
	Size     int    `xmlrpc:"size"`
}

//DiskFilter is used to filter the results DescribeDisks returns,
type DiskFilter struct {
	ID       int    `xmlrpc:"id"`
	RegionID int    `xmlrpc:"datacenter_id"`
	Name     string `xmlrpc:"name"`
	VMID     int    `xmlrpc:"vm_id"`
}

// DiskImage represents an image defined by Gandi
// with an OS, used to create new Disks and VMs
type DiskImage struct {
	ID       int    `xmlrpc:"id"`
	DiskID   int    `xmlrpc:"disk_id"`
	RegionID int    `xmlrpc:"datacenter_id"`
	Name     string `xmlrpc:"label"`
	Os       string `xmlrpc:"system"`
	Version  string `xmlrpc:"version"`
	Size     int    `xmlrpc:"size"`
}

// CreateDisk creates a disk either from an existing one if `src` is not nil,
// or an empty one if it is. `Name` is mandatory, `Size` defaults
// to 10GB and `Region` to FR_SD3
// Can't copy from any disk atm
func (h Hostingv4) CreateDisk(disk DiskSpec, src *DiskImage) (Disk, error) {
	var err error
	if disk.RegionID == 0 {
		disk.RegionID = defaultRegion
	}
	// Name is mandatory so we don't have to query the api again
	// to get the name generated
	if disk.Name == "" {
		return Disk{}, errors.New("Disk name required")
	}
	if disk.Size == 0 {
		disk.Size = defaultSize
	}

	var response = Operation{}
	request := []interface{}{disk}
	if src != nil {
		request = append(request, src.DiskID)
		err = h.Send("hosting.disk.create_from", request, &response)
	} else {
		err = h.Send("hosting.disk.create", request, &response)
	}
	if err != nil {
		return Disk{}, err
	}

	d := Disk{
		ID:       response.DiskID,
		Name:     disk.Name,
		Size:     disk.Size,
		RegionID: disk.RegionID,
		State:    "being_created",
		Type:     "data",
		BootDisk: false,
	}
	return d, nil
}

// DescribeDisks return a list of disks filtered with the options provided in `diskFilter`
func (h Hostingv4) DescribeDisks(diskfilter DiskFilter) ([]Disk, error) {
	var err error
	filter := diskFilterToMap(diskfilter)
	response := []Disk{}
	request := []interface{}{}
	if len(filter) > 0 {
		request = append(request, filter)
	}
	err = h.Send("hosting.disk.list", request, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteDisk deletes the Disk with ID `disk.ID`
func (h Hostingv4) DeleteDisk(disk *Disk) error {
	if disk.ID == 0 {
		return errors.New("Disk ID not provided, cannot delete")
	}

	response := Operation{}
	request := []interface{}{disk.ID}
	err := h.Send("hosting.disk.delete", request, &response)
	if err != nil {
		return err
	}

	disk.State = "deleted"
	return nil
}

// ExtendDisk extends `disk` size by `size` (original size + `size`)
// Disks cannot shrink in size
// `size` is in GB
func (h Hostingv4) ExtendDisk(disk *Disk, size uint) error {
	if disk.ID == 0 {
		return errors.New("Disk ID not provided")
	}
	// size has to be a multiple of 1024
	newSize := disk.Size + (int(size) * 1024)
	diskupdate := map[string]int{"size": newSize}

	response := Operation{}
	request := []interface{}{disk.ID, diskupdate}
	err := h.Send("hosting.disk.update", request, &response)
	if err != nil {
		return err
	}
	disk.Size = newSize
	return nil
}

// RenameDisk renames `disk` to `newName`
func (h Hostingv4) RenameDisk(disk *Disk, newName string) error {
	if disk.ID == 0 {
		return errors.New("Disk ID not provided")
	}
	diskupdate := map[string]string{"name": newName}

	response := Operation{}
	request := []interface{}{disk.ID, diskupdate}
	err := h.Send("hosting.disk.update", request, &response)
	if err != nil {
		return err
	}

	disk.Name = newName
	return nil
}

// ImageByName return the DiskImage exactly matching the name provided
func (h Hostingv4) ImageByName(name string, region Region) (DiskImage, error) {
	filter := map[string]interface{}{"label": name, "datacenter_id": region.ID}
	var res = []DiskImage{}
	request := []interface{}{filter}
	err := h.Send("hosting.image.list", request, &res)
	if err != nil {
		return DiskImage{}, err
	}

	if len(res) < 1 {
		return DiskImage{}, errors.New("Image not found")
	}
	return res[0], nil
}

// ImageByNameVersion returns the first DiskImage that matches both `os` and `version`
func (h Hostingv4) ImageByNameVersion(os string, version string, region Region) (DiskImage, error) {
	filter := map[string]interface{}{"system": os, "datacenter_id": region.ID}
	res := []DiskImage{}
	request := []interface{}{filter}
	err := h.Send("hosting.image.list", request, &res)
	if err != nil {
		return DiskImage{}, err
	}

	if len(res) < 1 {
		return DiskImage{}, errors.New("Image not found")
	}

	res = filterByVersion(res, version)
	if len(res) < 1 {
		return DiskImage{}, errors.New("Version not found for image")
	}
	return res[0], nil
}

// Filter a list of DiskImages by version
func filterByVersion(images []DiskImage, version string) []DiskImage {
	filtered := images[:0]
	for _, x := range images {
		if x.Version == version {
			filtered = append(filtered, x)
		}
	}
	return filtered
}

// Dirty, we should be able to do a generic function for structs with reflect
func diskFilterToMap(diskfilter DiskFilter) map[string]interface{} {
	res := map[string]interface{}{}

	if diskfilter.Name != "" {
		res["name"] = diskfilter.Name
	}
	if diskfilter.ID != 0 {
		res["id"] = diskfilter.ID
	}
	if diskfilter.RegionID != 0 {
		res["datacenter_id"] = diskfilter.RegionID
	}
	if diskfilter.VMID != 0 {
		res["vm_id"] = diskfilter.VMID
	}
	return res

}
