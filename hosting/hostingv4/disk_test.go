package hostingv4

import (
	"reflect"
	"testing"

	"github.com/PabloPie/Gandi-Go/mock"
	"github.com/golang/mock/gomock"
)

var (
	// default values
	defaultRegion = 3
	defaultSize   = 10240
	// expected params
	diskid   = 1
	diskname = "Disk1"
	disksize = 20480
	region   = 4
)

func TestCreateDiskWithNameSizeAndRegion(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	paramsDiskCreate := []interface{}{map[string]interface{}{
		"datacenter_id": region,
		"name":          diskname,
		"size":          disksize,
	}}
	responseDiskCreate := Operation{
		ID:     1,
		DiskID: 1,
	}
	creation := mockClient.EXPECT().Send("hosting.disk.create",
		paramsDiskCreate, gomock.Any()).SetArg(2, responseDiskCreate).Return(nil)

	paramsOperationInfo := []interface{}{responseDiskCreate.ID}
	responseWait := operationInfo{responseDiskCreate.ID, "DONE"}
	wait := mockClient.EXPECT().Send("operation.info",
		paramsOperationInfo, gomock.Any()).SetArg(2, responseWait).Return(nil).After(creation)

	paramsDiskInfo := []interface{}{responseDiskCreate.DiskID}
	responseDiskInfo := diskv4{1, diskname, disksize, region, "created", "data", []int{}, false}
	mockClient.EXPECT().Send("hosting.disk.info",
		paramsDiskInfo, gomock.Any()).SetArg(2, responseDiskInfo).Return(nil).After(wait)

	diskspec := DiskSpec{
		RegionID: "4",
		Name:     diskname,
		Size:     uint(disksize),
	}
	disk, _ := testHosting.CreateDisk(diskspec)

	expected := Disk{
		ID:       "1",
		Name:     diskname,
		Size:     uint(disksize),
		RegionID: "4",
		State:    "created",
		Type:     "data",
		BootDisk: false,
	}

	if !reflect.DeepEqual(disk, expected) {
		t.Errorf("Error, expected %+v, got instead %+v", expected, disk)
	}
}

//
// func TestCreateDiskFromSource(t *testing.T) {
// expected := hosting.Disk{
// ID:       diskid,
// Name:     diskname,
// Size:     defaultSize,
// RegionID: defaultRegion,
// State:    "being_created",
// Type:     "data",
// BootDisk: false,
// }
// diskspec := hosting.DiskSpec{
// Name: diskname,
// }
//TODO XXX, substitute DiskImage for ImageByNameVersion
// diskImage := hosting.DiskImage{DiskID: 5}
// disk, err := h.CreateDisk(diskspec, &diskImage)
// if err != nil {
// t.Errorf("Error creating disk: %v", err)
// }
//
// if !reflect.DeepEqual(disk, expected) {
// t.Errorf("Error, expected %+v, got instead %+v", expected, disk)
// }
// }
//
// func TestListDisksWithEmptyFilter(t *testing.T) {
// disks, err := h.DescribeDisks(hosting.DiskFilter{})
// if err != nil {
// t.Errorf("Error listing disks: %v", err)
// }
// if len(disks) < 1 {
// t.Errorf("Error, expected to get at least 1 Disk")
// }
// }
//
// func TestListDisksWithNameInFilter(t *testing.T) {
// expectedname := "disk3"
// disks, err := h.DescribeDisks(hosting.DiskFilter{Name: expectedname})
// if err != nil {
// t.Errorf("Error listing disks: %v", err)
// }
// if len(disks) != 1 {
// t.Errorf("Error, expected to get 1 Disk and got %d instead", len(disks))
// }
// if disks[0].Name != expectedname {
// t.Errorf("Error, expected to get Disk with name %s, got %s instead",
// expectedname, disks[0].Name)
// }
// }
//
// func TestListDisksWithVMIDInFilter(t *testing.T) {
// expectedregionid := 4
// disks, err := h.DescribeDisks(hosting.DiskFilter{RegionID: expectedregionid})
// if err != nil {
// t.Errorf("Error listing disks: %v", err)
// }
// for _, disk := range disks {
// if disk.RegionID != expectedregionid {
// t.Errorf("Error, expected to get Disk in region %d, got region %d instead",
// expectedregionid, disk.RegionID)
// }
// }
// }
//
// func TestDeleteDisk(t *testing.T) {
// disk := hosting.Disk{ID: 1}
// err := h.DeleteDisk(&disk)
// if err != nil {
// t.Errorf("Error deleting disk: %v", err)
// }
// if disk.State != "deleted" {
// t.Errorf("Error, disk state should be 'deleted' but is %s", disk.State)
// }
// }
//
// func TestExtendDisk2GB(t *testing.T) {
// disk := hosting.Disk{ID: 1, Size: defaultSize}
// err := h.ExtendDisk(&disk, 2)
// if err != nil {
// t.Errorf("Error extending disk: %v", err)
// }
// expectedSize := defaultSize + 2*1024
// if disk.Size != expectedSize {
// t.Errorf("Error, expected size of %d, got %d instead", expectedSize, disk.Size)
// }
// }
//
// func TestRenameDisk(t *testing.T) {
// expectedName := "NewName"
// disk := hosting.Disk{ID: 1, Name: diskname}
// err := h.RenameDisk(&disk, expectedName)
// if err != nil {
// t.Errorf("Error renaming disk: %v", err)
// }
// if disk.Name != expectedName {
// t.Errorf("Error, expected disk name to be %s, got %s instead", expectedName, disk.Name)
// }
// }
//
// func TestGetDebian9Image(t *testing.T) {
// expected := "Debian 9"
// region := hosting.Region{ID: region}
// diskimage, err := h.ImageByName(expected, region)
// if err != nil {
// t.Errorf("Error getting image: %v", err)
// }
// if diskimage.Name != expected {
// t.Errorf("Error, expected image %s, got %s instead", expected, diskimage.Name)
// }
// }
//
