package user_temp

import (
	"check-id-api/internal/logger"
	"check-id-api/internal/msg"
	"check-id-api/pkg/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type handlerUserTemp struct {
	DB   *sqlx.DB
	TxID string
}

// createUserTemp godoc
// @Summary Registra la sesión el usuario
// @Description Método para registra la sesión del usuario
// @tags Auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Success 200 {object} resAnny
// @Router /api/v1/user_temp/create [post]
func (h *handlerUserTemp) createUserTemp(c *fiber.Ctx) error {
	res := resAnny{Error: true}
	req := RqCreateUserTemp{}
	err := c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("no se pudo parsear la petición, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(15, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	usr, code, err := srvAuth.SrvUserTemp.GetUserTempByEmail(req.Email)
	if err != nil {
		logger.Error.Printf("No se pudo consultar el usuario, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if usr != nil {
		res.Error = false
		res.Code, res.Type, res.Msg = msg.GetByCode(5, h.DB, h.TxID)
		return c.Status(http.StatusOK).JSON(res)
	}

	_, code, err = srvAuth.SrvUserTemp.CreateUserTemp(uuid.New().String(), req.FullName, req.Surname, req.Name, req.Picture, req.Email, req.Domain)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el total del trabajo, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = "Datos registrados correctamente"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}
