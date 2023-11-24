package domain

import "fmt"

type RBACUser string
type RBACRole string
type RBACObject string
type RBACAction string

func NewUserObject(appUserID AppUserID) RBACUser {
	return RBACUser(fmt.Sprintf("user_%d", appUserID.Int()))
}
