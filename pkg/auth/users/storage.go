package users

import (
	"check-id-api/internal/logger"
	"check-id-api/internal/models"
	"github.com/jmoiron/sqlx"
)

const (
	Postgresql = "postgres"
)

type ServicesUsersRepository interface {
	create(m *Users) error
	update(m *Users) error
	delete(id string) error
	getByID(id string) (*Users, error)
	getAll() ([]*Users, error)
	getByEmail(email string) (*Users, error)
	getLasted(email string, limit, offset int) ([]*Users, error)
	getNotStarted() ([]*Users, error)
	getNoUploadFile(fileType int) ([]*Users, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesUsersRepository {
	var s ServicesUsersRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newUsersPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
