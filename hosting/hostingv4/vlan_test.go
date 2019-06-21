package hostingv4

import (
	"reflect"
	"testing"

	"github.com/PabloPie/go-gandi/hosting"
	"github.com/PabloPie/go-gandi/mock"
	"github.com/golang/mock/gomock"
)

func TestVlanCreation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	paramsVlanCreate := []interface{}{map[string]interface{}{
		"name":          "testvlan",
		"datacenter_id": region,
	}}
	responseVlanCreate := Operation{
		ID: 1,
	}
	creation := mockClient.EXPECT().Send("hosting.vlan.create",
		paramsVlanCreate, gomock.Any()).SetArg(2, responseVlanCreate).Return(nil)

	paramsWait := []interface{}{responseVlanCreate.ID}
	responseWait := operationInfo{responseVlanCreate.ID, "DONE"}
	wait := mockClient.EXPECT().Send("operation.info",
		paramsWait, gomock.Any()).SetArg(2, responseWait).Return(nil).After(creation)

	paramsVlanList := []interface{}{map[string]interface{}{
		"name": "testvlan",
	}}
	responseVlanList := []vlanv4{{ID: 2, Name: "testvlan", Gateway: "", Subnet: "192.168.0.0/24", RegionID: region}}
	mockClient.EXPECT().Send("hosting.vlan.list",
		paramsVlanList, gomock.Any()).SetArg(2, responseVlanList).Return(nil).After(wait)

	vlanspec := hosting.VlanSpec{Name: "testvlan", RegionID: regionstr}
	vlan, _ := testHosting.CreateVlan(vlanspec)

	expected := hosting.Vlan{
		ID:       "2",
		Name:     "testvlan",
		Subnet:   "192.168.0.0/24",
		RegionID: regionstr,
	}

	if !reflect.DeepEqual(vlan, expected) {
		t.Errorf("Error, expected %+v, got instead %+v", expected, vlan)
	}
}

func TestVlanDelete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	paramsVlanDelete := []interface{}{2}
	responseVlanDelete := Operation{
		ID: 1,
	}
	creation := mockClient.EXPECT().Send("hosting.vlan.delete",
		paramsVlanDelete, gomock.Any()).SetArg(2, responseVlanDelete).Return(nil)

	paramsWait := []interface{}{responseVlanDelete.ID}
	responseWait := operationInfo{responseVlanDelete.ID, "DONE"}
	mockClient.EXPECT().Send("operation.info",
		paramsWait, gomock.Any()).SetArg(2, responseWait).Return(nil).After(creation)

	vlan := hosting.Vlan{ID: "2"}
	err := testHosting.DeleteVlan(vlan)

	if err != nil {
		t.Errorf("Error, expected vlan to be deleted, got error '%v' instead", err)
	}
}

func TestVlanUpdateGW(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	vlanupdate := map[string]string{"gateway": "192.168.1.200"}
	paramsVlanUpdate := []interface{}{1, vlanupdate}
	responseVlanUpdate := Operation{
		ID: 1,
	}
	update := mockClient.EXPECT().Send("hosting.vlan.update",
		paramsVlanUpdate, gomock.Any()).SetArg(2, responseVlanUpdate).Return(nil)

	paramsWait := []interface{}{responseVlanUpdate.ID}
	responseWait := operationInfo{responseVlanUpdate.ID, "DONE"}
	wait := mockClient.EXPECT().Send("operation.info",
		paramsWait, gomock.Any()).SetArg(2, responseWait).Return(nil).After(update)

	paramsVlanList := []interface{}{map[string]interface{}{"id": []int{1}}}
	responseVlanList := []vlanv4{{ID: 1, RegionID: 5, Gateway: "192.168.1.200",
		Name: "testvlan", Subnet: "192.168.1.0/24"}}
	mockClient.EXPECT().Send("hosting.vlan.list",
		paramsVlanList, gomock.Any()).SetArg(2, responseVlanList).Return(nil).After(wait)

	expectedVlan := hosting.Vlan{ID: "1", RegionID: "5", Gateway: "192.168.1.200",
		Name: "testvlan", Subnet: "192.168.1.0/24"}

	oldvlan := hosting.Vlan{ID: "1", Gateway: "192.168.1.1"}
	vlan, _ := testHosting.UpdateVlanGW(oldvlan, "192.168.1.200")

	if !reflect.DeepEqual(vlan, expectedVlan) {
		t.Errorf("Error, expected %+v, got instead %+v", expectedVlan, vlan)
	}
}

func TestRenameVlan(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	vlanupdate := map[string]string{"name": "newvlanname"}
	paramsVlanUpdate := []interface{}{1, vlanupdate}
	responseVlanUpdate := Operation{
		ID: 1,
	}
	update := mockClient.EXPECT().Send("hosting.vlan.update",
		paramsVlanUpdate, gomock.Any()).SetArg(2, responseVlanUpdate).Return(nil)

	paramsWait := []interface{}{responseVlanUpdate.ID}
	responseWait := operationInfo{responseVlanUpdate.ID, "DONE"}
	wait := mockClient.EXPECT().Send("operation.info",
		paramsWait, gomock.Any()).SetArg(2, responseWait).Return(nil).After(update)

	paramsVlanList := []interface{}{map[string]interface{}{"name": "newvlanname"}}
	responseVlanList := []vlanv4{{ID: 1, RegionID: 5, Gateway: "192.168.1.1",
		Name: "newvlanname", Subnet: "192.168.1.0/24"}}
	mockClient.EXPECT().Send("hosting.vlan.list",
		paramsVlanList, gomock.Any()).SetArg(2, responseVlanList).Return(nil).After(wait)

	expectedVlan := hosting.Vlan{ID: "1", RegionID: "5", Gateway: "192.168.1.1",
		Name: "newvlanname", Subnet: "192.168.1.0/24"}

	oldvlan := hosting.Vlan{ID: "1", Name: "oldvlanname"}
	vlan, _ := testHosting.RenameVlan(oldvlan, "newvlanname")

	if !reflect.DeepEqual(vlan, expectedVlan) {
		t.Errorf("Error, expected %+v, got instead %+v", expectedVlan, vlan)
	}
}
