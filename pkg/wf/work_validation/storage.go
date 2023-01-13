package work_validation

import (
	"github.com/jmoiron/sqlx"

	"check-id-api/internal/logger"
	"check-id-api/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesWorkValidationRepository interface {
	create(m *WorkValidation) error
	update(m *WorkValidation) error
	delete(id int64) error
	getByID(id int64) (*WorkValidation, error)
	getAll() ([]*WorkValidation, error)
	getByUserId(userID string) (*WorkValidation, error)
	getByStatus(status string) ([]*WorkValidation, error)
	updateStatus(status string, userID string) error
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesWorkValidationRepository {
	var s ServicesWorkValidationRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newWorkValidationPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
