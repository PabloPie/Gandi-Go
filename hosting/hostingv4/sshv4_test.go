package hostingv4

import (
	"testing"

	"github.com/PabloPie/Gandi-Go/mock"
	"github.com/golang/mock/gomock"
)

var (
	keyid       = 1
	keyidstr    = "1"
	keyname     = "key1"
	keyvalue    = "ssh-rsa 12345 test@hosting"
	fingerprint = "11:22:33:44"
)

func TestCreateSSHKey(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	paramsCreateKey := []interface{}{
		map[string]string{
			"name":  keyname,
			"value": keyvalue,
		}}
	responseCreateKey := sshkeyv4{
		ID:          1,
		Name:        keyname,
		Fingerprint: fingerprint,
	}
	create := mockClient.EXPECT().Send("hosting.ssh.create",
		paramsCreateKey, gomock.Any()).SetArg(2, responseCreateKey).Return(nil)

	paramsKeyInfo := []interface{}{keyid}
	responseCreateKey.Value = keyvalue
	mockClient.EXPECT().Send("hosting.ssh.info",
		paramsKeyInfo, gomock.Any()).SetArg(2, responseCreateKey).Return(nil).After(create)

	expectedSSHKey := SSHKey{
		ID:          keyidstr,
		Name:        keyname,
		Fingerprint: fingerprint,
		Value:       keyvalue,
	}
	key, _ := testHosting.CreateKey(keyname, keyvalue)
	if key.Name != expectedSSHKey.Name {
		t.Errorf("Error, expected Key Name to be %s, got instead %s", key.Name, expectedSSHKey.Name)
	}
	if key.Value != expectedSSHKey.Value {
		t.Errorf("Error, expected Key Value to be %s, got instead %s", key.Value, expectedSSHKey.Value)
	}
}

func TestDeleteSSHKey(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	paramsDeleteKey := []interface{}{keyid}

	mockClient.EXPECT().Send("hosting.ssh.delete",
		paramsDeleteKey, gomock.Any()).SetArg(2, true).Return(nil)

	key := SSHKey{
		ID: keyidstr,
	}
	err := testHosting.DeleteKey(key)
	if err != nil {
		t.Errorf("Error, expected no errors when deleting key")
	}

}
func TestDeleteSSHKeyBadID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	key := SSHKey{
		ID: "thisisnotanint",
	}
	err := testHosting.DeleteKey(key)
	if err == nil {
		t.Errorf("Error, ID given was not an int, expected error")
	}
}

func TestKeyFromName(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	paramsListKey := []interface{}{
		map[string]string{
			"name": keyname,
		}}
	responseListKeys := []sshkeyv4{{
		ID:          keyid,
		Fingerprint: fingerprint,
		Name:        keyname,
	}}
	list := mockClient.EXPECT().Send("hosting.ssh.list",
		paramsListKey, gomock.Any()).SetArg(2, responseListKeys).Return(nil)

	paramsKeyInfo := []interface{}{keyid}
	responseListKeys[0].Value = keyvalue
	mockClient.EXPECT().Send("hosting.ssh.info",
		paramsKeyInfo, gomock.Any()).SetArg(2, responseListKeys[0]).Return(nil).After(list)

	expectedKey := SSHKey{
		ID:          keyidstr,
		Name:        keyname,
		Fingerprint: fingerprint,
		Value:       keyvalue,
	}
	key := testHosting.KeyfromName(keyname)

	if key != expectedKey {
		t.Errorf("Error, expected Key %+v, got instead %+v", expectedKey, key)
	}
}

func TestListKeys(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockV4Caller(mockCtrl)
	testHosting := Newv4Hosting(mockClient)

	paramsListKey := []interface{}{}
	responseListKeys := []sshkeyv4{{
		ID:          keyid,
		Fingerprint: fingerprint,
		Name:        keyname,
	}}
	mockClient.EXPECT().Send("hosting.ssh.list",
		paramsListKey, gomock.Any()).SetArg(2, responseListKeys).Return(nil)

	keys := testHosting.ListKeys()

	if len(keys) < 1 {
		t.Errorf("Error, expected at least a key, got %d instead", len(keys))
	}
}
