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
	user2 "check-id-api/pkg/auth/user"
	"check-id-api/pkg/cfg"
	"check-id-api/pkg/trx"
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

	err := c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("couldn't bind model start onboarding: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	user, code, err := srvAuth.SrvUser.GetUserByEmail(req.Email)
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

		if onboarding != nil && (onboarding.Status == "finished" || onboarding.Status == "pending") {
			// TODO validar el número máximo de las consultas de validación y el tiempo de vida
			ttl := time.Now().AddDate(0, 0, 1)
			validation, code, err := srvAuth.SrvLifeTest.CreateLifeTest(req.ClientId, 1, req.RequestId, ttl, user.ID, "pending")
			if err != nil {
				logger.Error.Printf("couldn't bind create validation request: %v", err)
				res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
				return c.Status(http.StatusAccepted).JSON(res)
			}

			// TODO validar el cifrado de datos
			res.Data = &Onboarding{
				Url:    e.OnlyOne.Url + e.OnlyOne.Onboarding + fmt.Sprintf("?process=validation&validation_id=%d&user_id=%s&email=%s", validation.ID, user.ID, req.Email),
				Method: "validation",
			}
			res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
			res.Error = false
			return c.Status(http.StatusOK).JSON(res)
		}
		if onboarding != nil && onboarding.Status == "started" {
			res.Data = &Onboarding{
				Url:    e.OnlyOne.Url + e.OnlyOne.Onboarding + fmt.Sprintf("?process=enroll&onboarding_id=%s&user_id=%s&email=%s", onboarding.ID, user.ID, req.Email),
				Method: "onboarding",
			}
			res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
			res.Error = false
			return c.Status(http.StatusOK).JSON(res)
		}

		if onboarding == nil {
			_, err = srvAuth.SrvUser.DeleteUser(user.ID)
			if err != nil {
				logger.Error.Printf("No se pudo eliminar el usuario, error: %v", err)
				res.Code, res.Type, res.Msg = 108, 1, "No se pudo eliminar el usuario, error: "+err.Error()
				return c.Status(http.StatusOK).JSON(res)
			}
		}
	}

	user, code, err = srvAuth.SrvUser.CreateUser(&user2.User{
		ID:             uuid.New().String(),
		Nickname:       req.Email,
		Email:          req.Email,
		Password:       req.DocumentNumber,
		FirstName:      req.FirstName,
		SecondName:     req.SecondName,
		FirstSurname:   req.FirstSurname,
		SecondSurname:  req.SecondSurname,
		DocumentNumber: req.DocumentNumber,
		Cellphone:      req.Cellphone,
		Nationality:    req.Nationality,
		RealIp:         c.IP(),
		StatusId:       0,
		FailedAttempts: 0,
		IsDeleted:      false,
	})
	if err != nil {
		logger.Error.Printf("couldn't create user, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvAuth.SrvUserRole.CreateUseRole(uuid.New().String(), user.ID, "14cbf8d2-485a-4fbe-baa6-16273c765f14")
	if err != nil {
		logger.Error.Printf("couldn't create user role, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvTrx.SrvTraceability.CreateTraceability("Registro", "info", "Registro de información básica", user.ID)
	if err != nil {
		logger.Error.Printf("couldn't create traceability, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	onboarding, code, err := srvAuth.SrvOnboarding.CreateOnboarding(uuid.New().String(), req.ClientId, req.RequestId, user.ID, "started", "")
	if err != nil {
		logger.Error.Printf("couldn't create onboarding, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = &Onboarding{
		Url:    e.OnlyOne.Url + e.OnlyOne.Onboarding + fmt.Sprintf("?process=enroll&onboarding_id=%s&user_id=%s&email=%s", onboarding.ID, user.ID, req.Email),
		Method: "onboarding",
	}
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
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if !resp {
		_, code, err = srvTrx.SrvTraceability.CreateTraceability("Validación de identidad", "error", "La validación de identidad fue rechazada", req.UserID)
		if err != nil {
			logger.Error.Printf("couldn't create traceability, error: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

		_, code, err = srvAuth.SrvOnboarding.UpdateOnboarding(onboarding.ID, onboarding.ClientId, onboarding.RequestId, onboarding.ID, "refused", "")
		if err != nil {
			logger.Error.Printf("couldn't update onboarding, error: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

		res.Error = false
		res.Code, res.Type, res.Msg = 109, 1, "Se ha procesado la información de enrolamiento"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	user, code, err := srvAuth.SrvUser.GetUserByID(req.UserID)
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

	_, code, err = srvCfg.SrvFiles.CreateFile(f.Path, f.FileName, 1, req.UserID)
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

	_, code, err = srvCfg.SrvFiles.CreateFile(f.Path, f.FileName, 2, req.UserID)
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

	_, code, err = srvCfg.SrvFiles.CreateFile(f.Path, f.FileName, 3, req.UserID)
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

	personSrv := persons.Persons{IdentityNumber: user.DocumentNumber}
	basicData, err := personSrv.GetPersonByIdentityNumber()
	if err != nil {
		logger.Error.Printf("couldn't get basic data of person, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	birthDate, _ := time.Parse("02-01-2006", basicData.BirthDate)
	age := int32(time.Now().Year() - birthDate.Year())
	nationality := "Colombia"

	user, code, err = srvAuth.SrvUser.UpdateUser(&user2.User{
		ID:             user.ID,
		Nickname:       user.Email,
		Email:          user.Email,
		Password:       strings.TrimSpace(strings.ToLower(basicData.FirstName) + user.DocumentNumber),
		FirstName:      &basicData.FirstName,
		SecondName:     &basicData.SecondName,
		FirstSurname:   &basicData.Surname,
		SecondSurname:  &basicData.SecondSurname,
		Age:            &age,
		DocumentNumber: user.DocumentNumber,
		Cellphone:      user.Cellphone,
		Gender:         &basicData.Gender,
		Nationality:    &nationality,
		RealIp:         user.RealIp,
		StatusId:       0,
		FailedAttempts: 0,
		BirthDate:      &birthDate,
		IsDeleted:      false,
	})
	if err != nil {
		logger.Error.Printf("couldn't update basic data of user, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvAuth.SrvUserRole.UpdateUseRoleByUserID(user.ID, "09ecd353-ee2b-42a0-88bc-45d118d7b65d")
	if err != nil {
		logger.Error.Printf("couldn't update user role, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvTrx.SrvTraceability.CreateTraceability("Actualización de datos", "info", "Actualización de los datos personales", req.UserID)
	if err != nil {
		logger.Error.Printf("couldn't create traceability, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	walletInfo, err := blockchain.CreateWallet(user)
	if err != nil {
		logger.Error.Printf("No se pudo crear el usuario en OnlyOne, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	description := "Validación de identidad del usuario"
	trxId, err := blockchain.CreateTransaction(user, "Validación de identidad", description, walletInfo.Id)
	if err != nil {
		logger.Error.Printf("No se pudo consultar el usuario, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvAuth.SrvOnboarding.UpdateOnboarding(onboarding.ID, onboarding.ClientId, onboarding.RequestId, onboarding.UserId, "pending", trxId)
	if err != nil {
		logger.Error.Printf("couldn't create onboarding, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvTrx.SrvTraceability.CreateTraceability("Validación de identidad", "success", "Validación de identidad aprobada", req.UserID)
	if err != nil {
		logger.Error.Printf("couldn't create traceability, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
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

// FinishOnboardingV2 godoc
// @Summary Método que permite terminar el enrolamiento de un usuario
// @Description Método que permite terminar el enrolamiento de un usuario que ha sido validado desde OnlyOne
// @tags Onboarding
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param RequestProcessOnboarding body RequestProcessOnboarding true "Datos para validar el enrolamiento del usuario"
// @Success 200 {object} ResProcessOnboarding
// @Router /api/v1/onboarding/process [post]
func (h *handlerOnboarding) FinishOnboardingV2(c *fiber.Ctx) error {
	res := ResProcessOnboarding{Error: true}
	req := RequestProcessOnboarding{}
	srvTrx := trx.NewServerTrx(h.DB, nil, h.TxID)
	srvCfg := cfg.NewServerCfg(h.DB, nil, h.TxID)
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

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

	user, code, err := srvAuth.SrvUser.GetUserByID(req.UserID)
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

	_, code, err = srvCfg.SrvFiles.CreateFile(f.Path, f.FileName, 1, req.UserID)
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

	_, code, err = srvCfg.SrvFiles.CreateFile(f.Path, f.FileName, 2, req.UserID)
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

	_, code, err = srvCfg.SrvFiles.CreateFile(f.Path, f.FileName, 3, req.UserID)
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

	_, code, err = srvAuth.SrvOnboarding.UpdateOnboarding(onboarding.ID, onboarding.ClientId, onboarding.RequestId, onboarding.UserId, "life-test", onboarding.TransactionId)
	if err != nil {
		logger.Error.Printf("couldn't update onboarding, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = "Información cargada correctamente"
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
// @Router /api/v1/onboarding/validate_identity [post]
func (h *handlerOnboarding) ValidateIdentity(c *fiber.Ctx) error {
	res := ResProcessOnboarding{Error: true}
	req := RequestValidationIdentity{}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)
	srvCfg := cfg.NewServerCfg(h.DB, nil, h.TxID)

	err := c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("couldn't bind model create wallets: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	validation, code, err := srvAuth.SrvLifeTest.GetLifeTestByID(req.ValidationId)
	if err != nil {
		logger.Error.Printf("couldn't bind get validation request: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if validation == nil {
		logger.Error.Printf("couldn't bind get validation request")
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if validation.Status != "pending" {
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

	fileSelfie, code, err := srvCfg.SrvFiles.GetFileByTypeAndUserID(1, req.UserID)
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
		_, code, err = srvAuth.SrvLifeTest.UpdateStatusLifeTest(req.ValidationId, "refused")
		if err != nil {
			logger.Error.Printf("couldn't bind update validation request: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

		res.Error = false
		res.Code, res.Type, res.Msg = 109, 1, "Se ha procesado la solicitud"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvAuth.SrvLifeTest.UpdateStatusLifeTest(req.ValidationId, "callback")
	if err != nil {
		logger.Error.Printf("couldn't bind update validation request: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}
