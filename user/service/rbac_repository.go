package service

import (
	"github.com/casbin/casbin/v2"

	"github.com/kujilabo/redstart/user/domain"
)

type RBACRepository interface {
	Init() error

	AddPolicy(domain domain.RBACDomain, subject domain.RBACSubject, action domain.RBACAction, object domain.RBACObject, effect domain.RBACEffect) error

	AddSubjectGroupingPolicy(domain domain.RBACDomain, subject domain.RBACUser, object domain.RBACRole) error
	AddObjectGroupingPolicy(domain domain.RBACDomain, child domain.RBACObject, parent domain.RBACObject) error

	NewEnforcerWithGroupsAndUsers(roles []domain.RBACRole, users []domain.RBACUser) (*casbin.Enforcer, error)
}
