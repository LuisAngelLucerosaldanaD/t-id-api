package onboarding_check_id

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// OnboardingCheckId  Model struct OnboardingCheckId
type OnboardingCheckId struct {
	ID        int64     `json:"id" db:"id" valid:"-"`
	UserId    string    `json:"user_id" db:"user_id" valid:"required"`
	Ip        string    `json:"ip" db:"ip" valid:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewOnboardingCheckId(id int64, userId string, ip string) *OnboardingCheckId {
	return &OnboardingCheckId{
		ID:     id,
		UserId: userId,
		Ip:     ip,
	}
}

func NewCreateOnboardingCheckId(userId string, ip string) *OnboardingCheckId {
	return &OnboardingCheckId{
		UserId: userId,
		Ip:     ip,
	}
}

func (m *OnboardingCheckId) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
