package service

import (
	"context"

	libdomain "github.com/kujilabo/redstart/lib/domain"
	liberrors "github.com/kujilabo/redstart/lib/errors"
	"github.com/kujilabo/redstart/user/domain"
)

type SystemOwner interface {
	domain.SystemOwnerModel

	GetOrganization(ctxc context.Context) (Organization, error)

	FindAppUserByID(ctx context.Context, id domain.AppUserID) (AppUser, error)

	FindAppUserByLoginID(ctx context.Context, loginID string) (AppUser, error)

	// AddAppUser(ctx context.Context, param AppUserAddParameter) (domain.AppUserID, error)

	AddFirstOwner(ctx context.Context, param FirstOwnerAddParameter) (domain.AppUserID, error)
}

type systemOwner struct {
	domain.SystemOwnerModel
	orgRepo       OrganizationRepository
	appUserRepo   AppUserRepository
	userGroupRepo UserGroupRepository
	// pairOfUserAndGroup PairOfUserAndGroupRepository
	// rbacRepo             RBACRepository
	authorizationManager AuthorizationManager
}

func NewSystemOwner(ctx context.Context, rf RepositoryFactory, systemOwnerModel domain.SystemOwnerModel) (SystemOwner, error) {
	orgRepo := rf.NewOrganizationRepository(ctx)
	appUserRepo := rf.NewAppUserRepository(ctx)
	userGroupRepo := rf.NewUserGroupRepository(ctx)
	// pairOfUserAndGroup := rf.NewPairOfUserAndGroupRepository(ctx)
	// rbacRepo := rf.NewRBACRepository(ctx)
	authorizationManager := rf.NewAuthorizationManager(ctx)

	m := &systemOwner{
		SystemOwnerModel: systemOwnerModel,
		orgRepo:          orgRepo,
		appUserRepo:      appUserRepo,
		userGroupRepo:    userGroupRepo,
		// pairOfUserAndGroup:   pairOfUserAndGroup,
		// rbacRepo:             rbacRepo,
		authorizationManager: authorizationManager,
	}

	if err := libdomain.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (m *systemOwner) GetOrganization(ctx context.Context) (Organization, error) {
	org, err := m.orgRepo.GetOrganization(ctx, m)
	if err != nil {
		return nil, liberrors.Errorf("m.orgRepo.GetOrganization. err: %w", err)
	}

	return org, nil
}

func (m *systemOwner) FindAppUserByID(ctx context.Context, id domain.AppUserID) (AppUser, error) {
	appUser, err := m.appUserRepo.FindAppUserByID(ctx, m, id)
	if err != nil {
		return nil, liberrors.Errorf("m.appUserRepo.FindAppUserByID. err: %w", err)
	}

	return appUser, nil
}

func (m *systemOwner) FindAppUserByLoginID(ctx context.Context, loginID string) (AppUser, error) {
	appUser, err := m.appUserRepo.FindAppUserByLoginID(ctx, m, loginID)
	if err != nil {
		return nil, liberrors.Errorf("m.appUserRepo.FindAppUserByLoginID. err: %w", err)
	}

	return appUser, nil
}

// func (m *systemOwner) AddAppUser(ctx context.Context, param AppUserAddParameter) (domain.AppUserID, error) {
// 	logger := liblog.GetLoggerFromContext(ctx, UserServiceContextKey)
// 	logger.InfoContext(ctx, "AddStudent")
// 	appUserID, err := m.appUserRepo.AddAppUser(ctx, m, param)
// 	if err != nil {
// 		return nil, liberrors.Errorf("m.appUserRepo.AddAppUser. err: %w", err)
// 	}

// 	return appUserID, nil
// }

func (m *systemOwner) AddFirstOwner(ctx context.Context, param FirstOwnerAddParameter) (domain.AppUserID, error) {
	// rbacAppUser := NewRBACAppUser(m.GetOrganizationID(), m.GetAppUserID())
	rbacAllUserRolesObject := NewRBACAllUserRolesObject(m.GetOrganizationID())

	// Can "the operator" "set" "all-user-roles" ?
	ok, err := m.authorizationManager.Authorize(ctx, m, RBACSetAction, rbacAllUserRolesObject)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, libdomain.ErrPermissionDenied
	}

	// add owner
	firstOwnerID, err := m.appUserRepo.AddAppUser(ctx, m, param)
	if err != nil {
		return nil, liberrors.Errorf("failed to AddFirstOwner. error: %w", err)
	}

	ownerGroup, err := m.userGroupRepo.FindUserGroupByKey(ctx, m, OwnerGroupKey)
	if err != nil {
		return nil, err
	}

	// add owner to owner-group
	if err := m.authorizationManager.AddUserToGroup(ctx, m, firstOwnerID, ownerGroup.GetUerGroupID()); err != nil {
		return nil, err
	}

	// add owner to owner-group
	// if err := m.pairOfUserAndGroup.AddPairOfUserAndGroup(ctx, m, ownerID, ownerGroup.GetUerGroupID()); err != nil {
	// 	return nil, err
	// }

	// rbacDomain := NewRBACOrganization(m.GetOrganizationID())

	// // "owner" "can" "set" "all-user-roles"
	// if err := m.rbacRepo.AddPolicy(rbacDomain, rbacAppUser, RBACSetAction, rbacAllUserRolesObject, RBACAllowEffect); err != nil {
	// 	return nil, liberrors.Errorf("Failed to AddNamedPolicy. priv: read, err: %w", err)
	// }

	// // "owner" "can" "unset" "all-user-roles"
	// if err := m.rbacRepo.AddPolicy(rbacDomain, rbacAppUser, RBACUnsetAction, rbacAllUserRolesObject, RBACAllowEffect); err != nil {
	// 	return nil, liberrors.Errorf("Failed to AddNamedPolicy. priv: read, err: %w", err)
	// }

	return firstOwnerID, nil
}
