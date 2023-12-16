package service

import (
	"context"

	libdomain "github.com/kujilabo/redstart/lib/domain"
	liberrors "github.com/kujilabo/redstart/lib/errors"
	"github.com/kujilabo/redstart/user/domain"
)

type UserGroupAddParameter interface {
	GetKey() string
	GetName() string
	GetDescription() string
}

type userGroupAddParameter struct {
	Key         string
	Name        string
	Description string
}

func NewUserGroupAddParameter(key, name, description string) (UserGroupAddParameter, error) {
	m := &userGroupAddParameter{
		Key:         key,
		Name:        name,
		Description: description,
	}
	if err := libdomain.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (p *userGroupAddParameter) GetKey() string {
	return p.Key
}
func (p *userGroupAddParameter) GetName() string {
	return p.Name
}
func (p *userGroupAddParameter) GetDescription() string {
	return p.Description
}

type UserGroupRepository interface {
	FindAllUserGroups(ctx context.Context, operator AppUserModelInterface) ([]*domain.UserGroupModel, error)

	FindSystemOwnerGroup(ctx context.Context, operator SystemAdminModelInterface, organizationID domain.OrganizationID) (*UserGroup, error)

	FindUserGroupByKey(ctx context.Context, operator AppUserModelInterface, key string) (*UserGroup, error)
	FindUserGroupByID(ctx context.Context, operator AppUserModelInterface, userGroupID domain.UserGroupID) (*UserGroup, error)
	AddOwnerGroup(ctx context.Context, operator SystemOwnerModelInterface, organizationID domain.OrganizationID) (domain.UserGroupID, error)

	AddSystemOwnerGroup(ctx context.Context, operator SystemAdminModelInterface, organizationID domain.OrganizationID) (domain.UserGroupID, error)

	AddUserGroup(ctx context.Context, operator OwnerModelInterface, parameter UserGroupAddParameter) (domain.UserGroupID, error)
}
