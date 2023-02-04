package clients

import (
	"github.com/jmoiron/sqlx"

	"check-id-api/internal/logger"
	"check-id-api/internal/models"
)

const (
	Postgresql = "postgres"
	SqlServer  = "sqlserver"
	Oracle     = "oci8"
)

type ServicesClientsRepository interface {
	create(m *Clients) error
	update(m *Clients) error
	delete(id int64) error
	getByID(id int64) (*Clients, error)
	getAll() ([]*Clients, error)
	getByNit(nit string) (*Clients, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesClientsRepository {
	var s ServicesClientsRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newClientsPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
