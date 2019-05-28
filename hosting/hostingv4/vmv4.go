package hostingv4

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/PabloPie/go-gandi/hosting"
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

// CreateVMWithExistingDiskAndIP creates a hosting.VM from a hosting.VMSpec if a valid hosting.IPAddress and hosting.Disk are given,
// that is, their IDs already exist.
//
// All 3 objects must reside in the same hosting.Region
// `hosting.VMSpec.RegionID` is the only mandatory parameter for the hosting.VM
func (h Hostingv4) CreateVMWithExistingDiskAndIP(vm hosting.VMSpec, ip hosting.IPAddress, disk hosting.Disk) (hosting.VM, hosting.IPAddress, hosting.Disk, error) {
	vmspecmap, ipid, diskid, _, err := h.checkParametersAndGetVMSpecMap("CreateVMWithExistingDiskAndIP", vm, &ip, &disk, nil)
	if err != nil {
		return hosting.VM{}, hosting.IPAddress{}, hosting.Disk{}, err
	}

	// call api to get the iface id that corresponds to the ip
	ifaceid, err := h.ifaceIDFromIPID(ipid)
	if err != nil {
		return hosting.VM{}, hosting.IPAddress{}, hosting.Disk{}, err
	}
	vmspecmap["iface_id"] = ifaceid
	vmspecmap["sys_disk_id"] = diskid

	// Call API, hosting.Disk and IP already exist, only one operation is returned
	vmid, err := h.createVMFromVMSpecMap(vmspecmap)
	if err != nil {
		return hosting.VM{}, hosting.IPAddress{}, hosting.Disk{}, err
	}

	vmRes, err := h.vmFromID(vmid)
	vmRes.SSHKeys = vm.SSHKeysID
	if err != nil {
		return hosting.VM{}, hosting.IPAddress{}, hosting.Disk{}, err
	}
	return vmRes, vmRes.Ips[0], vmRes.Disks[0], nil
}

// CreateVMWithExistingDisk creates a hosting.VM from a hosting.VMSpec if a valid hosting.Disk is given
//
// The disk must reside in the same hosting.Region as the hosting.VM
// An IP address will also be created in this region and attached to the hosting.VM
// `hosting.VMSpec.RegionID` is mandatory
func (h Hostingv4) CreateVMWithExistingDisk(vm hosting.VMSpec, version hosting.IPVersion, disk hosting.Disk) (hosting.VM, hosting.IPAddress, hosting.Disk, error) {
	vmspecmap, _, diskid, _, err := h.checkParametersAndGetVMSpecMap("CreateVMWithExistingDisk", vm, nil, &disk, nil)
	if err != nil {
		return hosting.VM{}, hosting.IPAddress{}, hosting.Disk{}, err
	}

	vmspecmap["sys_disk_id"] = diskid
	vmspecmap["ip_version"] = int(version)
	vmspecmap["bandwidth"] = hosting.DefaultBandwidth

	vmid, err := h.createVMFromVMSpecMap(vmspecmap)
	if err != nil {
		return hosting.VM{}, hosting.IPAddress{}, hosting.Disk{}, err
	}

	vmRes, err := h.vmFromID(vmid)
	vmRes.SSHKeys = vm.SSHKeysID
	if err != nil {
		return hosting.VM{}, hosting.IPAddress{}, hosting.Disk{}, err
	}
	return vmRes, vmRes.Ips[0], vmRes.Disks[0], nil
}

// CreateVMWithExistingIP creates a hosting.VM from a hosting.VMSpec if a valid hosting.IPAddress and hosting.DiskImage are given
//
// All three objects must be in the same hosting.Region, the new disk will be created in this region
// `hosting.VMSpec.RegionID` is mandatory
func (h Hostingv4) CreateVMWithExistingIP(vm hosting.VMSpec, image hosting.DiskImage, ip hosting.IPAddress, diskSize uint) (hosting.VM, hosting.IPAddress, hosting.Disk, error) {
	vmspecmap, ipid, _, imageid, err := h.checkParametersAndGetVMSpecMap("CreateVMWithExistingIP", vm, &ip, nil, &image)
	if err != nil {
		return hosting.VM{}, hosting.IPAddress{}, hosting.Disk{}, err
	}

	// Get the corresponding ifaceid of the ip
	ifaceid, err := h.ifaceIDFromIPID(ipid)
	if err != nil {
		return hosting.VM{}, hosting.IPAddress{}, hosting.Disk{}, err
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
	log.Printf("[INFO] Creating hosting.VM %s...", vmspecmap["hostname"])
	err = h.Send("hosting.vm.create_from", params, &response)
	if err != nil {
		return hosting.VM{}, hosting.IPAddress{}, hosting.Disk{}, err
	}

	vmop := response[2]
	// Wait for vm operation to finish, disk operation
	// will always end before
	if err = h.waitForOp(vmop); err != nil {
		return hosting.VM{}, hosting.IPAddress{}, hosting.Disk{}, err
	}
	log.Printf("[INFO] hosting.VM %s(ID: %d) created!", vmspecmap["hostname"], response[2].VMID)
	vmRes, err := h.vmFromID(vmop.ID)
	vmRes.SSHKeys = vm.SSHKeysID
	if err != nil {
		return hosting.VM{}, hosting.IPAddress{}, hosting.Disk{}, err
	}
	return vmRes, vmRes.Ips[0], vmRes.Disks[0], nil
}

// CreateVM creates a hosting.VM from scratch, creating also a system disk and an ip address
//
// `hosting.VMSpec.RegionID` is mandatory
func (h Hostingv4) CreateVM(vm hosting.VMSpec, image hosting.DiskImage, version hosting.IPVersion, diskSize uint) (hosting.VM, hosting.IPAddress, hosting.Disk, error) {
	vmspecmap, _, _, imageid, err := h.checkParametersAndGetVMSpecMap("CreateVMWithExistingIP", vm, nil, nil, &image)
	if err != nil {
		return hosting.VM{}, hosting.IPAddress{}, hosting.Disk{}, err
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
	log.Printf("[INFO] Creating hosting.VM %s...", vmspecmap["hostname"])
	err = h.Send("hosting.vm.create_from", params, &response)
	if err != nil {
		return hosting.VM{}, hosting.IPAddress{}, hosting.Disk{}, err
	}

	vmop := response[2]
	// Wait for vm operation to finish, disk operation
	// will always end before
	if err = h.waitForOp(vmop); err != nil {
		return hosting.VM{}, hosting.IPAddress{}, hosting.Disk{}, err
	}
	log.Printf("[INFO] hosting.VM %s(ID: %d) created!", vmspecmap["hostname"], response[2].VMID)
	vmRes, err := h.vmFromID(vmop.ID)
	vmRes.SSHKeys = vm.SSHKeysID
	if err != nil {
		return hosting.VM{}, hosting.IPAddress{}, hosting.Disk{}, err
	}
	return vmRes, vmRes.Ips[0], vmRes.Disks[0], nil
}

// AttachDisk attaches a hosting.Disk to a hosting.VM, both objects must already exist
// and be in the same hosting.Region
func (h Hostingv4) AttachDisk(vm hosting.VM, disk hosting.Disk) (hosting.VM, hosting.Disk, error) {
	var fn = "disk_attach"
	return h.diskAttachDetach(vm, disk, fn, -1)
}

// AttachDiskAtPosition attaches or swaps a hosting.Disk to a hosting.VM at the given position,
// both objects must already exist and be in the same hosting.Region
func (h Hostingv4) AttachDiskAtPosition(vm hosting.VM, disk hosting.Disk, position int) (hosting.VM, hosting.Disk, error) {
	var fn = "disk_attach"
	return h.diskAttachDetach(vm, disk, fn, position)
}

// DetachDisk detaches a hosting.Disk from a hosting.VM, will fail if it is a boot hosting.Disk
func (h Hostingv4) DetachDisk(vm hosting.VM, disk hosting.Disk) (hosting.VM, hosting.Disk, error) {
	var fn = "disk_detach"
	return h.diskAttachDetach(vm, disk, fn, -1)
}

// Attach and detach operations on a disk are almost identical, using a common function
// reduces significantly code size, the variable `op` determines which operation we are calling
func (h Hostingv4) diskAttachDetach(vm hosting.VM, disk hosting.Disk, op string, position int) (hosting.VM, hosting.Disk, error) {
	if vm.RegionID != disk.RegionID {
		return hosting.VM{}, hosting.Disk{}, &HostingError{op, "hosting.VM/hosting.Disk", "RegionID", ErrMismatch}
	}
	vmid, err := strconv.Atoi(vm.ID)
	if err != nil {
		return hosting.VM{}, hosting.Disk{}, internalParseError("hosting.VM", "ID")
	}
	diskid, err := strconv.Atoi(disk.ID)
	if err != nil {
		return hosting.VM{}, hosting.Disk{}, internalParseError("hosting.Disk", "ID")
	}

	params := []interface{}{vmid, diskid}

	if op == "disk_attach" && position >= 0 {
		params = append(params, map[string]interface{}{"position": position})
	}

	response := Operation{}
	err = h.Send("hosting.vm."+op, params, &response)
	if err != nil {
		return hosting.VM{}, hosting.Disk{}, err
	}
	if err = h.waitForOp(response); err != nil {
		return hosting.VM{}, hosting.Disk{}, err
	}

	vmRes, err := h.vmFromID(vmid)
	if err != nil {
		return hosting.VM{}, hosting.Disk{}, err
	}
	diskRes, err := h.diskFromID(diskid)
	if err != nil {
		return hosting.VM{}, hosting.Disk{}, err
	}
	return vmRes, diskRes, nil
}

// AttachIP attaches an IP to a hosting.VM, both objects must already exist
// and be in the same hosting.Region
func (h Hostingv4) AttachIP(vm hosting.VM, ip hosting.IPAddress) (hosting.VM, hosting.IPAddress, error) {
	var fn = "iface_attach"
	return h.ipAttachDetach(vm, ip, fn)
}

// DetachIP detaches an IP from a hosting.VM, meaning the IP will be free
// to be attached to another hosting.VM
func (h Hostingv4) DetachIP(vm hosting.VM, ip hosting.IPAddress) (hosting.VM, hosting.IPAddress, error) {
	var fn = "iface_detach"
	return h.ipAttachDetach(vm, ip, fn)
}

// Same as Disks, attach and detach operations are almost identical
func (h Hostingv4) ipAttachDetach(vm hosting.VM, ip hosting.IPAddress, op string) (hosting.VM, hosting.IPAddress, error) {
	if vm.RegionID != ip.RegionID {
		return hosting.VM{}, hosting.IPAddress{}, &HostingError{op, "hosting.VM/hosting.IPAddress", "RegionID", ErrMismatch}
	}
	vmid, err := strconv.Atoi(vm.ID)
	if err != nil {
		return hosting.VM{}, hosting.IPAddress{}, internalParseError("hosting.VM", "ID")
	}
	ipid, err := strconv.Atoi(ip.ID)
	if err != nil {
		return hosting.VM{}, hosting.IPAddress{}, internalParseError("hosting.IPAddress", "ID")
	}
	// Get corresponding iface id
	ifaceid, err := h.ifaceIDFromIPID(ipid)
	if err != nil {
		return hosting.VM{}, hosting.IPAddress{}, err
	}

	params := []interface{}{vmid, ifaceid}
	response := Operation{}
	err = h.Send("hosting.vm."+op, params, &response)
	if err != nil {
		return hosting.VM{}, hosting.IPAddress{}, err
	}
	if err = h.waitForOp(response); err != nil {
		return hosting.VM{}, hosting.IPAddress{}, err
	}

	vmRes, err := h.vmFromID(vmid)
	if err != nil {
		return hosting.VM{}, hosting.IPAddress{}, err
	}
	ipRes, err := h.ipFromID(ipid)
	if err != nil {
		return hosting.VM{}, hosting.IPAddress{}, err
	}
	return vmRes, ipRes, nil
}

// StartVM starts a stopped hosting.VM
func (h Hostingv4) StartVM(vm hosting.VM) error {
	var fn = "start"
	return h.opVM(vm, fn)
}

// StopVM stops a running hosting.VM
func (h Hostingv4) StopVM(vm hosting.VM) error {
	var fn = "stop"
	return h.opVM(vm, fn)
}

// RebootVM reboots a hosting.VM
func (h Hostingv4) RebootVM(vm hosting.VM) error {
	var fn = "reboot"
	return h.opVM(vm, fn)
}

// DeleteVM deletes a vm
//
// Add cascade option?
// Automatically stop vm before deleting?
func (h Hostingv4) DeleteVM(vm hosting.VM) error {
	var fn = "delete"
	return h.opVM(vm, fn)
}

// Common function for hosting.VM operations
func (h Hostingv4) opVM(vm hosting.VM, op string) error {
	if vm.ID == "" {
		return &HostingError{op, "hosting.VM", "ID", ErrNotProvided}
	}

	vmid, err := strconv.Atoi(vm.ID)
	if err != nil {
		return internalParseError("hosting.VM", "ID")
	}
	params := []interface{}{vmid}
	response := Operation{}
	err = h.Send("hosting.vm."+op, params, &response)
	if err != nil {
		return err
	}
	return h.waitForOp(response)
}

// ListVMs returns a list of VMs filtered with the options provided in `vmfilter`
func (h Hostingv4) ListVMs(vmfilter hosting.VMFilter) ([]hosting.VM, error) {
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
	err = h.Send("hosting.vm.list", params, &response)
	if err != nil {
		return nil, err
	}

	var vms []hosting.VM
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

// VMFromName is a helper function to get a hosting.VM given its name
//
// The function returns an error if the hosting.VM doesn't exist
func (h Hostingv4) VMFromName(name string) (hosting.VM, error) {
	if name == "" {
		return hosting.VM{}, &HostingError{"VMFromName", "-", "name", ErrNotProvided}
	}
	vms, err := h.ListVMs(hosting.VMFilter{Hostname: name})
	if err != nil {
		return hosting.VM{}, err
	}
	if len(vms) < 1 {
		return hosting.VM{}, fmt.Errorf("hosting.VM '%s' does not exist", name)
	}

	return vms[0], nil
}

// ListAllVMs lists every hosting.VM
func (h Hostingv4) ListAllVMs() ([]hosting.VM, error) {
	return h.ListVMs(hosting.VMFilter{})
}

// UpdateVMMemory updates the memory of a hosting.VM, new value can be higher
// or lower than the previous value
func (h Hostingv4) UpdateVMMemory(vm hosting.VM, memory int) (hosting.VM, error) {
	vmupdate := map[string]interface{}{"memory": memory}
	return h.updateVM(vm, vmupdate)

}

// UpdateVMCores updates the number of cores of a hosting.VM
func (h Hostingv4) UpdateVMCores(vm hosting.VM, cores int) (hosting.VM, error) {
	vmupdate := map[string]interface{}{"cores": cores}
	return h.updateVM(vm, vmupdate)

}

// RenameVM renames a hosting.VM
func (h Hostingv4) RenameVM(vm hosting.VM, newname string) (hosting.VM, error) {
	vmupdate := map[string]interface{}{"hostname": newname}
	return h.updateVM(vm, vmupdate)
}

// Common function for update operations
func (h Hostingv4) updateVM(vm hosting.VM, vmupdate map[string]interface{}) (hosting.VM, error) {
	var fn = "UpdateVM"
	if vm.ID == "" {
		return hosting.VM{}, &HostingError{fn, "hosting.VM", "ID", ErrNotProvided}
	}
	vmid, err := strconv.Atoi(vm.ID)
	if err != nil {
		return hosting.VM{}, &HostingError{fn, "hosting.VM", "ID", ErrParse}
	}

	response := Operation{}
	request := []interface{}{vmid, vmupdate}
	err = h.Send("hosting.vm.update", request, &response)
	if err != nil {
		return hosting.VM{}, err
	}
	err = h.waitForOp(response)
	if err != nil {
		return hosting.VM{}, err
	}

	return h.vmFromID(response.VMID)
}

// Helper functions

// vmFromID returns a global hosting.VM object from a v4 id
func (h Hostingv4) vmFromID(vmid int) (hosting.VM, error) {
	response := vmv4{}
	params := []interface{}{vmid}
	err := h.Send("hosting.vm.info", params, &response)
	if err != nil {
		return hosting.VM{}, err
	}
	vm := fromVMv4(response)
	return vm, nil
}

// Internal functions for creation

// Checks parameters of a hosting.VM creation
func (h Hostingv4) checkParametersAndGetVMSpecMap(fn string,
	vm hosting.VMSpec, ip *hosting.IPAddress, disk *hosting.Disk, image *hosting.DiskImage) (map[string]interface{}, int, int, int, error) {
	var ipid int
	var diskid int
	var imageid int
	var err error

	if vm.RegionID == "" {
		return nil, diskid, ipid, imageid, &HostingError{fn, "hosting.VMSpec", "RegionID", ErrNotProvided}
	}

	if disk != nil {
		if vm.RegionID != disk.RegionID {
			return nil, diskid, ipid, imageid, &HostingError{fn, "hosting.VMSpec/hosting.Disk", "RegionID", ErrMismatch}
		}
		if disk.ID == "" {
			return nil, diskid, ipid, imageid, &HostingError{fn, "hosting.Disk", "ID", ErrNotProvided}
		}
		diskid, err = strconv.Atoi(disk.ID)
		if err != nil {
			return nil, diskid, ipid, imageid, internalParseError("hosting.Disk", "ID")
		}
	}

	if ip != nil {
		if vm.RegionID != ip.RegionID {
			return nil, diskid, ipid, imageid, &HostingError{fn, "hosting.VMSpec/hosting.IPAddress", "RegionID", ErrMismatch}
		}
		if ip.ID == "" {
			return nil, diskid, ipid, imageid, &HostingError{fn, "hosting.IPAddress", "ID", ErrNotProvided}
		}
		ipid, err = strconv.Atoi(ip.ID)
		if err != nil {
			return nil, diskid, ipid, imageid, internalParseError("hosting.IPAddress", "ID")
		}
	}

	if image != nil {
		if vm.RegionID != image.RegionID {
			return nil, diskid, ipid, imageid, &HostingError{fn, "hosting.VMSpec/hosting.DiskImage", "RegionID", ErrMismatch}
		}
		if image.DiskID == "" {
			return nil, diskid, ipid, imageid, &HostingError{fn, "hosting.DiskImage", "ID", ErrNotProvided}
		}
		imageid, err = strconv.Atoi(image.DiskID)
		if err != nil {
			return nil, diskid, ipid, imageid, internalParseError("hosting.DiskImage", "ID")
		}
	}

	vmspec, err := h.toVMSpecv4(vm)
	if err != nil {
		return nil, diskid, ipid, imageid, err
	}

	vmspecmap, err := structToMap(vmspec)
	return vmspecmap, ipid, diskid, imageid, err
}

func (h Hostingv4) createVMFromVMSpecMap(vmspecmap map[string]interface{}) (int, error) {
	log.Printf("[INFO] Creating hosting.VM %s...", vmspecmap["hostname"])
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

	log.Printf("[INFO] hosting.VM %s(ID: %d) created!", vmspecmap["hostname"], operation.VMID)

	return operation.VMID, nil
}

// Internal functions for type conversion

// Hosting VMFilter -> VMFilter v4
func toVMFilterv4(vmfilter hosting.VMFilter) (vmFilterv4, error) {
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

// Hosting hosting.VMSpec -> hosting.VMSpec v4
func (h Hostingv4) toVMSpecv4(vm hosting.VMSpec) (vmSpecv4, error) {
	regionid, err := strconv.Atoi(vm.RegionID)
	if err != nil {
		return vmSpecv4{}, internalParseError("hosting.VMSpec", "RegionID")
	}
	var keys []int
	for _, key := range vm.SSHKeysID {
		sshkey := h.KeyFromName(key)
		keyid, err := strconv.Atoi(sshkey.ID)
		if err != nil {
			return vmSpecv4{}, errors.New("Key '" + key + "' does not exist")
		}
		keys = append(keys, keyid)
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

// vm v4 -> Hosting hosting.VM
func fromVMv4(vm vmv4) hosting.VM {
	id := strconv.Itoa(vm.ID)
	regionid := strconv.Itoa(vm.RegionID)
	var ips []hosting.IPAddress
	// v4 works with interfaces, extract the ips from them
	for _, iface := range vm.Ifaces {
		for _, ip := range iface.IPs {
			ips = append(ips, toIPAddress(ip))
		}
	}
	var disks []hosting.Disk
	for _, disk := range vm.Disks {
		disks = append(disks, fromDiskv4(disk))
	}
	return hosting.VM{
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
		State:       vm.State}
}
