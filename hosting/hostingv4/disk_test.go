package hostingv4

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/PabloPie/Gandi-Go/mock"
	"github.com/golang/mock/gomock"
)

var (
	// default values
	defaultRegion = 3
	defaultSize   = 10240
	// expected params
	diskid     = 1
	diskidstr  = "1"
	imageid    = 100
	imageidstr = "100"
	diskname   = "Disk1"
	disksize   = 20480
	region     = 4
	regionstr  = "4"
)

var disks = []diskv4{
	diskv4{1, "sys_disk1", 10240, 4, "created", "data", []int{1}, true},
	diskv4{4, "sys_disk3", 10240, 4, "created", "data", []int{3}, true},
	diskv4{2, "sys_disk2", 20480, 3, "created", "data", []int{2}, true},
	diskv4{3, "disk3", 1024, 3, "created", "data", []int{2}, false},
	diskv4{5, diskname, 10240, 3, "created", "data", []int{}, false},
}

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
		DiskID: diskid,
	}
	creation := mockClient.EXPECT().Send("hosting.disk.create",
		paramsDiskCreate, gomock.Any()).SetArg(2, responseDiskCreate).Return(nil)

	paramsWait := []interface{}{responseDiskCreate.ID}
	responseWait := operationInfo{responseDiskCreate.ID, "DONE"}
	wait := mockClient.EXPECT().Send("operation.info",
		paramsWait, gomock.Any()).SetArg(2, responseWait).Return(nil).After(creation)

	paramsDiskInfo := []interface{}{responseDiskCreate.DiskID}
	responseDiskInfo := diskv4{diskid, diskname, disksize, region, "created", "data", []int{}, false}
	mockClient.EXPECT().Send("hosting.disk.info",
		paramsDiskInfo, gomock.Any()).SetArg(2, responseDiskInfo).Return(nil).After(wait)

	diskspec := DiskSpec{
		RegionID: regionstr,
		Name:     diskname,
		Size:     uint(disksize),
	}
	disk, _ := testHosting.CreateDisk(diskspec)

	expected := Disk{
		ID:       diskidstr,
		Name:     diskname,
		Size:     uint(disksize),
		RegionID: regionstr,
		State:    "created",
		Type:     "data",
		BootDisk: false,
	}

	if !reflect.DeepEqual(disk, expected) {
		t.Errorf("Error, expected %+v, got instead %+v", expected, disk)
	}
}

func TestCreateDiskFromImageWithoutSize(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	paramsDiskCreate := []interface{}{map[string]interface{}{
		"datacenter_id": region,
		"name":          diskname,
	}, imageid}
	responseDiskCreate := Operation{
		ID:     1,
		DiskID: diskid,
	}
	creation := mockClient.EXPECT().Send("hosting.disk.create_from",
		paramsDiskCreate, gomock.Any()).SetArg(2, responseDiskCreate).Return(nil)

	paramsWait := []interface{}{responseDiskCreate.ID}
	responseWait := operationInfo{responseDiskCreate.ID, "DONE"}
	wait := mockClient.EXPECT().Send("operation.info",
		paramsWait, gomock.Any()).SetArg(2, responseWait).Return(nil).After(creation)

	paramsDiskInfo := []interface{}{responseDiskCreate.DiskID}
	responseDiskInfo := diskv4{diskid, diskname, 3072, region, "created", "data", []int{}, false}
	mockClient.EXPECT().Send("hosting.disk.info",
		paramsDiskInfo, gomock.Any()).SetArg(2, responseDiskInfo).Return(nil).After(wait)

	diskspec := DiskSpec{
		RegionID: regionstr,
		Name:     diskname,
	}
	diskimage := DiskImage{
		DiskID:   imageidstr,
		Size:     3072,
		Name:     "Debian 9",
		RegionID: regionstr,
	}
	disk, _ := testHosting.CreateDiskFromImage(diskspec, diskimage)

	expected := Disk{
		ID:       diskidstr,
		Name:     diskname,
		Size:     3072,
		RegionID: regionstr,
		State:    "created",
		Type:     "data",
		BootDisk: false,
	}

	if !reflect.DeepEqual(disk, expected) {
		t.Errorf("Error, expected %+v, got instead %+v", expected, disk)
	}
}

func TestListDisk(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	paramsDiskInfo := []interface{}{}
	responseDiskInfo := disks
	mockClient.EXPECT().Send("hosting.disk.list",
		paramsDiskInfo, gomock.Any()).SetArg(2, responseDiskInfo).Return(nil)

	disks, _ := testHosting.ListDisks()

	if len(disks) < 1 {
		t.Errorf("Error, expected to get at least 1 Disk")
	}
}

func TestListDiskWithEmptyFilterNodisks(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	paramsDiskInfo := []interface{}{}
	responseDiskInfo := []diskv4{}
	mockClient.EXPECT().Send("hosting.disk.list",
		paramsDiskInfo, gomock.Any()).SetArg(2, responseDiskInfo).Return(nil)

	diskfilter := DiskFilter{}
	disks, _ := testHosting.DescribeDisks(diskfilter)

	if len(disks) > 0 {
		t.Errorf("Error, expected to get no Disks")
	}
}

func TestListDiskWithNameInFilter(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	paramsDiskInfo := []interface{}{map[string]interface{}{
		"name": diskname,
	}}
	responseDiskInfo := disks[4:]
	mockClient.EXPECT().Send("hosting.disk.list",
		paramsDiskInfo, gomock.Any()).SetArg(2, responseDiskInfo).Return(nil)

	diskfilter := DiskFilter{Name: diskname}
	disks, _ := testHosting.DescribeDisks(diskfilter)

	if len(disks) != 1 {
		t.Errorf("Error, expected to get only 1 Disk and got %d instead", len(disks))
	}
	if disks[0].Name != diskname {
		t.Errorf("Error, expected to get Disk with name '%s', got '%s' instead",
			diskname, disks[0].Name)
	}
}

func TestDeleteDisk(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	paramsDiskDelete := []interface{}{diskid}
	responseDiskDelete := Operation{
		ID:     1,
		DiskID: diskid,
	}
	delete := mockClient.EXPECT().Send("hosting.disk.delete",
		paramsDiskDelete, gomock.Any()).SetArg(2, responseDiskDelete).Return(nil)

	paramsWait := []interface{}{responseDiskDelete.ID}
	responseWait := operationInfo{responseDiskDelete.ID, "DONE"}
	mockClient.EXPECT().Send("operation.info",
		paramsWait, gomock.Any()).SetArg(2, responseWait).Return(nil).After(delete)

	err := testHosting.DeleteDisk(Disk{ID: diskidstr})
	if err != nil {
		t.Errorf("Error, expected disk to be deleted, got error '%v' instead", err)
	}
}
func TestExtendDisk(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	var extendSize uint = 5
	var sizeInMB = 5120

	paramsDiskExtend := []interface{}{
		disks[0].ID,
		map[string]int{"size": disks[0].Size + sizeInMB},
	}
	responseDiskExtend := Operation{
		ID:     1,
		DiskID: disks[0].ID,
	}
	delete := mockClient.EXPECT().Send("hosting.disk.update",
		paramsDiskExtend, gomock.Any()).SetArg(2, responseDiskExtend).Return(nil)

	paramsWait := []interface{}{responseDiskExtend.ID}
	responseWait := operationInfo{responseDiskExtend.ID, "DONE"}
	wait := mockClient.EXPECT().Send("operation.info",
		paramsWait, gomock.Any()).SetArg(2, responseWait).Return(nil).After(delete)

	paramsDiskInfo := []interface{}{responseDiskExtend.DiskID}
	responseDiskInfo := disks[0]
	responseDiskInfo.Size += sizeInMB
	mockClient.EXPECT().Send("hosting.disk.info",
		paramsDiskInfo, gomock.Any()).SetArg(2, responseDiskInfo).Return(nil).After(wait)

	diskparam := Disk{ID: strconv.Itoa(disks[0].ID), Size: uint(disks[0].Size)}
	disk, _ := testHosting.ExtendDisk(diskparam, extendSize)
	if disk.Size != uint(disks[0].Size+sizeInMB) {
		t.Errorf("Error, expected disk size %dGB, got a size of %dGB instead",
			disks[0].Size+sizeInMB, disk.Size)
	}
}

func TestRenameDisk(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	var newName = "Disk_renamed"

	paramsDiskExtend := []interface{}{
		disks[0].ID,
		map[string]string{"name": newName},
	}
	responseDiskExtend := Operation{
		ID:     1,
		DiskID: disks[0].ID,
	}
	delete := mockClient.EXPECT().Send("hosting.disk.update",
		paramsDiskExtend, gomock.Any()).SetArg(2, responseDiskExtend).Return(nil)

	paramsWait := []interface{}{responseDiskExtend.ID}
	responseWait := operationInfo{responseDiskExtend.ID, "DONE"}
	wait := mockClient.EXPECT().Send("operation.info",
		paramsWait, gomock.Any()).SetArg(2, responseWait).Return(nil).After(delete)

	paramsDiskInfo := []interface{}{responseDiskExtend.DiskID}
	responseDiskInfo := disks[0]
	responseDiskInfo.Name = newName
	mockClient.EXPECT().Send("hosting.disk.info",
		paramsDiskInfo, gomock.Any()).SetArg(2, responseDiskInfo).Return(nil).After(wait)

	diskparam := Disk{ID: strconv.Itoa(disks[0].ID)}
	disk, _ := testHosting.RenameDisk(diskparam, newName)
	if disk.Name != newName {
		t.Errorf("Error, expected disk name to be %s, got '%s' instead",
			newName, disk.Name)
	}
}
