package hosting

import "time"

// VMManager is an interface embedded in the Hosting interface, not
// implemented by VM!
type VMManager interface {
	// VM creation operations
	CreateVM(vm VMSpec, image DiskImage, version IPVersion, diskSize uint) (VM, IPAddress, Disk, error)
	CreateVMWithExistingIP(vm VMSpec, image DiskImage, ip IPAddress, diskSize uint) (VM, IPAddress, Disk, error)
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

	VMFromName(name string) (VM, error)
	DescribeVM(vmfilter VMFilter) ([]VM, error)
	ListVMs() ([]VM, error)

	// VM update operations
	UpdateVMMemory(vm VM, memory int) (VM, error)
	UpdateVMCores(vm VM, cores int) (VM, error)
	RenameVM(vm VM, newname string) (VM, error)
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

// VMFilter contains filtering options for VMs
type VMFilter struct {
	RegionID string
	Farm     string
	Hostname string
	ID       string
	State    string
}
