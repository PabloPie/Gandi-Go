package hosting

type Operation struct {
	DiskID  int    `xmlrpc:"disk_id"`
	ID      int    `xmlrpc:"id"`
	IfaceID int    `xmlrpc:"iface_id"`
	IPID    int    `xmlrpc:"ip_id"`
	Step    string `xmlrpc:"step"`
	Type    string `xmlrpc:"type"`
	VMID    int    `xmlrpc:"vm_id"`
}
