package validation_users

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// ValidationUsers  Model struct ValidationUsers
type ValidationUsers struct {
	ID            string    `json:"id" db:"id" valid:"required,uuid"`
	TransactionId string    `json:"transaction_id" db:"transaction_id" valid:"required"`
	UserId        string    `json:"user_id" db:"user_id" valid:"required"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

func NewValidationUsers(id string, transactionId string, userId string) *ValidationUsers {
	return &ValidationUsers{
		ID:            id,
		TransactionId: transactionId,
		UserId:        userId,
	}
}

func (m *ValidationUsers) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
