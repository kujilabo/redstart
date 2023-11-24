package service

import "context"

type RepositoryFactory interface {
	NewOrganizationRepository(ctx context.Context) OrganizationRepository
	NewAppUserRepository(ctx context.Context) AppUserRepository
	NewUserRoleRepository(ctx context.Context) UserRoleRepository
	NewUserGroupRepository(ctx context.Context) UserGroupRepository

	NewPairOfUserAndGroupRepository(ctx context.Context) PairOfUserAndGroupRepository
	NewPairOfUserAndRoleRepository(ctx context.Context) PairOfUserAndRoleRepository

	NewRBACRepository(ctx context.Context) RBACRepository
}
