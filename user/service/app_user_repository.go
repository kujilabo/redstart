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

type AppUserAddParameter interface {
	GetLoginID() string
	GetUsername() string
	GetPassword() string
	// GetRoles() []string
	// GetDetails() string
}

type appUserAddParameter struct {
	LoginID  string
	Username string
	Password string
	// Roles    []string
	// Details  string
}

func NewAppUserAddParameter(loginID, username, password string,

// , roles []string, details string
) (AppUserAddParameter, error) {
	m := &appUserAddParameter{
		LoginID:  loginID,
		Username: username,
		Password: password,
		// Roles:    roles,
		// Details:  details,
	}
	if err := libdomain.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (p *appUserAddParameter) GetLoginID() string {
	return p.LoginID
}
func (p *appUserAddParameter) GetUsername() string {
	return p.Username
}
func (p *appUserAddParameter) GetPassword() string {
	return p.Password
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
	FindSystemOwnerByOrganizationID(ctx context.Context, operator SystemAdminModelInterface, organizationID *domain.OrganizationID) (*SystemOwner, error)

	FindSystemOwnerByOrganizationName(ctx context.Context, operator SystemAdminModelInterface, organizationName string, options ...Option) (*SystemOwner, error)

	FindAppUserByID(ctx context.Context, operator AppUserModelInterface, id *domain.AppUserID, options ...Option) (*AppUser, error)

	FindAppUserByLoginID(ctx context.Context, operator AppUserModelInterface, loginID string) (*AppUser, error)

	FindOwnerByLoginID(ctx context.Context, operator SystemOwnerModelInterface, loginID string) (*Owner, error)

	AddAppUser(ctx context.Context, operator OwnerModelInterface, param AppUserAddParameter) (*domain.AppUserID, error)

	AddSystemOwner(ctx context.Context, operator SystemAdminModelInterface, organizationID *domain.OrganizationID) (*domain.AppUserID, error)

	// AddFirstOwner(ctx context.Context, operator domain.SystemOwnerModel, param FirstOwnerAddParameter) (domain.AppUserID, error)

	// FindAppUserIDs(ctx context.Context, operator domain.SystemOwnerModel, pageNo, pageSize int) ([]domain.AppUserID, error)
}
