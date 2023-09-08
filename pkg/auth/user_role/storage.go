package user_role

import (
	"github.com/jmoiron/sqlx"

	"check-id-api/internal/logger"
	"check-id-api/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesUseRoleRepository interface {
	create(m *UseRole) error
	update(m *UseRole) error
	delete(id string) error
	getByID(id string) (*UseRole, error)
	getAll() ([]*UseRole, error)
	updateRoleByUserid(m *UseRole) error
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesUseRoleRepository {
	var s ServicesUseRoleRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newUseRolePsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
