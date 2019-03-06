Design proposal for Gandi Hosting GO driver

```
// Node represents a virtual machine
type VMInstance struct {
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

// this is our interface hosting,
// required for v4 and v5
```
type hosting interface {

// Sync or Async? Sync means a call has to wait up to 
// a minute for the creation to end, async means we can't
// abstract away "operations/events"
func createVM(VMSpec vm, []disk_spec, []ip_spec)  {}
func attach_disk(VMID int, DiskID int)            {}
func detach_disk(VMID int, DiskID int)            {}
func attach_ip(VMID int, IpID int)                {}
func detach_ip(VMID int, IpID int)                {}
func startVM(VMID int)                            {}
func stopVM(VMID int)                             {}
func rebootVM(VMID int)                           {}
func deleteVM(VMID int)                           {}
func infoVM(VMID int)                             {}
func updateVM(cores int, memory int)              {}
func listVMs()                                    {}



// Disk
func createDisk(DiskSpec)  {}
func deleteDisk()  {}
func extendDisk(int diskid, unsigned int size) {}
func renameDisk(int diskid, string name) {}

// IP
func createIP() {}
func deleteIP() {}

// Images
func listImages(regionID int) {}
func imageByName(name string, regionID int) {}

// Regions
func listRegions() {}
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
