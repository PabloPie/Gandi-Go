package hosting

import "time"

// VMManager is an interface embedded in the Hosting interface, not
// implemented by VM!
type VMManager interface {
	//Creates vm with a new disk of size `size` based on diskimage vm.image
	// CreateVM(vm VMSpec, image DiskImage, version IPVersion, diskSize uint) (VM, Disk, IPAddress, error)
	// CreateVMWithExistingIP(vm VMSpec, image DiskImage, ip IPAddress, diskSize uint) (VM, Disk, IPAddress, error)
	CreateVMWithExistingDisk(vm VMSpec, version IPVersion, disk Disk) (VM, IPAddress, Disk, error)
	CreateVMWithExistingDiskAndIP(vm VMSpec, ip IPAddress, disk Disk) (VM, IPAddress, Disk, error)

	AttachDisk(vm VM, disk Disk) (VM, Disk, error)
	DetachDisk(vm VM, disk Disk) (VM, Disk, error)
	AttachIP(vm VM, ip IPAddress) (VM, IPAddress, error)
	DetachIP(vm VM, ip IPAddress) (VM, IPAddress, error)

	// Operations on VM state
	StartVM(vm VM) error
	StopVM(vm VM) error
	RebootVM(vm VM) error
	DeleteVM(vm VM) error

	VMFromName(name string) VM
	DescribeVM(vmfilter VMFilter) ([]VM, error)
	ListVMs() ([]VM, error)

	// Updates vm memory to the value passed as parameter
	// UpdateMemoryVM(vm *VM, memory int) error
	// Updates the number of cores to the value passed as parameter
	// UpdateCoresVM(vm *VM, cores int) error
}

// VM represents a virtual machine
type VM struct {
	ID          string
	Hostname    string
	RegionID    string
	Farm        string
	Description string
	Cores       int
	Memory      int
	DateCreated time.Time
	Ips         []IPAddress
	Disks       []Disk
	SSHKeysID   []string
	State       string
}

// VMSpec gives the options available to create a VM
type VMSpec struct {
	RegionID  string
	Hostname  string
	Farm      string
	Memory    int
	Cores     int
	SSHKeysID []string
	Login     string
	Password  string
}

type VMFilter struct {
	RegionID string
	Farm     string
	Hostname string
	ID       string
	State    string
}
