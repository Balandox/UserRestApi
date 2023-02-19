package enum

import "errors"

type Role string

const (
	GUEST   Role = "GUEST"
	USER    Role = "USER"
	MANAGER Role = "MANAGER"
	ADMIN   Role = "ADMIN"
)

var (
	RoleIsNotValid = errors.New("role is not valid")
)

func (s Role) IsValid() bool {
	switch s {
	case GUEST, USER, MANAGER, ADMIN:
		return true
	}
	return false
}
