package status_request

import (
	"github.com/jmoiron/sqlx"

	"check-id-api/internal/logger"
	"check-id-api/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesStatusRequestRepository interface {
	create(m *StatusRequest) error
	update(m *StatusRequest) error
	delete(id int64) error
	getByID(id int64) (*StatusRequest, error)
	getAll() ([]*StatusRequest, error)
	getByUserId(userID string) (*StatusRequest, error)
	getByStatus(status string) ([]*StatusRequest, error)
	updateStatus(status string, description string, userID string) error
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesStatusRequestRepository {
	var s ServicesStatusRequestRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newStatusRequestPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
