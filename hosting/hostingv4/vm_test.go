package hostingv4

import (
	"log"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/PabloPie/go-gandi/hosting"
	"github.com/PabloPie/go-gandi/mock"
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

	paramsKeyList := []interface{}{
		map[string]string{
			"name": "key1",
		}}
	responseKeyList := []sshkeyv4{{ID: 1}}
	keylist := mockClient.EXPECT().Send("hosting.ssh.list",
		paramsKeyList, gomock.Any()).SetArg(2, responseKeyList).Return(nil)

	paramsKeyInfo := []interface{}{1}
	responseKeyInfo := sshkeyv4{ID: 1}
	keyinfo := mockClient.EXPECT().Send("hosting.ssh.info",
		paramsKeyInfo, gomock.Any()).SetArg(2, responseKeyInfo).Return(nil).After(keylist)

	paramsVMCreate := []interface{}{
		map[string]interface{}{
			"ip_version":    4,
			"bandwidth":     hosting.DefaultBandwidth,
			"datacenter_id": region,
			"hostname":      vmname,
			"keys":          []int{1},
		},
		map[string]interface{}{
			"datacenter_id": region,
			"size":          disksizeMB,
		}, imageid}
	responseVMCreate := []Operation{{}, {}, {ID: 1, VMID: vmid}}
	creation := mockClient.EXPECT().Send("hosting.vm.create_from",
		paramsVMCreate, gomock.Any()).SetArg(2, responseVMCreate).Return(nil).After(keyinfo)

	paramsWait := []interface{}{responseVMCreate[2].ID}
	responseWait := operationInfo{responseVMCreate[2].ID, "DONE"}
	wait := mockClient.EXPECT().Send("operation.info",
		paramsWait, gomock.Any()).SetArg(2, responseWait).Return(nil).After(creation)

	paramsVMInfo := []interface{}{responseVMCreate[2].VMID}
	ipsresponse := []iPAddressv4{{1, "192.168.1.1", region, 4, vmid, "used"}}
	ifaceresponse := []iface{{ipsresponse, region, 1, vmid}}
	diskresponse := []diskv4{{1, "sysdisk_1", disksizeMB, region, "created", "data", []int{vmid}, true}}
	responseVMInfo := vmv4{vmid, vmname, region, "", "", 1, 512, now, ifaceresponse, diskresponse, "running"}
	mockClient.EXPECT().Send("hosting.vm.info",
		paramsVMInfo, gomock.Any()).SetArg(2, responseVMInfo).Return(nil).After(wait)

	vmspec := hosting.VMSpec{
		RegionID:  regionstr,
		Hostname:  vmname,
		SSHKeysID: []string{"key1"},
	}
	diskimage := hosting.DiskImage{
		DiskID:   imageidstr,
		Size:     3,
		Name:     "Debian 9",
		RegionID: regionstr,
	}
	vm, _, _, err := testHosting.CreateVM(vmspec, diskimage, hosting.IPVersion(4), 20)
	if err != nil {
		log.Println(err)
	}
	expectedIPS := []hosting.IPAddress{{ID: "1", IP: "192.168.1.1", RegionID: regionstr,
		Version: hosting.IPVersion(4), VM: vmidstr, State: "used"}}
	expectedDisks := []hosting.Disk{{ID: diskidstr, Name: "sysdisk_1", Size: disksize, RegionID: regionstr,
		State: "created", Type: "data", VM: []string{vmidstr}, BootDisk: true}}
	expected := hosting.VM{
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
		SSHKeys:     []string{"key1"},
		State:       "running",
	}

	if !reflect.DeepEqual(vm, expected) {
		t.Errorf("Error, expected %+v, got instead %+v", expected, vm)
	}
}
func TestCreateVMWithExistingIP(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	now := time.Now()

	paramsIPInfo := []interface{}{1}
	responseIPInfo := Operation{IfaceID: 10}
	ipinfo := mockClient.EXPECT().Send("hosting.ip.info",
		paramsIPInfo, gomock.Any()).SetArg(2, responseIPInfo).Return(nil)

	paramsVMCreate := []interface{}{
		map[string]interface{}{
			"datacenter_id": region,
			"hostname":      vmname,
			"iface_id":      10,
		},
		map[string]interface{}{
			"datacenter_id": region,
			"size":          disksizeMB,
		}, imageid}
	responseVMCreate := []Operation{{}, {}, {ID: 1, VMID: vmid}}
	creation := mockClient.EXPECT().Send("hosting.vm.create_from",
		paramsVMCreate, gomock.Any()).SetArg(2, responseVMCreate).Return(nil).After(ipinfo)

	paramsWait := []interface{}{responseVMCreate[2].ID}
	responseWait := operationInfo{responseVMCreate[2].ID, "DONE"}
	wait := mockClient.EXPECT().Send("operation.info",
		paramsWait, gomock.Any()).SetArg(2, responseWait).Return(nil).After(creation)

	paramsVMInfo := []interface{}{responseVMCreate[2].VMID}
	ipsresponse := []iPAddressv4{{1, "192.168.1.1", region, 4, vmid, "used"}}
	ifaceresponse := []iface{{ipsresponse, region, 10, vmid}}
	diskresponse := []diskv4{{1, "sysdisk_1", disksizeMB, region, "created", "data", []int{vmid}, true}}
	responseVMInfo := vmv4{vmid, vmname, region, "", "", 1, 512, now, ifaceresponse, diskresponse, "running"}
	mockClient.EXPECT().Send("hosting.vm.info",
		paramsVMInfo, gomock.Any()).SetArg(2, responseVMInfo).Return(nil).After(wait)

	vmspec := hosting.VMSpec{
		RegionID: regionstr,
		Hostname: vmname,
	}
	diskimage := hosting.DiskImage{
		DiskID:   imageidstr,
		Size:     3,
		Name:     "Debian 9",
		RegionID: regionstr,
	}
	ip := hosting.IPAddress{
		ID:       "1",
		IP:       "192.168.1.1",
		RegionID: regionstr,
		Version:  hosting.IPVersion(4),
		State:    "created",
	}
	vm, _, _, err := testHosting.CreateVMWithExistingIP(vmspec, diskimage, ip, 20)
	if err != nil {
		log.Println(err)
	}

	expectedIPS := []hosting.IPAddress{{ID: "1", IP: "192.168.1.1", RegionID: regionstr,
		Version: hosting.IPVersion(4), VM: vmidstr, State: "used"}}
	expectedDisks := []hosting.Disk{{ID: diskidstr, Name: "sysdisk_1", Size: disksize, RegionID: regionstr,
		State: "created", Type: "data", VM: []string{vmidstr}, BootDisk: true}}
	expected := hosting.VM{
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
		State:       "running",
	}

	if !reflect.DeepEqual(vm, expected) {
		t.Errorf("Error, expected %+v, got instead %+v", expected, vm)
	}
}

func TestCreateVMWithExistingDiskAndIP(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	now := time.Now()

	paramsIPInfo := []interface{}{1}
	responseIPInfo := Operation{IfaceID: 10}
	ipinfo := mockClient.EXPECT().Send("hosting.ip.info",
		paramsIPInfo, gomock.Any()).SetArg(2, responseIPInfo).Return(nil)

	paramsVMCreate := []interface{}{
		map[string]interface{}{
			"datacenter_id": region,
			"sys_disk_id":   1,
			"hostname":      vmname,
			"iface_id":      10,
		}}
	responseVMCreate := []Operation{{}, {ID: 1, VMID: vmid}}
	creation := mockClient.EXPECT().Send("hosting.vm.create",
		paramsVMCreate, gomock.Any()).SetArg(2, responseVMCreate).Return(nil).After(ipinfo)

	paramsWait := []interface{}{responseVMCreate[1].ID}
	responseWait := operationInfo{responseVMCreate[1].ID, "DONE"}
	wait := mockClient.EXPECT().Send("operation.info",
		paramsWait, gomock.Any()).SetArg(2, responseWait).Return(nil).After(creation)

	paramsVMInfo := []interface{}{responseVMCreate[1].VMID}
	ipsresponse := []iPAddressv4{{1, "192.168.1.1", region, 4, vmid, "used"}}
	ifaceresponse := []iface{{ipsresponse, region, 10, vmid}}
	diskresponse := []diskv4{{1, "sysdisk_1", disksizeMB, region, "created", "data", []int{vmid}, true}}
	responseVMInfo := vmv4{vmid, vmname, region, "", "", 1, 512, now, ifaceresponse, diskresponse, "running"}
	mockClient.EXPECT().Send("hosting.vm.info",
		paramsVMInfo, gomock.Any()).SetArg(2, responseVMInfo).Return(nil).After(wait)

	vmspec := hosting.VMSpec{
		RegionID: regionstr,
		Hostname: vmname,
	}
	disk := hosting.Disk{
		ID:       "1",
		RegionID: regionstr,
	}
	ip := hosting.IPAddress{
		ID:       "1",
		IP:       "192.168.1.1",
		RegionID: regionstr,
		Version:  hosting.IPVersion(4),
		State:    "created",
	}
	vm, _, _, err := testHosting.CreateVMWithExistingDiskAndIP(vmspec, ip, disk)
	if err != nil {
		log.Println(err)
	}

	expectedIPS := []hosting.IPAddress{{ID: "1", IP: "192.168.1.1", RegionID: regionstr,
		Version: hosting.IPVersion(4), VM: vmidstr, State: "used"}}
	expectedDisks := []hosting.Disk{{ID: "1", Name: "sysdisk_1", Size: disksize, RegionID: regionstr,
		State: "created", Type: "data", VM: []string{vmidstr}, BootDisk: true}}
	expected := hosting.VM{
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
	responseVMInfo := vmv4{3, vmname, region, "", "", 1, 512, now, ifaceresponse, diskresponse, "running"}
	info := mockClient.EXPECT().Send("hosting.vm.info",
		paramsVMInfo, gomock.Any()).SetArg(2, responseVMInfo).Return(nil).After(wait)

	paramsDiskInfo := []interface{}{4}
	responseDiskInfo := diskv4{4, "disk2", 10240, region, "created", "data", []int{}, false}
	mockClient.EXPECT().Send("hosting.disk.info",
		paramsDiskInfo, gomock.Any()).SetArg(2, responseDiskInfo).Return(nil).After(info)

	disks := []hosting.Disk{
		{"4", "disk2", 10, regionstr, "created", "data", []string{"3"}, false},
		{"1", "sysdisk_1", disksize, regionstr, "created", "data", []string{"3"}, true},
	}
	vm := hosting.VM{ID: "3", Disks: disks}
	disk := hosting.Disk{ID: "4"}
	vmres, _, _ := testHosting.DetachDisk(vm, disk)

	expectedDisks := []hosting.Disk{{"4", "disk2", 10, regionstr, "created", "data", []string{"3"}, false}}
	expected := hosting.VM{
		ID:    "3",
		Disks: expectedDisks,
	}

	if len(vmres.Disks) != len(expected.Disks) {
		t.Errorf("Error, number of disks does not match, expected %v, got %v instead", expected.Disks, vmres.Disks)
	}
}

func TestDiskAttachAtPosition(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	vmid := 666
	vmidStr := strconv.Itoa(vmid)

	disks := []hosting.Disk{
		{"1", "d1", disksize, regionstr, "created", "data", []string{vmidStr}, true},
		{"2", "d2", disksize, regionstr, "created", "data", []string{vmidStr}, true},
		{"3", "d3", disksize, regionstr, "created", "data", []string{vmidStr}, true},
	}
	diskid := 3
	diskidStr := strconv.Itoa(diskid)
	disk := hosting.Disk{ID: diskidStr}

	position := 0

	vm := hosting.VM{ID: vmidStr, Disks: []hosting.Disk{disks[0], disks[1], disks[2]}}
	expectedVM := hosting.VM{ID: vmidstr, Disks: []hosting.Disk{disks[2], disks[1], disks[0]}}

	parameters := []interface{}{vmid, diskid, map[string]interface{}{"position": position}}
	response := Operation{ID: 1337, DiskID: diskid, VMID: vmid}

	attach := mockClient.EXPECT().Send("hosting.vm.disk_attach",
		parameters, gomock.Any()).SetArg(2, response).Return(nil)

	paramsWait := []interface{}{response.ID}
	responseWait := operationInfo{response.ID, "DONE"}

	wait := mockClient.EXPECT().Send("operation.info",
		paramsWait, gomock.Any()).SetArg(2, responseWait).Return(nil).After(attach)

	paramsVMInfo := []interface{}{response.VMID}

	diskresponse := []diskv4{{3, "d3", disksizeMB, region, "created", "data", []int{vmid}, true},
		{2, "d2", disksizeMB, region, "created", "data", []int{vmid}, true},
		{1, "d1", disksizeMB, region, "created", "data", []int{vmid}, true},
	}

	responseVMInfo := vmv4{Disks: diskresponse}
	info := mockClient.EXPECT().Send("hosting.vm.info",
		paramsVMInfo, gomock.Any()).SetArg(2, responseVMInfo).Return(nil).After(wait)

	mockClient.EXPECT().Send("hosting.disk.info",
		[]interface{}{3}, gomock.Any()).SetArg(2, diskresponse[0]).Return(nil).After(info)

	vmres, _, _ := testHosting.AttachDiskAtPosition(vm, disk, position)

	for i := range vmres.Disks {
		if expectedVM.Disks[i].ID != vmres.Disks[i].ID {
			t.Errorf("Error, disk %d does not match, expected %v, got %v instead", i, expectedVM.Disks[i], vmres.Disks[i])
		}
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
	responseVMInfo := vmv4{vmid, vmname, region, "", "", 1, 512, now, ifaceresponse, diskresponse, "running"}
	info := mockClient.EXPECT().Send("hosting.vm.info",
		paramsVMInfo, gomock.Any()).SetArg(2, responseVMInfo).Return(nil).After(wait)

	paramsIPInfo2 := []interface{}{3}
	responseIPInfo2 := iPAddressv4{3, "192.168.10.2", region, 4, 0, "created"}
	mockClient.EXPECT().Send("hosting.ip.info",
		paramsIPInfo2, gomock.Any()).SetArg(2, responseIPInfo2).Return(nil).After(info)

	ips := []hosting.IPAddress{
		{"2", "192.168.1.1", regionstr, hosting.IPVersion(4), "3", "used"},
		{"3", "192.168.10.2", regionstr, hosting.IPVersion(4), "3", "used"},
	}
	vm := hosting.VM{ID: "3", Ips: ips}
	ip := hosting.IPAddress{ID: "3"}
	vmres, _, _ := testHosting.DetachIP(vm, ip)

	expectedIPS := []hosting.IPAddress{{"2", "192.168.1.1", regionstr, hosting.IPVersion(4), "3", "used"}}
	expected := hosting.VM{
		ID:  "3",
		Ips: expectedIPS,
	}

	if len(vmres.Ips) != len(expected.Ips) {
		t.Errorf("Error, number of ips does not match, expected %v, got %v instead", expected.Ips, vmres.Ips)
	}
}

func TestVMStop(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	paramsVMStop := []interface{}{3}
	responseVMStop := Operation{ID: 5, VMID: 3}
	stop := mockClient.EXPECT().Send("hosting.vm.stop",
		paramsVMStop, gomock.Any()).SetArg(2, responseVMStop).Return(nil)

	paramsWait := []interface{}{responseVMStop.ID}
	responseWait := operationInfo{responseVMStop.ID, "DONE"}
	mockClient.EXPECT().Send("operation.info",
		paramsWait, gomock.Any()).SetArg(2, responseWait).Return(nil).After(stop)

	err := testHosting.StopVM(hosting.VM{ID: "3"})

	if err != nil {
		t.Errorf("Error, %s", err)
	}
}

func TestCreateVMWithExistingDisk(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	now := time.Now()

	paramsVMCreate := []interface{}{
		map[string]interface{}{
			"ip_version":    4,
			"sys_disk_id":   5,
			"bandwidth":     hosting.DefaultBandwidth,
			"datacenter_id": region,
			"hostname":      vmname,
		}}
	responseVMCreate := []Operation{{ID: 3, IfaceID: 5}, {ID: 1, VMID: vmid}}
	creation := mockClient.EXPECT().Send("hosting.vm.create",
		paramsVMCreate, gomock.Any()).SetArg(2, responseVMCreate).Return(nil)

	paramsWait := []interface{}{responseVMCreate[1].ID}
	responseWait := operationInfo{responseVMCreate[1].ID, "DONE"}
	wait := mockClient.EXPECT().Send("operation.info",
		paramsWait, gomock.Any()).SetArg(2, responseWait).Return(nil).After(creation)

	paramsVMInfo := []interface{}{responseVMCreate[1].VMID}
	ipsresponse := []iPAddressv4{{1, "192.168.1.1", region, 4, vmid, "used"}}
	ifaceresponse := []iface{{ipsresponse, region, 1, vmid}}
	diskresponse := []diskv4{{5, "sysdisk_1", disksizeMB, region, "created", "data", []int{vmid}, true}}
	responseVMInfo := vmv4{vmid, vmname, region, "", "", 1, 512, now, ifaceresponse, diskresponse, "running"}
	mockClient.EXPECT().Send("hosting.vm.info",
		paramsVMInfo, gomock.Any()).SetArg(2, responseVMInfo).Return(nil).After(wait)

	vmspec := hosting.VMSpec{
		RegionID: regionstr,
		Hostname: vmname,
	}
	disk := hosting.Disk{
		ID:       "5",
		Name:     "sysdisk_1",
		Size:     disksize,
		RegionID: regionstr,
		State:    "created",
		Type:     "data",
		VM:       []string{},
		BootDisk: false,
	}

	vm, _, _, err := testHosting.CreateVMWithExistingDisk(vmspec, hosting.IPVersion(4), disk)
	if err != nil {
		log.Println(err)
	}

	expectedIPS := []hosting.IPAddress{{"1", "192.168.1.1", regionstr, hosting.IPVersion(4), vmidstr, "used"}}
	expectedDisks := []hosting.Disk{{"5", "sysdisk_1", disksize, regionstr, "created", "data", []string{vmidstr}, true}}
	expected := hosting.VM{
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
		State:       "running",
	}

	if !reflect.DeepEqual(vm, expected) {
		t.Errorf("Error, expected %+v, got instead %+v", expected, vm)
	}
}

func TestVMFromName(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	now := time.Now()

	paramsVMList := []interface{}{map[string]interface{}{"hostname": vmname}}
	responseVMList := []vmv4{{ID: vmid, Hostname: vmname}}
	list := mockClient.EXPECT().Send("hosting.vm.list",
		paramsVMList, gomock.Any()).SetArg(2, responseVMList).Return(nil)

	paramsVMInfo := []interface{}{responseVMList[0].ID}
	ipsresponse := []iPAddressv4{{1, "192.168.1.1", region, 4, vmid, "used"}}
	ifaceresponse := []iface{{ipsresponse, region, 1, vmid}}
	diskresponse := []diskv4{{5, "sysdisk_1", disksizeMB, region, "created", "data", []int{vmid}, true}}
	responseVMInfo := vmv4{vmid, vmname, region, "", "", 1, 512, now, ifaceresponse, diskresponse, "running"}
	mockClient.EXPECT().Send("hosting.vm.info",
		paramsVMInfo, gomock.Any()).SetArg(2, responseVMInfo).Return(nil).After(list)

	vm, _ := testHosting.VMFromName(vmname)

	expectedIPS := []hosting.IPAddress{{"1", "192.168.1.1", regionstr, hosting.IPVersion(4), vmidstr, "used"}}
	expectedDisks := []hosting.Disk{{"5", "sysdisk_1", disksize, regionstr, "created", "data", []string{vmidstr}, true}}
	expected := hosting.VM{
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
		State:       "running",
	}

	if !reflect.DeepEqual(vm, expected) {
		t.Errorf("Error, expected %+v, got instead %+v", expected, vm)
	}
}

func TestRenameVM(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	now := time.Now()

	paramsVMUpdate := []interface{}{vmid, map[string]interface{}{"hostname": "NEWNAME"}}
	responseVMUpdate := Operation{ID: 5, VMID: vmid}
	update := mockClient.EXPECT().Send("hosting.vm.update",
		paramsVMUpdate, gomock.Any()).SetArg(2, responseVMUpdate).Return(nil)

	paramsWait := []interface{}{responseVMUpdate.ID}
	responseWait := operationInfo{responseVMUpdate.ID, "DONE"}
	wait := mockClient.EXPECT().Send("operation.info",
		paramsWait, gomock.Any()).SetArg(2, responseWait).Return(nil).After(update)

	paramsVMInfo := []interface{}{vmid}
	ipsresponse := []iPAddressv4{{1, "192.168.1.1", region, 4, vmid, "used"}}
	ifaceresponse := []iface{{ipsresponse, region, 1, vmid}}
	diskresponse := []diskv4{{5, "sysdisk_1", disksizeMB, region, "created", "data", []int{vmid}, true}}
	responseVMInfo := vmv4{vmid, "NEWNAME", region, "", "", 1, 512, now, ifaceresponse, diskresponse, "running"}
	mockClient.EXPECT().Send("hosting.vm.info",
		paramsVMInfo, gomock.Any()).SetArg(2, responseVMInfo).Return(nil).After(wait)

	vmreq := hosting.VM{ID: vmidstr}
	vm, _ := testHosting.RenameVM(vmreq, "NEWNAME")

	expectedIPS := []hosting.IPAddress{{"1", "192.168.1.1", regionstr, hosting.IPVersion(4), vmidstr, "used"}}
	expectedDisks := []hosting.Disk{{"5", "sysdisk_1", disksize, regionstr, "created", "data", []string{vmidstr}, true}}
	expected := hosting.VM{
		ID:          vmidstr,
		Hostname:    "NEWNAME",
		RegionID:    regionstr,
		Farm:        "",
		Description: "",
		Cores:       1,
		Memory:      512,
		DateCreated: now,
		Ips:         expectedIPS,
		Disks:       expectedDisks,
		State:       "running",
	}

	if !reflect.DeepEqual(vm, expected) {
		t.Errorf("Error, expected %+v, got instead %+v", expected, vm)
	}
}
