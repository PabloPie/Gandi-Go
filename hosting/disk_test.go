package hosting_test

import (
	"reflect"
	"testing"

	"github.com/PabloPie/Gandi-Go/hosting"
	"github.com/PabloPie/Gandi-Go/mock"
)

var (
	client = mock.NewMockClientv4()
	h      = hosting.Newv4Hosting(client)
	// default values
	defaultRegion = 3
	defaultSize   = 10240
	// expected params
	diskid   = 1
	diskname = "Disk1"
	disksize = 20480
	region   = 4
)

func TestCreateDiskWithSizeAndRegion(t *testing.T) {
	expected := hosting.Disk{
		ID:       diskid,
		Name:     diskname,
		Size:     disksize,
		RegionID: region,
		State:    "being_created",
		Type:     "data",
		BootDisk: false,
	}
	diskspec := hosting.DiskSpec{
		RegionID: region,
		Name:     diskname,
		Size:     disksize,
	}
	disk, err := h.CreateDisk(diskspec, nil)
	if err != nil {
		t.Errorf("Error creating disk: %v", err)
	}

	if !reflect.DeepEqual(disk, expected) {
		t.Errorf("Error, expected %+v, got instead %+v", expected, disk)
	}
}

func TestCreateDiskWithoutSizeAndRegion(t *testing.T) {
	expected := hosting.Disk{
		ID:       diskid,
		Name:     diskname,
		Size:     defaultSize,
		RegionID: defaultRegion,
		State:    "being_created",
		Type:     "data",
		BootDisk: false,
	}
	diskspec := hosting.DiskSpec{
		Name: diskname,
	}
	disk, err := h.CreateDisk(diskspec, nil)
	if err != nil {
		t.Errorf("Error creating disk: %v", err)
	}

	if !reflect.DeepEqual(disk, expected) {
		t.Errorf("Error, expected %+v, got instead %+v", expected, disk)
	}
}

func TestCreateDiskWithoutName(t *testing.T) {
	diskspec := hosting.DiskSpec{}
	_, err := h.CreateDisk(diskspec, nil)
	if err == nil {
		t.Errorf("No name provided")
	}
}

func TestCreateDiskFromSource(t *testing.T) {
	expected := hosting.Disk{
		ID:       diskid,
		Name:     diskname,
		Size:     defaultSize,
		RegionID: defaultRegion,
		State:    "being_created",
		Type:     "data",
		BootDisk: false,
	}
	diskspec := hosting.DiskSpec{
		Name: diskname,
	}
	// TODO XXX, substitute DiskImage for ImageByNameVersion
	diskImage := hosting.DiskImage{DiskID: 5}
	disk, err := h.CreateDisk(diskspec, &diskImage)
	if err != nil {
		t.Errorf("Error creating disk: %v", err)
	}

	if !reflect.DeepEqual(disk, expected) {
		t.Errorf("Error, expected %+v, got instead %+v", expected, disk)
	}
}

func TestListDisksWithEmptyFilter(t *testing.T) {
	disks, err := h.DescribeDisks(hosting.DiskFilter{})
	if err != nil {
		t.Errorf("Error listing disks: %v", err)
	}
	if len(disks) < 1 {
		t.Errorf("Error, expected to get at least 1 Disk")
	}
}

func TestDeleteDisk(t *testing.T) {
	disk := hosting.Disk{ID: 1}
	err := h.DeleteDisk(&disk)
	if err != nil {
		t.Errorf("Error deleting disk: %v", err)
	}
	if disk.State != "deleted" {
		t.Errorf("Error, disk state should be 'deleted' but is %s", disk.State)
	}
}

func TestExtendDisk2GB(t *testing.T) {
	disk := hosting.Disk{ID: 1, Size: defaultSize}
	err := h.ExtendDisk(&disk, 2)
	if err != nil {
		t.Errorf("Error extending disk: %v", err)
	}
	expectedSize := defaultSize + 2*1024
	if disk.Size != expectedSize {
		t.Errorf("Error, expected size of %d, got %d instead", expectedSize, disk.Size)
	}
}

func TestRenameDisk(t *testing.T) {
	expectedName := "NewName"
	disk := hosting.Disk{ID: 1, Name: diskname}
	err := h.RenameDisk(&disk, expectedName)
	if err != nil {
		t.Errorf("Error renaming disk: %v", err)
	}
	if disk.Name != expectedName {
		t.Errorf("Error, expected disk name to be %s, got %s instead", expectedName, disk.Name)
	}
}

func TestGetDebian9Image(t *testing.T) {
	expected := "Debian 9"
	region := hosting.Region{ID: region}
	diskimage, err := h.ImageByName(expected, region)
	if err != nil {
		t.Errorf("Error getting image: %v", err)
	}
	if diskimage.Name != expected {
		t.Errorf("Error, expected image %s, got %s instead", expected, diskimage.Name)
	}
}

func TestGetDebian9ImageNameVersion(t *testing.T) {
	expectedOS := "Debian"
	expectedVersion := "9"
	region := hosting.Region{ID: region}
	diskimage, err := h.ImageByNameVersion(expectedOS, expectedVersion, region)
	if err != nil {
		t.Errorf("Error getting image: %v", err)
	}
	if diskimage.Os != expectedOS || diskimage.Version != expectedVersion {
		t.Errorf("Error, expected os %s version %s, got %s %s instead",
			expectedOS, expectedVersion, diskimage.Os, diskimage.Version)
	}
}