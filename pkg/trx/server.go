package trx

import (
	"check-id-api/internal/models"
	"check-id-api/pkg/trx/traceability"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	SrvTraceability traceability.PortsServerTraceability
}

func NewServerTrx(db *sqlx.DB, user *models.User, txID string) *Server {

	repoTraceability := traceability.FactoryStorage(db, user, txID)
	srvTraceability := traceability.NewTraceabilityService(repoTraceability, user, txID)

	return &Server{
		SrvTraceability: srvTraceability,
	}
}
