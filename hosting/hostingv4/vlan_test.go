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
