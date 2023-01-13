package work_validation

import (
	"fmt"
	"github.com/asaskevich/govalidator"

	"check-id-api/internal/logger"
	"check-id-api/internal/models"
)

type PortsServerWorkValidation interface {
	CreateWorkValidation(status string, userId string) (*WorkValidation, int, error)
	UpdateWorkValidation(id int64, status string, userId string) (*WorkValidation, int, error)
	DeleteWorkValidation(id int64) (int, error)
	GetWorkValidationByID(id int64) (*WorkValidation, int, error)
	GetAllWorkValidation() ([]*WorkValidation, error)
	GetWorkValidationByUserId(userId string) (*WorkValidation, int, error)
	GetAllWorkValidationByStatus(status string) ([]*WorkValidation, error)
}

type service struct {
	repository ServicesWorkValidationRepository
	user       *models.User
	txID       string
}

func NewWorkValidationService(repository ServicesWorkValidationRepository, user *models.User, TxID string) PortsServerWorkValidation {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateWorkValidation(status string, userId string) (*WorkValidation, int, error) {
	m := NewCreateWorkValidation(status, userId)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create WorkValidation :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateWorkValidation(id int64, status string, userId string) (*WorkValidation, int, error) {
	m := NewWorkValidation(id, status, userId)
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return m, 15, fmt.Errorf("id is required")
	}
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update WorkValidation :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteWorkValidation(id int64) (int, error) {
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

func (s *service) GetWorkValidationByID(id int64) (*WorkValidation, int, error) {
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

func (s *service) GetAllWorkValidation() ([]*WorkValidation, error) {
	return s.repository.getAll()
}

func (s *service) GetWorkValidationByUserId(userId string) (*WorkValidation, int, error) {
	if !govalidator.IsUUID(userId) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("userId is not uuid"))
		return nil, 15, fmt.Errorf("userId is not uuid")
	}
	m, err := s.repository.getByUserId(userId)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByUserId row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetAllWorkValidationByStatus(status string) ([]*WorkValidation, error) {
	return s.repository.getByStatus(status)
}
