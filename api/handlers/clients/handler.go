package clients

import (
	"check-id-api/internal/logger"
	"check-id-api/internal/msg"
	"check-id-api/pkg/cfg"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
)

type handlerWork struct {
	DB   *sqlx.DB
	TxID string
}

// getDataClient godoc
// @Summary Obtiene la data del cliente
// @Description Método para obtener la información del cliente de CheckID
// @tags Client
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param nit path string true "NIT del cliente"
// @Success 200 {object} ResClient
// @Router /api/v1/clients/{nit} [get]
func (h *handlerWork) getDataClient(c *fiber.Ctx) error {
	res := ResClient{Error: true}
	srvCfg := cfg.NewServerCfg(h.DB, nil, h.TxID)
	nit := c.Params("nit")
	if nit == "" {
		logger.Error.Printf("El nit del cliente es requerido")
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	client, code, err := srvCfg.SrvClients.GetClientsByNit(nit)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el cliente, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if client == nil {
		logger.Error.Printf("No existe un cliente asociado al NIT")
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	bannerId, _ := strconv.ParseInt(client.Banner, 10, 64)
	smallId, _ := strconv.ParseInt(client.LogoSmall, 10, 64)

	fileBanner, code, err := srvCfg.SrvFiles.GetFilesByID(bannerId)
	if err != nil {
		logger.Error.Printf("No se pudo obtener la data del banner, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	fileS3Banner, code, err := srvCfg.SrvFilesS3.GetFileByPath(fileBanner.Path, fileBanner.Name)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el banner, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	fileSmall, code, err := srvCfg.SrvFiles.GetFilesByID(smallId)
	if err != nil {
		logger.Error.Printf("No se pudo obtener la data del logo small, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	fileS3Small, code, err := srvCfg.SrvFilesS3.GetFileByPath(fileSmall.Path, fileSmall.Name)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el logo small, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = Client{
		ID:          client.ID,
		FullName:    client.FullName,
		Nit:         client.Nit,
		Banner:      fileS3Banner.Encoding,
		LogoSmall:   fileS3Small.Encoding,
		MainColor:   client.MainColor,
		SecondColor: client.SecondColor,
		UrlRedirect: client.UrlRedirect,
		UrlApi:      client.UrlApi,
	}
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}
