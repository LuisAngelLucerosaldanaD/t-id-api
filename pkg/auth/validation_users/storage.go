package validation_users

import (
	"github.com/jmoiron/sqlx"

	"check-id-api/internal/logger"
	"check-id-api/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesValidationUsersRepository interface {
	create(m *ValidationUsers) error
	update(m *ValidationUsers) error
	delete(id string) error
	getByID(id string) (*ValidationUsers, error)
	getAll() ([]*ValidationUsers, error)
	getByUserID(userId string) (*ValidationUsers, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesValidationUsersRepository {
	var s ServicesValidationUsersRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newValidationUsersPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
