package auth

import (
	"check-id-api/internal/models"
	"check-id-api/pkg/auth/roles"
	"check-id-api/pkg/auth/users"
	"check-id-api/pkg/auth/users_rol"
	"check-id-api/pkg/auth/validation_users"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	SrvUser            users.PortsServerUsers
	SrvRoles           roles.PortsServerRoles
	SrvUsersRol        users_rol.PortsServerUsersRol
	SrvValidationUsers validation_users.PortsServerValidationUsers
}

func NewServerAuth(db *sqlx.DB, user *models.User, txID string) *Server {

	repoUser := users.FactoryStorage(db, user, txID)
	srvUser := users.NewUsersService(repoUser, user, txID)

	repoRoles := roles.FactoryStorage(db, user, txID)
	srvRoles := roles.NewRolesService(repoRoles, user, txID)

	repoUsersRol := users_rol.FactoryStorage(db, user, txID)
	srvUsersRol := users_rol.NewUsersRolService(repoUsersRol, user, txID)

	repoValidationUsers := validation_users.FactoryStorage(db, user, txID)
	srvValidationUsers := validation_users.NewValidationUsersService(repoValidationUsers, user, txID)

	return &Server{
		SrvUser:            srvUser,
		SrvRoles:           srvRoles,
		SrvUsersRol:        srvUsersRol,
		SrvValidationUsers: srvValidationUsers,
	}
}
