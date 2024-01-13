package service

import (
	"context"

	"github.com/casbin/casbin/v2"

	"github.com/kujilabo/redstart/user/domain"
)

type RBACRepository interface {
	Init() error

	AddPolicy(ctx context.Context, domain domain.RBACDomain, subject domain.RBACSubject, action domain.RBACAction, object domain.RBACObject, effect domain.RBACEffect) error

	AddSubjectGroupingPolicy(ctx context.Context, domain domain.RBACDomain, subject domain.RBACUser, object domain.RBACRole) error
	AddObjectGroupingPolicy(ctx context.Context, domain domain.RBACDomain, child domain.RBACObject, parent domain.RBACObject) error

	RemovePolicy(ctx context.Context, domain domain.RBACDomain, subject domain.RBACSubject, action domain.RBACAction, object domain.RBACObject, effect domain.RBACEffect) error
	// RemoveSubjectPolicy(domain domain.RBACDomain, subject domain.RBACSubject) error

	RemoveSubjectGroupingPolicy(ctx context.Context, domain domain.RBACDomain, subject domain.RBACUser, object domain.RBACRole) error
	RemoveObjectGroupingPolicy(ctx context.Context, domain domain.RBACDomain, child domain.RBACObject, parent domain.RBACObject) error

	NewEnforcerWithGroupsAndUsers(ctx context.Context, roles []domain.RBACRole, users []domain.RBACUser) (*casbin.Enforcer, error)
}
