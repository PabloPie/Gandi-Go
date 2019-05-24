package hostingv4

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/PabloPie/go-gandi/hosting"
)

type (
	Region = hosting.Region
)

type diskImagev4 struct {
	ID       int    `xmlrpc:"id"`
	DiskID   int    `xmlrpc:"disk_id"`
	RegionID int    `xmlrpc:"datacenter_id"`
	Name     string `xmlrpc:"label"`
	Size     int    `xmlrpc:"size"`
}

// ImageByName return the DiskImage with label `name`
func (h Hostingv4) ImageByName(name string, region Region) (DiskImage, error) {
	if region.ID == "" {
		return DiskImage{}, errors.New("Region provided does not have an ID")
	}

	regionid, err := strconv.Atoi(region.ID)
	if err != nil {
		return DiskImage{},
			fmt.Errorf("Error parsing RegionID '%s' from Region %v", region.ID, region)
	}
	filter := map[string]interface{}{"label": name, "datacenter_id": regionid}

	var response = []diskImagev4{}
	request := []interface{}{filter}
	err = h.Send("hosting.image.list", request, &response)
	if err != nil {
		return DiskImage{}, err
	}

	if len(response) < 1 {
		return DiskImage{}, errors.New("Image not found")
	}

	return fromDiskImagev4(response[0]), nil
}

// ListImagesInRegion returns the list of Images available in `region`
func (h Hostingv4) ListImagesInRegion(region Region) ([]DiskImage, error) {
	if region.ID == "" {
		return []DiskImage{}, errors.New("Region provided does not have an ID")
	}

	regionid, err := strconv.Atoi(region.ID)
	if err != nil {
		return nil,
			fmt.Errorf("Error parsing RegionID '%s' from Region %v", region.ID, region)
	}
	filter := map[string]interface{}{"datacenter_id": regionid}

	response := []diskImagev4{}
	request := []interface{}{filter}
	err = h.Send("hosting.image.list", request, &response)
	if err != nil {
		return []DiskImage{}, err
	}

	if len(response) < 1 {
		return []DiskImage{}, errors.New("No images")
	}
	var diskimages []DiskImage
	for _, image := range response {
		diskimages = append(diskimages, fromDiskImagev4(image))
	}

	return diskimages, nil
}

func fromDiskImagev4(image diskImagev4) DiskImage {
	id := strconv.Itoa(image.ID)
	diskid := strconv.Itoa(image.DiskID)
	regionid := strconv.Itoa(image.RegionID)
	return DiskImage{
		ID:       id,
		DiskID:   diskid,
		RegionID: regionid,
		Name:     image.Name,
		Size:     image.Size / 1024,
	}
}
