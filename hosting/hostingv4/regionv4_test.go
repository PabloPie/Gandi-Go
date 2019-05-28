package hostingv4

import (
	"errors"
	"reflect"
	"testing"

	"github.com/PabloPie/go-gandi/hosting"
	"github.com/PabloPie/go-gandi/mock"
	"github.com/golang/mock/gomock"
)

var regionsv4 = []regionv4{
	{ID: 12, Name: "FR-SD3", Country: "France"},
	{ID: 34, Name: "FR-SD6", Country: "France"},
	{ID: 56, Name: "EN-DC1", Country: "United Kingdom"},
	{ID: 78, Name: "ES-CD5", Country: "Espana"},
	{ID: 90, Name: "ES-CD4", Country: "Espana"},
}

func TestListRegions(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	mockClient.EXPECT().Send("hosting.datacenter.list",
		[]interface{}{},
		gomock.Any()).SetArg(2, regionsv4).Return(nil)

	var expected []hosting.Region
	for _, r := range regionsv4 {
		expected = append(expected, fromRegionv4(r))
	}

	regionsresult, _ := testHosting.ListRegions()

	if !reflect.DeepEqual(expected, regionsresult) {
		t.Errorf("Error, expected %+v, got instead %+v", expected, regionsresult)
	}
}

func TestAllRegionsByCode(t *testing.T) {
	for _, rv4 := range regionsv4 {
		region, _ := testRegionByCode(t, rv4.Name)
		regionexpected := fromRegionv4(rv4)

		if !reflect.DeepEqual(regionexpected, region) {
			t.Errorf("Error, expected %+v, got instead %+v", regionexpected, region)
		}
	}
}

func TestRegionByBadCode(t *testing.T) {
	r, err := testRegionByCode(t, "BAD-DC")
	expected := errors.New("hosting.Region not found")

	if !reflect.DeepEqual(expected, err) {
		t.Errorf("Error, expected %+v, got instead %+v (%+v)", expected, err, r)
	}
}

func testRegionByCode(t *testing.T, code string) (hosting.Region, error) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	var theRegion []regionv4
	for _, r := range regionsv4 {
		if r.Name == code {
			theRegion = append(theRegion, r)
			break
		}
	}

	mockClient.EXPECT().Send("hosting.datacenter.list",
		[]interface{}{map[string]string{"dc_code": code}},
		gomock.Any()).SetArg(2, theRegion).Return(nil)

	return testHosting.RegionbyCode(code)
}
