package service

import (
	"context"

	libdomain "github.com/kujilabo/redstart/lib/domain"
	liberrors "github.com/kujilabo/redstart/lib/errors"
	"github.com/kujilabo/redstart/user/domain"
)

type UserRoleAddParameter interface {
	GetKey() string
	GetName() string
	GetDescription() string
}

type userRoleAddParameter struct {
	Key         string
	Name        string
	Description string
}

func NewUserRoleAddParameter(key, name, description string) (UserRoleAddParameter, error) {
	m := &userRoleAddParameter{
		Key:         key,
		Name:        name,
		Description: description,
	}
	if err := libdomain.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (p *userRoleAddParameter) GetKey() string {
	return p.Key
}
func (p *userRoleAddParameter) GetName() string {
	return p.Name
}
func (p *userRoleAddParameter) GetDescription() string {
	return p.Description
}

type UserRoleRepository interface {
	FindSystemOwnerRole(ctx context.Context, operator domain.SystemAdminModel, organizationID domain.OrganizationID) (UserRole, error)

	FindUserRoleByKey(ctx context.Context, operator domain.AppUserModel, key string) (UserRole, error)

	// AddUserRole(ctx context.Context, operator domain.AppUserModel, parameter UserRoleAddParameter) (domain.UserRoleID, error)

	AddOwnerRole(ctx context.Context, operator domain.SystemAdminModel, organizationID domain.OrganizationID) (domain.UserRoleID, error)

	AddSystemOwnerRole(ctx context.Context, operator domain.SystemAdminModel, organizationID domain.OrganizationID) (domain.UserRoleID, error)
	// AddPersonalGroup(operator SystemOwner, studentID uint) (uint, error)
}
