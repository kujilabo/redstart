package service

import (
	"context"

	"github.com/kujilabo/redstart/user/domain"
)

type UserGroupRepository interface {
	FindPublicGroup(ctx context.Context, operator domain.SystemOwnerModel) (UserGroup, error)

	AddPublicGroup(ctx context.Context, operator domain.SystemOwnerModel) (domain.UserGroupID, error)
	// AddUserGroup(ctx context.Context, operator domain.AppUserModel, parameter UserGroupAddParameter) (domain.UserGroupID, error)
	// AddPersonalGroup(operator SystemOwner, studentID uint) (uint, error)
}
