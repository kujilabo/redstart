package service

import (
	libdomain "github.com/kujilabo/redstart/lib/domain"
	liberrors "github.com/kujilabo/redstart/lib/errors"
	"github.com/kujilabo/redstart/user/domain"
)

type UserRole interface {
	domain.UserRoleModel
}

type userRole struct {
	domain.UserRoleModel
}

// NewUserRole returns a new UserRole
func NewUserRole(userRoleModel domain.UserRoleModel) (UserRole, error) {
	m := &userRole{
		userRoleModel,
	}

	if err := libdomain.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return m, nil
}
