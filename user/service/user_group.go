package service

import (
	libdomain "github.com/kujilabo/redstart/lib/domain"
	liberrors "github.com/kujilabo/redstart/lib/errors"
	"github.com/kujilabo/redstart/user/domain"
)

type UserGroup interface {
	domain.UserGroupModel
}

type userGroup struct {
	domain.UserGroupModel
}

// NewUserGroup returns a new UserGroup
func NewUserGroup(userGroupModel domain.UserGroupModel) (UserGroup, error) {
	m := &userGroup{
		userGroupModel,
	}

	if err := libdomain.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return m, nil
}
