package user_role

import (
	"check-id-api/internal/logger"
	"check-id-api/internal/models"
	"fmt"
	"github.com/asaskevich/govalidator"
)

type PortsServerUseRole interface {
	CreateUseRole(id string, userId string, roleId string) (*UseRole, int, error)
	UpdateUseRole(id string, userId string, roleId string) (*UseRole, int, error)
	DeleteUseRole(id string) (int, error)
	GetUseRoleByID(id string) (*UseRole, int, error)
	GetAllUseRole() ([]*UseRole, error)
	UpdateUseRoleByUserID(userId string, roleId string) (*UseRole, int, error)
	GetUseRoleByUserID(id string) (*UseRole, int, error)
}

type service struct {
	repository ServicesUseRoleRepository
	user       *models.User
	txID       string
}

func NewUseRoleService(repository ServicesUseRoleRepository, user *models.User, TxID string) PortsServerUseRole {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateUseRole(id string, userId string, roleId string) (*UseRole, int, error) {
	m := NewUseRole(id, userId, roleId)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create UseRole :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateUseRole(id string, userId string, roleId string) (*UseRole, int, error) {
	m := NewUseRole(id, userId, roleId)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update UseRole :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteUseRole(id string) (int, error) {
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

func (s *service) GetUseRoleByID(id string) (*UseRole, int, error) {
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

func (s *service) GetUseRoleByUserID(id string) (*UseRole, int, error) {
	if !govalidator.IsUUID(id) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("user id isn't uuid"))
		return nil, 15, fmt.Errorf("user id isn't uuid")
	}
	m, err := s.repository.getByUseID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetAllUseRole() ([]*UseRole, error) {
	return s.repository.getAll()
}

func (s *service) UpdateUseRoleByUserID(userId string, roleId string) (*UseRole, int, error) {
	m := &UseRole{
		UserId: userId,
		RoleId: roleId,
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update UseRole :", err)
		return m, 18, err
	}
	return m, 29, nil
}
