//go:generate mockery --output mock --name AppUserModel
package domain

import (
	libdomain "github.com/kujilabo/redstart/lib/domain"
	liberrors "github.com/kujilabo/redstart/lib/errors"
)

// type AppUserID interface {
// 	Int() int
// 	IsAppUserID() bool
// }

type AppUserID struct {
	Value int `validate:"required,gte=0"`
}

func NewAppUserID(value int) (*AppUserID, error) {
	return &AppUserID{
		Value: value,
	}, nil
}

func (v *AppUserID) Int() int {
	return v.Value
}
func (v *AppUserID) IsAppUserID() bool {
	return true
}

// type AppUserModel interface {
// 	libdomain.BaseModel
// 	GetAppUserID() AppUserID
// 	GetOrganizationID() OrganizationID
// 	GetLoginID() string
// 	GetUsername() string
// 	GetUserGroups() []UserGroupModel
// }

type AppUserModel struct {
	libdomain.BaseModel
	AppUserID      *AppUserID
	OrganizationID *OrganizationID
	LoginID        string `validate:"required"`
	Username       string `validate:"required"`
	UserGroups     []*UserGroupModel
}

func NewAppUserModel(baseModel libdomain.BaseModel, appUserID *AppUserID, organizationID *OrganizationID, loginID, username string, userGroups []*UserGroupModel) (*AppUserModel, error) {
	m := &AppUserModel{
		BaseModel:      baseModel,
		AppUserID:      appUserID,
		OrganizationID: organizationID,
		LoginID:        loginID,
		Username:       username,
		UserGroups:     userGroups,
	}

	if err := libdomain.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return m, nil
}

// func (m *appUserModel) GetAppUserID() AppUserID {
// 	return m.AppUserID
// }

// func (m *appUserModel) GetOrganizationID() OrganizationID {
// 	return m.OrganizationID
// }

// func (m *appUserModel) GetLoginID() string {
// 	return m.LoginID
// }

// func (m *appUserModel) GetUsername() string {
// 	return m.Username
// }

// func (m *appUserModel) GetUserGroups() []UserGroupModel {
// 	return m.UserGroups
// }
