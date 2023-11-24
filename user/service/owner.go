package service

import (
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
