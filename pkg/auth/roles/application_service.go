package roles

import (
	"fmt"

	"check-id-api/internal/logger"
	"check-id-api/internal/models"
	"github.com/asaskevich/govalidator"
)

type PortsServerRoles interface {
	CreateRoles(id string, name string, description string) (*Roles, int, error)
	UpdateRoles(id string, name string, description string) (*Roles, int, error)
	DeleteRoles(id string) (int, error)
	GetRolesByID(id string) (*Roles, int, error)
	GetAllRoles() ([]*Roles, error)
}

type service struct {
	repository ServicesRolesRepository
	user       *models.User
	txID       string
}

func NewRolesService(repository ServicesRolesRepository, user *models.User, TxID string) PortsServerRoles {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateRoles(id string, name string, description string) (*Roles, int, error) {
	m := NewRoles(id, name, description)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create Roles :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateRoles(id string, name string, description string) (*Roles, int, error) {
	m := NewRoles(id, name, description)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update Roles :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteRoles(id string) (int, error) {
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

func (s *service) GetRolesByID(id string) (*Roles, int, error) {
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

func (s *service) GetAllRoles() ([]*Roles, error) {
	return s.repository.getAll()
}
