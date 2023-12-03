package service

import (
	"context"

	liberrors "github.com/kujilabo/redstart/lib/errors"
	"github.com/kujilabo/redstart/user/domain"
)

type Owner interface {
	domain.OwnerModel
}

type owner struct {
	rf RepositoryFactory
	domain.OwnerModel
}

func NewOwner(rf RepositoryFactory, ownerModel domain.OwnerModel) Owner {
	return &owner{
		rf:         rf,
		OwnerModel: ownerModel,
	}
}

func (m *owner) AddAppUser(ctx context.Context, param AppUserAddParameter) (domain.AppUserID, error) {
	appUserRepo := m.rf.NewAppUserRepository(ctx)
	appUserID, err := appUserRepo.AddAppUser(ctx, m, param)
	if err != nil {
		return nil, liberrors.Errorf("m.appUserRepo.AddAppUser. err: %w", err)
	}

	return appUserID, nil
}
