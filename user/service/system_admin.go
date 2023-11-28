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

	userRoleRepo := m.rf.NewUserRoleRepository(ctx)

	// add system-owner-role
	systemOwnerRoleID, err := userRoleRepo.AddSystemOwnerRole(ctx, m, organizationID)
	if err != nil {
		return nil, liberrors.Errorf("userRoleRepo.AddSystemOwnerRole. error: %w", err)
	}

	// add owner role
	ownerRoleID, err := userRoleRepo.AddOwnerRole(ctx, m, organizationID)
	if err != nil {
		return nil, err
	}

	systemOwner, err := m.initSystemOwner(ctx, systemOwnerRoleID, organizationID, param.GetName())
	if err != nil {
		return nil, liberrors.Errorf("m.initSystemOwner. error: %w", err)
	}

	owner, err := m.initFirstOwner(ctx, systemOwner, ownerRoleID, param.GetFirstOwner())
	if err != nil {
		return nil, liberrors.Errorf("m.initFirstOwner. error: %w", err)
	}

	// spaceRepo := m.rf.NewSpaceRepository(ctx)

	// // add default space
	// spaceID, err := spaceRepo.AddDefaultSpace(ctx, systemOwner)
	// if err != nil {
	// 	return nil, liberrors.Errorf("failed to AddDefaultSpace. error: %w", err)
	// }

	logger.InfoContext(ctx, fmt.Sprintf("SystemOwnerID:%d, owner: %+v", systemOwner.GetAppUserID().Int(), owner))
	// logger.Infof("SystemOwnerID:%d, SystemStudentID:%d, owner: %+v, spaceID: %d", systemOwnerID, systemStudentID, owner, spaceID)

	// // add personal group
	// personalGroupID, err := s.appUserGroupRepositor.AddPublicGroup(owner)
	// if err != nil {
	// 	return 0, fmt.Errorf("failed to AddPersonalGroup. error: %w", err)
	// }

	// // personal-group <-> owner
	// if err := s.groupUserRepository.AddGroupUser(systemOwner, personalGroupID, ownerID); err != nil {
	// 	return 0, fmt.Errorf("failed to AddGroupUser. error: %w", err)
	// }

	return organizationID, nil
}

func (m *systemAdmin) initSystemOwner(ctx context.Context, systemOwnerRoleID domain.UserRoleID, organizationID domain.OrganizationID, organizationName string) (SystemOwner, error) {
	// add system-owner
	if _, err := m.appUserRepo.AddSystemOwner(ctx, m, organizationID); err != nil {
		return nil, liberrors.Errorf("failed to AddSystemOwner. error: %w", err)
	}

	systemOwner, err := m.appUserRepo.FindSystemOwnerByOrganizationName(ctx, m, organizationName)
	if err != nil {
		return nil, liberrors.Errorf("failed to FindSystemOwnerByOrganizationName. error: %w", err)
	}

	// systen-owner <-> system-owner-role
	pairOfUserAndRole := m.rf.NewPairOfUserAndRoleRepository(ctx)

	if err := pairOfUserAndRole.AddPairOfUserAndRoleToSystemOwner(ctx, m, systemOwner, systemOwnerRoleID); err != nil {
		return nil, err
	}

	return systemOwner, nil
}

func (m *systemAdmin) initFirstOwner(ctx context.Context, systemOwner SystemOwner, ownerRoleID domain.UserRoleID, param FirstOwnerAddParameter) (Owner, error) {
	userGroupRepo := m.rf.NewUserGroupRepository(ctx)
	pairOfUserAndGroupRepo := m.rf.NewPairOfUserAndGroupRepository(ctx)

	// add owner
	ownerID, err := m.appUserRepo.AddFirstOwner(ctx, systemOwner, param)
	if err != nil {
		return nil, liberrors.Errorf("failed to AddFirstOwner. error: %w", err)
	}

	// owner <-> owner-role
	pairOfUserAndRole := m.rf.NewPairOfUserAndRoleRepository(ctx)
	if err := pairOfUserAndRole.AddPairOfUserAndRole(ctx, systemOwner, ownerID, ownerRoleID); err != nil {
		return nil, err
	}

	// add public group
	publicGroupID, err := userGroupRepo.AddPublicGroup(ctx, systemOwner)
	if err != nil {
		return nil, liberrors.Errorf("failed to AddPublicGroup. error: %w", err)
	}

	// public-group <-> owner
	if err := pairOfUserAndGroupRepo.AddPairOfUserAndGroup(ctx, systemOwner, publicGroupID, ownerID); err != nil {
		return nil, liberrors.Errorf("failed to AddGroupUser. error: %w", err)
	}

	//
	rbacRepo := m.rf.NewRBACRepository(ctx)
	rbacAppUser := NewRBACAppUser(ownerID)
	rbacAllUserRolesObject := NewRBACAllUserRoleObject()

	if err := rbacRepo.AddNamedPolicy(rbacAppUser, rbacAllUserRolesObject, RBACSetAction); err != nil {
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

//	func NewRBACUserRole(userRoleID domain.UserRoleID) domain.RBACRole {
//		return domain.NewRBACRole(fmt.Sprintf("role_%d", userRoleID.Int()))
//	}
func NewRBACUserRole(key string) domain.RBACRole {
	return domain.NewRBACRole(fmt.Sprintf("role_%s", key))
}

func NewRBACUserRoleObject(userRoleID domain.UserRoleID) domain.RBACObject {
	return domain.NewRBACObject(fmt.Sprintf("role_%d", userRoleID.Int()))
}

func NewRBACAllUserRoleObject() domain.RBACObject {
	return domain.NewRBACObject("role_*")
}

var RBACSetAction = domain.NewRBACAction("Set")
var RBACUnsetAction = domain.NewRBACAction("Unset")
