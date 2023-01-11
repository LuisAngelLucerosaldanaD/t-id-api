package traceability

import (
	"fmt"
	"github.com/asaskevich/govalidator"

	"check-id-api/internal/logger"
	"check-id-api/internal/models"
)

type PortsServerTraceability interface {
	CreateTraceability(action string, description string, userId string) (*Traceability, int, error)
	UpdateTraceability(id int64, action string, description string, userId string) (*Traceability, int, error)
	DeleteTraceability(id int64) (int, error)
	GetTraceabilityByID(id int64) (*Traceability, int, error)
	GetAllTraceability() ([]*Traceability, error)
	GetTraceabilityByUserID(userId string) ([]*Traceability, int, error)
}

type service struct {
	repository ServicesTraceabilityRepository
	user       *models.User
	txID       string
}

func NewTraceabilityService(repository ServicesTraceabilityRepository, user *models.User, TxID string) PortsServerTraceability {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateTraceability(action string, description string, userId string) (*Traceability, int, error) {
	m := NewCreateTraceability(action, description, userId)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create Traceability :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateTraceability(id int64, action string, description string, userId string) (*Traceability, int, error) {
	m := NewTraceability(id, action, description, userId)
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return m, 15, fmt.Errorf("id is required")
	}
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update Traceability :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteTraceability(id int64) (int, error) {
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

func (s *service) GetTraceabilityByID(id int64) (*Traceability, int, error) {
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

func (s *service) GetAllTraceability() ([]*Traceability, error) {
	return s.repository.getAll()
}

func (s *service) GetTraceabilityByUserID(userId string) ([]*Traceability, int, error) {
	if !govalidator.IsUUID(userId) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("userId isn't uuid"))
		return nil, 15, fmt.Errorf("userId isn't uuid")
	}
	m, err := s.repository.getByUserID(userId)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByUserID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}
