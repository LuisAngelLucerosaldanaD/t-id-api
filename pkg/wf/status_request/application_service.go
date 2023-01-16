package status_request

import (
	"fmt"
	"github.com/asaskevich/govalidator"

	"check-id-api/internal/logger"
	"check-id-api/internal/models"
)

type PortsServerStatusRequest interface {
	CreateStatusRequest(status string, description string, userId string) (*StatusRequest, int, error)
	UpdateStatusRequest(id int64, status string, description string, userId string) (*StatusRequest, int, error)
	DeleteStatusRequest(id int64) (int, error)
	GetStatusRequestByID(id int64) (*StatusRequest, int, error)
	GetAllStatusRequest() ([]*StatusRequest, error)
	GetStatusRequestByUserID(userId string) (*StatusRequest, int, error)
	GetStatusRequestByStatus(status string) ([]*StatusRequest, int, error)
	UpdateStatusRequestByUserId(status string, description string, userID string) (int, error)
}

type service struct {
	repository ServicesStatusRequestRepository
	user       *models.User
	txID       string
}

func NewStatusRequestService(repository ServicesStatusRequestRepository, user *models.User, TxID string) PortsServerStatusRequest {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateStatusRequest(status string, description string, userId string) (*StatusRequest, int, error) {
	m := NewCreateStatusRequest(status, description, userId)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create StatusRequest :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateStatusRequest(id int64, status string, description string, userId string) (*StatusRequest, int, error) {
	m := NewStatusRequest(id, status, description, userId)
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return m, 15, fmt.Errorf("id is required")
	}
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update StatusRequest :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteStatusRequest(id int64) (int, error) {
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

func (s *service) GetStatusRequestByID(id int64) (*StatusRequest, int, error) {
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

func (s *service) GetAllStatusRequest() ([]*StatusRequest, error) {
	return s.repository.getAll()
}

func (s *service) GetStatusRequestByUserID(userId string) (*StatusRequest, int, error) {
	if !govalidator.IsUUID(userId) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is not uuid valid"))
		return nil, 15, fmt.Errorf("id is not uuid valid")
	}
	m, err := s.repository.getByUserId(userId)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByUserId row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetStatusRequestByStatus(status string) ([]*StatusRequest, int, error) {
	m, err := s.repository.getByStatus(status)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByStatus row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) UpdateStatusRequestByUserId(status string, description string, userID string) (int, error) {
	if err := s.repository.updateStatus(status, description, userID); err != nil {
		logger.Error.Println(s.txID, " - couldn't update WorkValidation :", err)
		return 18, err
	}
	return 29, nil
}
