package service

import "context"

type RepositoryFactory interface {
	NewOrganizationRepository(ctx context.Context) OrganizationRepository
	NewAppUserRepository(ctx context.Context) AppUserRepository
	NewUserGroupRepository(ctx context.Context) UserGroupRepository

	// NewPairOfUserAndGroupRepository(ctx context.Context) PairOfUserAndGroupRepository

	// NewRBACRepository(ctx context.Context) RBACRepository

	NewAuthorizationManager(ctx context.Context) AuthorizationManager
}
