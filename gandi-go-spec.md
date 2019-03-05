Design proposal for Gandi Hosting GO driver

```
// Node represents a virtual machine
type VMInstance struct {
	ID           int
	hostname     string
	datacenterID int
	farm         string
	description  string
	cores        int
	memory       int
	dateCreated  dateTime.iso8601
	dateUpdated  dateTime.iso8601
	ips          []int
	disks        []int
	disksID      []int
	sshKeys      []string
	state        string
}
```
```
// VMCreateRequest
type VMSpec struct {
	DcID       int
	Hostname   string
	Memory     int
	Cores      int
	IPVersion  int
	SSHKey     string
}
```

```
type RegionInstance struct {
	id					int
	name				string
	state				int
}
```

```
type ImageInstance struct {
	id					int
	disk_id				int
	ragion_id			int
	name				string
	os					string
	size				int
	state				int
}
```


```
type DiskInstance struct {
	id					int
	name				string
	size				int
	datacenter			int
	state				string
	type				string
	vm					int[]
	is_bootable			bool
//	is_migrating		bool
	can_snapshot		bool
}
```

```
type IpInstance struct {
	id					int
	ip					string
	region_id			int
	version				int
	vm_id				int
	state				string
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
func createVMs(vm_spec, []disk_spec, []ip_spec) {}
func attach_disk()                              {}
func detach_disk()                              {}
func attach_ip()                                {}
func detach_ip()                                {}
func startVM()                                  {}
func stopVM()                                   {}
func rebootVM()                                 {}
func deleteVM()                                 {}
func infoVM()                                   {}
func update()                                   {}
func migrateVM()                                {}



// Disk
func createDisk()								{}
func deleteDisk()								{}
func upgradeDisk()								{}

// IP
func createIP()									{}
func destroyIP()								{}

// Images
func listImage(region)								{}

// Regions
func listRegions()								{}
}
```



## Disk image
	- Image identifiée par un int dans l'APIv4
	- une image = sur un DC
## Disk
	- Est-ce vraiment utile de proposer un countDisk() ?
	- Il y a une histoire de noyaux dispo dans tel ou tel datacenter dans l'APIv4
	- Migration d'un disque d'un DC à l'autre ? (en bonus ?)
