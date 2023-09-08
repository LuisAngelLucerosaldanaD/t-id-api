package role

import (
	"github.com/jmoiron/sqlx"

	"check-id-api/internal/logger"
	"check-id-api/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesRoleRepository interface {
	create(m *Role) error
	update(m *Role) error
	delete(id string) error
	getByID(id string) (*Role, error)
	getAll() ([]*Role, error)
	getByUserID(userID string) (*Role, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesRoleRepository {
	var s ServicesRoleRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newRolePsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
