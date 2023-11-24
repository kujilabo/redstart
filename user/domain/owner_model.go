package domain

type OwnerModel interface {
	AppUserModel
	IsOwnerModel() bool
}

type ownerModel struct {
	AppUserModel
}

func NewOwnerModel(appUser AppUserModel) (OwnerModel, error) {
	return &ownerModel{
		AppUserModel: appUser,
	}, nil
}

func (m *ownerModel) IsOwnerModel() bool {
	return true
}
