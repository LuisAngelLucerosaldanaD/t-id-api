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
// @Router /api/v1/client/{nit} [get]
func (h *handlerWork) getDataClient(c *fiber.Ctx) error {
	res := ResClient{Error: true}
	srvCfg := cfg.NewServerCfg(h.DB, nil, h.TxID)
	nit := c.Params("nit")
	if nit == "" {
		logger.Error.Printf("El nit del cliente es requerido")
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	client, code, err := srvCfg.SrvClients.GetClientByNit(nit)
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

	fileBanner, code, err := srvCfg.SrvFiles.GetFileByID(bannerId)
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

	fileSmall, code, err := srvCfg.SrvFiles.GetFileByID(smallId)
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

// CreateClient godoc
// @Summary Crea el cliente en el sistema
// @Description Método para crear el cliente en el sistema
// @tags Client
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param Client body Client true "Datos del cliente a crear"
// @Success 200 {object} ResAnny
// @Router /api/v1/client [post]
func (h *handlerWork) CreateClient(c *fiber.Ctx) error {
	res := ResAnny{Error: true}
	srvCfg := cfg.NewServerCfg(h.DB, nil, h.TxID)
	req := Client{}
	err := c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("No se pudo parsear el cuerpo de la petición, erro: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	client, code, err := srvCfg.SrvClients.GetClientByNit(req.Nit)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el cliente, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if client != nil {
		logger.Error.Printf("Ya existe un cliente creado con la información proporcionada")
		res.Code, res.Type, res.Msg = msg.GetByCode(5, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	fileS3Banner, err := srvCfg.SrvFilesS3.UploadFile(req.Nit, "banner.png", req.Banner)
	if err != nil {
		logger.Error.Printf("No se pudo cargar el banner, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	fileBanner, code, err := srvCfg.SrvFiles.CreateFile(fileS3Banner.Path, fileS3Banner.FileName, 6, "2435b1d2-6e0a-4541-a3a5-810b22e961d1")
	if err != nil {
		logger.Error.Printf("No se pudo crear la data del banner, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	fileS3Small, err := srvCfg.SrvFilesS3.UploadFile(req.Nit, "logo_small.png", req.LogoSmall)
	if err != nil {
		logger.Error.Printf("No se pudo cargar el logo small, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(5, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	fileSmall, code, err := srvCfg.SrvFiles.CreateFile(fileS3Small.Path, fileS3Small.FileName, 7, "2435b1d2-6e0a-4541-a3a5-810b22e961d1")
	if err != nil {
		logger.Error.Printf("No se pudo cargar la data del logo small, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvCfg.SrvClients.CreateClient(req.FullName, req.Nit, strconv.FormatInt(fileBanner.ID, 10), strconv.FormatInt(fileSmall.ID, 10), req.MainColor, req.SecondColor, req.UrlRedirect, req.UrlApi)
	if err != nil {
		logger.Error.Printf("No se pudo crear al cliente, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = "Cliente creado correctamente"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}
