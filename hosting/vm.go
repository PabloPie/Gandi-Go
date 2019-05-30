package hosting

import "time"

// VMManager represents a service capable of manipulation virtual machine objects in Gandi's platform
type VMManager interface {

	// CreateVM creates a VM, an IP and a Disk from the DiskImage provided with the size given
	//
	// It returns all the objects created
	// It is used when one doesn't have any IPs or valid Disks already created
	CreateVM(vm VMSpec, image DiskImage, version IPVersion, diskSize uint) (VM, IPAddress, Disk, error)

	// CreateVMWithExistingIP creates a VM and a Disk from the DiskImage provided with the size given
	//
	// It returns the three objects, VM and Disk being new, and IPAdress updated
	// Useful if one already has an IPAddress that he wants to keep for a new VM
	CreateVMWithExistingIP(vm VMSpec, image DiskImage, ip IPAddress, diskSize uint) (VM, IPAddress, Disk, error)

	// CreateVMWithExistingDisk creates a VM and an IP, when a Disk that can be used as boot is given
	//
	// It returns the three objects, VM and IPAddress being new, and Disk updated
	CreateVMWithExistingDisk(vm VMSpec, version IPVersion, disk Disk) (VM, IPAddress, Disk, error)

	// CreateVMWithExistingDiskAndIP creates a VM, when a Disk that can be used as boot and a valid
	// IP are given
	//
	// It returns the three objects, VM being new, and IPAddress and Disk updated
	// This function is used when the three objects are created separately
	CreateVMWithExistingDiskAndIP(vm VMSpec, ip IPAddress, disk Disk) (VM, IPAddress, Disk, error)

	// AttachDisk attaches a Disk to a VM, if they are in the same Region
	//
	// It returns the updated objects
	// The Disk is attached at the next available position
	AttachDisk(vm VM, disk Disk) (VM, Disk, error)

	// AttachDiskAtPosition attaches a Disk to a VM at the position specified
	// if they are in the same Region
	//
	// It returns the updated objects
	AttachDiskAtPosition(vm VM, disk Disk, position int) (VM, Disk, error)

	// DetachDisk detaches a Disk from a VM
	DetachDisk(vm VM, disk Disk) (VM, Disk, error)

	// AttachIP attaches an IPAddress to a VM if both are in the same Region
	AttachIP(vm VM, ip IPAddress) (VM, IPAddress, error)

	// DetachIP detaches an IPAddress from a VM
	DetachIP(vm VM, ip IPAddress) (VM, IPAddress, error)

	// Operations on VM state
	StartVM(vm VM) error
	StopVM(vm VM) error
	RebootVM(vm VM) error

	// DeleteVM deletes a VM
	//
	// To be able to delete a VM it is first needed to call StopVM
	// The first IP and the Disk used as boot disk are deleted with
	// the VM if they are not detached after stopping the VM
	DeleteVM(vm VM) error

	// VMFromName returns a VM given a name
	//
	// If a VM with name provided does not exist,
	// an error is returned
	VMFromName(name string) (VM, error)

	// ListVMs return a list of VMs, filtered with the options
	// given in the VMFilter
	//
	// An unset field in `vmfilter` is ignored when making the
	// request
	ListVMs(vmfilter VMFilter) ([]VM, error)

	// ListAllVMs is a helper function to list every VM
	//
	// It is equivalent to calling ListVMs with an empty filter
	ListAllVMs() ([]VM, error)

	// UpdateVMMemory updates the memory of a VM, this value
	// can increase and decrease
	UpdateVMMemory(vm VM, memory int) (VM, error)

	// UpdateVMCores updates the number of cores of the VM
	// this value can increase and decrease
	UpdateVMCores(vm VM, cores int) (VM, error)

	// RenameVM renames a VM
	RenameVM(vm VM, newname string) (VM, error)
}

// VM represents a virtual machine
type VM struct {

	// ID of the object
	ID string

	// Name of the VM
	Hostname string

	// ID of the Region the VM is in
	RegionID string

	// Farm tag
	Farm string

	// We actually forgot about
	// this one during the creation
	// so it will be empty
	Description string

	// Number of cores
	Cores int

	// Memory in MB
	Memory int

	// Time of creation of the VM
	DateCreated time.Time

	// List of IPAddresses of the VM
	Ips []IPAddress

	// List of Disks of the VM
	// Disk at position 0 is the boot disk
	Disks []Disk

	// List of SSHKeys the VM was created with
	SSHKeys []string

	// State of the VM:
	// paused, running, halted, locked,
	// being_created, deleted, being_migrated
	State string
}

// VMSpec contains the parameters
// specified by the user to create a VM
type VMSpec struct {

	// ID of the Region the VM
	// will be in
	RegionID string

	// Name of the VM
	Hostname string

	// Farm is an optional parameter to
	// group a set of VMs
	// It's just a tag...
	Farm string

	// Memory in MB
	Memory int

	// Number of cores
	Cores int

	// List of SSHKey names to be copied
	// inside the VM on creation
	SSHKeysID []string

	// Optional login and password to login
	// into the VM
	Login    string
	Password string
}

// VMFilter is used to list virtual machines,
// filtered with the options provided
type VMFilter struct {
	RegionID string
	Farm     string
	Hostname string
	ID       string
	State    string
}
