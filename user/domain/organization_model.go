package domain

import (
	libdomain "github.com/kujilabo/redstart/lib/domain"
	liberrors "github.com/kujilabo/redstart/lib/errors"
)

type OrganizationID interface {
	Int() int
}

type organizationID struct {
	Value int `validate:"required,gte=1"`
}

func NewOrganizationID(value int) (OrganizationID, error) {
	return &organizationID{
		Value: value,
	}, nil
}

func (v *organizationID) Int() int {
	return v.Value
}

type OrganizationModel interface {
	libdomain.BaseModel
	GetID() OrganizationID
	GetName() string
}

type organizationModel struct {
	libdomain.BaseModel
	OrganizationID OrganizationID
	Name           string `validate:"required"`
}

func NewOrganizationModel(basemodel libdomain.BaseModel, organizationID OrganizationID, name string) (OrganizationModel, error) {
	m := &organizationModel{
		BaseModel:      basemodel,
		OrganizationID: organizationID,
		Name:           name,
	}
	if err := libdomain.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (m *organizationModel) GetID() OrganizationID {
	return m.OrganizationID
}

func (m *organizationModel) GetName() string {
	return m.Name
}
