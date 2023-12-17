//go:generate mockery --output mock --name OrganizationRepository
package service

import (
	"context"
	"errors"

	libdomain "github.com/kujilabo/redstart/lib/domain"
	liberrors "github.com/kujilabo/redstart/lib/errors"
	"github.com/kujilabo/redstart/user/domain"
)

type AppUserModelInterface interface {
	AppUserID() *domain.AppUserID
	OrganizationID() *domain.OrganizationID
	LoginID() string
	Username() string
	// GetUserGroups() []domain.UserGroupModel
}
type OwnerModelInterface interface {
	AppUserModelInterface
	IsOwner() bool
	// GetUserGroups() []domain.UserGroupModel
}
type SystemOwnerModelInterface interface {
	OwnerModelInterface
	IsSystemOwner() bool
	// GetUserGroups() []domain.UserGroupModel
}
type SystemAdminModelInterface interface {
	AppUserID() *domain.AppUserID
	IsSystemAdmin() bool
	// GetUserGroups() []domain.UserGroupModel
}

var ErrOrganizationNotFound = errors.New("organization not found")
var ErrOrganizationAlreadyExists = errors.New("organization already exists")

// type FirstOwnerAddParameter interface {
// 	GetLoginID() string
// 	GetUsername() string
// 	GetPassword() string
// }

// type firstOwnerAddParameter struct {
// 	LoginID  string `validate:"required"`
// 	Username string `validate:"required"`
// 	Password string `validate:"required"`
// }

// func NewFirstOwnerAddParameter(loginID, username, password string) (FirstOwnerAddParameter, error) {
// 	m := &firstOwnerAddParameter{
// 		LoginID:  loginID,
// 		Username: username,
// 		Password: password,
// 	}

// 	if err := libdomain.Validator.Struct(m); err != nil {
// 		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
// 	}

// 	return m, nil
// }

// func (p *firstOwnerAddParameter) GetLoginID() string {
// 	return p.LoginID
// }
// func (p *firstOwnerAddParameter) GetUsername() string {
// 	return p.Username
// }
// func (p *firstOwnerAddParameter) GetPassword() string {
// 	return p.Password
// }

type OrganizationAddParameter interface {
	GetName() string
	GetFirstOwner() AppUserAddParameterInterface
}

type organizationAddParameter struct {
	Name       string `validate:"required"`
	FirstOwner AppUserAddParameterInterface
}

func NewOrganizationAddParameter(name string, firstOwner AppUserAddParameterInterface) (OrganizationAddParameter, error) {
	m := &organizationAddParameter{
		Name:       name,
		FirstOwner: firstOwner,
	}
	if err := libdomain.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (p *organizationAddParameter) GetName() string {
	return p.Name
}
func (p *organizationAddParameter) GetFirstOwner() AppUserAddParameterInterface {
	return p.FirstOwner
}

type OrganizationRepository interface {
	GetOrganization(ctx context.Context, operator AppUserModelInterface) (*Organization, error)

	FindOrganizationByName(ctx context.Context, operator SystemAdminModelInterface, name string) (*Organization, error)

	FindOrganizationByID(ctx context.Context, operator SystemAdminModelInterface, id *domain.OrganizationID) (*Organization, error)

	AddOrganization(ctx context.Context, operator SystemAdminModelInterface, param OrganizationAddParameter) (*domain.OrganizationID, error)

	// FindOrganizationByName(ctx context.Context, operator SystemAdmin, name string) (Organization, error)
	// FindOrganization(ctx context.Context, operator AppUser) (Organization, error)
}
