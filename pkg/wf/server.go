package wf

import (
	"check-id-api/internal/models"
	"check-id-api/pkg/wf/status_request"
	"check-id-api/pkg/wf/work_validation"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	SrvWork      work_validation.PortsServerWorkValidation
	SrvStatusReq status_request.PortsServerStatusRequest
}

func NewServerWf(db *sqlx.DB, user *models.User, txID string) *Server {

	repoWork := work_validation.FactoryStorage(db, user, txID)
	srvWork := work_validation.NewWorkValidationService(repoWork, user, txID)

	repoStatusReq := status_request.FactoryStorage(db, user, txID)
	srvStatusReq := status_request.NewStatusRequestService(repoStatusReq, user, txID)

	return &Server{
		SrvWork:      srvWork,
		SrvStatusReq: srvStatusReq,
	}
}
