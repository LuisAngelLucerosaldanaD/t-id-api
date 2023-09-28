package onboarding_check_id

import (
	"fmt"

	"check-id-api/internal/logger"
	"check-id-api/internal/models"
)

type PortsServerOnboardingCheckId interface {
	CreateOnboardingCheckId(userId string, ip string) (*OnboardingCheckId, int, error)
	UpdateOnboardingCheckId(id int64, userId string, ip string) (*OnboardingCheckId, int, error)
	DeleteOnboardingCheckId(id int64) (int, error)
	GetOnboardingCheckIdByID(id int64) (*OnboardingCheckId, int, error)
	GetAllOnboardingCheckId() ([]*OnboardingCheckId, error)
}

type service struct {
	repository ServicesOnboardingCheckIdRepository
	user       *models.User
	txID       string
}

func NewOnboardingCheckIdService(repository ServicesOnboardingCheckIdRepository, user *models.User, TxID string) PortsServerOnboardingCheckId {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateOnboardingCheckId(userId string, ip string) (*OnboardingCheckId, int, error) {
	m := NewCreateOnboardingCheckId(userId, ip)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create OnboardingCheckId :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateOnboardingCheckId(id int64, userId string, ip string) (*OnboardingCheckId, int, error) {
	m := NewOnboardingCheckId(id, userId, ip)
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return m, 15, fmt.Errorf("id is required")
	}
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update OnboardingCheckId :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteOnboardingCheckId(id int64) (int, error) {
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return 15, fmt.Errorf("id is required")
	}

	if err := s.repository.delete(id); err != nil {
		if err.Error() == "ecatch:108" {
			return 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't update row:", err)
		return 20, err
	}
	return 28, nil
}

func (s *service) GetOnboardingCheckIdByID(id int64) (*OnboardingCheckId, int, error) {
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return nil, 15, fmt.Errorf("id is required")
	}
	m, err := s.repository.getByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetAllOnboardingCheckId() ([]*OnboardingCheckId, error) {
	return s.repository.getAll()
}
