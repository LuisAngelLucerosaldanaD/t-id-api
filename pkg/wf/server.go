package wf

import (
	"check-id-api/internal/models"
	"check-id-api/pkg/wf/work_validation"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	SrvWork work_validation.PortsServerWorkValidation
}

func NewServerWf(db *sqlx.DB, user *models.User, txID string) *Server {

	repoWork := work_validation.FactoryStorage(db, user, txID)
	srvWork := work_validation.NewWorkValidationService(repoWork, user, txID)

	return &Server{
		SrvWork: srvWork,
	}
}
