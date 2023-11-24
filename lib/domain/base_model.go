package domain

import (
	"time"

	liberrors "github.com/kujilabo/redstart/lib/errors"
)

type BaseModel interface {
	GetVersion() int
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	GetCreatedBy() int
	GetUpdatedBy() int
}

type baseModel struct {
	Version   int `validate:"required,gte=1"`
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy int `validate:"gte=0"`
	UpdatedBy int `validate:"gte=0"`
}

func NewBaseModel(version int, createdAt, updatedAt time.Time, createdBy, updatedBy int) (BaseModel, error) {
	m := &baseModel{
		Version:   version,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		CreatedBy: createdBy,
		UpdatedBy: updatedBy,
	}

	if err := Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (m *baseModel) GetVersion() int {
	return m.Version
}

func (m *baseModel) GetCreatedAt() time.Time {
	return m.CreatedAt
}

func (m *baseModel) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}

func (m *baseModel) GetCreatedBy() int {
	return m.CreatedBy
}

func (m *baseModel) GetUpdatedBy() int {
	return m.UpdatedBy
}
