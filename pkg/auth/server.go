package auth

import (
	"check-id-api/internal/models"
	"check-id-api/pkg/auth/life_test"
	"check-id-api/pkg/auth/onboarding"
	"check-id-api/pkg/auth/onboarding_check_id"
	"check-id-api/pkg/auth/role"
	"check-id-api/pkg/auth/user"
	"check-id-api/pkg/auth/user_role"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	SrvUser              user.PortsServerUser
	SrvOnboarding        onboarding.PortsServerOnboarding
	SrvUserRole          user_role.PortsServerUseRole
	SrvRole              role.PortsServerRole
	SrvLifeTest          life_test.PortsServerLifeTest
	SrvOnboardingCheckId onboarding_check_id.PortsServerOnboardingCheckId
}

func NewServerAuth(db *sqlx.DB, userModel *models.User, txID string) *Server {

	repoUser := user.FactoryStorage(db, userModel, txID)
	srvUser := user.NewUsersService(repoUser, userModel, txID)

	repoOnboarding := onboarding.FactoryStorage(db, userModel, txID)
	srvOnboarding := onboarding.NewOnboardingService(repoOnboarding, userModel, txID)

	repoUserRole := user_role.FactoryStorage(db, userModel, txID)
	srvUserRole := user_role.NewUseRoleService(repoUserRole, userModel, txID)

	repoRole := role.FactoryStorage(db, userModel, txID)
	srvRole := role.NewRoleService(repoRole, userModel, txID)

	repoLifeTest := life_test.FactoryStorage(db, userModel, txID)
	srvLifeTest := life_test.NewLifeTestService(repoLifeTest, userModel, txID)

	repoOnboardingCheckId := onboarding_check_id.FactoryStorage(db, userModel, txID)
	srvOnboardingCheckId := onboarding_check_id.NewOnboardingCheckIdService(repoOnboardingCheckId, userModel, txID)

	return &Server{
		SrvUser:              srvUser,
		SrvOnboarding:        srvOnboarding,
		SrvUserRole:          srvUserRole,
		SrvRole:              srvRole,
		SrvLifeTest:          srvLifeTest,
		SrvOnboardingCheckId: srvOnboardingCheckId,
	}
}
