package user_temp

import (
	"fmt"

	"check-id-api/internal/logger"
	"check-id-api/internal/models"
	"github.com/asaskevich/govalidator"
)

type PortsServerUserTemp interface {
	CreateUserTemp(id string, fullName string, surname string, name string, picture string, email string, domain string) (*UserTemp, int, error)
	UpdateUserTemp(id string, fullName string, surname string, name string, picture string, email string, domain string) (*UserTemp, int, error)
	DeleteUserTemp(id string) (int, error)
	GetUserTempByID(id string) (*UserTemp, int, error)
	GetAllUserTemp() ([]*UserTemp, error)
	GetUserTempByEmail(email string) (*UserTemp, int, error)
	DeleteUserTempByEmail(email string) (int, error)
}

type service struct {
	repository ServicesUserTempRepository
	user       *models.User
	txID       string
}

func NewUserTempService(repository ServicesUserTempRepository, user *models.User, TxID string) PortsServerUserTemp {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateUserTemp(id string, fullName string, surname string, name string, picture string, email string, domain string) (*UserTemp, int, error) {
	m := NewUserTemp(id, fullName, surname, name, picture, email, domain)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create UserTemp :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateUserTemp(id string, fullName string, surname string, name string, picture string, email string, domain string) (*UserTemp, int, error) {
	m := NewUserTemp(id, fullName, surname, name, picture, email, domain)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update UserTemp :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteUserTemp(id string) (int, error) {
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

func (s *service) GetUserTempByID(id string) (*UserTemp, int, error) {
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

func (s *service) GetAllUserTemp() ([]*UserTemp, error) {
	return s.repository.getAll()
}

func (s *service) GetUserTempByEmail(email string) (*UserTemp, int, error) {
	if !govalidator.IsEmail(email) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("emailRq isn't email valid"))
		return nil, 15, fmt.Errorf("emailRq isn't email valid")
	}
	m, err := s.repository.getByEmail(email)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByEmail row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) DeleteUserTempByEmail(email string) (int, error) {
	if !govalidator.IsEmail(email) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("emailRq isn't email valid"))
		return 15, fmt.Errorf("emailRq isn't email valid")
	}

	if err := s.repository.deleteByEmail(email); err != nil {
		if err.Error() == "ecatch:108" {
			return 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't update row:", err)
		return 20, err
	}
	return 28, nil
}
