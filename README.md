# Gandi-Go Hosting library

[![Actions Status](https://wdp9fww0r9.execute-api.us-west-2.amazonaws.com/production/badge/PabloPie/go-gandi)](https://wdp9fww0r9.execute-api.us-west-2.amazonaws.com/production/results/PabloPie/go-gandi) [![codecov](https://codecov.io/gh/PabloPie/go-gandi/branch/master/graph/badge.svg)](https://codecov.io/gh/PabloPie/go-gandi)

Go library to interact with [Gandi](https://www.gandi.net/en)'s Hosting API and manage Virtual Machines, Disks, IPs and Vlans. Currently working on implementing the driver for the [XMLRPC API](https://doc.rpc.gandi.net/overview.html) while the development of the new API, [V5](https://docs.gandi.net/en/cloud/index.html), finishes. The API this library provides will not change once the driver for V5 is implemented.

## Usage Example

This example shows how to create an IP, a Disk and a VM. To use the library you need to [get your API Key](https://v4.gandi.net/admin/api_key).

```go
// This example uses the current api v4

import (
	"github.com/PabloPie/go-gandi/client"
	"github.com/PabloPie/go-gandi/hosting"
	"github.com/PabloPie/go-gandi/hosting/hostingv4"
)

apikey := "MYAPIKEY"
// Using an empty URL sets it to the default URL for v4
c, _ := client.NewClientv4("", apikey)
// We create a v4 Hosting driver with our client
h := hostingv4.Newv4Hosting(c)

region, _ := h.RegionbyCode("FR-SD6")
// We search for image Debian 9 in Region FR-SD6
image, _ := h.ImageByName("Debian 9", region)

// We create an IPv4 in Region FR-SD6
ip, _ := h.CreateIP(region, 4)

// We define a diskspec for a disk in region FR-SD6
// named Disk1 with size 15GB
diskspec := hosting.DiskSpec{
    RegionID: region.ID,
    Name:     "Disk1",
    Size:     15,
}
// We create our disk from an Image Debian9
disk, _ := h.CreateDiskFromImage(diskspec, image)

// We create a new SSH Key for our VM
key, _ := h.CreateKey("key1", "<publicsshkey>")

// We define our vmspec to create a VM in region FR-SD6
vmspec := hosting.VMSpec{
    RegionID:  region.ID,
    Hostname:  "VM1",
    Memory:    1024,
    Cores:     2,
    SSHKeysID: []string{key.Name},
}

// We create our VM, given that our IPAddress and Disk already existed
// we use the function CreateVMWithExistingDiskAndIP (yeah, naming is hard)
vm, ip, disk, _ := h.CreateVMWithExistingDiskAndIP(vmspec, ip, disk)

```

