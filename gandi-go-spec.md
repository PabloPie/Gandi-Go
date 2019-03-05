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

//
func createDisk() {}

//
func createIP() {}
}
```
