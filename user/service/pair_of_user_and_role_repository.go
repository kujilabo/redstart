package service

import (
	"context"

	"github.com/kujilabo/redstart/user/domain"
)

type PairOfUserAndRoleRepository interface {
	AddPairOfUserAndRoleToSystemOwner(ctx context.Context, operator domain.SystemAdminModel, systemOwner domain.SystemOwnerModel, userRoleID domain.UserRoleID) error

	AddPairOfUserAndRole(ctx context.Context, operator domain.AppUserModel, appUserID domain.AppUserID, userRoleID domain.UserRoleID) error

	FindUserRolesByUserID(ctx context.Context, operator domain.AppUserModel, appUserID domain.AppUserID) ([]domain.UserRoleModel, error)
}
