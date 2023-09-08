package life_test

import (
	"github.com/jmoiron/sqlx"

	"check-id-api/internal/logger"
	"check-id-api/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesLifeTestRepository interface {
	create(m *LifeTest) error
	update(m *LifeTest) error
	delete(id int64) error
	getByID(id int64) (*LifeTest, error)
	getAll() ([]*LifeTest, error)
	getByClientIDAndRequestID(clientIid int64, requestID string) (*LifeTest, error)
	updateStatus(m *LifeTest) error
	getAllByUserId(userID string) ([]*LifeTest, error)
	getByUserID(id string) (*LifeTest, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesLifeTestRepository {
	var s ServicesLifeTestRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newLifeTestPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
