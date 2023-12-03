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

	// add system-owner-group
	systemOwnerGroupID, err := userGroupRepo.AddSystemOwnerGroup(ctx, m, organizationID)
	if err != nil {
		return nil, liberrors.Errorf("userGroupRepo.AddSystemOwnerRole. error: %w", err)
	}

	// add system-owner
	if _, err := m.appUserRepo.AddSystemOwner(ctx, m, organizationID); err != nil {
		return nil, liberrors.Errorf("failed to AddSystemOwner. error: %w", err)
	}

	systemOwner, err := m.appUserRepo.FindSystemOwnerByOrganizationName(ctx, m, param.GetName())
	if err != nil {
		return nil, liberrors.Errorf("failed to FindSystemOwnerByOrganizationName. error: %w", err)
	}

	pairOfUserAndGroup := m.rf.NewPairOfUserAndGroupRepository(ctx)

	// systen-owner belongs to system-owner-group
	if err := pairOfUserAndGroup.AddPairOfUserAndGroupToSystemOwner(ctx, m, systemOwner, systemOwnerGroupID); err != nil {
		return nil, err
	}

	// add owner group
	if _, err := userGroupRepo.AddOwnerGroup(ctx, m, organizationID); err != nil {
		return nil, err
	}

	// add first owner
	ownerID, err := systemOwner.AddFirstOwner(ctx, param.GetFirstOwner())
	if err != nil {
		return nil, liberrors.Errorf("m.initFirstOwner. error: %w", err)
	}

	logger.InfoContext(ctx, fmt.Sprintf("SystemOwnerID:%d, ownerID: %d", systemOwner.GetAppUserID().Int(), ownerID.Int()))

	return organizationID, nil
}

func NewRBACOrganization(organizationID domain.OrganizationID) domain.RBACDomain {
	return domain.NewRBACDomain(fmt.Sprintf("domain:%d", organizationID.Int()))
}

func NewRBACAppUser(organizationID domain.OrganizationID, appUserID domain.AppUserID) domain.RBACUser {
	return domain.NewRBACUser(fmt.Sprintf("user:%d", appUserID.Int()))
}

//	func NewRBACUserRole(userRoleID domain.UserGroupID) domain.RBACRole {
//		return domain.NewRBACRole(fmt.Sprintf("role_%d", userRoleID.Int()))
//	}
func NewRBACUserRole(organizationID domain.OrganizationID, userGroupID domain.UserGroupID) domain.RBACRole {
	return domain.NewRBACRole(fmt.Sprintf("domain:%d_role:%d", organizationID.Int(), userGroupID.Int()))
}

func NewRBACUserRoleObject(organizationID domain.OrganizationID, userRoleID domain.UserGroupID) domain.RBACObject {
	return domain.NewRBACObject(fmt.Sprintf("domain:%d_role:%d", organizationID.Int(), userRoleID.Int()))
}

func NewRBACAllUserRolesObject(organizationID domain.OrganizationID) domain.RBACObject {
	return domain.NewRBACObject(fmt.Sprintf("domain:%d_role:*", organizationID.Int()))
}

var RBACSetAction = domain.NewRBACAction("Set")
var RBACUnsetAction = domain.NewRBACAction("Unset")

var RBACAllowEffect = domain.NewRBACEffect("allow")
var RBACDenyEffect = domain.NewRBACEffect("deny")
