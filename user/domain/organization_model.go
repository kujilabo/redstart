package domain

import (
	libdomain "github.com/kujilabo/redstart/lib/domain"
	liberrors "github.com/kujilabo/redstart/lib/errors"
)

type OrganizationID struct {
	Value int `validate:"required,gte=1"`
}

func NewOrganizationID(value int) (*OrganizationID, error) {
	return &OrganizationID{
		Value: value,
	}, nil
}

func (v *OrganizationID) Int() int {
	return v.Value
}
func (v *OrganizationID) IsOrganizationID() bool {
	return true
}

type OrganizationModel struct {
	*libdomain.BaseModel
	OrganizationID *OrganizationID
	Name           string `validate:"required"`
}

func NewOrganizationModel(basemodel *libdomain.BaseModel, organizationID *OrganizationID, name string) (*OrganizationModel, error) {
	m := &OrganizationModel{
		BaseModel:      basemodel,
		OrganizationID: organizationID,
		Name:           name,
	}
	if err := libdomain.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return m, nil
}
