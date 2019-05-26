package hosting

// SSHKeyManager represents a service capable of manipulating
// SSH Keys in Gandi's platform
type SSHKeyManager interface {
	CreateKey(name string, value string) (SSHKey, error)
	DeleteKey(key SSHKey) error
	KeyFromName(name string) SSHKey
	ListKeys() []SSHKey
}

// SSHKey represents an ssh key
type SSHKey struct {
	Fingerprint string
	ID          string
	Name        string
	Value       string
}
