package hostingv4

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/PabloPie/go-gandi/hosting"
)

// internal representation of a hosting.Disk Image for API v4
type diskImagev4 struct {
	ID       int    `xmlrpc:"id"`
	DiskID   int    `xmlrpc:"disk_id"`
	RegionID int    `xmlrpc:"datacenter_id"`
	Name     string `xmlrpc:"label"`
	Size     int    `xmlrpc:"size"`
}

// ImageByName returns the hosting.DiskImage with label `name` found in `region`
func (h Hostingv4) ImageByName(name string, region hosting.Region) (hosting.DiskImage, error) {
	if region.ID == "" {
		return hosting.DiskImage{}, errors.New("hosting.Region provided does not have an ID")
	}

	regionid, err := strconv.Atoi(region.ID)
	if err != nil {
		return hosting.DiskImage{},
			fmt.Errorf("Error parsing RegionID '%s' from hosting.Region %v", region.ID, region)
	}
	filter := map[string]interface{}{"label": name, "datacenter_id": regionid}

	var response = []diskImagev4{}
	request := []interface{}{filter}
	err = h.Send("hosting.image.list", request, &response)
	if err != nil {
		return hosting.DiskImage{}, err
	}

	if len(response) < 1 {
		return hosting.DiskImage{}, errors.New("Image not found")
	}

	return fromDiskImagev4(response[0]), nil
}

// ListImagesInRegion returns the list of Images available in `region`
func (h Hostingv4) ListImagesInRegion(region hosting.Region) ([]hosting.DiskImage, error) {
	if region.ID == "" {
		return []hosting.DiskImage{}, errors.New("hosting.Region provided does not have an ID")
	}

	regionid, err := strconv.Atoi(region.ID)
	if err != nil {
		return nil,
			fmt.Errorf("Error parsing RegionID '%s' from hosting.Region %v", region.ID, region)
	}
	filter := map[string]interface{}{"datacenter_id": regionid}

	response := []diskImagev4{}
	request := []interface{}{filter}
	err = h.Send("hosting.image.list", request, &response)
	if err != nil {
		return []hosting.DiskImage{}, err
	}

	if len(response) < 1 {
		return []hosting.DiskImage{}, errors.New("No images")
	}
	var diskimages []hosting.DiskImage
	for _, image := range response {
		diskimages = append(diskimages, fromDiskImagev4(image))
	}

	return diskimages, nil
}

// diskImagev4 -> Hosting hosting.DiskImage
func fromDiskImagev4(image diskImagev4) hosting.DiskImage {
	id := strconv.Itoa(image.ID)
	diskid := strconv.Itoa(image.DiskID)
	regionid := strconv.Itoa(image.RegionID)
	return hosting.DiskImage{
		ID:       id,
		DiskID:   diskid,
		RegionID: regionid,
		Name:     image.Name,
		Size:     image.Size / 1024,
	}
}
