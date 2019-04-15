package hosting

type SSHKeyManager interface {
	CreateKey(name string, value string) (SSHKey, error)
	DeleteKey(key SSHKey) error
	// Get the key with name `name`
	KeyfromName(name string) SSHKey
	ListKeys() []SSHKey
}

type SSHKey struct {
	Fingerprint string
	ID          string
	Name        string
	Value       string
}
