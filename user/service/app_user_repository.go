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
}

type AppUserAddParameter struct {
	LoginIDInternal  string
	UsernameInternal string
	PasswordInternal string
}

func NewAppUserAddParameter(loginID, username, password string) (*AppUserAddParameter, error) {
	m := &AppUserAddParameter{
		LoginIDInternal:  loginID,
		UsernameInternal: username,
		PasswordInternal: password,
	}
	if err := libdomain.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (p *AppUserAddParameter) LoginID() string {
	return p.LoginIDInternal
}
func (p *AppUserAddParameter) Username() string {
	return p.UsernameInternal
}
func (p *AppUserAddParameter) Password() string {
	return p.PasswordInternal
}

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
