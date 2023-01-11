package files

import (
	"github.com/jmoiron/sqlx"

	"check-id-api/internal/logger"
	"check-id-api/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesFilesRepository interface {
	create(m *Files) error
	update(m *Files) error
	delete(id int64) error
	getByID(id int64) (*Files, error)
	getAll() ([]*Files, error)
	getByUserID(userID string) ([]*Files, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesFilesRepository {
	var s ServicesFilesRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newFilesPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
