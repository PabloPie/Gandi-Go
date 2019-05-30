package hosting

// SSHKeyManager represents a service capable of manipulating
// SSH Keys in Gandi's platform
type SSHKeyManager interface {

	// CreateKey creates an SSH Key to use when creating VMs
	//
	// Returns an object containing the fingerprint and the ID
	CreateKey(name string, value string) (SSHKey, error)

	// DeleteKey deletes the SSHKey given
	DeleteKey(key SSHKey) error

	// KeyFromName return the SSHKey with name `name`,
	// if no key with such name exists, an empty
	// object is returned instead
	KeyFromName(name string) SSHKey

	// ListKeys lists every SSHKey created
	ListKeys() []SSHKey
}

// SSHKey represents an ssh key
type SSHKey struct {
	Fingerprint string
	ID          string
	Name        string
	Value       string
}
