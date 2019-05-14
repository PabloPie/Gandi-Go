package hostingv4

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/PabloPie/Gandi-Go/hosting"
)

type (
	VM       = hosting.VM
	VMSpec   = hosting.VMSpec
	VMFilter = hosting.VMFilter
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

type vmFilterv4 struct {
	RegionID int    `xmlrpc:"datacenter_id"`
	Farm     string `xmlrpc:"farm"`
	Hostname string `xmlrpc:"hostname"`
	ID       int    `xmlrpc:"id"`
	State    string `xmlrpc:"state"`
}

// CreateVMWithExistingDiskAndIP creates a VM from a VMSpec if a valid IPAddress and Disk are given,
// that is, their IDs already exist. All 3 objects must reside in the same Region.
// `VMSpec.RegionID` is the only mandatory parameter for the VM.
func (h Hostingv4) CreateVMWithExistingDiskAndIP(vm VMSpec, ip IPAddress, disk Disk) (VM, IPAddress, Disk, error) {
	vmspecmap, ipid, diskid, _, err := checkParametersAndGetVMSpecMap("CreateVMWithExistingDiskAndIP", vm, &ip, &disk, nil)
	if err != nil {
		return VM{}, IPAddress{}, Disk{}, err
	}

	// call api to get the iface id that corresponds to the ip
	ifaceid, err := h.ifaceIDFromIPID(ipid)
	if err != nil {
		return VM{}, IPAddress{}, Disk{}, err
	}
	vmspecmap["iface_id"] = ifaceid
	vmspecmap["sys_disk_id"] = diskid

	// Call API, Disk and IP already exist, only one operation is returned
	vmid, err := h.createVMFromVMSpecMap(vmspecmap)
	if err != nil {
		return VM{}, IPAddress{}, Disk{}, err
	}

	vmRes, err := h.vmFromID(vmid)
	if err != nil {
		return VM{}, IPAddress{}, Disk{}, err
	}
	return vmRes, vmRes.Ips[0], vmRes.Disks[0], nil
}

// CreateVMWithExistingDisk creates a VM from a VMSpec if a valid Disk is given,
// The disk must reside in the same Region as the VM and an IP address
// will be created in this region.
// `VMSpec.RegionID` is the only mandatory parameter for the VM.
func (h Hostingv4) CreateVMWithExistingDisk(vm VMSpec, version hosting.IPVersion, disk Disk) (VM, IPAddress, Disk, error) {
	vmspecmap, _, diskid, _, err := checkParametersAndGetVMSpecMap("CreateVMWithExistingDisk", vm, nil, &disk, nil)
	if err != nil {
		return VM{}, IPAddress{}, Disk{}, err
	}

	vmspecmap["sys_disk_id"] = diskid
	vmspecmap["ip_version"] = int(version)
	vmspecmap["bandwidth"] = hosting.DefaultBandwidth

	vmid, err := h.createVMFromVMSpecMap(vmspecmap)
	if err != nil {
		return VM{}, IPAddress{}, Disk{}, err
	}

	vmRes, err := h.vmFromID(vmid)
	if err != nil {
		return VM{}, IPAddress{}, Disk{}, err
	}
	return vmRes, vmRes.Ips[0], vmRes.Disks[0], nil
}

// CreateVMWithExistingIP creates a VM from a VMSpec if valid IPAddress and DiskImage are given
// All three objects must be in the same Region, the new disk will be created in the same region
// `VMSpec.RegionID` is the only mandatory parameter for VM creation
func (h Hostingv4) CreateVMWithExistingIP(vm VMSpec, image DiskImage, ip IPAddress, diskSize uint) (VM, IPAddress, Disk, error) {
	vmspecmap, ipid, _, imageid, err := checkParametersAndGetVMSpecMap("CreateVMWithExistingIP", vm, &ip, nil, &image)
	if err != nil {
		return VM{}, IPAddress{}, Disk{}, err
	}

	// Get the corresponding ifaceid of the ip
	ifaceid, err := h.ifaceIDFromIPID(ipid)
	if err != nil {
		return VM{}, IPAddress{}, Disk{}, err
	}
	vmspecmap["iface_id"] = ifaceid
	diskspec := diskSpecv4{
		// Docs say datacenter_id is an optional parameter
		RegionID: vmspecmap["datacenter_id"].(int),
		Size:     int(diskSize) * 1024,
	}
	diskparam, _ := structToMap(diskspec)

	params := []interface{}{vmspecmap, diskparam, imageid}
	response := []Operation{}
	log.Printf("[INFO] Creating VM %s...", vmspecmap["hostname"])
	err = h.Send("hosting.vm.create_from", params, &response)
	if err != nil {
		return VM{}, IPAddress{}, Disk{}, err
	}

	vmop := response[2]
	// Wait for vm operation to finish, disk operation
	// will always end before
	if err = h.waitForOp(vmop); err != nil {
		return VM{}, IPAddress{}, Disk{}, err
	}
	log.Printf("[INFO] VM %s(ID: %d) created!", vmspecmap["hostname"], response[2].VMID)
	vmRes, err := h.vmFromID(vmop.ID)
	if err != nil {
		return VM{}, IPAddress{}, Disk{}, err
	}
	return vmRes, vmRes.Ips[0], vmRes.Disks[0], nil
}

// CreateVM creates a VM together with its system disk and an ip address
// `VMSpec.RegionID` is the only mandatory parameter for VM creation
func (h Hostingv4) CreateVM(vm VMSpec, image DiskImage, version hosting.IPVersion, diskSize uint) (VM, IPAddress, Disk, error) {
	vmspecmap, _, _, imageid, err := checkParametersAndGetVMSpecMap("CreateVMWithExistingIP", vm, nil, nil, &image)
	if err != nil {
		return VM{}, IPAddress{}, Disk{}, err
	}

	vmspecmap["ip_version"] = int(version)
	vmspecmap["bandwidth"] = hosting.DefaultBandwidth
	diskspec := diskSpecv4{
		// Docs say datacenter_id is an optional parameter
		RegionID: vmspecmap["datacenter_id"].(int),
		Size:     int(diskSize) * 1024,
	}
	diskparam, _ := structToMap(diskspec)
	params := []interface{}{vmspecmap, diskparam, imageid}
	response := []Operation{}
	log.Printf("[INFO] Creating VM %s...", vmspecmap["hostname"])
	err = h.Send("hosting.vm.create_from", params, &response)
	if err != nil {
		return VM{}, IPAddress{}, Disk{}, err
	}

	vmop := response[2]
	// Wait for vm operation to finish, disk operation
	// will always end before
	if err = h.waitForOp(vmop); err != nil {
		return VM{}, IPAddress{}, Disk{}, err
	}
	log.Printf("[INFO] VM %s(ID: %d) created!", vmspecmap["hostname"], response[2].VMID)
	vmRes, err := h.vmFromID(vmop.ID)
	if err != nil {
		return VM{}, IPAddress{}, Disk{}, err
	}
	return vmRes, vmRes.Ips[0], vmRes.Disks[0], nil
}

// AttachDisk attaches a Disk to a VM, both objects must already exist
func (h Hostingv4) AttachDisk(vm VM, disk Disk) (VM, Disk, error) {
	var fn = "disk_attach"
	return h.diskAttachDetach(vm, disk, fn, false)
}

// ChangeBootDisk attaches or swaps a Disk to a VM as the boot disk, both objects must already exist
func (h Hostingv4) ChangeBootDisk(vm VM, disk Disk) (VM, Disk, error) {
	var fn = "disk_attach"
	return h.diskAttachDetach(vm, disk, fn, true)
}

// DetachDisk detaches a Disk from a VM, will fail if it is a boot Disk
func (h Hostingv4) DetachDisk(vm VM, disk Disk) (VM, Disk, error) {
	var fn = "disk_detach"
	return h.diskAttachDetach(vm, disk, fn, false)
}

// Attach and detach operations on a disk are almost identical, using a common function
// reduces significantly code size, the variable `op` determines which operation we are calling
func (h Hostingv4) diskAttachDetach(vm VM, disk Disk, op string, is_boot bool) (VM, Disk, error) {
	if vm.RegionID != disk.RegionID {
		return VM{}, Disk{}, &HostingError{op, "VM/Disk", "RegionID", ErrMismatch}
	}
	vmid, err := strconv.Atoi(vm.ID)
	if err != nil {
		return VM{}, Disk{}, internalParseError("VM", "ID")
	}
	diskid, err := strconv.Atoi(disk.ID)
	if err != nil {
		return VM{}, Disk{}, internalParseError("Disk", "ID")
	}

	params := []interface{}{vmid, diskid}
	
	if op == "disk_attach" && is_boot {
		params = append(params,map[string]interface{}{"position": 0})
	}
	
	response := Operation{}
	err = h.Send("hosting.vm."+op, params, &response)
	if err != nil {
		return VM{}, Disk{}, err
	}
	if err = h.waitForOp(response); err != nil {
		return VM{}, Disk{}, err
	}

	vmRes, err := h.vmFromID(vmid)
	if err != nil {
		return VM{}, Disk{}, err
	}
	diskRes, err := h.diskFromID(diskid)
	if err != nil {
		return VM{}, Disk{}, err
	}
	return vmRes, diskRes, nil
}

// AttachIP attaches an IP to a VM, both objects must already exist
func (h Hostingv4) AttachIP(vm VM, ip IPAddress) (VM, IPAddress, error) {
	var fn = "iface_attach"
	return h.ipAttachDetach(vm, ip, fn)
}

// DetachIP detaches an IP from a VM, meaning the IP will be free
// to be attached to another VM
func (h Hostingv4) DetachIP(vm VM, ip IPAddress) (VM, IPAddress, error) {
	var fn = "iface_detach"
	return h.ipAttachDetach(vm, ip, fn)
}

// Same as Disks, attach and detach operations are almost identical
func (h Hostingv4) ipAttachDetach(vm VM, ip IPAddress, op string) (VM, IPAddress, error) {
	if vm.RegionID != ip.RegionID {
		return VM{}, IPAddress{}, &HostingError{op, "VM/IPAddress", "RegionID", ErrMismatch}
	}
	vmid, err := strconv.Atoi(vm.ID)
	if err != nil {
		return VM{}, IPAddress{}, internalParseError("VM", "ID")
	}
	ipid, err := strconv.Atoi(ip.ID)
	if err != nil {
		return VM{}, IPAddress{}, internalParseError("IPAddress", "ID")
	}
	// Get corresponding iface id
	ifaceid, err := h.ifaceIDFromIPID(ipid)
	if err != nil {
		return VM{}, IPAddress{}, err
	}

	params := []interface{}{vmid, ifaceid}
	response := Operation{}
	err = h.Send("hosting.vm."+op, params, &response)
	if err != nil {
		return VM{}, IPAddress{}, err
	}
	if err = h.waitForOp(response); err != nil {
		return VM{}, IPAddress{}, err
	}

	vmRes, err := h.vmFromID(vmid)
	if err != nil {
		return VM{}, IPAddress{}, err
	}
	ipRes, err := h.ipFromID(ipid)
	if err != nil {
		return VM{}, IPAddress{}, err
	}
	return vmRes, ipRes, nil
}

// StartVM starts a stopped VM
func (h Hostingv4) StartVM(vm VM) error {
	var fn = "start"
	return h.opVM(vm, fn)
}

// StopVM stops a running VM
func (h Hostingv4) StopVM(vm VM) error {
	var fn = "stop"
	return h.opVM(vm, fn)
}

// RebootVM reboots a VM
func (h Hostingv4) RebootVM(vm VM) error {
	var fn = "reboot"
	return h.opVM(vm, fn)
}

// DeleteVM deletes a vm
// Add cascade option?
// Automatically stop vm before deleting?
func (h Hostingv4) DeleteVM(vm VM) error {
	var fn = "delete"
	return h.opVM(vm, fn)
}

// Common function for VM operation
func (h Hostingv4) opVM(vm VM, op string) error {
	if vm.ID == "" {
		return &HostingError{op, "VM", "ID", ErrNotProvided}
	}

	vmid, err := strconv.Atoi(vm.ID)
	if err != nil {
		return internalParseError("VM", "ID")
	}
	params := []interface{}{vmid}
	response := Operation{}
	err = h.Send("hosting.vm."+op, params, &response)
	if err != nil {
		return err
	}
	return h.waitForOp(response)
}

// DescribeVM returns a list of VMs filtered with the options provided in `vmfilter`
func (h Hostingv4) DescribeVM(vmfilter VMFilter) ([]VM, error) {
	filterv4, err := toVMFilterv4(vmfilter)
	if err != nil {
		return nil, err
	}
	filter, _ := structToMap(filterv4)

	response := []vmv4{}
	params := []interface{}{}
	if len(filter) > 0 {
		params = append(params, filter)
	}
	// disk.list and disk.info return the same information
	err = h.Send("hosting.vm.list", params, &response)
	if err != nil {
		return nil, err
	}

	var vms []VM
	for _, vmv4 := range response {
		// vm list does not a contain the full description
		// call vm info to get a vm's interfaces and disks
		vm, err := h.vmFromID(vmv4.ID)
		if err != nil {
			log.Printf("[WARN] Error getting %s (ID: %s) information, excluded from list: %s", vm.Hostname, vm.ID, err)
			continue
		}
		vms = append(vms, vm)
	}
	return vms, nil
}

// VMFromName is a helper function to get a VM given its name
func (h Hostingv4) VMFromName(name string) (VM, error) {
	disks, err := h.DescribeVM(VMFilter{Hostname: name})
	if err != nil {
		return VM{}, err
	}
	if len(disks) < 1 {
		return VM{}, fmt.Errorf("VM '%s' does not exist", name)
	}

	return disks[0], nil
}

// ListVMs lists every VM
func (h Hostingv4) ListVMs() ([]VM, error) {
	return h.DescribeVM(VMFilter{})
}

// UpdateVMMemory updates the memory of a VM, new value can be higher
// or lower than the previous value
func (h Hostingv4) UpdateVMMemory(vm VM, memory int) (VM, error) {
	vmupdate := map[string]interface{}{"memory": memory}
	return h.updateVM(vm, vmupdate)

}

// UpdateVMCores updates the number of cores of a VM
func (h Hostingv4) UpdateVMCores(vm VM, cores int) (VM, error) {
	vmupdate := map[string]interface{}{"cores": cores}
	return h.updateVM(vm, vmupdate)

}

// RenameVM renames a VM
func (h Hostingv4) RenameVM(vm VM, newname string) (VM, error) {
	vmupdate := map[string]interface{}{"hostname": newname}
	return h.updateVM(vm, vmupdate)
}

func (h Hostingv4) updateVM(vm VM, vmupdate map[string]interface{}) (VM, error) {
	var fn = "UpdateVM"
	if vm.ID == "" {
		return VM{}, &HostingError{fn, "VM", "ID", ErrNotProvided}
	}
	vmid, err := strconv.Atoi(vm.ID)
	if err != nil {
		return VM{}, &HostingError{fn, "VM", "ID", ErrParse}
	}

	response := Operation{}
	request := []interface{}{vmid, vmupdate}
	err = h.Send("hosting.vm.update", request, &response)
	if err != nil {
		return VM{}, err
	}
	err = h.waitForOp(response)
	if err != nil {
		return VM{}, err
	}

	return h.vmFromID(response.VMID)
}

// vmFromID returns a global VM object from a v4 id
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

// Internal functions for creation

func checkParametersAndGetVMSpecMap(fn string,
	vm VMSpec, ip *IPAddress, disk *Disk, image *DiskImage) (map[string]interface{}, int, int, int, error) {
	var ipid int
	var diskid int
	var imageid int
	var err error

	if vm.RegionID == "" {
		return nil, diskid, ipid, imageid, &HostingError{fn, "VMSpec", "RegionID", ErrNotProvided}
	}

	if disk != nil {
		if vm.RegionID != disk.RegionID {
			return nil, diskid, ipid, imageid, &HostingError{fn, "VMSpec/Disk", "RegionID", ErrMismatch}
		}
		if disk.ID == "" {
			return nil, diskid, ipid, imageid, &HostingError{fn, "Disk", "ID", ErrNotProvided}
		}
		diskid, err = strconv.Atoi(disk.ID)
		if err != nil {
			return nil, diskid, ipid, imageid, internalParseError("Disk", "ID")
		}
	}

	if ip != nil {
		if vm.RegionID != ip.RegionID {
			return nil, diskid, ipid, imageid, &HostingError{fn, "VMSpec/IPAddress", "RegionID", ErrMismatch}
		}
		if ip.ID == "" {
			return nil, diskid, ipid, imageid, &HostingError{fn, "IPAddress", "ID", ErrNotProvided}
		}
		ipid, err = strconv.Atoi(ip.ID)
		if err != nil {
			return nil, diskid, ipid, imageid, internalParseError("IPAddress", "ID")
		}
	}

	if image != nil {
		if vm.RegionID != image.RegionID {
			return nil, diskid, ipid, imageid, &HostingError{fn, "VMSpec/DiskImage", "RegionID", ErrMismatch}
		}
		if image.DiskID == "" {
			return nil, diskid, ipid, imageid, &HostingError{fn, "DiskImage", "ID", ErrNotProvided}
		}
		imageid, err = strconv.Atoi(image.DiskID)
		if err != nil {
			return nil, diskid, ipid, imageid, internalParseError("DiskImage", "ID")
		}
	}

	vmspec, err := toVMSpecv4(vm)
	if err != nil {
		return nil, diskid, ipid, imageid, err
	}

	vmspecmap, err := structToMap(vmspec)
	return vmspecmap, ipid, diskid, imageid, err
}

func (h Hostingv4) createVMFromVMSpecMap(vmspecmap map[string]interface{}) (int, error) {
	log.Printf("[INFO] Creating VM %s...", vmspecmap["hostname"])
	request := []interface{}{vmspecmap}
	response := []Operation{}
	if err := h.Send("hosting.vm.create", request, &response); err != nil {
		return -1, err
	}

	operation := response[0]
	if len(response) > 1 {
		operation = response[1]
	}

	if err := h.waitForOp(operation); err != nil {
		return -1, err
	}

	log.Printf("[INFO] VM %s(ID: %d) created!", vmspecmap["hostname"], operation.VMID)

	return operation.VMID, nil
}

// Internal functions for type conversion

func toVMFilterv4(vmfilter VMFilter) (vmFilterv4, error) {
	region := toInt(vmfilter.RegionID)
	if region == -1 {
		return vmFilterv4{}, internalParseError("VMFilter", "RegionID")
	}

	id := toInt(vmfilter.ID)
	if id == -1 {
		return vmFilterv4{}, internalParseError("VMFilter", "ID")
	}

	return vmFilterv4{
		RegionID: region,
		ID:       id,
		Hostname: vmfilter.Hostname,
		Farm:     vmfilter.Farm,
		State:    vmfilter.State,
	}, nil
}

func toVMSpecv4(vm VMSpec) (vmSpecv4, error) {
	regionid, err := strconv.Atoi(vm.RegionID)
	if err != nil {
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
	// v4 works with interfaces, extract the ips from them
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
