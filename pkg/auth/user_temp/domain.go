package user_temp

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// UserTemp  Model struct UserTemp
type UserTemp struct {
	ID        string    `json:"id" db:"id" valid:"required,uuid"`
	FullName  string    `json:"full_name" db:"full_name" valid:"required"`
	Surname   string    `json:"surname" db:"surname" valid:"required"`
	Name      string    `json:"name" db:"name" valid:"required"`
	Picture   string    `json:"picture" db:"picture" valid:"required"`
	Email     string    `json:"email" db:"email" valid:"required"`
	Domain    string    `json:"domain" db:"domain" valid:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewUserTemp(id string, fullName string, surname string, name string, picture string, email string, domain string) *UserTemp {
	return &UserTemp{
		ID:       id,
		FullName: fullName,
		Surname:  surname,
		Name:     name,
		Picture:  picture,
		Email:    email,
		Domain:   domain,
	}
}

func (m *UserTemp) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
