package domain

// type SystemAdminModel interface {
// 	GetAppUserID() AppUserID
// 	IsSystemAdminModel() bool
// }

type SystemAdminModel struct {
	AppUserID *AppUserID
}

func NewSystemAdminModel() *SystemAdminModel {
	return &SystemAdminModel{
		AppUserID: SystemAdminID,
	}
}

// func (s *SystemAdminModel) AppUserID() *AppUserID {
// 	return SystemAdminID
// }

// func (s *systemAdminModel) IsSystemAdminModel() bool {
// 	return true
// }
