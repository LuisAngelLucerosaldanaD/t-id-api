package onboarding

import (
	"check-id-api/internal/aws_ia"
	"check-id-api/internal/blockchain"
	"check-id-api/internal/env"
	"check-id-api/internal/logger"
	"check-id-api/internal/msg"
	"check-id-api/internal/persons"
	"check-id-api/internal/send_grid"
	"check-id-api/internal/template"
	"check-id-api/pkg/auth"
	"check-id-api/pkg/cfg"
	"check-id-api/pkg/trx"
	"check-id-api/pkg/wf"
	"encoding/base64"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"net/http"
	"strings"
	"time"
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
	srvCfg := cfg.NewServerCfg(h.DB, nil, h.TxID)

	err := c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("couldn't bind model start onboarding: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	user, code, err := srvAuth.SrvUser.GetUsersByEmail(req.Email)
	if err != nil {
		logger.Error.Printf("couldn't bind get user by email: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if user != nil {
		onboarding, code, err := srvAuth.SrvOnboarding.GetOnboardingByUserID(user.ID)
		if err != nil {
			logger.Error.Printf("couldn't bind get onboarding by user id: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

		if onboarding != nil && onboarding.Status == "finished" {
			// TODO validar el número máximo de las consultas de validación y el tiempo de vida
			ttl := time.Now().AddDate(0, 0, 3)
			validation, code, err := srvCfg.SrvValidationRequest.CreateValidationRequest(req.ClientId, 3, req.RequestId, ttl, user.DocumentNumber, "pending")
			if err != nil {
				logger.Error.Printf("couldn't bind create validation request: %v", err)
				res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
				return c.Status(http.StatusAccepted).JSON(res)
			}

			// TODO validar el cifrado de datos
			res.Data = e.OnlyOne.Url + e.OnlyOne.Onboarding + fmt.Sprintf("?validation_id=%d&user_id=%s&email=%s&process=validation", validation.ID, user.ID, req.Email)
			res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
			res.Error = false
			return c.Status(http.StatusOK).JSON(res)
		}
		if onboarding != nil && onboarding.Status == "pending" {
			res.Data = e.OnlyOne.Url + e.OnlyOne.Onboarding + fmt.Sprintf("?onboarding_id=%s&user_id=%s&email=%s&process=enrolamiento", onboarding.ID, user.ID, req.Email)
			res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
			res.Error = false
			return c.Status(http.StatusOK).JSON(res)
		}
	}

	user, code, err = srvAuth.SrvUser.CreateUsers(uuid.New().String(), nil, req.DocumentNumber, nil, req.Email, req.FirstName, req.SecondName, req.SecondSurname, nil, nil, req.Nationality, nil, req.FirstSurname, nil, nil, nil, nil, c.IP(), req.Cellphone)
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

	res.Data = e.OnlyOne.Url + e.OnlyOne.Onboarding + fmt.Sprintf("?onboarding_id=%s&user_id=%s&email=%s&process=enrolamiento", onboarding.ID, user.ID, req.Email)
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
	e := env.NewConfiguration()

	err := c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("couldn't bind model create wallets: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
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

	if onboarding.Status != "started" {
		res.Code, res.Type, res.Msg = 22, 1, "El usuario ya ha finalizado el proceso de enrolamiento"
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

	user, code, err := srvAuth.SrvUser.GetUsersByID(req.UserID)
	if err != nil {
		logger.Error.Printf("couldn't get user by identity number, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if user == nil {
		logger.Error.Printf("couldn't get user by identity number")
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	f, err := srvCfg.SrvFilesS3.UploadFile(req.UserID, req.UserID+"_selfie.jpg", req.Selfie)
	if err != nil {
		logger.Error.Printf("couldn't upload file s3: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvCfg.SrvFiles.CreateFiles(f.Path, f.FileName, 1, req.UserID)
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

	_, code, err = srvAuth.SrvOnboarding.UpdateOnboarding(onboarding.ID, onboarding.ClientId, onboarding.RequestId, onboarding.ID, "pending")
	if err != nil {
		logger.Error.Printf("couldn't create onboarding, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	personSrv := persons.Persons{IdentityNumber: user.DocumentNumber}
	basicData, err := personSrv.GetPersonByIdentityNumber()
	if err != nil {
		logger.Error.Printf("couldn't get basic data of person, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	birthDate, _ := time.Parse("02/01/2006", basicData.BirthDate)
	expeditionDate, _ := time.Parse("02/01/2006", basicData.ExpeditionDate)
	age := int32(time.Now().Year() - birthDate.Year())

	user, code, err = srvAuth.SrvUser.UpdateUsers(user.ID, nil, user.DocumentNumber, &expeditionDate,
		user.Email, &basicData.FirstName, &basicData.SecondName, &basicData.SecondSurname, &age, &basicData.Gender,
		user.Nationality, nil, &basicData.Surname, &birthDate, nil, nil, nil, user.RealIp, user.Cellphone)
	if err != nil {
		logger.Error.Printf("couldn't update basic data of user, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	walletInfo, err := blockchain.CreateAccountAndWallet(user, req.Selfie, req.UserID+"_selfie.jpg")
	if err != nil {
		logger.Error.Printf("No se pudo crear el usuario en OnlyOne, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	description := "Validación de identidad del usuario"
	trxId, err := blockchain.CreateTransactionV2(user, "Validación de identidad", description, walletInfo.Id, "")
	if err != nil {
		logger.Error.Printf("No se pudo consultar el usuario, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvAuth.SrvValidationUsers.CreateValidationUsers(uuid.New().String(), trxId, req.UserID)
	if err != nil {
		logger.Error.Printf("No se pudo registrar el id de la transacción, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
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

	fullName := strings.TrimSpace(*user.FirstName + " " + *user.SecondName + " " + *user.FirstSurname + " " + *user.SecondSurname)
	param := make(map[string]string)
	var mailAttachment []*mail.Attachment
	param["TEMPLATE-PATH"] = e.Template.WalletMail
	param["FULL_NAME"] = fullName
	param["WALLET_ID"] = walletInfo.Id

	body, err := template.GenerateTemplateMail(param)
	if err != nil {
		logger.Error.Printf("couldn't generate body in NotificationEmail: %v", err)
	}

	filePrivate := mail.NewAttachment()
	filePrivate.SetContent(base64.StdEncoding.EncodeToString([]byte(walletInfo.Private)))
	filePrivate.SetType("text/plain")
	filePrivate.SetFilename("private.pem")
	filePrivate.SetDisposition("attachment")
	mailAttachment = append(mailAttachment, filePrivate)

	filePublic := mail.NewAttachment()
	filePublic.SetContent(base64.StdEncoding.EncodeToString([]byte(walletInfo.Public)))
	filePublic.SetType("text/plain")
	filePublic.SetFilename("public.pem")
	filePublic.SetDisposition("attachment")
	mailAttachment = append(mailAttachment, filePublic)

	fileMnemonic := mail.NewAttachment()
	fileMnemonic.SetContent(base64.StdEncoding.EncodeToString([]byte(walletInfo.Mnemonic)))
	fileMnemonic.SetType("text/plain")
	fileMnemonic.SetFilename("mnemonic.txt")
	fileMnemonic.SetDisposition("attachment")
	mailAttachment = append(mailAttachment, fileMnemonic)

	emailSd := send_grid.Model{
		FromMail: e.SendGrid.FromMail,
		FromName: e.SendGrid.FromName,
		Tos: []send_grid.To{
			{
				Name: fullName,
				Mail: user.Email,
			},
		},
		Subject:     "Certificados públicos y privados OnlyOne",
		HTMLContent: body,
		Attachments: mailAttachment,
	}

	err = emailSd.SendMail()
	if err != nil {
		logger.Error.Println(h.TxID, " - error al enviar el correo con las credenciales de la wallet: %s", err.Error())
	}

	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// ValidateIdentity godoc
// @Summary Método que permite finalizar la validación de identidad de un usuario
// @Description Método que permite finalizar la validación de identidad de un usuario por la aplicación de OnlyOne
// @tags Onboarding
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param RequestValidationIdentity body RequestValidationIdentity true "Datos para validar la identidad del usuario"
// @Success 200 {object} ResProcessOnboarding
// @Router /api/v1/onboarding/validate-identity [post]
func (h *handlerOnboarding) ValidateIdentity(c *fiber.Ctx) error {
	res := ResProcessOnboarding{Error: true}
	req := RequestValidationIdentity{}
	srvCfg := cfg.NewServerCfg(h.DB, nil, h.TxID)

	err := c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("couldn't bind model create wallets: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	onboarding, code, err := srvCfg.SrvValidationRequest.GetValidationRequestByID(req.ValidationId)
	if err != nil {
		logger.Error.Printf("couldn't bind get validation request: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if onboarding == nil {
		logger.Error.Printf("couldn't bind get validation request")
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if onboarding.Status != "pending" {
		logger.Error.Printf("La validación de identidad ya ha sido realizada o no existe")
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	selfieBytes, err := base64.StdEncoding.DecodeString(req.FaceImage)
	if err != nil {
		logger.Error.Printf("couldn't decode selfie: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	fileSelfie, code, err := srvCfg.SrvFiles.GetFilesByTypeAndUserID(1, req.UserID)
	if err != nil {
		logger.Error.Printf("no se pudo obtener la imagen del documento de identidad, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if fileSelfie == nil {
		res.Code, res.Type, res.Msg = 22, 1, "El usuario no ha cargado su documento de identidad aun"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	documentB64, code, err := srvCfg.SrvFilesS3.GetFileByPath(fileSelfie.Path, fileSelfie.Name)
	if err != nil {
		logger.Error.Printf("no se pudo obtener la imagen del documento de identidad, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	selfieStorageBytes, err := base64.StdEncoding.DecodeString(documentB64.Encoding)
	if err != nil {
		logger.Error.Printf("couldn't decode document front: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	resp, err := aws_ia.CompareFacesV2(selfieBytes, selfieStorageBytes)
	if err != nil {
		logger.Error.Printf("couldn't decode identity: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if !resp {
		res.Code, res.Type, res.Msg = 109, 1, "La persona no coincide con su documento de identidad"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvCfg.SrvValidationRequest.UpdateStatusValidationRequest(req.ValidationId, "callback")
	if err != nil {
		logger.Error.Printf("couldn't bind update validation request: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}
