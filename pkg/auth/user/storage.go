package user

import (
	"check-id-api/internal/logger"
	"check-id-api/internal/models"
	"github.com/jmoiron/sqlx"
)

const (
	Postgresql = "postgres"
)

type ServicesUserRepository interface {
	create(m *User) error
	update(m *User) error
	delete(id string) error
	getByID(id string) (*User, error)
	getAll() ([]*User, error)
	getByEmail(email string) (*User, error)
	getLasted(email string, limit, offset int) ([]*User, error)
	getNotStarted() ([]*User, error)
	getNoUploadFile(fileType int) ([]*User, error)
	getByIdentityNumber(identityNumber string) (*User, error)
	getByDniAndEmail(dni string, email string) (*User, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesUserRepository {
	var s ServicesUserRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newUserPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
