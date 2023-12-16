package domain

// type SystemAdminModel interface {
// 	GetAppUserID() AppUserID
// 	IsSystemAdminModel() bool
// }

type SystemAdminModel struct {
}

func NewSystemAdminModel() *SystemAdminModel {
	return &SystemAdminModel{}
}

func (s *SystemAdminModel) AppUserID() AppUserID {
	return SystemAdminID
}

// func (s *systemAdminModel) IsSystemAdminModel() bool {
// 	return true
// }
