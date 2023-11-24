package service

import (
	libdomain "github.com/kujilabo/redstart/lib/domain"
	liberrors "github.com/kujilabo/redstart/lib/errors"
	"github.com/kujilabo/redstart/user/domain"
)

type Organization interface {
	domain.OrganizationModel
}

type organization struct {
	domain.OrganizationModel
}

func NewOrganization(organizationModel domain.OrganizationModel) (Organization, error) {
	m := &organization{
		organizationModel,
	}

	if err := libdomain.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return m, nil
}
