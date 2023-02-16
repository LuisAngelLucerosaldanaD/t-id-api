package validation_request

import (
	"github.com/jmoiron/sqlx"

	"check-id-api/internal/logger"
	"check-id-api/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesValidationRequestRepository interface {
	create(m *ValidationRequest) error
	update(m *ValidationRequest) error
	delete(id int64) error
	getByID(id int64) (*ValidationRequest, error)
	getAll() ([]*ValidationRequest, error)
	getByClientIDAndRequestID(clientIid int64, requestID string) (*ValidationRequest, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesValidationRequestRepository {
	var s ServicesValidationRequestRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newValidationRequestPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
