package hostingv4

import (
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/PabloPie/Gandi-Go/hosting"
	"github.com/PabloPie/Gandi-Go/mock"
	"github.com/golang/mock/gomock"
)

var (
	vmid    = 1
	vmidstr = "1"
	vmname  = "TestVM"
)

func TestCreateVM(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	now := time.Now()

	paramsVMCreate := []interface{}{
		map[string]interface{}{
			"ip_version":    4,
			"bandwidth":     hosting.DefaultBandwidth,
			"datacenter_id": region,
			"hostname":      vmname,
		},
		map[string]interface{}{
			"datacenter_id": region,
			"size":          disksizeMB,
		}, imageid}
	responseVMCreate := []Operation{{}, {}, {ID: 1, VMID: vmid}}
	creation := mockClient.EXPECT().Send("hosting.vm.create_from",
		paramsVMCreate, gomock.Any()).SetArg(2, responseVMCreate).Return(nil)

	paramsWait := []interface{}{responseVMCreate[2].ID}
	responseWait := operationInfo{responseVMCreate[2].ID, "DONE"}
	wait := mockClient.EXPECT().Send("operation.info",
		paramsWait, gomock.Any()).SetArg(2, responseWait).Return(nil).After(creation)

	paramsVMInfo := []interface{}{responseVMCreate[2].VMID}
	ipsresponse := []iPAddressv4{{1, "192.168.1.1", region, 4, vmid, "used"}}
	ifaceresponse := []iface{{ipsresponse, region, 1, vmid}}
	diskresponse := []diskv4{{1, "sysdisk_1", disksizeMB, region, "created", "data", []int{vmid}, true}}
	responseVMInfo := vmv4{vmid, vmname, region, "", "", 1, 512, now, ifaceresponse, diskresponse, []int{1, 2, 3}, "running"}
	mockClient.EXPECT().Send("hosting.vm.info",
		paramsVMInfo, gomock.Any()).SetArg(2, responseVMInfo).Return(nil).After(wait)

	vmspec := VMSpec{
		RegionID: regionstr,
		Hostname: vmname,
	}
	diskimage := DiskImage{
		DiskID:   imageidstr,
		Size:     3,
		Name:     "Debian 9",
		RegionID: regionstr,
	}
	vm, _, _, err := testHosting.CreateVM(vmspec, diskimage, hosting.IPVersion(4), 20)
	if err != nil {
		log.Println(err)
	}

	expectedIPS := []IPAddress{{"1", "192.168.1.1", regionstr, hosting.IPVersion(4), vmidstr, "used"}}
	expectedDisks := []Disk{{diskidstr, "sysdisk_1", disksize, regionstr, "created", "data", []string{vmidstr}, true}}
	expected := VM{
		ID:          vmidstr,
		Hostname:    vmname,
		RegionID:    regionstr,
		Farm:        "",
		Description: "",
		Cores:       1,
		Memory:      512,
		DateCreated: now,
		Ips:         expectedIPS,
		Disks:       expectedDisks,
		SSHKeysID:   []string{"1", "2", "3"},
		State:       "running",
	}

	if !reflect.DeepEqual(vm, expected) {
		t.Errorf("Error, expected %+v, got instead %+v", expected, vm)
	}
}

func TestDiskDetach(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	now := time.Now()

	paramsDiskDetach := []interface{}{3, 4}
	responseDiskDetach := Operation{ID: 5, DiskID: 4, VMID: 3}
	attach := mockClient.EXPECT().Send("hosting.vm.disk_detach",
		paramsDiskDetach, gomock.Any()).SetArg(2, responseDiskDetach).Return(nil)

	paramsWait := []interface{}{responseDiskDetach.ID}
	responseWait := operationInfo{responseDiskDetach.ID, "DONE"}
	wait := mockClient.EXPECT().Send("operation.info",
		paramsWait, gomock.Any()).SetArg(2, responseWait).Return(nil).After(attach)

	paramsVMInfo := []interface{}{responseDiskDetach.VMID}
	ipsresponse := []iPAddressv4{{2, "192.168.1.1", region, 4, 3, "used"}}
	ifaceresponse := []iface{{ipsresponse, region, 2, 3}}
	diskresponse := []diskv4{{1, "sysdisk_1", disksizeMB, region, "created", "data", []int{3}, true}}
	responseVMInfo := vmv4{3, vmname, region, "", "", 1, 512, now, ifaceresponse, diskresponse, []int{1, 2, 3}, "running"}
	info := mockClient.EXPECT().Send("hosting.vm.info",
		paramsVMInfo, gomock.Any()).SetArg(2, responseVMInfo).Return(nil).After(wait)

	paramsDiskInfo := []interface{}{4}
	responseDiskInfo := diskv4{4, "disk2", 10240, region, "created", "data", []int{}, false}
	mockClient.EXPECT().Send("hosting.disk.info",
		paramsDiskInfo, gomock.Any()).SetArg(2, responseDiskInfo).Return(nil).After(info)

	disks := []Disk{
		{"4", "disk2", 10, regionstr, "created", "data", []string{"3"}, false},
		{"1", "sysdisk_1", disksize, regionstr, "created", "data", []string{"3"}, true},
	}
	vm := VM{ID: "3", Disks: disks}
	disk := Disk{ID: "4"}
	vmres, _, _ := testHosting.DetachDisk(vm, disk)

	expectedDisks := []Disk{{"4", "disk2", 10, regionstr, "created", "data", []string{"3"}, false}}
	expected := VM{
		ID:    "3",
		Disks: expectedDisks,
	}

	if len(vmres.Disks) != len(expected.Disks) {
		t.Errorf("Error, number of disks does not match, expected %v, got %v instead", expected.Disks, vmres.Disks)
	}
}

func TestIPDetach(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	now := time.Now()

	paramsIPInfo := []interface{}{3}
	responseIPInfo := Operation{IfaceID: 2}
	ipinfo := mockClient.EXPECT().Send("hosting.ip.info",
		paramsIPInfo, gomock.Any()).SetArg(2, responseIPInfo).Return(nil)

	paramsIPDetach := []interface{}{3, responseIPInfo.IfaceID}
	responseIPDetach := Operation{ID: 5, IfaceID: responseIPInfo.IfaceID, VMID: 3}
	attach := mockClient.EXPECT().Send("hosting.vm.iface_detach",
		paramsIPDetach, gomock.Any()).SetArg(2, responseIPDetach).Return(nil).After(ipinfo)

	paramsWait := []interface{}{responseIPDetach.ID}
	responseWait := operationInfo{responseIPDetach.ID, "DONE"}
	wait := mockClient.EXPECT().Send("operation.info",
		paramsWait, gomock.Any()).SetArg(2, responseWait).Return(nil).After(attach)

	paramsVMInfo := []interface{}{responseIPDetach.VMID}
	ipsresponse := []iPAddressv4{{2, "192.168.1.1", region, 4, 3, "used"}}
	ifaceresponse := []iface{{ipsresponse, region, 2, vmid}}
	diskresponse := []diskv4{{1, "sysdisk_1", disksizeMB, region, "created", "data", []int{vmid}, true}}
	responseVMInfo := vmv4{vmid, vmname, region, "", "", 1, 512, now, ifaceresponse, diskresponse, []int{1, 2, 3}, "running"}
	info := mockClient.EXPECT().Send("hosting.vm.info",
		paramsVMInfo, gomock.Any()).SetArg(2, responseVMInfo).Return(nil).After(wait)

	paramsIPInfo2 := []interface{}{3}
	responseIPInfo2 := iPAddressv4{3, "192.168.10.2", region, 4, 0, "created"}
	mockClient.EXPECT().Send("hosting.ip.info",
		paramsIPInfo2, gomock.Any()).SetArg(2, responseIPInfo2).Return(nil).After(info)

	ips := []IPAddress{
		{"2", "192.168.1.1", regionstr, hosting.IPVersion(4), "3", "used"},
		{"3", "192.168.10.2", regionstr, hosting.IPVersion(4), "3", "used"},
	}
	vm := VM{ID: "3", Ips: ips}
	ip := IPAddress{ID: "3"}
	vmres, _, _ := testHosting.DetachIP(vm, ip)

	expectedIPS := []IPAddress{{"2", "192.168.1.1", regionstr, hosting.IPVersion(4), "3", "used"}}
	expected := VM{
		ID:  "3",
		Ips: expectedIPS,
	}

	if len(vmres.Ips) != len(expected.Ips) {
		t.Errorf("Error, number of ips does not match, expected %v, got %v instead", expected.Ips, vmres.Ips)
	}
}
