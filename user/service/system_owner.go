package service

import (
	"context"

	libdomain "github.com/kujilabo/redstart/lib/domain"
	liberrors "github.com/kujilabo/redstart/lib/errors"
	liblog "github.com/kujilabo/redstart/lib/log"
	"github.com/kujilabo/redstart/user/domain"
)

type SystemOwner interface {
	domain.SystemOwnerModel

	GetOrganization(ctxc context.Context) (Organization, error)

	FindAppUserByID(ctx context.Context, id domain.AppUserID) (AppUser, error)

	FindAppUserByLoginID(ctx context.Context, loginID string) (AppUser, error)

	AddAppUser(ctx context.Context, param AppUserAddParameter) (domain.AppUserID, error)
}

type systemOwner struct {
	domain.SystemOwnerModel
	orgRepo                OrganizationRepository
	userGroupRepo          UserGroupRepository
	appUserRepo            AppUserRepository
	pairOfUserAndGroupRepo PairOfUserAndGroupRepository
	rbacRepo               RBACRepository
}

func NewSystemOwner(ctx context.Context, rf RepositoryFactory, systemOwnerModel domain.SystemOwnerModel) (SystemOwner, error) {
	orgRepo := rf.NewOrganizationRepository(ctx)
	appUserRepo := rf.NewAppUserRepository(ctx)
	userGroupRepo := rf.NewUserGroupRepository(ctx)
	pairOfUserAndGroupRepo := rf.NewPairOfUserAndGroupRepository(ctx)
	rbacRepo := rf.NewRBACRepository(ctx)

	m := &systemOwner{
		SystemOwnerModel:       systemOwnerModel,
		orgRepo:                orgRepo,
		userGroupRepo:          userGroupRepo,
		appUserRepo:            appUserRepo,
		pairOfUserAndGroupRepo: pairOfUserAndGroupRepo,
		rbacRepo:               rbacRepo,
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

func (m *systemOwner) AddAppUser(ctx context.Context, param AppUserAddParameter) (domain.AppUserID, error) {
	logger := liblog.GetLoggerFromContext(ctx, UserServiceContextKey)
	logger.InfoContext(ctx, "AddStudent")
	appUserID, err := m.appUserRepo.AddAppUser(ctx, m, param)
	if err != nil {
		return nil, liberrors.Errorf("m.appUserRepo.AddAppUser. err: %w", err)
	}
	appUser, err := m.appUserRepo.FindAppUserByID(ctx, m, appUserID)
	if err != nil {
		return nil, liberrors.Errorf("m.appUserRepo.FindAppUserByID. err: %w", err)
	}

	// personalGroupID, err := m.rf.NewAppUserGroupRepository().AddPersonalGroup(m, studentID)
	// if err != nil {
	// 	return 0, err
	// }

	publicGroup, err := m.userGroupRepo.FindPublicGroup(ctx, m)
	if err != nil {
		return nil, liberrors.Errorf("m.userGroupRepo.FindPublicGroup. err: %w", err)
	}
	if err := m.pairOfUserAndGroupRepo.AddPairOfUserAndGroup(ctx, m, publicGroup.GetUerGroupID(), appUser.GetAppUserID()); err != nil {
		return nil, liberrors.Errorf("m.groupUserRepo.AddGroupUser. err: %w", err)
	}

	// spaceID, err := m.spaceRepo.AddPersonalSpace(ctx, m, appUser)
	// if err != nil {
	// 	return 0, liberrors.Errorf("m.spaceRepo.AddPersonalSpace. err: %w", err)
	// }

	// logger.Infof("Personal spaceID: %d", spaceID)

	// spaceWriter := domain.NewSpaceWriterRole(spaceID)
	// spaceObject := domain.NewSpaceObject(spaceID)
	// userSubject := domain.NewUserObject(appUserID)

	// if err := m.rbacRepo.AddNamedPolicy(spaceWriter, spaceObject, "read"); err != nil {
	// 	return 0, liberrors.Errorf("problemRepo.AddNamedPolicy(read). err: %w", err)
	// }

	// if err := m.rbacRepo.AddNamedPolicy(spaceWriter, spaceObject, "write"); err != nil {
	// 	return 0, liberrors.Errorf("problemRepo.AddNamedPolicy(write). err: %w", err)
	// }

	// if err := m.rbacRepo.AddNamedGroupingPolicy(userSubject, spaceWriter); err != nil {
	// 	return 0, liberrors.Errorf("problemRepo.AddNamedGroupingPolicy. err: %w", err)
	// }

	// defaultSpace, err := m.rf.NewSpaceRepository().FindDefaultSpace(ctx, s)
	// if err != nil {
	// 	return 0, err
	// }

	// if err := m.rf.NewUserSpaceRepository().Add(ctx, appUser, SpaceID(defaultSpace.GetID())); err != nil {
	// 	return 0, err
	// }

	return appUserID, nil
}

// func (m *systemOwner) AddSystemSpace(ctx context.Context) (domain.SpaceID, error) {
// 	logger := liblog.GetLoggerFromContext(ctx, UserServiceContextKey)
// 	logger.Infof("AddSystemSpace")

// 	spaceID, err := m.spaceRepo.AddSystemSpace(ctx, m)
// 	if err != nil {
// 		return 0, liberrors.Errorf("m.spaceRepo.AddSystemSpace. err: %w", err)
// 	}
// 	return spaceID, nil
// }
