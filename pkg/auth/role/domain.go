package role

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Role  Model struct Role
type Role struct {
	ID          string    `json:"id" db:"id" valid:"required,uuid"`
	Name        string    `json:"name" db:"name" valid:"required"`
	Description string    `json:"description" db:"description" valid:"required"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func NewRole(id string, name string, description string) *Role {
	return &Role{
		ID:          id,
		Name:        name,
		Description: description,
	}
}

func (m *Role) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
