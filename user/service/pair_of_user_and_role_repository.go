package service

import (
	"context"

	"github.com/kujilabo/redstart/user/domain"
)

type PairOfUserAndRoleRepository interface {
	AddPairOfUserAndRole(ctx context.Context, operator domain.AppUserModel, userRoleID domain.UserGroupID, appUserID domain.AppUserID) error
}
