package onboarding

import (
	"check-id-api/internal/logger"
	"check-id-api/internal/models"
	"github.com/jmoiron/sqlx"
)

const (
	Postgresql = "postgres"
)

type ServicesOnboardingRepository interface {
	create(m *Onboarding) error
	update(m *Onboarding) error
	delete(id string) error
	getByID(id string) (*Onboarding, error)
	getAll() ([]*Onboarding, error)
	getByUserID(userId string) (*Onboarding, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesOnboardingRepository {
	var s ServicesOnboardingRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newOnboardingPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
