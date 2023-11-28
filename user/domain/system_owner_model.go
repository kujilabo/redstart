package domain

import (
	libdomain "github.com/kujilabo/redstart/lib/domain"
	liberrors "github.com/kujilabo/redstart/lib/errors"
)

const SystemOwnerID = 2

type SystemOwnerModel interface {
	OwnerModel
	IsSystemOwnerModel() bool
}

type systemOwnerModel struct {
	OwnerModel
	AppUserID AppUserID
}

func NewSystemOwnerModel(appUser OwnerModel) (SystemOwnerModel, error) {
	m := &systemOwnerModel{
		OwnerModel: appUser,
	}

	if err := libdomain.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (m *systemOwnerModel) IsSystemOwnerModel() bool {
	return true
}
