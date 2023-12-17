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

type OrganizationAddParameterInterface interface {
	Name() string
	FirstOwner() AppUserAddParameterInterface
}

type OrganizationAddParameter struct {
	Name_       string `validate:"required"`
	FirstOwner_ AppUserAddParameterInterface
}

func NewOrganizationAddParameter(name string, firstOwner AppUserAddParameterInterface) (*OrganizationAddParameter, error) {
	m := &OrganizationAddParameter{
		Name_:       name,
		FirstOwner_: firstOwner,
	}
	if err := libdomain.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (p *OrganizationAddParameter) Name() string {
	return p.Name_
}
func (p *OrganizationAddParameter) FirstOwner() AppUserAddParameterInterface {
	return p.FirstOwner_
}

type OrganizationRepository interface {
	GetOrganization(ctx context.Context, operator AppUserModelInterface) (*Organization, error)

	FindOrganizationByName(ctx context.Context, operator SystemAdminModelInterface, name string) (*Organization, error)

	FindOrganizationByID(ctx context.Context, operator SystemAdminModelInterface, id *domain.OrganizationID) (*Organization, error)

	AddOrganization(ctx context.Context, operator SystemAdminModelInterface, param OrganizationAddParameterInterface) (*domain.OrganizationID, error)

	// FindOrganizationByName(ctx context.Context, operator SystemAdmin, name string) (Organization, error)
	// FindOrganization(ctx context.Context, operator AppUser) (Organization, error)
}
