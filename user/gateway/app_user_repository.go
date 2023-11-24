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
	AppUserTableName = "app_user"

	SystemAdminLoginID = "__system_admin"
	SystemOwnerLoginID = "__system_owner"

	SystemOwnerRoleKey = "__system_owner"
	OwnerRoleKey       = "__owner"
	// SystemStudentLoginID = "system-student"
	// GuestLoginID         = "guest"

	// AdministratorRole = "Administrator"
	OwnerRole = "_owner"
	// ManagerRole       = "Manager"
	// UserRole          = "User"
	// GuestRole         = "Guest"
	// UnknownRole       = "Unknown"
)

type appUserRepository struct {
	db *gorm.DB
	rf service.RepositoryFactory
}

type appUserEntity struct {
	BaseModelEntity
	ID                   int
	OrganizationID       int
	LoginID              string
	Username             string
	HashedPassword       string
	Provider             string
	ProviderID           string
	ProviderAccessToken  string
	ProviderRefreshToken string
	Removed              bool
}

func (e *appUserEntity) TableName() string {
	return AppUserTableName
}

func (e *appUserEntity) toSystemOwner(ctx context.Context, rf service.RepositoryFactory) (service.SystemOwner, error) {
	if e.LoginID != SystemOwnerLoginID {
		return nil, liberrors.Errorf("invalid system owner. loginID: %s", e.LoginID)
	}

	baseModel, err := e.toBaseModel()
	if err != nil {
		return nil, liberrors.Errorf("e.toBaseModel. err: %w", err)
	}

	appUserID, err := domain.NewAppUserID(e.ID)
	if err != nil {
		return nil, liberrors.Errorf("domain.NewAppUserModel. err: %w", err)
	}

	organizationID, err := domain.NewOrganizationID(e.OrganizationID)
	if err != nil {
		return nil, liberrors.Errorf("domain.NewOrganizationID. err: %w", err)
	}

	appUserModel, err := domain.NewAppUserModel(baseModel, appUserID, organizationID, e.LoginID, e.Username)
	if err != nil {
		return nil, liberrors.Errorf("domain.NewAppUserModel. err: %w", err)
	}

	ownerModel, err := domain.NewOwnerModel(appUserModel)
	if err != nil {
		return nil, liberrors.Errorf("domain.NewOwnerModel. err: %w", err)
	}

	systemOwnerModel, err := domain.NewSystemOwnerModel(ownerModel)
	if err != nil {
		return nil, liberrors.Errorf("domain.NewSystemOwnerModel. err: %w", err)
	}

	systemOwner, err := service.NewSystemOwner(ctx, rf, systemOwnerModel)
	if err != nil {
		return nil, liberrors.Errorf("service.NewSystemOwner. err: %w", err)
	}

	return systemOwner, nil
}

func (e *appUserEntity) toOwner(rf service.RepositoryFactory) (service.Owner, error) {
	appUserModel, err := e.toAppUserModel()
	if err != nil {
		return nil, liberrors.Errorf("e.toAppUserModel. err: %w", err)
	}

	ownerModel, err := domain.NewOwnerModel(appUserModel)
	if err != nil {
		return nil, liberrors.Errorf("domain.NewOwnerModel. err: %w", err)
	}

	return service.NewOwner(rf, ownerModel), nil
}

func (e *appUserEntity) toAppUserModel() (domain.AppUserModel, error) {
	baseModel, err := e.toBaseModel()
	if err != nil {
		return nil, liberrors.Errorf("e.toModel. err: %w", err)
	}

	appUserID, err := domain.NewAppUserID(e.ID)
	if err != nil {
		return nil, liberrors.Errorf("domain.NewAppUserModel. err: %w", err)
	}

	organizationID, err := domain.NewOrganizationID(e.OrganizationID)
	if err != nil {
		return nil, liberrors.Errorf("domain.NewOrganizationID. err: %w", err)
	}

	appUserModel, err := domain.NewAppUserModel(baseModel, appUserID, organizationID, e.LoginID, e.Username)
	if err != nil {
		return nil, liberrors.Errorf("domain.NewAppUserModel. err: %w", err)
	}

	return appUserModel, nil
}

func (e *appUserEntity) toAppUser(ctx context.Context, rf service.RepositoryFactory) (service.AppUser, error) {
	appUserModel, err := e.toAppUserModel()
	if err != nil {
		return nil, liberrors.Errorf("e.toAppUserModel. err: %w", err)
	}

	appUser, err := service.NewAppUser(ctx, rf, appUserModel)
	if err != nil {
		return nil, liberrors.Errorf("service.NewAppUser. err: %w", err)
	}

	return appUser, nil
}

// func (e *appUserEntity) toAppUser(ctx context.Context, rf service.RepositoryFactory) (service.AppUser, error) {
// 	appUserModel, err := e.toAppUserModel()
// 	if err != nil {
// 		return nil, liberrors.Errorf("e.toAppUserModel. err: %w", err)
// 	}

// 	appUser, err := service.NewAppUser(ctx, rf, appUserModel)
// 	if err != nil {
// 		return nil, liberrors.Errorf("service.NewAppUser. err: %w", err)
// 	}

// 	return appUser, nil
// }

// func NewAppUserRepository(ctx context.Context, rf service.RepositoryFactory, db *gorm.DB) service.AppUserRepository {
// 	if rf == nil {
// 		panic(errors.New("rf is nil"))
// 	} else if db == nil {
// 		panic(errors.New("db is nil"))
// 	}
// 	return &appUserRepository{
// 		rf: rf,
// 		db: db,
// 	}
// }

func NewAppUserRepository(ctx context.Context, driverName string, db *gorm.DB, rf service.RepositoryFactory) service.AppUserRepository {
	return &appUserRepository{
		db: db,
		rf: rf,
	}
}

func (r *appUserRepository) FindSystemOwnerByOrganizationID(ctx context.Context, operator domain.SystemAdminModel, organizationID domain.OrganizationID) (service.SystemOwner, error) {
	_, span := tracer.Start(ctx, "appUserRepository.FindSystemOwnerByOrganizationID")
	defer span.End()

	appUser := appUserEntity{}
	if result := r.db.Where("organization_id = ?", organizationID).
		Where("login_id = ? and removed = 0", SystemOwnerLoginID).
		First(&appUser); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, liberrors.Errorf("system owner not found. organization ID: %d, err: %w", organizationID, service.ErrSystemOwnerNotFound)
		}
		return nil, result.Error
	}
	return appUser.toSystemOwner(ctx, r.rf)
}

func (r *appUserRepository) FindSystemOwnerByOrganizationName(ctx context.Context, operator domain.SystemAdminModel, organizationName string) (service.SystemOwner, error) {
	_, span := tracer.Start(ctx, "appUserRepository.FindSystemOwnerByOrganizationName")
	defer span.End()

	appUser := appUserEntity{}
	if result := r.db.Table("organization").Select("app_user.*").
		Where("organization.name = ? and app_user.removed = 0", organizationName).
		Where("login_id = ?", SystemOwnerLoginID).
		Joins("inner join app_user on organization.id = app_user.organization_id").
		First(&appUser); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, liberrors.Errorf("system owner not found. organization name: %s, err: %w", organizationName, service.ErrSystemOwnerNotFound)
		}

		return nil, result.Error
	}
	return appUser.toSystemOwner(ctx, r.rf)
}

func (r *appUserRepository) FindAppUserByID(ctx context.Context, operator domain.AppUserModel, id domain.AppUserID) (service.AppUser, error) {
	_, span := tracer.Start(ctx, "appUserRepository.FindAppUserByID")
	defer span.End()

	appUser := appUserEntity{}
	if result := r.db.Where("organization_id = ?", operator.GetOrganizationID().Int()).
		Where("id = ? and removed = 0", id.Int()).
		First(&appUser); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrAppUserNotFound
		}

		return nil, result.Error
	}

	return appUser.toAppUser(ctx, r.rf)
}

func (r *appUserRepository) FindAppUserByLoginID(ctx context.Context, operator domain.AppUserModel, loginID string) (service.AppUser, error) {
	_, span := tracer.Start(ctx, "appUserRepository.FindAppUserByLoginID")
	defer span.End()

	appUser := appUserEntity{}
	if result := r.db.Where("organization_id = ?", operator.GetOrganizationID().Int()).
		Where("login_id = ? and removed = 0", loginID).
		First(&appUser); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrAppUserNotFound
		}

		return nil, result.Error
	}

	return appUser.toAppUser(ctx, r.rf)
}

func (r *appUserRepository) FindOwnerByLoginID(ctx context.Context, operator domain.SystemOwnerModel, loginID string) (service.Owner, error) {
	_, span := tracer.Start(ctx, "appUserRepository.FindOwnerByLoginID")
	defer span.End()

	appUser := appUserEntity{}
	if result := r.db.Where(&appUserEntity{
		OrganizationID: operator.GetOrganizationID().Int(),
		LoginID:        loginID,
	}).First(&appUser); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrAppUserNotFound
		}

		return nil, result.Error
	}

	return appUser.toOwner(r.rf)
}

func (r *appUserRepository) addAppUser(ctx context.Context, appUserEntity *appUserEntity) (domain.AppUserID, error) {
	if result := r.db.Create(appUserEntity); result.Error != nil {
		return nil, liberrors.Errorf("db.Create. err: %w", libgateway.ConvertDuplicatedError(result.Error, service.ErrAppUserAlreadyExists))
	}

	appUserID, err := domain.NewAppUserID(appUserEntity.ID)
	if err != nil {
		return nil, err
	}

	return appUserID, nil
}

func (r *appUserRepository) AddAppUser(ctx context.Context, operator domain.OwnerModel, param service.AppUserAddParameter) (domain.AppUserID, error) {
	_, span := tracer.Start(ctx, "appUserRepository.AddAppUser")
	defer span.End()

	hashedPassword := ""
	if len(param.GetPassword()) != 0 {
		hashedPasswordTmp, err := libgateway.HashPassword(param.GetPassword())
		if err != nil {
			return nil, liberrors.Errorf("libgateway.HashPassword. err: %w", err)
		}

		hashedPassword = hashedPasswordTmp
	}

	appUserEntity := appUserEntity{
		BaseModelEntity: BaseModelEntity{
			Version:   1,
			CreatedBy: operator.GetAppUserID().Int(),
			UpdatedBy: operator.GetAppUserID().Int(),
		},
		OrganizationID: operator.GetAppUserID().Int(),
		LoginID:        param.GetLoginID(),
		Username:       param.GetUsername(),
		HashedPassword: hashedPassword,
	}

	appUserID, err := r.addAppUser(ctx, &appUserEntity)
	if err != nil {
		return nil, err
	}

	return appUserID, nil
}

func (r *appUserRepository) AddSystemOwner(ctx context.Context, operator domain.SystemAdminModel, organizationID domain.OrganizationID) (domain.AppUserID, error) {
	_, span := tracer.Start(ctx, "appUserRepository.AddSystemOwner")
	defer span.End()

	appUserEntity := appUserEntity{
		BaseModelEntity: BaseModelEntity{
			Version:   1,
			CreatedBy: operator.GetAppUserID().Int(),
			UpdatedBy: operator.GetAppUserID().Int(),
		},
		OrganizationID: organizationID.Int(),
		LoginID:        SystemOwnerLoginID,
		Username:       "SystemOwner",
	}

	appUserID, err := r.addAppUser(ctx, &appUserEntity)
	if err != nil {
		return nil, err
	}

	return appUserID, nil
}

func (r *appUserRepository) AddFirstOwner(ctx context.Context, operator domain.SystemOwnerModel, param service.FirstOwnerAddParameter) (domain.AppUserID, error) {
	_, span := tracer.Start(ctx, "appUserRepository.AddFirstOwner")
	defer span.End()

	hashedPassword, err := libgateway.HashPassword(param.GetPassword())
	if err != nil {
		return nil, liberrors.Errorf("passwordhelper.HashPassword. err: %w", err)
	}

	appUserEntity := appUserEntity{
		BaseModelEntity: BaseModelEntity{
			Version:   1,
			CreatedBy: operator.GetAppUserID().Int(),
			UpdatedBy: operator.GetAppUserID().Int(),
		},
		OrganizationID: operator.GetAppUserID().Int(),
		LoginID:        SystemOwnerLoginID,
		Username:       "SystemOwner",
		HashedPassword: hashedPassword,
	}

	appUserID, err := r.addAppUser(ctx, &appUserEntity)
	if err != nil {
		return nil, err
	}

	return appUserID, nil
}
