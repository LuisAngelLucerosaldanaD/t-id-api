package work

import (
	"check-id-api/internal/logger"
	"check-id-api/internal/msg"
	"check-id-api/pkg/auth"
	"check-id-api/pkg/trx"
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

	wfOK, code, err := srvWf.SrvStatusReq.GetStatusRequestByStatus("validado")
	if err != nil {
		logger.Error.Printf("No se pudo obtener el total del trabajo validado, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	wfPending, code, err := srvWf.SrvStatusReq.GetStatusRequestByStatus("pendiente")
	if err != nil {
		logger.Error.Printf("No se pudo obtener el total del trabajo pendiente, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	wfRefused, code, err := srvWf.SrvStatusReq.GetStatusRequestByStatus("rechazado")
	if err != nil {
		logger.Error.Printf("No se pudo obtener el total del trabajo rechazado, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	wfExpired, code, err := srvWf.SrvStatusReq.GetStatusRequestByStatus("expirado")
	if err != nil {
		logger.Error.Printf("No se pudo obtener el total del trabajo expirado, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	usersTemp, err := srvAuth.SrvUser.GetAllNotStarted()
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

// acceptUserData godoc
// @Summary Acepta la información registrada
// @Description Método para aceptar la data registrada de un usuario por parte del administrador
// @tags Work
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param ReqAccept body ReqAccept true "Datos de solicitud para la aceptación"
// @Success 200 {object} resAnny
// @Router /api/v1/work/accept [post]
func (h *handlerWork) acceptUserData(c *fiber.Ctx) error {
	res := resAnny{Error: true}
	req := ReqAccept{}
	err := c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("el id del usuario es requerido")
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	srvWf := wf.NewServerWf(h.DB, nil, h.TxID)
	srvTrx := trx.NewServerTrx(h.DB, nil, h.TxID)

	code, err := srvWf.SrvWork.UpdateWorkValidationStatus("validado", req.UserID)
	if err != nil {
		logger.Error.Printf("No se pudo actualizar el registro, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvTrx.SrvTraceability.CreateTraceability("Validación de datos", "success", "Los datos registrados fueron validados y aceptados", req.UserID)
	if err != nil {
		logger.Error.Printf("couldn't create traceability, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = "Datos registrados correctamente"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// refusedUserData godoc
// @Summary Rechaza la información registrada
// @Description Método para rechazar la data registrada de un usuario por parte del administrador
// @tags Work
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param ReqAccept body ReqAccept true "Datos de solicitud para la aceptación"
// @Success 200 {object} resAnny
// @Router /api/v1/work/refused [post]
func (h *handlerWork) refusedUserData(c *fiber.Ctx) error {
	res := resAnny{Error: true}
	req := ReqAccept{}
	err := c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("el id del usuario es requerido")
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	srvWf := wf.NewServerWf(h.DB, nil, h.TxID)
	srvTrx := trx.NewServerTrx(h.DB, nil, h.TxID)

	code, err := srvWf.SrvWork.UpdateWorkValidationStatus("rechazado", req.UserID)
	if err != nil {
		logger.Error.Printf("No se pudo actualizar el registro, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvTrx.SrvTraceability.CreateTraceability("Validación de datos", "error", "Los datos fueron rechazados por un administrador", req.UserID)
	if err != nil {
		logger.Error.Printf("couldn't create traceability, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = "Datos registrados correctamente"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}
