package capabilities

// Permission describes a runtime permission required by a capability.
type Permission string

const (
	PermissionRead       Permission = "read"
	PermissionWrite      Permission = "write"
	PermissionExecute    Permission = "execute"
	PermissionNetwork    Permission = "network"
	PermissionFilesystem Permission = "filesystem"
	PermissionSecrets    Permission = "secrets"
)

// PermissionSet is a small helper for matching capabilities to policies.
type PermissionSet map[Permission]struct{}

func NewPermissionSet(values ...Permission) PermissionSet {
	set := make(PermissionSet, len(values))
	for _, value := range values {
		set[value] = struct{}{}
	}
	return set
}

func (s PermissionSet) Has(value Permission) bool {
	_, ok := s[value]
	return ok
}

func (s PermissionSet) Allows(required ...Permission) bool {
	if len(required) == 0 {
		return true
	}
	for _, value := range required {
		if !s.Has(value) {
			return false
		}
	}
	return true
}

func (s PermissionSet) AllowsSet(required PermissionSet) bool {
	if len(required) == 0 {
		return true
	}
	for value := range required {
		if !s.Has(value) {
			return false
		}
	}
	return true
}
