package service

import (
	"context"

	libdomain "github.com/kujilabo/redstart/lib/domain"
	liberrors "github.com/kujilabo/redstart/lib/errors"
	"github.com/kujilabo/redstart/user/domain"
)

// type AppUser interface {
// 	// domain.AppUserModel
// }

type AppUser struct {
	*domain.AppUserModel
}

func NewAppUser(ctx context.Context, rf RepositoryFactory, appUserModel *domain.AppUserModel) (*AppUser, error) {
	if rf == nil {
		return nil, liberrors.Errorf("rf is nil. err: %w", libdomain.ErrInvalidArgument)
	}
	if appUserModel == nil {
		return nil, liberrors.Errorf("appUserModel is nil. err: %w", libdomain.ErrInvalidArgument)
	}

	m := &AppUser{
		AppUserModel: appUserModel,
	}

	if err := libdomain.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (m *AppUser) AppUserID() *domain.AppUserID {
	return m.AppUserModel.AppUserID
}
func (m *AppUser) OrganizationID() *domain.OrganizationID {
	return m.AppUserModel.OrganizationID
}
func (m *AppUser) LoginID() string {
	return m.AppUserModel.LoginID
}
func (m *AppUser) Username() string {
	return m.AppUserModel.Username
}
