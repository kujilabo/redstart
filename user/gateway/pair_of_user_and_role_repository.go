package gateway

import (
	"context"

	"gorm.io/gorm"

	liberrors "github.com/kujilabo/redstart/lib/errors"
	libgateway "github.com/kujilabo/redstart/lib/gateway"
	"github.com/kujilabo/redstart/user/domain"
	"github.com/kujilabo/redstart/user/service"
)

var (
	PairOfUserAndRoleTableName = "user_n_role"
)

type pairOfUserAndRoleRepository struct {
	db *gorm.DB
}

type pairOfUserAndRoleEntity struct {
	JunctionModelEntity
	OrganizationID int
	AppUserID      int
	UserGroupID    int
}

func (u *pairOfUserAndRoleEntity) TableName() string {
	return PairOfUserAndRoleTableName
}

func NewPairOfUserAndRoleRepository(ctx context.Context, db *gorm.DB) service.PairOfUserAndRoleRepository {
	return &pairOfUserAndRoleRepository{
		db: db,
	}
}

func (r *pairOfUserAndRoleRepository) AddPairOfUserAndRole(ctx context.Context, operator domain.AppUserModel, userGroupID domain.UserGroupID, appUserID domain.AppUserID) error {
	_, span := tracer.Start(ctx, "pairOfUserAndRoleRepository.AddPairOfUserAndRole")
	defer span.End()

	pairOfUserAndRole := pairOfUserAndRoleEntity{
		JunctionModelEntity: JunctionModelEntity{
			CreatedBy: operator.GetAppUserID().Int(),
		},
		OrganizationID: operator.GetOrganizationID().Int(),
		AppUserID:      appUserID.Int(),
		UserGroupID:    userGroupID.Int(),
	}
	if result := r.db.Create(&pairOfUserAndRole); result.Error != nil {
		return liberrors.Errorf(". err: %w", libgateway.ConvertDuplicatedError(result.Error, service.ErrAppUserAlreadyExists))
	}
	return nil
}
