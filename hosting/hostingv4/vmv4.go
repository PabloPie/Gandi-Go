package hostingv4

import (
	"log"
	"strconv"
	"time"

	"github.com/PabloPie/Gandi-Go/hosting"
)

type (
	VM     = hosting.VM
	VMSpec = hosting.VMSpec
)

type vmv4 struct {
	ID          int       `xmlrpc:"id"`
	Hostname    string    `xmlrpc:"hostname"`
	RegionID    int       `xmlrpc:"datacenter_id"`
	Farm        string    `xmlrpc:"farm"`
	Description string    `xmlrpc:"description"`
	Cores       int       `xmlrpc:"cores"`
	Memory      int       `xmlrpc:"memory"`
	DateCreated time.Time `xmlrpc:"date_created"`
	Ifaces      []iface   `xmlrpc:"ifaces"`
	Disks       []diskv4  `xmlrpc:"disks"`
	SSHKeysID   []int     `xmlrpc:"keys"`
	State       string    `xmlrpc:"state"`
}

type vmSpecv4 struct {
	RegionID  int    `xmlrpc:"datacenter_id"`
	Hostname  string `xmlrpc:"hostname"`
	Farm      string `xmlrpc:"farm"`
	Memory    int    `xmlrpc:"memory"`
	Cores     int    `xmlrpc:"cores"`
	SSHKeysID []int  `xmlrpc:"keys"`
	Login     string `xmlrpc:"login"`
	Password  string `xmlrpc:"password"`
}

// CreateVMWithExistingDiskAndIP creates a VM from a VMSpec if a valid IPAddress and Disk are given,
// that is, their IDs already exist. All 3 objects must reside in the same Region.
// `VMSpec.RegionID` is the only mandatory parameter for the VM.
func (h Hostingv4) CreateVMWithExistingDiskAndIP(vm VMSpec, ip IPAddress, disk Disk) (VM, IPAddress, Disk, error) {
	var fn = "CreateVMWithExistingDiskAndIP"
	if vm.RegionID == "" {
		return VM{}, IPAddress{}, Disk{}, &HostingError{fn, "VMSpec", "RegionID", ErrNotProvided}
	}
	if vm.RegionID != disk.RegionID || vm.RegionID != ip.RegionID {
		return VM{}, IPAddress{}, Disk{}, &HostingError{fn, "VMSpec/IPAddress/Disk", "RegionID", ErrMismatch}
	}
	if disk.ID == "" || ip.ID == "" {
		return VM{}, IPAddress{}, Disk{}, &HostingError{fn, "Disk/IPAddress", "ID", ErrNotProvided}
	}

	// Parsing and casting to v4 structs
	vmspec, err := toVMSpecv4(vm)
	if err != nil {
		return VM{}, IPAddress{}, Disk{}, err
	}
	vmspecmap, _ := structToMap(vmspec)

	diskid, err := strconv.Atoi(disk.ID)
	if err != nil {
		return VM{}, IPAddress{}, Disk{}, internalParseError("Disk", "ID")
	}
	ipid, err := strconv.Atoi(ip.ID)
	if err != nil {
		return VM{}, IPAddress{}, Disk{}, internalParseError("IP", "ID")
	}
	// call api to get the iface id that corresponds to the ip
	ifaceid, err := h.ifaceIDFromIPID(ipid)
	if err != nil {
		return VM{}, IPAddress{}, Disk{}, err
	}
	vmspecmap["iface_id"] = ifaceid
	vmspecmap["sys_disk_id"] = diskid

	// Call API, Disk and IP already exist, only one operation is returned
	log.Printf("Creating VM %s...", vm.Hostname)
	request := []interface{}{vmspecmap}
	response := []Operation{}
	err = h.Send("hosting.vm.create", request, &response)
	if err != nil {
		return VM{}, IPAddress{}, Disk{}, err
	}
	err = h.waitForOp(response[0])
	log.Printf("VM %s(ID: %d) created!", vm.Hostname, response[0].VMID)

	vmRes, err := h.vmFromID(response[0].VMID)
	return vmRes, vmRes.Ips[0], vmRes.Disks[0], nil
}

func toVMSpecv4(vm VMSpec) (vmSpecv4, error) {
	regionid, err := strconv.Atoi(vm.RegionID)
	if !isIgnorableErr(err) {
		return vmSpecv4{}, internalParseError("VMSpec", "RegionID")
	}
	var keys []int
	var errkey bool
	for _, key := range vm.SSHKeysID {
		keyid, err := strconv.Atoi(key)
		if err != nil {
			errkey = true
			break
		}
		keys = append(keys, keyid)
	}
	if errkey {
		return vmSpecv4{}, internalParseError("VMSpec", "SSHKeysID")
	}
	return vmSpecv4{
		RegionID:  regionid,
		Hostname:  vm.Hostname,
		Farm:      vm.Farm,
		Cores:     vm.Cores,
		Memory:    vm.Memory,
		SSHKeysID: keys,
		Login:     vm.Login,
		Password:  vm.Password,
	}, nil
}

func fromVMv4(vm vmv4) VM {
	id := strconv.Itoa(vm.ID)
	regionid := strconv.Itoa(vm.RegionID)
	var ips []IPAddress
	// v4 works with interfaces, extract the ips from those
	for _, iface := range vm.Ifaces {
		for _, ip := range iface.IPs {
			ips = append(ips, toIPAddress(ip))
		}
	}
	var disks []Disk
	for _, disk := range vm.Disks {
		disks = append(disks, fromDiskv4(disk))
	}
	var keys []string
	for _, key := range vm.SSHKeysID {
		keys = append(keys, strconv.Itoa(key))
	}
	return VM{
		ID:          id,
		Hostname:    vm.Hostname,
		RegionID:    regionid,
		Farm:        vm.Farm,
		Description: vm.Description,
		Cores:       vm.Cores,
		Memory:      vm.Memory,
		DateCreated: vm.DateCreated,
		Ips:         ips,
		Disks:       disks,
		SSHKeysID:   keys,
		State:       vm.State}
}

func (h Hostingv4) vmFromID(vmid int) (VM, error) {
	response := vmv4{}
	params := []interface{}{vmid}
	err := h.Send("hosting.vm.info", params, &response)
	if err != nil {
		return VM{}, err
	}
	vm := fromVMv4(response)
	return vm, nil
}
