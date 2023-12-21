package service

import (
	"context"
	"errors"

	libdomain "github.com/kujilabo/redstart/lib/domain"
	liberrors "github.com/kujilabo/redstart/lib/errors"
	"github.com/kujilabo/redstart/user/domain"
)

var ErrAppUserNotFound = errors.New("AppUser not found")
var ErrAppUserAlreadyExists = errors.New("AppUser already exists")

var ErrSystemOwnerNotFound = errors.New("SystemOwner not found")

type AppUserAddParameterInterface interface {
	LoginID() string
	Username() string
	Password() string
	// GetRoles() []string
	// GetDetails() string
}

type AppUserAddParameter struct {
	LoginID_  string
	Username_ string
	Password_ string
	// Roles    []string
	// Details  string
}

func NewAppUserAddParameter(loginID, username, password string) (*AppUserAddParameter, error) {
	m := &AppUserAddParameter{
		LoginID_:  loginID,
		Username_: username,
		Password_: password,
		// Roles:    roles,
		// Details:  details,
	}
	if err := libdomain.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (p *AppUserAddParameter) LoginID() string {
	return p.LoginID_
}
func (p *AppUserAddParameter) Username() string {
	return p.Username_
}
func (p *AppUserAddParameter) Password() string {
	return p.Password_
}

// func (p *appUserAddParameter) GetRoles() []string {
// 	return p.Roles
// }
// func (p *appUserAddParameter) GetDetails() string {
// 	return p.Details
// }

type Option string

var IncludeGroups Option = "IncludeGroups"

type AppUserRepository interface {
	FindSystemOwnerByOrganizationID(ctx context.Context, operator SystemAdminInterface, organizationID *domain.OrganizationID) (*SystemOwner, error)

	FindSystemOwnerByOrganizationName(ctx context.Context, operator SystemAdminInterface, organizationName string, options ...Option) (*SystemOwner, error)

	FindAppUserByID(ctx context.Context, operator AppUserInterface, id *domain.AppUserID, options ...Option) (*AppUser, error)

	FindAppUserByLoginID(ctx context.Context, operator AppUserInterface, loginID string) (*AppUser, error)

	FindOwnerByLoginID(ctx context.Context, operator SystemOwnerInterface, loginID string) (*Owner, error)

	AddAppUser(ctx context.Context, operator OwnerModelInterface, param AppUserAddParameterInterface) (*domain.AppUserID, error)

	AddSystemOwner(ctx context.Context, operator SystemAdminInterface, organizationID *domain.OrganizationID) (*domain.AppUserID, error)

	// AddFirstOwner(ctx context.Context, operator domain.SystemOwnerModel, param FirstOwnerAddParameter) (domain.AppUserID, error)

	// FindAppUserIDs(ctx context.Context, operator domain.SystemOwnerModel, pageNo, pageSize int) ([]domain.AppUserID, error)
}
