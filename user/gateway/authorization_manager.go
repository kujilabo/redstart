package gateway

import (
	"context"

	liberrors "github.com/kujilabo/redstart/lib/errors"
	"github.com/kujilabo/redstart/user/domain"
	"github.com/kujilabo/redstart/user/service"
	"gorm.io/gorm"
)

type authorizationManager struct {
	db *gorm.DB
	rf service.RepositoryFactory
}

func NewAuthorizationManager(ctx context.Context, db *gorm.DB, rf service.RepositoryFactory) service.AuthorizationManager {
	return &authorizationManager{
		db: db,
		rf: rf,
	}
}

func (m *authorizationManager) AddUserToGroupBySystemAdmin(ctx context.Context, operator service.SystemAdminModelInterface, organizationID domain.OrganizationID, appUserID domain.AppUserID, userGroupID domain.UserGroupID) error {
	pairOfUserAndGroupRepo := NewPairOfUserAndGroupRepository(ctx, m.db, m.rf)

	if err := pairOfUserAndGroupRepo.AddPairOfUserAndGroupBySystemAdmin(ctx, operator, organizationID, appUserID, userGroupID); err != nil {
		return err
	}

	rbacRepo := newRBACRepository(ctx, m.db)
	rbacAppUser := service.NewRBACAppUser(organizationID, appUserID)
	rbacUserRole := service.NewRBACUserRole(organizationID, userGroupID)
	rbacDomain := service.NewRBACOrganization(organizationID)

	// app-user belongs to user-role
	if err := rbacRepo.AddSubjectGroupingPolicy(rbacDomain, rbacAppUser, rbacUserRole); err != nil {
		return liberrors.Errorf("rbacRepo.AddNamedGroupingPolicy. err: %w", err)
	}

	return nil
}
func (m *authorizationManager) AddUserToGroup(ctx context.Context, operator service.AppUserModelInterface, appUserID domain.AppUserID, userGroupID domain.UserGroupID) error {
	pairOfUserAndGroupRepo := NewPairOfUserAndGroupRepository(ctx, m.db, m.rf)

	if err := pairOfUserAndGroupRepo.AddPairOfUserAndGroup(ctx, operator, appUserID, userGroupID); err != nil {
		return err
	}

	organizationID := operator.OrganizationID()

	rbacRepo := newRBACRepository(ctx, m.db)
	rbacAppUser := service.NewRBACAppUser(organizationID, appUserID)
	rbacUserRole := service.NewRBACUserRole(organizationID, userGroupID)
	rbacDomain := service.NewRBACOrganization(organizationID)

	// app-user belongs to user-role
	if err := rbacRepo.AddSubjectGroupingPolicy(rbacDomain, rbacAppUser, rbacUserRole); err != nil {
		return liberrors.Errorf("rbacRepo.AddNamedGroupingPolicy. err: %w", err)
	}

	return nil
}

func (m *authorizationManager) AddPolicyToUser(ctx context.Context, operator service.AppUserModelInterface, subject domain.RBACSubject, action domain.RBACAction, object domain.RBACObject, effect domain.RBACEffect) error {
	rbacRepo := newRBACRepository(ctx, m.db)
	rbacDomain := service.NewRBACOrganization(operator.OrganizationID())

	if err := rbacRepo.AddPolicy(rbacDomain, subject, action, object, effect); err != nil {
		return liberrors.Errorf("Failed to AddNamedPolicy. priv: read, err: %w", err)
	}

	return nil
}

func (m *authorizationManager) AddPolicyToUserBySystemAdmin(ctx context.Context, operator service.SystemAdminModelInterface, organizationID domain.OrganizationID, subject domain.RBACSubject, action domain.RBACAction, object domain.RBACObject, effect domain.RBACEffect) error {
	rbacRepo := newRBACRepository(ctx, m.db)
	rbacDomain := service.NewRBACOrganization(organizationID)

	if err := rbacRepo.AddPolicy(rbacDomain, subject, action, object, effect); err != nil {
		return liberrors.Errorf("Failed to AddNamedPolicy. priv: read, err: %w", err)
	}

	return nil
}

func (m *authorizationManager) AddPolicyToGroup(ctx context.Context, operator service.AppUserModelInterface, subject domain.RBACSubject, action domain.RBACAction, object domain.RBACObject, effect domain.RBACEffect) error {
	rbacRepo := newRBACRepository(ctx, m.db)
	rbacDomain := service.NewRBACOrganization(operator.OrganizationID())

	if err := rbacRepo.AddPolicy(rbacDomain, subject, action, object, effect); err != nil {
		return liberrors.Errorf("Failed to AddNamedPolicy. priv: read, err: %w", err)
	}

	return nil
}

func (m *authorizationManager) AddPolicyToGroupBySystemAdmin(ctx context.Context, operator service.SystemAdminModelInterface, organizationID domain.OrganizationID, subject domain.RBACSubject, action domain.RBACAction, object domain.RBACObject, effect domain.RBACEffect) error {
	rbacRepo := newRBACRepository(ctx, m.db)
	rbacDomain := service.NewRBACOrganization(organizationID)

	if err := rbacRepo.AddPolicy(rbacDomain, subject, action, object, effect); err != nil {
		return liberrors.Errorf("Failed to AddNamedPolicy. priv: read, err: %w", err)
	}

	return nil
}

func (m *authorizationManager) Authorize(ctx context.Context, operator service.AppUserModelInterface, rbacAction domain.RBACAction, rbacObject domain.RBACObject) (bool, error) {
	rbacDomain := service.NewRBACOrganization(operator.OrganizationID())

	userGroupRepo := m.rf.NewUserGroupRepository(ctx)
	userGroups, err := userGroupRepo.FindAllUserGroups(ctx, operator)
	if err != nil {
		return false, err
	}

	rbacRoles := make([]domain.RBACRole, 0)
	for _, userGroup := range userGroups {
		rbacRoles = append(rbacRoles, service.NewRBACUserRole(operator.OrganizationID(), userGroup.GetUerGroupID()))
	}

	rbacRepo := newRBACRepository(ctx, m.db)
	rbacOperator := service.NewRBACAppUser(operator.OrganizationID(), operator.AppUserID())
	e, err := rbacRepo.NewEnforcerWithGroupsAndUsers(rbacRoles, []domain.RBACUser{rbacOperator})
	if err != nil {
		return false, err
	}

	ok, err := e.Enforce(rbacOperator.Subject(), rbacObject.Object(), rbacAction.Action(), rbacDomain.Domain())
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}

	return false, nil
}
