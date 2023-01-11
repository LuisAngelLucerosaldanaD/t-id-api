package validation_users

import (
	"fmt"

	"check-id-api/internal/logger"
	"check-id-api/internal/models"
	"github.com/asaskevich/govalidator"
)

type PortsServerValidationUsers interface {
	CreateValidationUsers(id string, transactionId string, userId string) (*ValidationUsers, int, error)
	UpdateValidationUsers(id string, transactionId string, userId string) (*ValidationUsers, int, error)
	DeleteValidationUsers(id string) (int, error)
	GetValidationUsersByID(id string) (*ValidationUsers, int, error)
	GetAllValidationUsers() ([]*ValidationUsers, error)
	GetValidationUsersByUserID(userId string) (*ValidationUsers, int, error)
}

type service struct {
	repository ServicesValidationUsersRepository
	user       *models.User
	txID       string
}

func NewValidationUsersService(repository ServicesValidationUsersRepository, user *models.User, TxID string) PortsServerValidationUsers {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateValidationUsers(id string, transactionId string, userId string) (*ValidationUsers, int, error) {
	m := NewValidationUsers(id, transactionId, userId)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create ValidationUsers :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateValidationUsers(id string, transactionId string, userId string) (*ValidationUsers, int, error) {
	m := NewValidationUsers(id, transactionId, userId)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update ValidationUsers :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteValidationUsers(id string) (int, error) {
	if !govalidator.IsUUID(id) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return 15, fmt.Errorf("id isn't uuid")
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

func (s *service) GetValidationUsersByID(id string) (*ValidationUsers, int, error) {
	if !govalidator.IsUUID(id) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return nil, 15, fmt.Errorf("id isn't uuid")
	}
	m, err := s.repository.getByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetAllValidationUsers() ([]*ValidationUsers, error) {
	return s.repository.getAll()
}

func (s *service) GetValidationUsersByUserID(userId string) (*ValidationUsers, int, error) {
	if !govalidator.IsUUID(userId) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return nil, 15, fmt.Errorf("id isn't uuid")
	}
	m, err := s.repository.getByUserID(userId)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByUserID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}
