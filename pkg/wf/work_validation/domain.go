package work_validation

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// WorkValidation  Model struct WorkValidation
type WorkValidation struct {
	ID        int64     `json:"id" db:"id" valid:"-"`
	Status    string    `json:"status" db:"status" valid:"required"`
	UserId    string    `json:"user_id" db:"user_id" valid:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewWorkValidation(id int64, status string, userId string) *WorkValidation {
	return &WorkValidation{
		ID:     id,
		Status: status,
		UserId: userId,
	}
}

func NewCreateWorkValidation(status string, userId string) *WorkValidation {
	return &WorkValidation{
		Status: status,
		UserId: userId,
	}
}

func (m *WorkValidation) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
