//go:generate mockery --output mock --name RBACRepository
package service

import (
	"github.com/casbin/casbin/v2"

	"github.com/kujilabo/redstart/user/domain"
)

type RBACRepository interface {
	Init() error

	AddNamedPolicy(subject domain.RBACRole, object domain.RBACObject, action domain.RBACAction) error

	AddNamedGroupingPolicy(subject domain.RBACUser, object domain.RBACRole) error

	NewEnforcerWithRolesAndUsers(roles []domain.RBACRole, users []domain.RBACUser) (*casbin.Enforcer, error)
}
