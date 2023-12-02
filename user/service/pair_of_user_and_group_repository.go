package service

import (
	"context"

	"github.com/kujilabo/redstart/user/domain"
)

type PairOfUserAndGroupRepository interface {
	AddPairOfUserAndGroupToSystemOwner(ctx context.Context, operator domain.SystemAdminModel, systemOwner domain.SystemOwnerModel, userGroupID domain.UserGroupID) error

	AddPairOfUserAndGroup(ctx context.Context, operator domain.AppUserModel, appUserID domain.AppUserID, userGroupID domain.UserGroupID) error

	FindUserGroupsByUserID(ctx context.Context, operator domain.AppUserModel, appUserID domain.AppUserID) ([]domain.UserGroupModel, error)
}
