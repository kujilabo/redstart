package service

import (
	"context"

	liberrors "github.com/kujilabo/redstart/lib/errors"
	"github.com/kujilabo/redstart/user/domain"
)

// type Owner interface {
// 	// domain.OwnerModel
// 	AddAppUser(ctx context.Context, param AppUserAddParameter) (domain.AppUserID, error)
// }

type Owner struct {
	rf RepositoryFactory
	*domain.OwnerModel
}

func NewOwner(rf RepositoryFactory, ownerModel *domain.OwnerModel) *Owner {
	m := &Owner{
		rf:         rf,
		OwnerModel: ownerModel,
	}

	return m
}

func (m *Owner) AddAppUser(ctx context.Context, param AppUserAddParameter) (*domain.AppUserID, error) {
	appUserRepo := m.rf.NewAppUserRepository(ctx)
	appUserID, err := appUserRepo.AddAppUser(ctx, m, param)
	if err != nil {
		return nil, liberrors.Errorf("m.appUserRepo.AddAppUser. err: %w", err)
	}

	return appUserID, nil
}

func (m *Owner) AppUserID() *domain.AppUserID {
	return m.AppUserModel.AppUserID
}
func (m *Owner) OrganizationID() *domain.OrganizationID {
	return m.AppUserModel.OrganizationID
}
func (m *Owner) LoginID() string {
	return m.AppUserModel.LoginID
}
func (m *Owner) Username() string {
	return m.AppUserModel.Username
}
func (m *Owner) IsOwner() bool {
	return true
}
