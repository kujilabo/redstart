package domain

type SystemAdminModel interface {
	GetAppUserID() AppUserID
	IsSystemAdminModel() bool
}

type systemAdminModel struct {
}

func NewSystemAdminModel() SystemAdminModel {
	return &systemAdminModel{}
}

func (s *systemAdminModel) GetAppUserID() AppUserID {
	return SystemAdminID
}

func (s *systemAdminModel) IsSystemAdminModel() bool {
	return true
}
