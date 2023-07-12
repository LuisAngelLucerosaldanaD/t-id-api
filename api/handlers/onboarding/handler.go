package onboarding

import (
	"check-id-api/internal/aws_ia"
	"check-id-api/internal/env"
	"check-id-api/internal/logger"
	"check-id-api/internal/msg"
	"check-id-api/pkg/auth"
	"check-id-api/pkg/cfg"
	"check-id-api/pkg/trx"
	"check-id-api/pkg/wf"
	"encoding/base64"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type handlerOnboarding struct {
	DB   *sqlx.DB
	TxID string
}

// Onboarding godoc
// @Summary Método que permite iniciar el enrolamiento de un usuario
// @Description Método que permite iniciar el enrolamiento de un usuario que puede ser desde un tercero o desde el mismo sistema
// @tags Onboarding
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param resCreateOnboarding body resCreateOnboarding true "Datos para el enrolamiento del usuario"
// @Success 200 {object} resCreateOnboarding
// @Router /api/v1/onboarding/ [post]
func (h *handlerOnboarding) Onboarding(c *fiber.Ctx) error {
	e := env.NewConfiguration()
	res := resCreateOnboarding{Error: true}
	req := requestCreateOnboarding{}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)
	srvTrx := trx.NewServerTrx(h.DB, nil, h.TxID)
	srvWf := wf.NewServerWf(h.DB, nil, h.TxID)

	err := c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("couldn't bind model create wallets: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	user, code, err := srvAuth.SrvUser.CreateUsers(uuid.New().String(), nil, req.DocumentNumber, nil, req.Email, req.FirstName, req.SecondName, req.SecondSurname, nil, nil, req.Nationality, nil, req.FirstSurname, nil, nil, nil, nil, c.IP(), req.Cellphone)
	if err != nil {
		logger.Error.Printf("couldn't create user, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvTrx.SrvTraceability.CreateTraceability("Registro", "info", "Registro de información básica", user.ID)
	if err != nil {
		logger.Error.Printf("couldn't create traceability, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvWf.SrvStatusReq.CreateStatusRequest("pendiente", "Pendiente por validación de identidad", user.ID)
	if err != nil {
		logger.Error.Printf("couldn't create status request, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvWf.SrvWork.CreateWorkValidation("pending", user.ID)
	if err != nil {
		logger.Error.Printf("couldn't start work, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	onboarding, code, err := srvAuth.SrvOnboarding.CreateOnboarding(uuid.New().String(), req.ClientId, req.RequestId, user.ID, "started")
	if err != nil {
		logger.Error.Printf("couldn't create onboarding, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = e.OnlyOne.Url + e.OnlyOne.Onboarding + fmt.Sprintf("%s/%s/%s", onboarding.ID, user.ID, req.Email)
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// FinishOnboarding godoc
// @Summary Método que permite terminar el enrolamiento de un usuario
// @Description Método que permite terminar el enrolamiento de un usuario que ha sido validado desde OnlyOne
// @tags Onboarding
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param RequestProcessOnboarding body RequestProcessOnboarding true "Datos para validar el enrolamiento del usuario"
// @Success 200 {object} ResProcessOnboarding
// @Router /api/v1/onboarding/process [post]
func (h *handlerOnboarding) FinishOnboarding(c *fiber.Ctx) error {
	res := ResProcessOnboarding{Error: true}
	req := RequestProcessOnboarding{}
	srvTrx := trx.NewServerTrx(h.DB, nil, h.TxID)
	srvWf := wf.NewServerWf(h.DB, nil, h.TxID)
	srvCfg := cfg.NewServerCfg(h.DB, nil, h.TxID)
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	err := c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("couldn't bind model create wallets: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	selfieBytes, err := base64.StdEncoding.DecodeString(req.Selfie)
	if err != nil {
		logger.Error.Printf("couldn't decode selfie: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	documentFrontBytes, err := base64.StdEncoding.DecodeString(req.DocumentFront)
	if err != nil {
		logger.Error.Printf("couldn't decode document front: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	resp, err := aws_ia.CompareFacesV2(selfieBytes, documentFrontBytes)
	if err != nil {
		logger.Error.Printf("couldn't decode identity: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if !resp {
		res.Code, res.Type, res.Msg = 109, 1, "La persona no coincide con su documento de identidad"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	f, err := srvCfg.SrvFilesS3.UploadFile(req.UserID, req.UserID+"_selfie.jpg", req.Selfie)
	if err != nil {
		logger.Error.Printf("couldn't upload file s3: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err := srvCfg.SrvFiles.CreateFiles(f.Path, f.FileName, 1, req.UserID)
	if err != nil {
		logger.Error.Printf("couldn't create selfie image, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		res.Msg = err.Error()
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvTrx.SrvTraceability.CreateTraceability("Carga de Selfie", "info", "Carga de la imagen de Selfie", req.UserID)
	if err != nil {
		logger.Error.Printf("couldn't create traceability, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	f, err = srvCfg.SrvFilesS3.UploadFile(req.UserID, req.UserID+"_doc_front.jpg", req.DocumentFront)
	if err != nil {
		logger.Error.Printf("couldn't upload file document front to s3: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvCfg.SrvFiles.CreateFiles(f.Path, f.FileName, 2, req.UserID)
	if err != nil {
		logger.Error.Printf("couldn't create file document front, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		res.Msg = err.Error()
		return c.Status(http.StatusAccepted).JSON(res)
	}

	f, err = srvCfg.SrvFilesS3.UploadFile(req.UserID, req.UserID+"_doc_back.jpg", req.DocumentBack)
	if err != nil {
		logger.Error.Printf("couldn't upload file document back to s3: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvCfg.SrvFiles.CreateFiles(f.Path, f.FileName, 3, req.UserID)
	if err != nil {
		logger.Error.Printf("couldn't create file document back, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		res.Msg = err.Error()
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvTrx.SrvTraceability.CreateTraceability("Carga del documento", "info", "Carga de documento de identidad", req.UserID)
	if err != nil {
		logger.Error.Printf("couldn't create traceability, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvTrx.SrvTraceability.CreateTraceability("Validación de identidad", "success", "Validación de identidad aprobada", req.UserID)
	if err != nil {
		logger.Error.Printf("couldn't create traceability, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvWf.SrvWork.CreateWorkValidation("finished", req.UserID)
	if err != nil {
		logger.Error.Printf("couldn't start work, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		res.Msg = err.Error()
		return c.Status(http.StatusAccepted).JSON(res)
	}

	onboarding, code, err := srvAuth.SrvOnboarding.GetOnboardingByID(req.Onboarding)
	if err != nil {
		logger.Error.Printf("couldn't get onboarding, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if onboarding == nil {
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvAuth.SrvOnboarding.UpdateOnboarding(onboarding.ID, onboarding.ClientId, onboarding.RequestId, onboarding.ID, "pending")
	if err != nil {
		logger.Error.Printf("couldn't create onboarding, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}
