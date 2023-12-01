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
	UserRoleTableName = "user_role"
)

type userRoleEntity struct {
	BaseModelEntity
	ID             int
	OrganizationID int
	Key            string
	Name           string
	Description    string
	Removed        bool
}

func (e *userRoleEntity) TableName() string {
	return UserRoleTableName
}

func (e *userRoleEntity) toUserRoleModel() (domain.UserRoleModel, error) {
	baseModel, err := e.toBaseModel()
	if err != nil {
		return nil, liberrors.Errorf("toBaseModel. err: %w", err)
	}

	userRoleID, err := domain.NewUserRoleID(e.ID)
	if err != nil {
		return nil, liberrors.Errorf("domain.NewAppUserModel. err: %w", err)
	}

	organizationID, err := domain.NewOrganizationID(e.OrganizationID)
	if err != nil {
		return nil, liberrors.Errorf("domain.NewOrganizationID. err: %w", err)
	}

	userRoleModel, err := domain.NewUserRoleModel(baseModel, userRoleID, organizationID, e.Key, e.Name, e.Description)
	if err != nil {
		return nil, liberrors.Errorf("domain.NewAppUserRole. err: %w", err)
	}

	return userRoleModel, nil
}

func (e *userRoleEntity) toUserRole() (service.UserRole, error) {
	userRoleModel, err := e.toUserRoleModel()
	if err != nil {
		return nil, liberrors.Errorf("e.toUserRoleModel. err: %w", err)
	}

	userRole, err := service.NewUserRole(userRoleModel)
	if err != nil {
		return nil, liberrors.Errorf("service.NewAppUserRole. err: %w", err)
	}

	return userRole, nil
}

type userRoleRepository struct {
	db *gorm.DB
}

func NewUserRoleRepository(ctx context.Context, db *gorm.DB) service.UserRoleRepository {
	return &userRoleRepository{
		db: db,
	}
}

func (r *userRoleRepository) FindSystemOwnerRole(ctx context.Context, operator domain.SystemAdminModel, organizationID domain.OrganizationID) (service.UserRole, error) {
	_, span := tracer.Start(ctx, "userRoleRepository.FindUserRoleByKey")
	defer span.End()

	userRole := userRoleEntity{}
	if result := r.db.Where(&userRoleEntity{
		OrganizationID: organizationID.Int(),
		Key:            SystemOwnerRoleKey,
	}).Find(&userRole); result.Error != nil {
		return nil, result.Error
	}
	return userRole.toUserRole()
}

func (r *userRoleRepository) FindUserRoleByID(ctx context.Context, operator domain.AppUserModel, userRoleID domain.UserRoleID) (service.UserRole, error) {
	_, span := tracer.Start(ctx, "userRoleRepository.FindUserRoleByID")
	defer span.End()

	userRole := userRoleEntity{}
	if result := r.db.Where("organization_id = ?", operator.GetOrganizationID().Int()).
		Where("id = ? and removed = 0", userRoleID.Int()).
		Find(&userRole); result.Error != nil {
		return nil, result.Error
	}
	return userRole.toUserRole()
}

func (r *userRoleRepository) FindUserRoleByKey(ctx context.Context, operator domain.AppUserModel, key string) (service.UserRole, error) {
	_, span := tracer.Start(ctx, "userRoleRepository.FindUserRoleByKey")
	defer span.End()

	userRole := userRoleEntity{}
	if result := r.db.Where("`organization_id` = ?", operator.GetOrganizationID().Int()).
		Where("`key` = ? and `removed` = 0", key).
		Find(&userRole); result.Error != nil {
		return nil, result.Error
	}
	return userRole.toUserRole()
}

func (r *userRoleRepository) AddSystemOwnerRole(ctx context.Context, operator domain.SystemAdminModel, organizationID domain.OrganizationID) (domain.UserRoleID, error) {
	_, span := tracer.Start(ctx, "userRoleRepository.AddSystemOwnerRole")
	defer span.End()

	userRole := userRoleEntity{
		BaseModelEntity: BaseModelEntity{
			Version:   1,
			CreatedBy: operator.GetAppUserID().Int(),
			UpdatedBy: operator.GetAppUserID().Int(),
		},
		OrganizationID: organizationID.Int(),
		Key:            SystemOwnerRoleKey,
		Name:           SystemOwnerRoleName,
	}
	if result := r.db.Create(&userRole); result.Error != nil {
		return nil, liberrors.Errorf(". err: %w", libgateway.ConvertDuplicatedError(result.Error, service.ErrAppUserAlreadyExists))
	}

	userRoleID, err := domain.NewUserRoleID(userRole.ID)
	if err != nil {
		return nil, err
	}

	return userRoleID, nil
}

func (r *userRoleRepository) AddOwnerRole(ctx context.Context, operator domain.SystemAdminModel, organizationID domain.OrganizationID) (domain.UserRoleID, error) {
	_, span := tracer.Start(ctx, "userRoleRepository.AddOwnerRole")
	defer span.End()

	userRole := userRoleEntity{
		BaseModelEntity: BaseModelEntity{
			Version:   1,
			CreatedBy: operator.GetAppUserID().Int(),
			UpdatedBy: operator.GetAppUserID().Int(),
		},
		OrganizationID: organizationID.Int(),
		Key:            OwnerRoleKey,
		Name:           OwnerRoleName,
	}
	if result := r.db.Create(&userRole); result.Error != nil {
		return nil, liberrors.Errorf(". err: %w", libgateway.ConvertDuplicatedError(result.Error, service.ErrAppUserAlreadyExists))
	}

	userRoleID, err := domain.NewUserRoleID(userRole.ID)
	if err != nil {
		return nil, err
	}

	return userRoleID, nil
}

func (r *userRoleRepository) AddUserRole(ctx context.Context, operator domain.OwnerModel, parameter service.UserRoleAddParameter) (domain.UserRoleID, error) {
	_, span := tracer.Start(ctx, "userRoleRepository.AddUserRole")
	defer span.End()

	userRole := userRoleEntity{
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
	if result := r.db.Create(&userRole); result.Error != nil {
		return nil, liberrors.Errorf(". err: %w", libgateway.ConvertDuplicatedError(result.Error, service.ErrAppUserAlreadyExists))
	}

	userRoleID, err := domain.NewUserRoleID(userRole.ID)
	if err != nil {
		return nil, err
	}

	return userRoleID, nil
}

// func (r *userRoleRepository) AddPersonalGroup(ctx context.Context, operator domain.AppUserModel) (domain.UserRoleID, error) {
// 	_, span := tracer.Start(ctx, "userRoleRepository.AddPersonalGroup")
// 	defer span.End()

// 	userRole := userRoleEntity{
// 		BaseModelEntity: BaseModelEntity{
// 			Version:   1,
// 			CreatedBy: operator.GetAppUserID().Int(),
// 			UpdatedBy: operator.GetAppUserID().Int(),
// 		},
// 		OrganizationID: operator.GetOrganizationID().Int(),
// 		Key:            "#" + operator.GetLoginID(),
// 		Name:           "Personal group",
// 	}
// 	if result := r.db.Create(&userRole); result.Error != nil {
// 		return nil, liberrors.Errorf(". err: %w", libgateway.ConvertDuplicatedError(result.Error, service.ErrAppUserAlreadyExists))
// 	}

// 	userRoleID, err := domain.NewUserRoleID(userRole.ID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return userRoleID, nil
// }
