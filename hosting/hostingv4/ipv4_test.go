package hostingv4

import (
	"errors"
	"reflect"
	"strconv"
	"testing"

	"github.com/PabloPie/go-gandi/client"

	"github.com/PabloPie/go-gandi/hosting"
	"github.com/PabloPie/go-gandi/mock"
	"github.com/golang/mock/gomock"
)

var (
	nbipv6              = 4
	nbipv4              = 4
	testRegionID        = 789
	testVersion         = hosting.IPv4
	nbTestRegionVersion = 2
	ipsv4               = []iPAddressv4{
		iPAddressv4{ID: 100, IP: "192.168.0.1", RegionID: 123, Version: 4, VM: 0, State: "created"},
		iPAddressv4{ID: 102, IP: "2001:4b98::DEAD", RegionID: 123, Version: 6, VM: 0, State: "created"},
		iPAddressv4{ID: 154, IP: "2001:4b98::BABE", RegionID: 456, Version: 6, VM: 0, State: "being_created"},
		iPAddressv4{ID: 154, IP: "2001:B00B::1337", RegionID: 456, Version: 6, VM: 0, State: "free"},
		iPAddressv4{ID: 111, IP: "192.168.0.10", RegionID: 123, Version: 4, VM: 0, State: "created"},
		iPAddressv4{ID: 115, IP: "192.168.5.20", RegionID: 789, Version: 4, VM: 0, State: "created"},
		iPAddressv4{ID: 116, IP: "192.168.44.30", RegionID: 789, Version: 4, VM: 0, State: "created"},
		iPAddressv4{ID: 117, IP: "2001:B00B::BOOB", RegionID: 789, Version: 6, VM: 0, State: "created"},
	}
)

var regions = []Region{
	Region{ID: "123", Name: "Datacentre 123", Country: "France"},
	Region{ID: "456", Name: "Datacenter 456", Country: "United Kingdom"},
	Region{ID: "789", Name: "Centro de datos 789", Country: "Espana"},
}

/* CreateIP */

func TestCreateIPv6(t *testing.T) {
	testCreateIP(t, hosting.IPv6, "fe80::DEAD:BABE:DEAD:BEEF", regions[1])
}

func TestCreateIPv4(t *testing.T) {
	testCreateIP(t, hosting.IPv4, "92.243.17.196", regions[2])
}

func testCreateIP(t *testing.T, version hosting.IPVersion, theIP string, region Region) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	myOp := Operation{ID: 1, IPID: 666}
	regionIDInt, _ := strconv.Atoi(region.ID)

	creation := mockClient.EXPECT().Send("hosting.iface.create",
		[]interface{}{map[string]interface{}{
			"datacenter_id": regionIDInt,
			"ip_version":    int(version),
			"bandwidth":     hosting.DefaultBandwidth,
		}},
		gomock.Any()).SetArg(2, myOp).Return(nil)

	wait := mockClient.EXPECT().Send("operation.info",
		[]interface{}{myOp.ID},
		gomock.Any()).SetArg(2, operationInfo{myOp.ID, "DONE"}).Return(nil).After(creation)

	ipaddressv4 := iPAddressv4{
		ID:       1337,
		IP:       theIP,
		RegionID: regionIDInt,
		Version:  int(version),
		VM:       0,
		State:    "created",
	}

	ipexpected := IPAddress{
		ID:       "1337",
		IP:       theIP,
		RegionID: region.ID,
		Version:  version,
		VM:       "0",
		State:    "created",
	}

	mockClient.EXPECT().Send("hosting.ip.info",
		[]interface{}{myOp.IPID},
		gomock.Any()).SetArg(2, ipaddressv4).Return(nil).After(wait)

	ipresult, _ := testHosting.CreateIP(region, version)

	if !reflect.DeepEqual(ipexpected, ipresult) {
		t.Errorf("Error, expected %+v, got instead %+v", ipexpected, ipresult)
	}
}

func TestCreateIPbadVersion(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	expected := errors.New("Bad IP version")
	_, err := testHosting.CreateIP(regions[0], 1234)

	if !reflect.DeepEqual(expected, err) {
		t.Errorf("Error, expected %+v, got instead %+v", expected, err)
	}
}

func TestCreateIPCreationFailed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	myOp := Operation{ID: 212, IPID: 99}
	region := regions[0]
	version := hosting.IPv4
	regionIDInt, _ := strconv.Atoi(region.ID)

	creation := mockClient.EXPECT().Send("hosting.iface.create",
		[]interface{}{map[string]interface{}{
			"datacenter_id": regionIDInt,
			"ip_version":    int(version),
			"bandwidth":     hosting.DefaultBandwidth,
		}},
		gomock.Any()).SetArg(2, myOp).Return(nil)

	wait1 := mockClient.EXPECT().Send("operation.info",
		[]interface{}{myOp.ID},
		gomock.Any()).SetArg(2, operationInfo{myOp.ID, "WAIT"}).Return(nil).After(creation)

	mockClient.EXPECT().Send("operation.info",
		[]interface{}{myOp.ID},
		gomock.Any()).SetArg(2, operationInfo{myOp.ID, "ERROR"}).Return(nil).After(wait1)

	_, err := testHosting.CreateIP(region, version)

	expected := errors.New("Bad operation status for 212 : ERROR")
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("Error, expected %+v, got instead %+v", expected, err)
	}
}

/* DeleteIP */

func TestDeleteIP(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	ip := toIPAddress(ipsv4[1])
	ipIDInt, _ := strconv.Atoi(ip.ID)
	opIPInfo := Operation{ID: 123654, IfaceID: 666}
	opWait := Operation{ID: 123777}

	infos := mockClient.EXPECT().Send("hosting.ip.info",
		[]interface{}{ipIDInt},
		gomock.Any()).SetArg(2, opIPInfo).Return(nil)

	delete := mockClient.EXPECT().Send("hosting.iface.delete",
		[]interface{}{opIPInfo.IfaceID},
		gomock.Any()).SetArg(2, opWait).Return(nil).After(infos)

	mockClient.EXPECT().Send("operation.info",
		[]interface{}{opWait.ID},
		gomock.Any()).SetArg(2, operationInfo{opWait.ID, "DONE"}).Return(nil).After(delete)

	err := testHosting.DeleteIP(ip)

	if !reflect.DeepEqual(nil, err) {
		t.Errorf("Error, expected %+v, got instead %+v", nil, err)
	}
}

/* DescribeIP */

func TestDescribeAllIP(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	filter := IPFilter{}
	ipmap, _ := ipFilterToMap(filter)

	mockClient.EXPECT().Send("hosting.ip.list",
		[]interface{}{ipmap},
		gomock.Any()).SetArg(2, ipsv4).Return(nil)

	ipsresult, _ := testHosting.DescribeIP(filter)

	var expected []IPAddress
	for _, ip := range ipsv4 {
		expected = append(expected, toIPAddress(ip))
	}

	if !reflect.DeepEqual(expected, ipsresult) {
		t.Errorf("Error, expected %+v, got instead %+v", expected, ipsresult)
	}
}

func TestDescribeAllIPv4(t *testing.T) {
	ipsresult, _ := testDescribeIPByVersion(t, hosting.IPv4)

	var expected []IPAddress
	for _, ip := range ipsv4 {
		if ip.Version == 4 {
			expected = append(expected, toIPAddress(ip))
		}
	}

	if !reflect.DeepEqual(expected, ipsresult) {
		t.Errorf("Error, expected %+v, got instead %+v", expected, ipsresult)
	}
}

func TestDescribeAllIPv6(t *testing.T) {
	ipsresult, _ := testDescribeIPByVersion(t, hosting.IPv6)

	var expected []IPAddress
	for _, ip := range ipsv4 {
		if ip.Version == 6 {
			expected = append(expected, toIPAddress(ip))
		}
	}

	if len(ipsresult) != 4 {
		t.Errorf("Error, expected 4 IPv6, got %+v !", len(ipsresult))
	}

	if !reflect.DeepEqual(expected, ipsresult) {
		t.Errorf("Error, expected %+v, got instead %+v", expected, ipsresult)
	}
}

func TestDescribeIPBadVersion(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)
	_, err := testHosting.DescribeIP(IPFilter{Version: 5})
	expected := internalParseError("IPFilter", "Version")

	if !reflect.DeepEqual(expected, err) {
		t.Errorf("Error, expected %+v, got instead %+v", expected, err)
	}
}

func testDescribeIPByVersion(t *testing.T, version hosting.IPVersion) ([]IPAddress, error) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	filter := IPFilter{Version: version}
	ipmap, _ := ipFilterToMap(filter)

	var ipsv4version []iPAddressv4
	versionInt := int(version)
	for _, ip := range ipsv4 {
		if ip.Version == versionInt {
			ipsv4version = append(ipsv4version, ip)
		}
	}

	mockClient.EXPECT().Send("hosting.ip.list",
		[]interface{}{ipmap},
		gomock.Any()).SetArg(2, ipsv4version).Return(nil)

	return testHosting.DescribeIP(filter)
}

func TestDescribeIPByIP(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	ip := ipsv4[0]
	filter := IPFilter{IP: ip.IP}
	ipmap, _ := ipFilterToMap(filter)

	mockClient.EXPECT().Send("hosting.ip.list",
		[]interface{}{ipmap},
		gomock.Any()).SetArg(2, []iPAddressv4{ip}).Return(nil)

	ipsresult, _ := testHosting.DescribeIP(filter)

	var expected []IPAddress
	for _, iip := range ipsv4 {
		if iip.IP == ip.IP {
			expected = append(expected, toIPAddress(iip))
		}
	}

	if len(ipsresult) != 1 {
		t.Errorf("Error, expected 1 IPs, got %+v !", len(ipsresult))
	}

	if !reflect.DeepEqual(expected, ipsresult) {
		t.Errorf("Error, expected %+v, got instead %+v", expected, ipsresult)
	}
}

func TestDescribeIPByID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	ip := ipsv4[1]
	idString := strconv.Itoa(ip.ID)
	filter := IPFilter{ID: idString}
	ipmap, _ := ipFilterToMap(filter)

	mockClient.EXPECT().Send("hosting.ip.list",
		[]interface{}{ipmap},
		gomock.Any()).SetArg(2, []iPAddressv4{ip}).Return(nil)

	ipsresult, _ := testHosting.DescribeIP(filter)

	var expected []IPAddress
	for _, iip := range ipsv4 {
		if iip.ID == ip.ID {
			expected = append(expected, toIPAddress(iip))
		}
	}

	if len(ipsresult) != 1 {
		t.Errorf("Error, expected 1 IPs, got %+v !", len(ipsresult))
	}

	if !reflect.DeepEqual(expected, ipsresult) {
		t.Errorf("Error, expected %+v, got instead %+v", expected, ipsresult)
	}
}

func TestDescribeIPByRegionID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	regionID := testRegionID
	regionIDString := strconv.Itoa(testRegionID)
	filter := IPFilter{RegionID: regionIDString}
	ipmap, _ := ipFilterToMap(filter)

	var ipsv4region []iPAddressv4
	var expected []IPAddress
	for _, iip := range ipsv4 {
		if iip.RegionID == regionID {
			ipsv4region = append(ipsv4region, iip)
			expected = append(expected, toIPAddress(iip))
		}
	}

	mockClient.EXPECT().Send("hosting.ip.list",
		[]interface{}{ipmap},
		gomock.Any()).SetArg(2, ipsv4region).Return(nil)

	ipsresult, _ := testHosting.DescribeIP(filter)

	if !reflect.DeepEqual(expected, ipsresult) {
		t.Errorf("Error, expected %+v, got instead %+v", expected, ipsresult)
	}
}

func TestDescribeIPByRegionIDAndVersion(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	regionID := testRegionID
	regionIDString := strconv.Itoa(testRegionID)
	version := testVersion
	versionInt := int(testVersion)

	filter := IPFilter{RegionID: regionIDString, Version: version}
	ipmap, _ := ipFilterToMap(filter)

	var ipsv4region []iPAddressv4
	var expected []IPAddress
	for _, iip := range ipsv4 {
		if iip.RegionID == regionID && iip.Version == versionInt {
			ipsv4region = append(ipsv4region, iip)
			expected = append(expected, toIPAddress(iip))
		}
	}

	mockClient.EXPECT().Send("hosting.ip.list",
		[]interface{}{ipmap},
		gomock.Any()).SetArg(2, ipsv4region).Return(nil)

	ipsresult, _ := testHosting.DescribeIP(filter)

	if len(ipsresult) != nbTestRegionVersion {
		t.Errorf("Error, expected %+v IPs, got %+v !", nbTestRegionVersion, len(ipsresult))
	}

	if !reflect.DeepEqual(expected, ipsresult) {
		t.Errorf("Error, expected %+v, got instead %+v", expected, ipsresult)
	}
}

func TestDeleteIPBadID(t *testing.T) {
	cl, err := client.NewClientv4("", "1234")
	testHosting := Newv4Hosting(cl)

	ip := IPAddress{
		ID: "ThisisnotAnID",
	}
	err = testHosting.DeleteIP(ip)
	if err == nil {
		t.Errorf("Error, expected error when parsing ID")
	}
}

func TestCreateIPBadRegionID(t *testing.T) {
	cl, err := client.NewClientv4("", "1234")
	testHosting := Newv4Hosting(cl)

	region := Region{
		ID: "ThisisnotAnID",
	}
	_, err = testHosting.CreateIP(region, hosting.IPVersion(4))
	if err == nil {
		t.Errorf("Error, expected error when parsing ID")
	}
}

func TestFilterBadID(t *testing.T) {
	cl, err := client.NewClientv4("", "1234")
	testHosting := Newv4Hosting(cl)

	filter := IPFilter{
		ID: "ThisisnotAnID",
	}
	_, err = testHosting.DescribeIP(filter)
	if err == nil {
		t.Errorf("Error, expected error when parsing ID")
	}
}

func TestFilterBadRegionID(t *testing.T) {
	cl, err := client.NewClientv4("", "1234")
	testHosting := Newv4Hosting(cl)

	filter := IPFilter{
		RegionID: "ThisisnotAnID",
	}
	_, err = testHosting.DescribeIP(filter)
	if err == nil {
		t.Errorf("Error, expected error when parsing ID")
	}
}
