package files

import (
	"fmt"
	"github.com/asaskevich/govalidator"

	"check-id-api/internal/logger"
	"check-id-api/internal/models"
)

type PortsServerFiles interface {
	CreateFiles(path string, name string, typeFile int32, userId string) (*Files, int, error)
	UpdateFiles(id int64, path string, name string, typeFile int32, userId string) (*Files, int, error)
	DeleteFiles(id int64) (int, error)
	GetFilesByID(id int64) (*Files, int, error)
	GetAllFiles() ([]*Files, error)
	GetFilesByUserID(userId string) ([]*Files, int, error)
	DeleteFilesByUserID(userId string) (int, error)
}

type service struct {
	repository ServicesFilesRepository
	user       *models.User
	txID       string
}

func NewFilesService(repository ServicesFilesRepository, user *models.User, TxID string) PortsServerFiles {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateFiles(path string, name string, typeFile int32, userId string) (*Files, int, error) {
	m := NewCreateFiles(path, name, typeFile, userId)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create Files :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateFiles(id int64, path string, name string, typeFile int32, userId string) (*Files, int, error) {
	m := NewFiles(id, path, name, typeFile, userId)
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return m, 15, fmt.Errorf("id is required")
	}
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update Files :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteFiles(id int64) (int, error) {
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

func (s *service) GetFilesByID(id int64) (*Files, int, error) {
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

func (s *service) GetAllFiles() ([]*Files, error) {
	return s.repository.getAll()
}

func (s *service) GetFilesByUserID(userId string) ([]*Files, int, error) {
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

func (s *service) DeleteFilesByUserID(userId string) (int, error) {
	if !govalidator.IsUUID(userId) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("userId isn't uuid"))
		return 15, fmt.Errorf("userId isn't uuid")
	}
	if err := s.repository.deleteByUserId(userId); err != nil {
		if err.Error() == "ecatch:108" {
			return 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't delete row:", err)
		return 20, err
	}
	return 28, nil
}
