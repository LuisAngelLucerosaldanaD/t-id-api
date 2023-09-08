package client

import (
	"github.com/jmoiron/sqlx"

	"check-id-api/internal/logger"
	"check-id-api/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesClientRepository interface {
	create(m *Client) error
	update(m *Client) error
	delete(id int64) error
	getByID(id int64) (*Client, error)
	getAll() ([]*Client, error)
	getByNit(nit string) (*Client, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesClientRepository {
	var s ServicesClientRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newClientPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
