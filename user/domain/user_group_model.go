package domain

import (
	"strings"

	libdomain "github.com/kujilabo/redstart/lib/domain"
	liberrors "github.com/kujilabo/redstart/lib/errors"
)

type UserGroupID interface {
	Int() int
	IsUserGroupID() bool
}

type userGroupID struct {
	Value int
}

func NewUserGroupID(value int) (UserGroupID, error) {
	return &userGroupID{
		Value: value,
	}, nil
}

func (v *userGroupID) Int() int {
	return v.Value
}
func (v *userGroupID) IsUserGroupID() bool {
	return true
}

type UserGroupModel interface {
	libdomain.BaseModel
	GetUerGroupID() UserGroupID
	GetOrganizationID() OrganizationID
	GetKey() string
	GetName() string
	GetDescription() string
	IsSystemGroup() bool
}

type userGroupModel struct {
	libdomain.BaseModel
	UserGroupID    UserGroupID
	OrganizationID OrganizationID
	Key            string `validate:"required"`
	Name           string `validate:"required"`
	Description    string
}

// NewUserGroupModel returns a new UserGroupModel
func NewUserGroupModel(baseModel libdomain.BaseModel, userGroupID UserGroupID, organizationID OrganizationID, key, name, description string) (UserGroupModel, error) {
	m := &userGroupModel{
		BaseModel:      baseModel,
		UserGroupID:    userGroupID,
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

func (m *userGroupModel) GetUerGroupID() UserGroupID {
	return m.UserGroupID
}

func (m *userGroupModel) GetOrganizationID() OrganizationID {
	return m.OrganizationID
}

func (m *userGroupModel) GetKey() string {
	return m.Key
}

func (m *userGroupModel) GetName() string {
	return m.Name
}

func (m *userGroupModel) GetDescription() string {
	return m.Description
}

func (m *userGroupModel) IsSystemGroup() bool {
	return strings.HasPrefix(m.Key, "__")
}