package user_temp

import (
	"github.com/jmoiron/sqlx"

	"check-id-api/internal/logger"
	"check-id-api/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesUserTempRepository interface {
	create(m *UserTemp) error
	update(m *UserTemp) error
	delete(id string) error
	getByID(id string) (*UserTemp, error)
	getAll() ([]*UserTemp, error)
	getByEmail(email string) (*UserTemp, error)
	deleteByEmail(email string) error
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesUserTempRepository {
	var s ServicesUserTempRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newUserTempPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
