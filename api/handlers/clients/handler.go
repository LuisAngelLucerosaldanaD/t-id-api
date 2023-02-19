package clients

import (
	"check-id-api/internal/logger"
	"check-id-api/internal/msg"
	"check-id-api/pkg/cfg"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	"time"
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

// CreateClient godoc
// @Summary Crea el cliente en el sistema
// @Description Método para crear el cliente en el sistema
// @tags Client
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param Client body Client true "Datos del cliente a crear"
// @Success 200 {object} ResAnny
// @Router /api/v1/clients [post]
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

	client, code, err := srvCfg.SrvClients.GetClientsByNit(req.Nit)
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

	fileBanner, code, err := srvCfg.SrvFiles.CreateFiles(fileS3Banner.Path, fileS3Banner.FileName, 6, "2435b1d2-6e0a-4541-a3a5-810b22e961d1")
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

	fileSmall, code, err := srvCfg.SrvFiles.CreateFiles(fileS3Small.Path, fileS3Small.FileName, 7, "2435b1d2-6e0a-4541-a3a5-810b22e961d1")
	if err != nil {
		logger.Error.Printf("No se pudo cargar la data del logo small, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvCfg.SrvClients.CreateClients(req.FullName, req.Nit, strconv.FormatInt(fileBanner.ID, 10), strconv.FormatInt(fileSmall.ID, 10), req.MainColor, req.SecondColor, req.UrlRedirect, req.UrlApi)
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

// GetValidationRequest godoc
// @Summary Obtiene el flujo de validación de un usuario
// @Description Método para obtener el flujo de validación de identidad de un usuario
// @tags Client
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param nit path string true "NIT del cliente"
// @Param request_id path string true "Número de solicitud"
// @Param document_number path string true "Número de identificación del usuario"
// @Success 200 {object} ResAnny
// @Router /api/v1/validation-workflow/{nit}/{request_id}/{document_number} [get]
func (h *handlerWork) GetValidationRequest(c *fiber.Ctx) error {
	res := ResAnny{Error: true}
	srvCfg := cfg.NewServerCfg(h.DB, nil, h.TxID)
	nit := c.Params("nit")
	requestID := c.Params("request_id")
	documentNumber := c.Params("document_number")
	if nit == "" || requestID == "" {
		logger.Error.Printf("No se pudo parasear los valores de busqueda")
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
		logger.Error.Printf("No se pudo obtener el cliente")
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	validationRequest, code, err := srvCfg.SrvValidationRequest.GetValidationRequestByClientIDAndRequestID(client.ID, requestID)
	if err != nil {
		logger.Error.Printf("No se pudo obtener la configuracion de la validación de identidad, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if validationRequest == nil {
		logger.Error.Printf("No se pudo obtener la data del flujo de validación de identidad")
		res.Code, res.Type, res.Msg = 404, 1, "No se encontró información sobre el flujo de validación de identidad con los datos ingresados"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if validationRequest.UserIdentification != documentNumber {
		res.Code, res.Type, res.Msg = 22, 1, "Este usuario no tiene configurado una solicitud de validación de identidad en el flujo"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	dateExpired := validationRequest.ExpiredAt.Sub(time.Now())
	if dateExpired.Minutes() <= 0 {
		res.Code, res.Type, res.Msg = 22, 1, "La fecha para validar la identidad a caducado"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if validationRequest.MaxNumValidation == 0 {
		res.Code, res.Type, res.Msg = 22, 1, "Se ha superado la cantidad máximas de consultas configuradas para este flujo de validación de identidad"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = validationRequest
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}
