package hostingv4

import (
	"errors"
	"reflect"
	"testing"

	"github.com/PabloPie/go-gandi/mock"
	"github.com/golang/mock/gomock"
)

func TestWaitForOpSucceededFirstCall(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Hostingv4{mockClient}

	myOp := Operation{ID: 1337}

	mockClient.EXPECT().Send("operation.info",
		[]interface{}{myOp.ID},
		gomock.Any()).SetArg(2, operationInfo{myOp.ID, "DONE"}).Return(nil)

	err := testHosting.waitForOp(myOp)

	if !reflect.DeepEqual(nil, err) {
		t.Errorf("Error, expected '%+v', got instead '%+v'", nil, err)
	}
}

func TestWaitForOpSucceeded2Calls(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Hostingv4{mockClient}

	myOp := Operation{ID: 1337, Step: "NULL"}

	call1 := mockClient.EXPECT().Send("operation.info",
		[]interface{}{myOp.ID},
		gomock.Any()).SetArg(2, operationInfo{myOp.ID, "WAIT"}).Return(nil)

	mockClient.EXPECT().Send("operation.info",
		[]interface{}{myOp.ID},
		gomock.Any()).SetArg(2, operationInfo{myOp.ID, "DONE"}).Return(nil).After(call1)

	err := testHosting.waitForOp(myOp)

	if !reflect.DeepEqual(nil, err) {
		t.Errorf("Error, expected '%+v', got instead '%+v'", nil, err)
	}
}

func TestWaitForOpFailedFirstCall(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Hostingv4{mockClient}

	myOp := Operation{ID: 1337}

	expected := errors.New("Cannot access Operation '1337'")

	mockClient.EXPECT().Send("operation.info", []interface{}{myOp.ID},
		gomock.Any()).Return(expected)

	err := testHosting.waitForOp(myOp)

	if !reflect.DeepEqual(expected, err) {
		t.Errorf("Error, expected '%+v', got instead '%+v'", expected, err)
	}
}

func TestWaitForOpFailed2Calls(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Hostingv4{mockClient}

	myOp := Operation{ID: 1337}

	call1 := mockClient.EXPECT().Send("operation.info",
		[]interface{}{myOp.ID},
		gomock.Any()).SetArg(2, operationInfo{myOp.ID, "BILL"}).Return(nil)

	mockClient.EXPECT().Send("operation.info",
		[]interface{}{myOp.ID},
		gomock.Any()).SetArg(2, operationInfo{myOp.ID, "ERROR"}).Return(nil).After(call1)

	err := testHosting.waitForOp(myOp)
	expected := errors.New("Bad operation status for 1337 : ERROR")

	if !reflect.DeepEqual(expected, err) {
		t.Errorf("Error, expected '%+v', got instead '%+v'", expected, err)
	}
}

func TestWaitForOpFailed2CallsSendError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Hostingv4{mockClient}

	myOp := Operation{ID: 1337}

	call1 := mockClient.EXPECT().Send("operation.info",
		[]interface{}{myOp.ID},
		gomock.Any()).SetArg(2, operationInfo{myOp.ID, "BILL"}).Return(nil)

	expected := errors.New("Cannot access Operation '1337'")

	mockClient.EXPECT().Send("operation.info",
		[]interface{}{myOp.ID},
		gomock.Any()).SetArg(2, operationInfo{myOp.ID, "WAIT"}).Return(expected).After(call1)

	err := testHosting.waitForOp(myOp)

	if !reflect.DeepEqual(expected, err) {
		t.Errorf("Error, expected '%+v', got instead '%+v'", expected, err)
	}
}
