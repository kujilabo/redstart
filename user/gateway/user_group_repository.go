package gateway

import (
	"context"
	"errors"

	"gorm.io/gorm"

	liberrors "github.com/kujilabo/redstart/lib/errors"
	libgateway "github.com/kujilabo/redstart/lib/gateway"
	"github.com/kujilabo/redstart/user/domain"
	"github.com/kujilabo/redstart/user/service"
)

var (
	UserGroupTableName = "user_group"
	PublicGroupKey     = "__public"
)

type userGroupEntity struct {
	BaseModelEntity
	ID             int
	OrganizationID int
	Key            string
	Name           string
	Description    string
	Removed        bool
}

func (e *userGroupEntity) TableName() string {
	return UserGroupTableName
}

func (e *userGroupEntity) toAppUserGroup() (service.UserGroup, error) {
	baseModel, err := e.toBaseModel()
	if err != nil {
		return nil, liberrors.Errorf("toAppUserGroup. err: %w", err)
	}

	userGroupID, err := domain.NewUserGroupID(e.ID)
	if err != nil {
		return nil, liberrors.Errorf("domain.NewAppUserModel. err: %w", err)
	}

	organizationID, err := domain.NewOrganizationID(e.OrganizationID)
	if err != nil {
		return nil, liberrors.Errorf("domain.NewOrganizationID. err: %w", err)
	}

	userGroupMdoel, err := domain.NewUserGroupModel(baseModel, userGroupID, organizationID, e.Key, e.Name, e.Description)
	if err != nil {
		return nil, liberrors.Errorf("domain.NewAppUserGroup. err: %w", err)
	}

	userGroup, err := service.NewUserGroup(userGroupMdoel)
	if err != nil {
		return nil, liberrors.Errorf("service.NewAppUserGroup. err: %w", err)
	}

	return userGroup, nil
}

type userGroupRepository struct {
	db *gorm.DB
}

func NewUserGroupRepository(ctx context.Context, db *gorm.DB) service.UserGroupRepository {
	if db == nil {
		panic(errors.New("db is nil"))
	}

	return &userGroupRepository{
		db: db,
	}
}

func (r *userGroupRepository) FindPublicGroup(ctx context.Context, operator domain.SystemOwnerModel) (service.UserGroup, error) {
	_, span := tracer.Start(ctx, "userGroupRepository.FindPublicGroup")
	defer span.End()

	userGroup := userGroupEntity{}
	if result := r.db.Where(&userGroupEntity{
		OrganizationID: operator.GetOrganizationID().Int(),
		Key:            PublicGroupKey,
	}).Find(&userGroup); result.Error != nil {
		return nil, result.Error
	}
	return userGroup.toAppUserGroup()
}

func (r *userGroupRepository) AddPublicGroup(ctx context.Context, operator domain.SystemOwnerModel) (domain.UserGroupID, error) {
	_, span := tracer.Start(ctx, "userGroupRepository.AddPublicGroup")
	defer span.End()

	userGroup := userGroupEntity{
		BaseModelEntity: BaseModelEntity{
			Version:   1,
			CreatedBy: operator.GetAppUserID().Int(),
			UpdatedBy: operator.GetAppUserID().Int(),
		},
		OrganizationID: operator.GetOrganizationID().Int(),
		Key:            PublicGroupKey,
		Name:           "Public group",
	}
	if result := r.db.Create(&userGroup); result.Error != nil {
		return nil, liberrors.Errorf(". err: %w", libgateway.ConvertDuplicatedError(result.Error, service.ErrAppUserAlreadyExists))
	}

	userGroupID, err := domain.NewUserGroupID(userGroup.ID)
	if err != nil {
		return nil, err
	}

	return userGroupID, nil
}

func (r *userGroupRepository) AddPersonalGroup(ctx context.Context, operator domain.AppUserModel) (domain.UserGroupID, error) {
	_, span := tracer.Start(ctx, "userGroupRepository.AddPersonalGroup")
	defer span.End()

	userGroup := userGroupEntity{
		BaseModelEntity: BaseModelEntity{
			Version:   1,
			CreatedBy: operator.GetAppUserID().Int(),
			UpdatedBy: operator.GetAppUserID().Int(),
		},
		OrganizationID: operator.GetOrganizationID().Int(),
		Key:            "#" + operator.GetLoginID(),
		Name:           "Personal group",
	}
	if result := r.db.Create(&userGroup); result.Error != nil {
		return nil, liberrors.Errorf(". err: %w", libgateway.ConvertDuplicatedError(result.Error, service.ErrAppUserAlreadyExists))
	}

	userGroupID, err := domain.NewUserGroupID(userGroup.ID)
	if err != nil {
		return nil, err
	}

	return userGroupID, nil
}
