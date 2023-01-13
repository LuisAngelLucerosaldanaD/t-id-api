package work

import (
	"check-id-api/internal/logger"
	"check-id-api/internal/msg"
	"check-id-api/pkg/auth"
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
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	wfOK, err := srvWf.SrvWork.GetAllWorkValidationByStatus("validado")
	if err != nil {
		logger.Error.Printf("No se pudo obtener el total del trabajo validado, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	wfPending, err := srvWf.SrvWork.GetAllWorkValidationByStatus("pendiente")
	if err != nil {
		logger.Error.Printf("No se pudo obtener el total del trabajo pendiente, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	wfRefused, err := srvWf.SrvWork.GetAllWorkValidationByStatus("rechazado")
	if err != nil {
		logger.Error.Printf("No se pudo obtener el total del trabajo rechazado, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	wfExpired, err := srvWf.SrvWork.GetAllWorkValidationByStatus("expirado")
	if err != nil {
		logger.Error.Printf("No se pudo obtener el total del trabajo expirado, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	usersTemp, err := srvAuth.SrvUserTemp.GetAllUserTemp()
	if err != nil {
		logger.Error.Printf("No se pudo obtener el total del trabajo no iniciado, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = Status{
		Valid:     len(wfOK),
		Pending:   len(wfPending),
		Refused:   len(wfRefused),
		Total:     len(wfOK) + len(wfPending) + len(wfRefused) + len(wfExpired) + len(usersTemp),
		Expired:   len(wfExpired),
		NotStated: len(usersTemp),
	}
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}
