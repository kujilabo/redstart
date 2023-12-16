package service

import (
	"context"

	"github.com/kujilabo/redstart/user/domain"
)

type PairOfUserAndGroupRepository interface {
	AddPairOfUserAndGroupBySystemAdmin(ctx context.Context, operator SystemAdminModelInterface, organizationID domain.OrganizationID, appUserID domain.AppUserID, userGroupID domain.UserGroupID) error

	AddPairOfUserAndGroup(ctx context.Context, operator AppUserModelInterface, appUserID domain.AppUserID, userGroupID domain.UserGroupID) error

	RemovePairOfUserAndGroup(ctx context.Context, operator AppUserModelInterface, appUserID domain.AppUserID, userGroupID domain.UserGroupID) error

	FindUserGroupsByUserID(ctx context.Context, operator AppUserModelInterface, appUserID domain.AppUserID) ([]*domain.UserGroupModel, error)
}
