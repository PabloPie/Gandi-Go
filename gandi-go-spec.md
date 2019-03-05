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
```
// VMDescription, used for vm creation
type VMSpec struct {
	DatacenterID int
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
// VMListOptions, used for filtering in VMList()
type VMFilter struct {
	DatacenterID   int[]
	Farm           string[]
	Hostname       string[]
	State          string
	ID             int[]
}
```

```
type Region struct {
	id      int
	name    string
	state   int
}
```

```
type DiskImage struct {
	id        int
	disk_id   int
	region_id int
	name      string
	os        string
	size      int
	state     int
}
```


```
type Disk struct {
	id           int
	name         string
	size         int
	datacenter   int
	state        string
	type         string
	vm           int[]
	is_bootable  bool
//	is_migrating bool
//	can_snapshot bool
}
```

```
type DiskSpec struct {

}
```

```
type DiskFilter struct {

}
```

```
type IpAddress struct {
	ID        int
	IP        string
	RegionID int
	Version   int
	VMID     int
	State     string
}
```

```
type IpAdressSpec struct {

}
```

```
type hosting interface {

// vm_spec, disk_spec and ip_spec
// must be common between v4 and v5
// this is our interface hosting,
// "implemented" by v4 and v5

func listVM() VMInstance                        {}

// Allow multiple interfaces | disks?
// Sync or Async? Sync means a call has to wait up to 
// a minute for the creation to end, async means we can't
// abstract away "operations/events"
func createVM(vm_spec) VMInstance               {}

// Another vm creation option
func createVMs(VMSpec vm, []disk_spec, []ip_spec) {}
func attach_disk(VMID int, DiskID int)            {}
func detach_disk(VMID int, DiskID int)            {}
func attach_ip(VMID int, IpID int)                {}
func detach_ip(VMID int, IpID int)                {}
func startVM(VMID int)                            {}
func stopVM(VMID int)                             {}
func rebootVM(VMID int)                           {}
func deleteVM(VMID int)                           {}
func infoVM(VMID int)                             {}
func update()                                     {}
func migrateVM(VMID int, DCID int)                {}
func listVMs(VMFilter filter)                     {}



// Disk
func createDisk()  {}
func deleteDisk()  {}
func upgradeDisk() {}

// IP
func createIP()  {}
func deleteIP() {}

// Images
func listImages(regionID int) {}

// Regions
func listRegions() {}
}
```



## Disk image
	- Image identifiée par un int dans l'APIv4
	- une image = sur un DC
## Disk
	- Est-ce vraiment utile de proposer un countDisk() ?
	- Il y a une histoire de noyaux dispo dans tel ou tel datacenter dans l'APIv4
	- Migration d'un disque d'un DC à l'autre ? (en bonus ?)
