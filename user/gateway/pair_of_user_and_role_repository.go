package gateway

import (
	"context"

	"gorm.io/gorm"

	libdomain "github.com/kujilabo/redstart/lib/domain"
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
	rf service.RepositoryFactory
}

type pairOfUserAndRoleEntity struct {
	JunctionModelEntity
	OrganizationID int
	AppUserID      int
	UserRoleID     int
}

func (u *pairOfUserAndRoleEntity) TableName() string {
	return PairOfUserAndRoleTableName
}

func NewPairOfUserAndRoleRepository(ctx context.Context, db *gorm.DB, rf service.RepositoryFactory) service.PairOfUserAndRoleRepository {
	return &pairOfUserAndRoleRepository{
		db: db,
		rf: rf,
	}
}

func (r *pairOfUserAndRoleRepository) AddPairOfUserAndRoleToSystemOwner(ctx context.Context, operator domain.SystemAdminModel, systemOwner domain.SystemOwnerModel, userRoleID domain.UserRoleID) error {
	_, span := tracer.Start(ctx, "pairOfUserAndRoleRepository.AddPairOfUserAndRole")
	defer span.End()

	if err := r.add(ctx, operator.GetAppUserID(), systemOwner.GetOrganizationID(), systemOwner.GetAppUserID(), userRoleID, SystemOwnerRoleKey); err != nil {
		return nil
	}

	return nil
}

func (r *pairOfUserAndRoleRepository) AddPairOfUserAndRole(ctx context.Context, operator domain.AppUserModel, appUserID domain.AppUserID, userRoleID domain.UserRoleID) error {
	_, span := tracer.Start(ctx, "pairOfUserAndRoleRepository.AddPairOfUserAndRole")
	defer span.End()

	rbacAllUserRoleObject := service.NewRBACAllUserRoleObject()
	rbacUserRoleObject := service.NewRBACUserRoleObject(userRoleID)

	ok, err := r.enforce(ctx, operator, []domain.RBACObject{rbacAllUserRoleObject, rbacUserRoleObject}, service.RBACSetAction)
	if err != nil {
		return err
	}
	if !ok {
		return libdomain.ErrPermissionDenied
	}

	userRoleRepo := r.rf.NewUserRoleRepository(ctx)
	userRole, err := userRoleRepo.FindUserRoleByID(ctx, operator, userRoleID)
	if err != nil {
		return err
	}

	if err := r.add(ctx, operator.GetAppUserID(), operator.GetOrganizationID(), appUserID, userRoleID, userRole.GetKey()); err != nil {
		return err
	}
	return nil
}

func (r *pairOfUserAndRoleRepository) add(ctx context.Context, operatorID domain.AppUserID, organizationID domain.OrganizationID, appUserID domain.AppUserID, userRoleID domain.UserRoleID, userRoleKey string) error {
	// add pairOfOuserAndRole
	pairOfUserAndRole := pairOfUserAndRoleEntity{
		JunctionModelEntity: JunctionModelEntity{
			CreatedBy: operatorID.Int(),
		},
		OrganizationID: organizationID.Int(),
		AppUserID:      appUserID.Int(),
		UserRoleID:     userRoleID.Int(),
	}
	if result := r.db.Create(&pairOfUserAndRole); result.Error != nil {
		return liberrors.Errorf(". err: %w", libgateway.ConvertDuplicatedError(result.Error, service.ErrAppUserAlreadyExists))
	}

	rbacRepo := r.rf.NewRBACRepository(ctx)
	rbacAppUser := service.NewRBACAppUser(appUserID)
	rbacUserRole := service.NewRBACUserRole(userRoleKey)

	// add namedGroupingPolicy
	if err := rbacRepo.AddNamedGroupingPolicy(rbacAppUser, rbacUserRole); err != nil {
		return liberrors.Errorf("rbacRepo.AddNamedGroupingPolicy. err: %w", err)
	}

	return nil
}

func (r *pairOfUserAndRoleRepository) enforce(ctx context.Context, operator domain.AppUserModel, rbacObjects []domain.RBACObject, rbacAction domain.RBACAction) (bool, error) {

	appUserRepo := r.rf.NewAppUserRepository(ctx)
	appUser, err := appUserRepo.FindAppUserByID(ctx, operator, operator.GetAppUserID(), service.IncludeRoles)
	if err != nil {
		return false, err
	}

	rbacRoles := make([]domain.RBACRole, 0)
	for _, userRole := range appUser.GetUserRoles() {
		rbacRoles = append(rbacRoles, service.NewRBACUserRole(userRole.GetKey()))
	}

	rbacRepo := r.rf.NewRBACRepository(ctx)
	rbacOperator := service.NewRBACAppUser(operator.GetAppUserID())
	e, err := rbacRepo.NewEnforcerWithGroupsAndUsers(rbacRoles, []domain.RBACUser{rbacOperator})
	if err != nil {
		return false, err
	}

	for _, object := range rbacObjects {
		ok, err := e.Enforce(rbacOperator.Subject(), object.Object(), rbacAction.Action())
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}

	return false, nil
}

func (r *pairOfUserAndRoleRepository) FindUserRolesByUserID(ctx context.Context, operator domain.AppUserModel, appUserID domain.AppUserID) ([]domain.UserRoleModel, error) {
	userRoles := []userRoleEntity{}
	if result := r.db.Table("user_role").Select("user_role.*").
		Where("user_role.organization_id = ?", operator.GetOrganizationID().Int()).
		Where("user_role.removed = 0").
		Where("app_user.organization_id = ?", operator.GetOrganizationID().Int()).
		Where("app_user.id = ? and app_user.removed = 0", appUserID.Int()).
		Joins("inner join user_n_role on user_role.id = user_n_role.user_role_id").
		Joins("inner join app_user on user_n_role.app_user_id = app_user.id").
		Order("`user_role`.`key`").
		Find(&userRoles); result.Error != nil {
		return nil, result.Error
	}

	userRoleModels := make([]domain.UserRoleModel, len(userRoles))
	for i, e := range userRoles {
		m, err := e.toUserRoleModel()
		if err != nil {
			return nil, err
		}
		userRoleModels[i] = m
	}

	return userRoleModels, nil
}
