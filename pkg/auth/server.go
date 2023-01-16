package auth

import (
	"check-id-api/internal/models"
	"check-id-api/pkg/auth/users"
	"check-id-api/pkg/auth/validation_users"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	SrvUser            users.PortsServerUsers
	SrvValidationUsers validation_users.PortsServerValidationUsers
}

func NewServerAuth(db *sqlx.DB, user *models.User, txID string) *Server {

	repoUser := users.FactoryStorage(db, user, txID)
	srvUser := users.NewUsersService(repoUser, user, txID)

	repoValidationUsers := validation_users.FactoryStorage(db, user, txID)
	srvValidationUsers := validation_users.NewValidationUsersService(repoValidationUsers, user, txID)

	return &Server{
		SrvUser:            srvUser,
		SrvValidationUsers: srvValidationUsers,
	}
}
