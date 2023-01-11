package users_rol

import (
	"github.com/jmoiron/sqlx"

	"check-id-api/internal/logger"
	"check-id-api/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesUsersRolRepository interface {
	create(m *UsersRol) error
	update(m *UsersRol) error
	delete(id string) error
	getByID(id string) (*UsersRol, error)
	getAll() ([]*UsersRol, error)
	getByUserID(userID string) (*UsersRol, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesUsersRolRepository {
	var s ServicesUsersRolRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newUsersRolPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
