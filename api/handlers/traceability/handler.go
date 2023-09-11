package traceability

import (
	"check-id-api/internal/logger"
	"check-id-api/internal/middleware"
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
// @Summary Obtención de los datos de trazabilidad de un usuario por su id
// @Description Método para obtención de los datos de trazabilidad de un usuario por su id
// @tags Traceability
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param userID path string true "ID del usuario"
// @Success 200 {object} resTraceability
// @Router /api/v1/traceability [get]
func (h *handlerTraceability) getTraceability(c *fiber.Ctx) error {
	res := resTraceability{Error: true}
	srvTrx := trx.NewServerTrx(h.DB, nil, h.TxID)

	userToken, err := middleware.GetUser(c)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el usuario del token, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(95, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	traceability, code, err := srvTrx.SrvTraceability.GetTraceabilityByUserID(userToken.ID)
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
