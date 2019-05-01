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
