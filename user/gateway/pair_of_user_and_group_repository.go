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
	PairOfUserAndGroupTableName = "user_n_group"
)

type pairOfUserAndGroupRepository struct {
	db *gorm.DB
}

type pairOfUserAndGroupEntity struct {
	JunctionModelEntity
	OrganizationID int
	AppUserID      int
	UserGroupID    int
}

func (u *pairOfUserAndGroupEntity) TableName() string {
	return PairOfUserAndGroupTableName
}

func NewPairOfUserAndGroupRepository(ctx context.Context, db *gorm.DB) service.PairOfUserAndGroupRepository {
	return &pairOfUserAndGroupRepository{
		db: db,
	}
}

func (r *pairOfUserAndGroupRepository) AddPairOfUserAndGroup(ctx context.Context, operator domain.AppUserModel, userGroupID domain.UserGroupID, appUserID domain.AppUserID) error {
	_, span := tracer.Start(ctx, "pairOfUserAndGroupRepository.AddPairOfUserAndGroup")
	defer span.End()

	pairOfUserAndGroup := pairOfUserAndGroupEntity{
		JunctionModelEntity: JunctionModelEntity{
			CreatedBy: operator.GetAppUserID().Int(),
		},
		OrganizationID: operator.GetOrganizationID().Int(),
		AppUserID:      appUserID.Int(),
		UserGroupID:    userGroupID.Int(),
	}
	if result := r.db.Create(&pairOfUserAndGroup); result.Error != nil {
		return liberrors.Errorf(". err: %w", libgateway.ConvertDuplicatedError(result.Error, service.ErrAppUserAlreadyExists))
	}
	return nil
}
