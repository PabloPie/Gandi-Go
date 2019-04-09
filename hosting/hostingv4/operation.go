package hostingv4

import (
	"fmt"
	"time"
)

type operationInfo struct {
	ID     int    `xmlrpc:"id"`
	Status string `xmlrpc:"step"`
}

// Operation is an operation in Gandi v4 API
type Operation struct {
	ID      int    `xmlrpc:"id"`
	VMID    int    `xmlrpc:"vm_id"`
	DiskID  int    `xmlrpc:"disk_id"`
	IfaceID int    `xmlrpc:"iface_id"`
	IPID    int    `xmlrpc:"ip_id"`
	Step    string `xmlrpc:"step"`
	Type    string `xmlrpc:"type"`
}

func (h Hostingv4) waitForOp(op Operation) error {
	res := operationInfo{}
	params := []interface{}{op.ID}
	err := h.Send("operation.info", params, &res)
	if err != nil {
		return err
	}
	for res.Status != "DONE" {
		time.Sleep(2 * time.Second)
		err := h.Send("operation.info", params, &res)
		if err != nil {
			return err
		}
		if res.Status == "DONE" {
			return nil
		}
		if res.Status != "BILL" && res.Status != "WAIT" && res.Status != "RUN" {
			return fmt.Errorf("Bad operation status for %d : %s", op.ID, res.Status)
		}
	}
	return nil
}
