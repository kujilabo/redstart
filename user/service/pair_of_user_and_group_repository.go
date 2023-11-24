package service

import (
	"context"

	"github.com/kujilabo/redstart/user/domain"
)

type PairOfUserAndGroupRepository interface {
	AddPairOfUserAndGroup(ctx context.Context, operator domain.AppUserModel, userGroupID domain.UserGroupID, appUserID domain.AppUserID) error
}
