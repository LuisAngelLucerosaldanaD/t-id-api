package users_rol

import (
	"fmt"

	"check-id-api/internal/logger"
	"check-id-api/internal/models"
	"github.com/asaskevich/govalidator"
)

type PortsServerUsersRol interface {
	CreateUsersRol(id string, userId string, roleId string) (*UsersRol, int, error)
	UpdateUsersRol(id string, userId string, roleId string) (*UsersRol, int, error)
	DeleteUsersRol(id string) (int, error)
	GetUsersRolByID(id string) (*UsersRol, int, error)
	GetAllUsersRol() ([]*UsersRol, error)
	GetUsersRolByUserId(userId string) (*UsersRol, int, error)
}

type service struct {
	repository ServicesUsersRolRepository
	user       *models.User
	txID       string
}

func NewUsersRolService(repository ServicesUsersRolRepository, user *models.User, TxID string) PortsServerUsersRol {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateUsersRol(id string, userId string, roleId string) (*UsersRol, int, error) {
	m := NewUsersRol(id, userId, roleId)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create UsersRol :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateUsersRol(id string, userId string, roleId string) (*UsersRol, int, error) {
	m := NewUsersRol(id, userId, roleId)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update UsersRol :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteUsersRol(id string) (int, error) {
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

func (s *service) GetUsersRolByID(id string) (*UsersRol, int, error) {
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

func (s *service) GetAllUsersRol() ([]*UsersRol, error) {
	return s.repository.getAll()
}

func (s *service) GetUsersRolByUserId(userId string) (*UsersRol, int, error) {
	if !govalidator.IsUUID(userId) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("userId isn't uuid"))
		return nil, 15, fmt.Errorf("userId isn't uuid")
	}
	m, err := s.repository.getByUserID(userId)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}
