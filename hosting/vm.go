package hosting

import "time"

// VMManager is an interface embedded in the Hosting interface, not
// implemented by VM!
type VMManager interface {
	//Creates vm with a new disk of size `size` based on diskimage vm.image
	CreateVMD(vm VMSpec, size int) (VM, Disk, IPAddress, error)

	//Creates vm using an already existing bootable disk as system disk
	CreateVM(vm VMSpec, disk *Disk) (VM, IPAddress, error)

	AttachDisk(vm *VM, disk *Disk) error
	DetachDisk(vm *VM, disk *Disk) error
	AttachIP(vm *VM, ip *IPAddress) error
	DetachIP(vm *VM, ip *IPAddress) error

	// Operations on VM state
	StartVM(vm *VM) error
	StopVM(vm *VM) error
	RebootVM(vm *VM) error
	DeleteVM(vm *VM) error

	// filter function? currently this function is of no use
	InfoVM(vmid int) (VM, error)
	// Updates vm memory to the value passed as parameter
	UpdateMemoryVM(vm *VM, memory int) error

	// Updates the number of cores to the value passed as parameter
	UpdateCoresVM(vm *VM, cores int) error

	ListVMs() ([]VM, error)
}

// VM represents a virtual machine, in any version of the API
type VM struct {
	ID           int
	Hostname     string
	DatacenterID int
	Farm         string
	Description  string
	Cores        int
	Memory       int
	DateCreated  time.Time
	Ips          []IPAddress
	Disks        []Disk
	SSHKeys      []string
	State        string
}

// Export struct or provide function to generate struct?
// VMSpec is used in v4 of the API for vm creation
type VMSpec struct {
	RegionID  int
	Hostname  string
	Farm      string
	Memory    int
	Cores     int
	IPVersion int
	Image     DiskImage
	SSHKey    string
	Login     string
	Password  string
}
