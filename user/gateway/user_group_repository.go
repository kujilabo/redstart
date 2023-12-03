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
	UserGroupTableName = "user_group"
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

func (e *userGroupEntity) toUserGroupModel() (domain.UserGroupModel, error) {
	baseModel, err := e.toBaseModel()
	if err != nil {
		return nil, liberrors.Errorf("toBaseModel. err: %w", err)
	}

	userGroupID, err := domain.NewUserGroupID(e.ID)
	if err != nil {
		return nil, liberrors.Errorf("domain.NewAppUserModel. err: %w", err)
	}

	organizationID, err := domain.NewOrganizationID(e.OrganizationID)
	if err != nil {
		return nil, liberrors.Errorf("domain.NewOrganizationID. err: %w", err)
	}

	userGroupModel, err := domain.NewUserGroupModel(baseModel, userGroupID, organizationID, e.Key, e.Name, e.Description)
	if err != nil {
		return nil, liberrors.Errorf("domain.NewUserGroupModel. err: %w", err)
	}

	return userGroupModel, nil
}

func (e *userGroupEntity) toUserGroup() (service.UserGroup, error) {
	userGroupModel, err := e.toUserGroupModel()
	if err != nil {
		return nil, liberrors.Errorf("e.touserGroupModel. err: %w", err)
	}

	userGroup, err := service.NewUserGroup(userGroupModel)
	if err != nil {
		return nil, liberrors.Errorf("service.NewUserGroup. err: %w", err)
	}

	return userGroup, nil
}

type userGroupRepository struct {
	db *gorm.DB
}

func NewUserGroupRepository(ctx context.Context, db *gorm.DB) service.UserGroupRepository {
	return &userGroupRepository{
		db: db,
	}
}

func (r *userGroupRepository) FindSystemOwnerGroup(ctx context.Context, operator domain.SystemAdminModel, organizationID domain.OrganizationID) (service.UserGroup, error) {
	_, span := tracer.Start(ctx, "userGroupRepository.FindSystemOwnerGroup")
	defer span.End()

	userGroup := userGroupEntity{}
	if result := r.db.Where(&userGroupEntity{
		OrganizationID: organizationID.Int(),
		Key:            service.SystemOwnerGroupKey,
	}).Find(&userGroup); result.Error != nil {
		return nil, result.Error
	}
	return userGroup.toUserGroup()
}

func (r *userGroupRepository) FindUserGroupByID(ctx context.Context, operator domain.AppUserModel, userGroupID domain.UserGroupID) (service.UserGroup, error) {
	_, span := tracer.Start(ctx, "userGroupRepository.FindUserGroupByID")
	defer span.End()

	userGroup := userGroupEntity{}
	if result := r.db.Where("organization_id = ?", operator.GetOrganizationID().Int()).
		Where("id = ? and removed = 0", userGroupID.Int()).
		Find(&userGroup); result.Error != nil {
		return nil, result.Error
	}
	return userGroup.toUserGroup()
}

func (r *userGroupRepository) FindUserGroupByKey(ctx context.Context, operator domain.AppUserModel, key string) (service.UserGroup, error) {
	_, span := tracer.Start(ctx, "userGroupRepository.FindUserGroupByKey")
	defer span.End()

	userGroup := userGroupEntity{}
	if result := r.db.Where("`organization_id` = ?", operator.GetOrganizationID().Int()).
		Where("`key` = ? and `removed` = 0", key).
		Find(&userGroup); result.Error != nil {
		return nil, result.Error
	}
	return userGroup.toUserGroup()
}

func (r *userGroupRepository) AddSystemOwnerGroup(ctx context.Context, operator domain.SystemAdminModel, organizationID domain.OrganizationID) (domain.UserGroupID, error) {
	_, span := tracer.Start(ctx, "userGroupRepository.AddSystemOwnerGroup")
	defer span.End()

	userGroup := userGroupEntity{
		BaseModelEntity: BaseModelEntity{
			Version:   1,
			CreatedBy: operator.GetAppUserID().Int(),
			UpdatedBy: operator.GetAppUserID().Int(),
		},
		OrganizationID: organizationID.Int(),
		Key:            service.SystemOwnerGroupKey,
		Name:           service.SystemOwnerGroupName,
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

func (r *userGroupRepository) AddOwnerGroup(ctx context.Context, operator domain.SystemAdminModel, organizationID domain.OrganizationID) (domain.UserGroupID, error) {
	_, span := tracer.Start(ctx, "userGroupRepository.AddOwnerGroup")
	defer span.End()

	userGroup := userGroupEntity{
		BaseModelEntity: BaseModelEntity{
			Version:   1,
			CreatedBy: operator.GetAppUserID().Int(),
			UpdatedBy: operator.GetAppUserID().Int(),
		},
		OrganizationID: organizationID.Int(),
		Key:            service.OwnerGroupKey,
		Name:           service.OwnerGroupName,
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

func (r *userGroupRepository) AddUserGroup(ctx context.Context, operator domain.OwnerModel, parameter service.UserGroupAddParameter) (domain.UserGroupID, error) {
	_, span := tracer.Start(ctx, "userGroupRepository.AddUserGroup")
	defer span.End()

	userGroup := userGroupEntity{
		BaseModelEntity: BaseModelEntity{
			Version:   1,
			CreatedBy: operator.GetAppUserID().Int(),
			UpdatedBy: operator.GetAppUserID().Int(),
		},
		OrganizationID: operator.GetOrganizationID().Int(),
		Key:            parameter.GetKey(),
		Name:           parameter.GetName(),
		Description:    parameter.GetDescription(),
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
