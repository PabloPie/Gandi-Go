Design proposal for Gandi Hosting GO driver

```
// Node represents a virtual machine
type VM struct {
	ID           int
	Hostname     string
	DatacenterID int
	Farm         string
	Description  string
	Cores        int
	Memory       int
	DateCreated  dateTime.iso8601
	Ips          []IPAddress
	Disks        []Disk
	SSHKeys      []string
	State        string
}
```
// vm_spec, disk_spec and ip_spec
// must be common between v4 and v5
```
// VMDescription, used for vm creation
type VMSpec struct {
	RegionID     int
	Hostname     string
	Farm         string
	Memory       int
	Cores        int
	IPVersion    int
	Image        id
	SSHKey       string
	Login        string
	Password     string
}
```

```
type Region struct {
	ID      int
	Name    string
	State   int
}
```

```
type DiskImage struct {
	ID        int
	DiskID    int
	RegionID  int
	Name      string
	Os        string
	Size      int
	State     int
}
```

```
type Disk struct {
	ID           int
	Name         string
	Size         int
	RegionID     int
	State        string
	Type         string
	Vm           int[]
	BootDisk     bool
}
```

```
type DiskSpec struct {
	RegionID   int
	Name       string
	Size       int
}
```

```
type IpAddress struct {
	ID        int
	IP        string
	RegionID  int
	Version   int
	VMID      int
	State     string
	Vlan      int // if private
}
```

```
type IpAddressSpec struct {
	RegionID  int
	Bandwidth int
	Version   int
	Vlan      int
	IP        int // if vlan defined
}
```

```
type Vlan struct {
	ID       int
	Name     String
	Gateway  String
	Subnet   String
	RegionID int
}
```

```
type VlanSpec struct {
	Name     String
	Gateway  String
	Subnet   String
	RegionID int
}
```

## Hosting Interface

This is our interface hosting, with all the functions that are going to be implemented, required for hostingv4 and hostingv5

```
type hosting interface {

//Creates vm with a new disk of size `size` based on diskimage vm.image
func createVMD(VMSpec vm, int size) VM, Disk, IPAddress             {}
//Creates vm using an already existing bootable disk as system disk
func createVM(VMSpec vm, diskid int) VM, IPAddress                  {}
func attachDisk(vmid int, diskid int) err                           {}
func detachDisk(vmid int, diskid int) err                           {}
func attachIP(vmid int, ipid int) err                               {}
func detachIP(vmid int, ipid int) err                               {}
func startVM(vmid int) err                                          {}
func stopVM(vmid int) err                                           {}
func rebootVM(vmid int) err                                         {}
func deleteVM(vmid int) err                                         {}
func infoVM(vmid int) VMInstance                                    {}
// Updates vm memory to the value passed as parameter
func updateMemoryVM(memory int) VMInstance                          {}
// Updates the number of cores to the value passed as parameter
func updateCoresVM(cores int) VMInstance                            {}
func listVMs() VMInstance[]                                         {}

// Disk
func createDisk(DiskSpec disk) Disk                   {}
func infoDisk(int diskid) Disk                        {}
func listDisk() Disk[]                                {}
func deleteDisk(int diskid) err                       {}
// Extends `diskid` size by `size` (original size + `size`)
func extendDisk(int diskid, unsigned int size) err    {}
func renameDisk(int diskid, string name) err          {}

// IP
func createIP(int version) IPAddress {}
func infoIP(int ipid) IPAddress      {}
func listIP() IPAddress[]            {}
func deleteIP(int ipid) err          {}

// Images
func listImages(regionid int) DiskImages[]      {}
func imageByName(name string, regionid int) int {}

// SSH
func createKey(string name, string value) SSHKey      {}
func deleteKey(string name) err                       {}
// Get the key with name `name`
func keyfromName(string name) SSHKey                  {}
func listKeys() SSHKeys[]                             {}

// Vlan
func createVlan(VlanSpec) Vlan {}
func infoVlan(int vlanid) Vlan {}
func listVlan() Vlan[]         {}
func updateVlanGW() err        {}
func renameVlan() err          {}

// Regions
func listRegions() Region[] {}
}
```



## Disk image
	- Image identifiée par un int dans l'APIv4
	- une image = sur un DC
	- utilite de list images? normalement l'utilisateur devrait avoir une liste des images dispos sans devoir faire un appel sur l'api
## Disk
	- Est-ce vraiment utile de proposer un countDisk() ?
	- Il y a une histoire de noyaux dispo dans tel ou tel datacenter dans l'APIv4
	- Migration d'un disque d'un DC à l'autre ? (en bonus ?)


# Travail
	- 1ere partie IP+DISK+VM -> Tests
	- 2eme partie SSH+VLAN -> Tests
	- 3eme partie Terraform -> Tests
	- 4eme partie Doc 
