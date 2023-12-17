package service

import (
	"context"

	"github.com/kujilabo/redstart/user/domain"
)

type AuthorizationManager interface {
	Init(ctx context.Context) error

	AddUserToGroup(ctx context.Context, operator AppUserModelInterface, appUserID *domain.AppUserID, userGroupID *domain.UserGroupID) error

	AddUserToGroupBySystemAdmin(ctx context.Context, operator SystemAdminModelInterface, organizationID *domain.OrganizationID, appUserID *domain.AppUserID, userGroupID *domain.UserGroupID) error

	// RemoveUserFromGroup()

	// AddGroupToGroup(ctx context.Context, operator domain.AppUserModel, src domain.UserGroupID, dst domain.UserGroupID) error

	// RemoveGroupFromGroup()

	// AddObjectToObject()

	// RemoveObjectFromObject()

	AddPolicyToUser(ctx context.Context, operator AppUserModelInterface, subject domain.RBACSubject, action domain.RBACAction, object domain.RBACObject, effect domain.RBACEffect) error

	AddPolicyToUserBySystemAdmin(ctx context.Context, operator SystemAdminModelInterface, organizationID *domain.OrganizationID, subject domain.RBACSubject, action domain.RBACAction, object domain.RBACObject, effect domain.RBACEffect) error

	AddPolicyToGroup(ctx context.Context, operator AppUserModelInterface, subject domain.RBACSubject, action domain.RBACAction, object domain.RBACObject, effect domain.RBACEffect) error

	AddPolicyToGroupBySystemAdmin(ctx context.Context, operator SystemAdminModelInterface, organizationID *domain.OrganizationID, subject domain.RBACSubject, action domain.RBACAction, object domain.RBACObject, effect domain.RBACEffect) error

	// AddPolicyToGroup()

	// RemovePolicyToGroup()

	Authorize(ctx context.Context, operator AppUserModelInterface, rbacAction domain.RBACAction, rbacObject domain.RBACObject) (bool, error)
}
