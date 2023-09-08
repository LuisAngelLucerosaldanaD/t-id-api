package role

import (
	"fmt"

	"check-id-api/internal/logger"
	"check-id-api/internal/models"
	"github.com/asaskevich/govalidator"
)

type PortsServerRole interface {
	CreateRole(id string, name string, description string) (*Role, int, error)
	UpdateRole(id string, name string, description string) (*Role, int, error)
	DeleteRole(id string) (int, error)
	GetRoleByID(id string) (*Role, int, error)
	GetAllRole() ([]*Role, error)
	GetRoleByUserID(userID string) (*Role, int, error)
}

type service struct {
	repository ServicesRoleRepository
	user       *models.User
	txID       string
}

func NewRoleService(repository ServicesRoleRepository, user *models.User, TxID string) PortsServerRole {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateRole(id string, name string, description string) (*Role, int, error) {
	m := NewRole(id, name, description)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create Role :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateRole(id string, name string, description string) (*Role, int, error) {
	m := NewRole(id, name, description)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update Role :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteRole(id string) (int, error) {
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

func (s *service) GetRoleByID(id string) (*Role, int, error) {
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

func (s *service) GetAllRole() ([]*Role, error) {
	return s.repository.getAll()
}

func (s *service) GetRoleByUserID(userID string) (*Role, int, error) {
	if !govalidator.IsUUID(userID) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("user id isn't uuid"))
		return nil, 15, fmt.Errorf("user id isn't uuid")
	}
	m, err := s.repository.getByUserID(userID)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByUserID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}
