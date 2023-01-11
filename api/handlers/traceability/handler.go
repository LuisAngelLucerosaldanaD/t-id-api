package traceability

import (
	"check-id-api/internal/logger"
	"check-id-api/internal/msg"
	"check-id-api/pkg/trx"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type handlerTraceability struct {
	DB   *sqlx.DB
	TxID string
}

// getTraceabilitySession godoc
// @Summary Obtención de los datos de trazabilidad
// @Description Método para obtener la trazabilidad registrada para el proceso de verificación de identidad
// @tags Traceability
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param userID path string true "ID del usuario"
// @Success 200 {object} resTraceability
// @Router /api/v1/traceability/user-session/{userID} [get]
func (h *handlerTraceability) getTraceabilitySession(c *fiber.Ctx) error {
	res := resTraceability{Error: true}
	srvTrx := trx.NewServerTrx(h.DB, nil, h.TxID)

	userID := c.Params("userID")
	if userID == "" {
		logger.Error.Printf("No se pudo obtener el id del usuario")
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	traceability, code, err := srvTrx.SrvTraceability.GetTraceabilityByUserID(userID)
	if err != nil {
		logger.Error.Printf("No se pudo obtener la trazabilidad del usuario, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = traceability
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}
