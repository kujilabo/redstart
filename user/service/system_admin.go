package service

import (
	"context"
	"fmt"

	liberrors "github.com/kujilabo/redstart/lib/errors"
	liblog "github.com/kujilabo/redstart/lib/log"
	"github.com/kujilabo/redstart/user/domain"
)

type SystemAdmin interface {
	domain.SystemAdminModel

	FindSystemOwnerByOrganizationID(ctx context.Context, organizationID domain.OrganizationID) (SystemOwner, error)

	FindSystemOwnerByOrganizationName(ctx context.Context, organizationName string) (SystemOwner, error)

	FindOrganizationByName(ctx context.Context, name string) (Organization, error)

	AddOrganization(ctx context.Context, parma OrganizationAddParameter) (domain.OrganizationID, error)
}

type systemAdmin struct {
	domain.SystemAdminModel
	rf          RepositoryFactory
	orgRepo     OrganizationRepository
	appUserRepo AppUserRepository
}

func NewSystemAdmin(ctx context.Context, rf RepositoryFactory) (SystemAdmin, error) {
	orgRepo := rf.NewOrganizationRepository(ctx)
	appUserRepo := rf.NewAppUserRepository(ctx)

	return &systemAdmin{
		SystemAdminModel: domain.NewSystemAdminModel(),
		rf:               rf,
		orgRepo:          orgRepo,
		appUserRepo:      appUserRepo,
	}, nil
}

func (m *systemAdmin) FindSystemOwnerByOrganizationID(ctx context.Context, organizationID domain.OrganizationID) (SystemOwner, error) {
	sysOwner, err := m.appUserRepo.FindSystemOwnerByOrganizationID(ctx, m, organizationID)
	if err != nil {
		return nil, liberrors.Errorf("m.appUserRepo.FindSystemOwnerByOrganizationID. error: %w", err)
	}

	return sysOwner, nil
}

func (m *systemAdmin) FindSystemOwnerByOrganizationName(ctx context.Context, organizationName string) (SystemOwner, error) {
	sysOwner, err := m.appUserRepo.FindSystemOwnerByOrganizationName(ctx, m, organizationName)
	if err != nil {
		return nil, liberrors.Errorf("m.appUserRepo.FindSystemOwnerByOrganizationName. error: %w", err)
	}

	return sysOwner, nil
}

func (m *systemAdmin) FindOrganizationByName(ctx context.Context, name string) (Organization, error) {
	org, err := m.orgRepo.FindOrganizationByName(ctx, m, name)
	if err != nil {
		return nil, liberrors.Errorf("m.orgRepo.FindOrganizationByName. error: %w", err)
	}

	return org, nil
}

func (m *systemAdmin) AddOrganization(ctx context.Context, param OrganizationAddParameter) (domain.OrganizationID, error) {
	logger := liblog.GetLoggerFromContext(ctx, UserServiceContextKey)

	// add organization
	organizationID, err := m.orgRepo.AddOrganization(ctx, m, param)
	if err != nil {
		return nil, liberrors.Errorf("failed to AddOrganization. error: %w", err)
	}

	userGroupRepo := m.rf.NewUserGroupRepository(ctx)

	// add system-owner-role
	systemOwnerGroupID, err := userGroupRepo.AddSystemOwnerGroup(ctx, m, organizationID)
	if err != nil {
		return nil, liberrors.Errorf("userGroupRepo.AddSystemOwnerRole. error: %w", err)
	}

	// add owner role
	ownerGroupID, err := userGroupRepo.AddOwnerGroup(ctx, m, organizationID)
	if err != nil {
		return nil, err
	}

	systemOwner, err := m.initSystemOwner(ctx, systemOwnerGroupID, organizationID, param.GetName())
	if err != nil {
		return nil, liberrors.Errorf("m.initSystemOwner. error: %w", err)
	}

	owner, err := m.initFirstOwner(ctx, systemOwner, ownerGroupID, param.GetFirstOwner())
	if err != nil {
		return nil, liberrors.Errorf("m.initFirstOwner. error: %w", err)
	}

	logger.InfoContext(ctx, fmt.Sprintf("SystemOwnerID:%d, owner: %+v", systemOwner.GetAppUserID().Int(), owner))

	return organizationID, nil
}

func (m *systemAdmin) initSystemOwner(ctx context.Context, systemOwnerGroupID domain.UserGroupID, organizationID domain.OrganizationID, organizationName string) (SystemOwner, error) {
	// add system-owner
	if _, err := m.appUserRepo.AddSystemOwner(ctx, m, organizationID); err != nil {
		return nil, liberrors.Errorf("failed to AddSystemOwner. error: %w", err)
	}

	systemOwner, err := m.appUserRepo.FindSystemOwnerByOrganizationName(ctx, m, organizationName)
	if err != nil {
		return nil, liberrors.Errorf("failed to FindSystemOwnerByOrganizationName. error: %w", err)
	}

	// systen-owner <-> system-owner-role
	pairOfUserAndGroup := m.rf.NewPairOfUserAndGroupRepository(ctx)

	if err := pairOfUserAndGroup.AddPairOfUserAndGroupToSystemOwner(ctx, m, systemOwner, systemOwnerGroupID); err != nil {
		return nil, err
	}

	return systemOwner, nil
}

func (m *systemAdmin) initFirstOwner(ctx context.Context, systemOwner SystemOwner, ownerGroupID domain.UserGroupID, param FirstOwnerAddParameter) (Owner, error) {

	// add owner
	ownerID, err := m.appUserRepo.AddFirstOwner(ctx, systemOwner, param)
	if err != nil {
		return nil, liberrors.Errorf("failed to AddFirstOwner. error: %w", err)
	}

	// owner <-> owner-role
	pairOfUserAndGroup := m.rf.NewPairOfUserAndGroupRepository(ctx)
	if err := pairOfUserAndGroup.AddPairOfUserAndGroup(ctx, systemOwner, ownerID, ownerGroupID); err != nil {
		return nil, err
	}

	//
	rbacRepo := m.rf.NewRBACRepository(ctx)
	rbacAppUser := NewRBACAppUser(ownerID)
	rbacAllUserRolesObject := NewRBACAllUserRoleObject()

	if err := rbacRepo.AddNamedPolicy(rbacAppUser, rbacAllUserRolesObject, RBACSetAction, RBACAllowEffect); err != nil {
		return nil, liberrors.Errorf("Failed to AddNamedPolicy. priv: read, err: %w", err)
	}

	owner, err := m.appUserRepo.FindOwnerByLoginID(ctx, systemOwner, param.GetLoginID())
	if err != nil {
		return nil, liberrors.Errorf("failed to FindOwnerByLoginID. error: %w", err)
	}

	return owner, nil
}

func NewRBACAppUser(appUserID domain.AppUserID) domain.RBACUser {
	return domain.NewRBACUser(fmt.Sprintf("user_%d", appUserID.Int()))
}

//	func NewRBACUserRole(userRoleID domain.UserGroupID) domain.RBACRole {
//		return domain.NewRBACRole(fmt.Sprintf("role_%d", userRoleID.Int()))
//	}
func NewRBACUserRole(key string) domain.RBACRole {
	return domain.NewRBACRole(fmt.Sprintf("role_%s", key))
}

func NewRBACUserRoleObject(userRoleID domain.UserGroupID) domain.RBACObject {
	return domain.NewRBACObject(fmt.Sprintf("role_%d", userRoleID.Int()))
}

func NewRBACAllUserRoleObject() domain.RBACObject {
	return domain.NewRBACObject("role_*")
}

var RBACSetAction = domain.NewRBACAction("Set")
var RBACUnsetAction = domain.NewRBACAction("Unset")

var RBACAllowEffect = domain.NewRBACEffect("allow")
var RBACDenyEffect = domain.NewRBACEffect("deny")
