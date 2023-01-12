package work

import (
	"check-id-api/internal/logger"
	"check-id-api/internal/msg"
	"check-id-api/pkg/wf"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type handlerWork struct {
	DB   *sqlx.DB
	TxID string
}

// getTotalWork godoc
// @Summary Trae la totalidad del trabajo existente
// @Description Método para obtener la totalidad del trabajo registrado por lo usuarios
// @tags Work
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Success 200 {object} resAllWork
// @Router /api/v1/work/all [get]
func (h *handlerWork) getTotalWork(c *fiber.Ctx) error {
	res := resAllWork{Error: true}
	srvWf := wf.NewServerWf(h.DB, nil, h.TxID)

	wfWork, err := srvWf.SrvWork.GetAllWorkValidation()
	if err != nil {
		logger.Error.Printf("No se pudo obtener el total del trabajo, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = wfWork
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}
