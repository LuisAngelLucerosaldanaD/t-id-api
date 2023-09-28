package onboarding_check_id

import (
	"github.com/jmoiron/sqlx"

	"check-id-api/internal/logger"
	"check-id-api/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesOnboardingCheckIdRepository interface {
	create(m *OnboardingCheckId) error
	update(m *OnboardingCheckId) error
	delete(id int64) error
	getByID(id int64) (*OnboardingCheckId, error)
	getAll() ([]*OnboardingCheckId, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesOnboardingCheckIdRepository {
	var s ServicesOnboardingCheckIdRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newOnboardingCheckIdPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
