package domain

import (
	"strings"

	libdomain "github.com/kujilabo/redstart/lib/domain"
	liberrors "github.com/kujilabo/redstart/lib/errors"
)

type UserRoleID interface {
	Int() int
}

type userRoleID struct {
	Value int
}

func NewUserRoleID(value int) (UserRoleID, error) {
	return &userRoleID{
		Value: value,
	}, nil
}

func (v *userRoleID) Int() int {
	return v.Value
}

type UserRoleModel interface {
	libdomain.BaseModel
	GetUerRoleID() UserRoleID
	GetOrganizationID() OrganizationID
	GetKey() string
	GetName() string
	GetDescription() string
	IsSystemRole() bool
}

type userRoleModel struct {
	libdomain.BaseModel
	UserRoleID     UserRoleID
	OrganizationID OrganizationID
	Key            string `validate:"required"`
	Name           string `validate:"required"`
	Description    string
}

// NewUserRoleModel returns a new UserRoleModel
func NewUserRoleModel(baseModel libdomain.BaseModel, userRoleID UserRoleID, organizationID OrganizationID, key, name, description string) (UserRoleModel, error) {
	m := &userRoleModel{
		BaseModel:      baseModel,
		UserRoleID:     userRoleID,
		OrganizationID: organizationID,
		Key:            key,
		Name:           name,
		Description:    description,
	}

	if err := libdomain.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (m *userRoleModel) GetUerRoleID() UserRoleID {
	return m.UserRoleID
}

func (m *userRoleModel) GetOrganizationID() OrganizationID {
	return m.OrganizationID
}

func (m *userRoleModel) GetKey() string {
	return m.Key
}

func (m *userRoleModel) GetName() string {
	return m.Name
}

func (m *userRoleModel) GetDescription() string {
	return m.Description
}

func (m *userRoleModel) IsSystemRole() bool {
	return strings.HasPrefix(m.Key, "__")
}
