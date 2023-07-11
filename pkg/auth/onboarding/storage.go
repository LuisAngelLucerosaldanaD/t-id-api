package onboarding

import (
	"github.com/jmoiron/sqlx"

	"onlyone_smc/internal/logger"
	"onlyone_smc/internal/models"
)

const (
	Postgresql = "postgres"
	SqlServer  = "sqlserver"
	Oracle     = "oci8"
)

type ServicesOnboardingRepository interface {
	create(m *Onboarding) error
	update(m *Onboarding) error
	delete(id string) error
	getByID(id string) (*Onboarding, error)
	getAll() ([]*Onboarding, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesOnboardingRepository {
	var s ServicesOnboardingRepository
	engine := db.DriverName()
	switch engine {
	case SqlServer:
		return newOnboardingSqlServerRepository(db, user, txID)
	case Postgresql:
		return newOnboardingPsqlRepository(db, user, txID)
	case Oracle:
		return newOnboardingOrclRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
