package hostingv4

import (
	"errors"
	"reflect"
	"testing"
	"strconv"

	"github.com/PabloPie/Gandi-Go/mock"
	"github.com/golang/mock/gomock"
)

var (
	images4 = []diskImagev4{
		diskImagev4{ID: 1, DiskID : 10, RegionID : 123, Name : "Debian", Size : 2048},
		diskImagev4{ID: 2, DiskID : 11, RegionID : 456, Name : "Debian", Size : 2048},
		diskImagev4{ID: 4, DiskID : 44, RegionID : 123, Name : "Ubuntu", Size : 8192},
		diskImagev4{ID: 8, DiskID : 54, RegionID : 456, Name : "Ubuntu", Size : 8192},
		diskImagev4{ID: 11, DiskID : 98, RegionID : 456, Name : "MS-DOS", Size : 64},
	}
)

/* ImageByName */

func TestImageByNameEmptyRegion(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	_, err := testHosting.ImageByName("Debian",Region{})
	expected := errors.New("Region provided does not have an ID")

	if !reflect.DeepEqual(expected, err) {
		t.Errorf("Error, expected '%+v', got instead '%+v'", expected, err)
	}
}

func TestImageByNameBadRegionID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	_, err := testHosting.ImageByName("Debian",Region{ID: "badid"})
	expected := errors.New("Error parsing RegionID 'badid' from Region {badid  }")

	if !reflect.DeepEqual(expected, err) {
		t.Errorf("Error, expected '%+v', got instead '%+v'", expected, err)
	}
}

func TestImageByNameNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	mockClient.EXPECT().Send("hosting.image.list",
		gomock.Any(),
		gomock.Any()).SetArg(2, []diskImagev4{}).Return(nil)

	_, err := testHosting.ImageByName("Debian",regions[1])
	expected := errors.New("Image not found")

	if !reflect.DeepEqual(expected, err) {
		t.Errorf("Error, expected '%+v', got instead '%+v'", expected, err)
	}
}

func TestImageByNameSucceeded(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	mockClient.EXPECT().Send("hosting.image.list",
		[]interface{}{map[string]interface{}{"label": images4[0].Name, "datacenter_id": images4[0].RegionID}},
		gomock.Any()).SetArg(2, []diskImagev4{images4[0]}).Return(nil)

	image, _ := testHosting.ImageByName(images4[0].Name, Region{ID:strconv.Itoa(images4[0].RegionID)})
	expected := fromDiskImagev4(images4[0])

	if !reflect.DeepEqual(expected, image) {
		t.Errorf("Error, expected '%+v', got instead '%+v'", expected, image)
	}
}

/* ListImagesInRegion */

func TestListImagesInRegionBadRegionID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	_, err := testHosting.ListImagesInRegion(Region{ID: "badid"})
	expected := errors.New("Error parsing RegionID 'badid' from Region {badid  }")

	if !reflect.DeepEqual(expected, err) {
		t.Errorf("Error, expected '%+v', got instead '%+v'", expected, err)
	}
}

func TestListImagesInRegionNoImage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	mockClient.EXPECT().Send("hosting.image.list",
		gomock.Any(),
		gomock.Any()).SetArg(2, []diskImagev4{}).Return(nil)

	_, err := testHosting.ListImagesInRegion(regions[1])
	expected := errors.New("No images")

	if !reflect.DeepEqual(expected, err) {
		t.Errorf("Error, expected '%+v', got instead '%+v'", expected, err)
	}
}

func TestListImagesInRegionSucceeded(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)
	
	regionID := images4[1].RegionID
	var theImages []diskImagev4
	var expected []DiskImage
	for _, i := range images4 {
		if i.RegionID == regionID {
			theImages = append(theImages, i)
			expected = append(expected, fromDiskImagev4(i))
		}
	}

	mockClient.EXPECT().Send("hosting.image.list",
		[]interface{}{map[string]interface{}{"datacenter_id": regionID}},
		gomock.Any()).SetArg(2, theImages).Return(nil)

	images, _ := testHosting.ListImagesInRegion(regions[1])

	if !reflect.DeepEqual(expected, images) {
		t.Errorf("Error, expected '%+v', got instead '%+v'", expected, images)
	}
}
