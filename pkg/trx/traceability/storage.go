package traceability

import (
	"github.com/jmoiron/sqlx"

	"check-id-api/internal/logger"
	"check-id-api/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesTraceabilityRepository interface {
	create(m *Traceability) error
	update(m *Traceability) error
	delete(id int64) error
	getByID(id int64) (*Traceability, error)
	getAll() ([]*Traceability, error)
	getByUserID(userId string) ([]*Traceability, error)
	deleteByUserID(userId string) error
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesTraceabilityRepository {
	var s ServicesTraceabilityRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newTraceabilityPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
