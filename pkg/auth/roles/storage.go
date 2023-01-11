package roles

import (
	"github.com/jmoiron/sqlx"

	"check-id-api/internal/logger"
	"check-id-api/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesRolesRepository interface {
	create(m *Roles) error
	update(m *Roles) error
	delete(id string) error
	getByID(id string) (*Roles, error)
	getAll() ([]*Roles, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesRolesRepository {
	var s ServicesRolesRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newRolesPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
