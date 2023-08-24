package work

import (
	"check-id-api/internal/blockchain"
	"check-id-api/internal/env"
	"check-id-api/internal/logger"
	"check-id-api/internal/msg"
	"check-id-api/internal/send_grid"
	"check-id-api/internal/template"
	"check-id-api/pkg/auth"
	"check-id-api/pkg/cfg"
	"check-id-api/pkg/trx"
	"check-id-api/pkg/wf"
	"encoding/base64"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"net/http"
	"strconv"
	"strings"
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
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)
	srvCfg := cfg.NewServerCfg(h.DB, nil, h.TxID)

	allOnboaring, err := srvAuth.SrvOnboarding.GetAllOnboarding()
	if err != nil {
		logger.Error.Printf("No se pudo obtener el total del trabajo validado, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	allValidation, err := srvCfg.SrvValidationRequest.GetAllValidationRequest()
	if err != nil {
		logger.Error.Printf("No se pudo obtener el total del trabajo pendiente, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	usersTemp, err := srvAuth.SrvUser.GetAllNotStarted()
	if err != nil {
		logger.Error.Printf("No se pudo obtener el total del trabajo no iniciado, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	var works []*Work

	for _, onboarding := range allOnboaring {
		user, code, err := srvAuth.SrvUser.GetUsersByID(onboarding.UserId)
		if err != nil {
			logger.Error.Printf("No se pudo obtener al usaurio de la solicitud de onboarding, error: %s", err.Error())
			res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

		if user == nil {
			logger.Info.Printf("No se pudo obtener al usaurio de la solicitud de onboarding, onboarding_id: " + onboarding.ID)
			continue
		}
		works = append(works, &Work{
			Process: "Onboarding",
			User: User{
				ID:             user.ID,
				TypeDocument:   user.TypeDocument,
				DocumentNumber: user.DocumentNumber,
				ExpeditionDate: user.ExpeditionDate,
				Email:          user.Email,
				FirstName:      user.FirstName,
				SecondName:     user.SecondName,
				SecondSurname:  user.SecondSurname,
				Age:            user.Age,
				Gender:         user.Gender,
				Nationality:    user.Nationality,
				CivilStatus:    user.CivilStatus,
				FirstSurname:   user.FirstSurname,
				BirthDate:      user.BirthDate,
				Country:        user.Country,
				Department:     user.Department,
				Cellphone:      user.Cellphone,
				City:           user.City,
				RealIp:         user.RealIp,
				CreatedAt:      user.CreatedAt,
				UpdatedAt:      user.UpdatedAt,
			},
			ClientID:  onboarding.ClientId,
			RequestID: onboarding.RequestId,
			Status:    onboarding.Status,
			ExpiredAt: "-",
			CreateAt:  onboarding.CreatedAt.String(),
		})
	}

	for _, validation := range allValidation {
		user, code, err := srvAuth.SrvUser.GetUsersByID(validation.UserID)
		if err != nil {
			logger.Error.Printf("No se pudo obtener al usaurio de la solicitud de onboarding, error: %s", err.Error())
			res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

		if user == nil {
			logger.Info.Printf("No se pudo obtener al usaurio de la solicitud de validacion de identidad, validation_id: " + validation.RequestId)
			continue
		}
		works = append(works, &Work{
			Process: "Validación de identidad",
			User: User{
				ID:             user.ID,
				TypeDocument:   user.TypeDocument,
				DocumentNumber: user.DocumentNumber,
				ExpeditionDate: user.ExpeditionDate,
				Email:          user.Email,
				FirstName:      user.FirstName,
				SecondName:     user.SecondName,
				SecondSurname:  user.SecondSurname,
				Age:            user.Age,
				Gender:         user.Gender,
				Nationality:    user.Nationality,
				CivilStatus:    user.CivilStatus,
				FirstSurname:   user.FirstSurname,
				BirthDate:      user.BirthDate,
				Country:        user.Country,
				Department:     user.Department,
				Cellphone:      user.Cellphone,
				City:           user.City,
				RealIp:         user.RealIp,
				CreatedAt:      user.CreatedAt,
				UpdatedAt:      user.UpdatedAt,
			},
			ClientID:  validation.ClientId,
			RequestID: validation.RequestId,
			Status:    validation.Status,
			ExpiredAt: "-",
			CreateAt:  validation.CreatedAt.String(),
		})
	}

	for _, user := range usersTemp {
		works = append(works, &Work{
			Process: "System",
			User: User{
				ID:             user.ID,
				TypeDocument:   user.TypeDocument,
				DocumentNumber: user.DocumentNumber,
				ExpeditionDate: user.ExpeditionDate,
				Email:          user.Email,
				FirstName:      user.FirstName,
				SecondName:     user.SecondName,
				SecondSurname:  user.SecondSurname,
				Age:            user.Age,
				Gender:         user.Gender,
				Nationality:    user.Nationality,
				CivilStatus:    user.CivilStatus,
				FirstSurname:   user.FirstSurname,
				BirthDate:      user.BirthDate,
				Country:        user.Country,
				Department:     user.Department,
				Cellphone:      user.Cellphone,
				City:           user.City,
				RealIp:         user.RealIp,
				CreatedAt:      user.CreatedAt,
				UpdatedAt:      user.UpdatedAt,
			},
			ClientID:  -1,
			RequestID: "-",
			Status:    "No iniciado",
			ExpiredAt: "-",
			CreateAt:  user.CreatedAt.String(),
		})
	}

	res.Data = works
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
	param := make(map[string]string)
	var mailAttachment []*mail.Attachment
	e := env.NewConfiguration()
	err := c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("el id del usuario es requerido")
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	srvWf := wf.NewServerWf(h.DB, nil, h.TxID)
	srvTrx := trx.NewServerTrx(h.DB, nil, h.TxID)
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)
	srvCfg := cfg.NewServerCfg(h.DB, nil, h.TxID)

	status, err := srvWf.SrvWork.GetAllWorkValidationByStatus("ok")
	if err != nil {
		logger.Error.Printf("No se pudo consultar el estado, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if status == nil {
		res.Code, res.Type, res.Msg = msg.GetByCode(116, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	user, code, err := srvAuth.SrvUser.GetUsersByID(req.UserID)
	if err != nil {
		logger.Error.Printf("No se pudo consultar el usuario, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	description := "Validación de identidad del usuario"
	identifier := []blockchain.Identifier{
		{
			Name: "Información básica",
			Attributes: []blockchain.Attribute{
				{
					Id:    1,
					Name:  "Primer Nombre",
					Value: strings.TrimSpace(*user.FirstName),
				},
				{
					Id:    2,
					Name:  "Segundo Nombre",
					Value: strings.TrimSpace(*user.SecondName),
				},
				{
					Id:    3,
					Name:  "Primer Apellido",
					Value: strings.TrimSpace(*user.FirstSurname),
				},
				{
					Id:    4,
					Name:  "Segundo Apellido",
					Value: strings.TrimSpace(*user.SecondSurname),
				},
				{
					Id:    5,
					Name:  "Tipo de Documento",
					Value: *user.TypeDocument,
				},
				{
					Id:    6,
					Name:  "Número de Documento",
					Value: user.DocumentNumber,
				},
				{
					Id:    7,
					Name:  "Correo Electrónico",
					Value: user.Email,
				},
				{
					Id:    8,
					Name:  "Edad",
					Value: strconv.Itoa(int(*user.Age)),
				},
				{
					Id:    9,
					Name:  "Sexo",
					Value: *user.Gender,
				},
				{
					Id:    10,
					Name:  "Fecha de Nacimiento",
					Value: user.BirthDate.String(),
				},
				{
					Id:    11,
					Name:  "Fecha de Expedición del Documento",
					Value: user.ExpeditionDate.UTC().String(),
				},
				{
					Id:    12,
					Name:  "Estado Civil",
					Value: *user.CivilStatus,
				},
				{
					Id:    13,
					Name:  "IP de Dispositivo",
					Value: user.RealIp,
				},
				{
					Id:    14,
					Name:  "Nacionalidad",
					Value: *user.Nationality,
				},
				{
					Id:    15,
					Name:  "Fecha de Creación",
					Value: user.CreatedAt.UTC().String(),
				},
			},
		},
	}

	file, code, err := srvCfg.SrvFiles.GetFilesByTypeAndUserID(1, user.ID)
	if err != nil {
		logger.Error.Printf("No se pudo obtener la foto del usuario, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	fileS3, code, err := srvCfg.SrvFilesS3.GetFileByPath(file.Path, file.Name)
	if err != nil {
		logger.Error.Printf("No se pudo obtener la foto del usuario de S3, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	walletInfo, err := blockchain.CreateAccountAndWallet(user, fileS3.Encoding, fileS3.NameDocument)
	if err != nil {
		logger.Error.Printf("No se pudo crear el usuario en OnlyOne, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	trxId, err := blockchain.CreateTransaction(identifier, "Validación de identidad", description, walletInfo.Id, "")
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

	code, err = srvWf.SrvStatusReq.UpdateStatusRequestByUserId("validado", "La información del usuario ha sido validado correctamente", req.UserID)
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

	fullName := strings.TrimSpace(*user.FirstName + " " + *user.SecondName + " " + *user.FirstSurname + " " + *user.SecondSurname)

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
// @Param ReqRefused body ReqRefused true "Datos de solicitud para el rechazo"
// @Success 200 {object} resAnny
// @Router /api/v1/work/refused [post]
func (h *handlerWork) refusedUserData(c *fiber.Ctx) error {
	res := resAnny{Error: true}
	req := ReqRefused{}
	err := c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("el id del usuario es requerido")
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	srvWf := wf.NewServerWf(h.DB, nil, h.TxID)
	srvTrx := trx.NewServerTrx(h.DB, nil, h.TxID)

	code, err := srvWf.SrvStatusReq.UpdateStatusRequestByUserId("rechazado", req.Motivo, req.UserID)
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
