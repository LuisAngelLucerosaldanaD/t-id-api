package users

import (
	"check-id-api/internal/aws_ia"
	"check-id-api/internal/logger"
	"check-id-api/internal/msg"
	"check-id-api/internal/ws"
	"check-id-api/pkg/auth"
	"check-id-api/pkg/auth/users"
	"check-id-api/pkg/cfg"
	"check-id-api/pkg/trx"
	"check-id-api/pkg/wf"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	"time"
)

type handlerUser struct {
	DB   *sqlx.DB
	TxID string
}

// uploadSelfie godoc
// @Summary Carga de selfie del usuario
// @Description Método para cargar la selfie del usuario
// @tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param UploadSelfie body reqUploadSelfie true "Selfie del usuario"
// @Success 200 {object} responseAnny
// @Router /api/v1/user/upload-selfie [post]
func (h *handlerUser) uploadSelfie(c *fiber.Ctx) error {
	res := responseAnny{Error: true}
	req := reqUploadSelfie{}
	srvCfg := cfg.NewServerCfg(h.DB, nil, h.TxID)
	srvTrx := trx.NewServerTrx(h.DB, nil, h.TxID)
	srvWf := wf.NewServerWf(h.DB, nil, h.TxID)

	err := c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("couldn't bind model create wallets: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	fileSelfie, code, err := srvCfg.SrvFiles.GetFilesByTypeAndUserID(1, req.UserID)
	if err != nil {
		logger.Error.Printf("couldn't get user file, error validation user, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if fileSelfie != nil {
		res.Code, res.Type, res.Msg = msg.GetByCode(5, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	f, err := srvCfg.SrvFilesS3.UploadFile(req.UserID, req.UserID+"_selfie.jpg", req.SelfieImg)
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

	_, code, err = srvWf.SrvStatusReq.CreateStatusRequest("pendiente", "Pendiente por carga de documento de identidad", req.UserID)
	if err != nil {
		logger.Error.Printf("couldn't create status request, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = "Selfie cargada correctamente"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// uploadDocuments godoc
// @Summary Carga del documento de identidad
// @Description Método para cargar el documento de identidad
// @tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param uploadDocument body reqUploadDocument true "Documento de identidad"
// @Success 200 {object} responseAnny
// @Router /api/v1/user/upload-documents [post]
func (h *handlerUser) uploadDocuments(c *fiber.Ctx) error {
	res := responseAnny{Error: true}
	req := reqUploadDocument{}
	srvCfg := cfg.NewServerCfg(h.DB, nil, h.TxID)
	srvTrx := trx.NewServerTrx(h.DB, nil, h.TxID)
	srvWf := wf.NewServerWf(h.DB, nil, h.TxID)

	err := c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("couldn't bind model create wallets: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	fileTemp, code, err := srvCfg.SrvFiles.GetFilesByTypeAndUserID(2, req.UserID)
	if err != nil {
		logger.Error.Printf("couldn't upload documents, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if fileTemp != nil {
		res.Code, res.Type, res.Msg = msg.GetByCode(5, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	f, err := srvCfg.SrvFilesS3.UploadFile(req.UserID, req.UserID+"_doc_front.jpg", req.DocumentFrontImg)
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

	f, err = srvCfg.SrvFilesS3.UploadFile(req.UserID, req.UserID+"_doc_back.jpg", req.DocumentBackImg)
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

	_, code, err = srvWf.SrvStatusReq.CreateStatusRequest("pendiente", "Pendiente por carga de información básica", req.UserID)
	if err != nil {
		logger.Error.Printf("couldn't create status request, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = "Documento cargado correctamente"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// registerBasicInformation godoc
// @Summary Registro de información básica
// @Description Método para el registro de los datos básicos de una persona
// @tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param BasicInformation body requestValidateIdentity true "request of validate user identity"
// @Success 200 {object} resCreateUser
// @Router /api/v1/user/basic-information [post]
func (h *handlerUser) registerBasicInformation(c *fiber.Ctx) error {
	res := resCreateUser{Error: true}
	req := requestValidateIdentity{}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)
	srvTrx := trx.NewServerTrx(h.DB, nil, h.TxID)
	srvWf := wf.NewServerWf(h.DB, nil, h.TxID)

	err := c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("couldn't bind model create wallets: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	user, code, err := srvAuth.SrvUser.UpdateUsers(req.Id, req.TypeDocument, req.DocumentNumber, req.ExpeditionDate, req.Email, req.FirstName, req.SecondName, req.SecondSurname, req.Age, req.Gender, req.Nationality, req.CivilStatus, req.FirstSurname, req.BirthDate, req.Country, req.Department, req.City, c.IP())
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
		res.Msg = err.Error()
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = user
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// createUser godoc
// @Summary Creación de un usuario
// @Description Método para crear el usuario
// @tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param BasicInformation body requestValidateIdentity true "request of validate user identity"
// @Success 200 {object} responseAnny
// @Router /api/v1/user/create [post]
func (h *handlerUser) createUser(c *fiber.Ctx) error {
	res := responseAnny{Error: true}
	req := requestValidateIdentity{}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)
	srvTrx := trx.NewServerTrx(h.DB, nil, h.TxID)
	srvWf := wf.NewServerWf(h.DB, nil, h.TxID)

	err := c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("couldn't bind model create wallets: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	user, code, err := srvAuth.SrvUser.CreateUsers(uuid.New().String(), req.TypeDocument, req.DocumentNumber, req.ExpeditionDate, req.Email, req.FirstName, req.SecondName, req.SecondSurname, req.Age, req.Gender, req.Nationality, req.CivilStatus, req.FirstSurname, req.BirthDate, req.Country, req.Department, req.City, c.IP())
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

	_, code, err = srvWf.SrvWork.CreateWorkValidation("pendiente", user.ID)
	if err != nil {
		_, _ = srvAuth.SrvUser.DeleteUsers(user.ID)
		_, _ = srvTrx.SrvTraceability.DeleteTraceabilityByUserID(user.ID)
		logger.Error.Printf("couldn't start work, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		res.Msg = err.Error()
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = "Datos registrados correctamente"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// getUserSession godoc
// @Summary Obtiene los datos registrados del usuario por su email o su id
// @Description Método para el obtener la información del usuario en sesión por su email o id
// @tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param identifier path string true "Identificador para la búsqueda del usuario"
// @Success 200 {object} resGetUserSession
// @Router /api/v1/user/user-session/{identifier} [get]
func (h *handlerUser) getUserSession(c *fiber.Ctx) error {
	res := resGetUserSession{}
	identifier := c.Params("identifier")
	var user *users.Users
	if identifier == "" {
		logger.Error.Printf("el identifier es requerido")
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)
	srvCfg := cfg.NewServerCfg(h.DB, nil, h.TxID)

	if !govalidator.IsUUID(identifier) && !govalidator.IsEmail(identifier) {
		logger.Error.Printf("el identifier no es un parámetro valido de búsqueda")
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if govalidator.IsUUID(identifier) {
		userTmp, code, err := srvAuth.SrvUser.GetUsersByID(identifier)
		if err != nil {
			logger.Error.Printf("No se pudo obtener el usuario, error: %s", err.Error())
			res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}
		user = userTmp
	} else {
		userTmp, code, err := srvAuth.SrvUser.GetUsersByEmail(identifier)
		if err != nil {
			logger.Error.Printf("No se pudo obtener el usuario, error: %s", err.Error())
			res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}
		user = userTmp
	}

	if user == nil {
		res.Error = false
		res.Code, res.Type, res.Msg = msg.GetByCode(95, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	userValidation, code, err := srvAuth.SrvValidationUsers.GetValidationUsersByUserID(user.ID)
	if err != nil {
		logger.Error.Printf("No se pudo obtener la validación del usuario, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	transactionID := ""
	selfieImg := ""
	frontDocument := ""
	backDocument := ""

	if userValidation != nil {
		transactionID = userValidation.TransactionId
	}

	files, code, err := srvCfg.SrvFiles.GetFilesByUserID(user.ID)
	if err != nil {
		logger.Error.Printf("No se pudo obtener la validación del usuario, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if files != nil {
		for _, file := range files {
			fileS3, code, err := srvCfg.SrvFilesS3.GetFileByPath(file.Path, file.Name)
			if err != nil {
				logger.Error.Printf("No se pudo descargar el archivo, error: %s", err.Error())
				res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
				return c.Status(http.StatusAccepted).JSON(res)
			}

			switch file.Type {
			case 1:
				selfieImg = fileS3.Encoding
				break
			case 2:
				frontDocument = fileS3.Encoding
				break
			default:
				backDocument = fileS3.Encoding
				break
			}
		}
	}

	res.Data = &UserValidation{
		ID:               user.ID,
		TypeDocument:     user.TypeDocument,
		DocumentNumber:   user.DocumentNumber,
		ExpeditionDate:   user.ExpeditionDate,
		Email:            user.Email,
		FirstName:        user.FirstName,
		SecondName:       user.SecondName,
		SecondSurname:    user.SecondSurname,
		Age:              user.Age,
		Gender:           user.Gender,
		Nationality:      user.Nationality,
		CivilStatus:      user.CivilStatus,
		FirstSurname:     user.FirstSurname,
		BirthDate:        user.BirthDate,
		Country:          user.Country,
		TransactionId:    transactionID,
		Department:       user.Department,
		City:             user.City,
		SelfieImg:        selfieImg,
		BackDocumentImg:  backDocument,
		FrontDocumentImg: frontDocument,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
	}
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// getLastedUsers godoc
// @Summary Obtiene los registros de usuarios
// @Description Método para el obtener los registros de los usuarios
// @tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param email path string true "Correo electrónico del usuario"
// @Param limit path string true "Cantidad de registros por consulta"
// @Param offset path string true "Inicio del conteo de los registros por consulta"
// @Success 200 {object} resGetUsersLasted
// @Router /api/v1/user/users-lasted/{email}/{limit}/{offset} [get]
func (h *handlerUser) getLastedUsers(c *fiber.Ctx) error {
	res := resGetUsersLasted{}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)
	srvWf := wf.NewServerWf(h.DB, nil, h.TxID)
	var limit, offset int

	email := c.Params("email")
	if email == "" {
		logger.Error.Printf("el email del usuario es requerido")
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	limitStr := c.Params("limit")
	if limitStr == "" {
		limit = 10
	} else {
		limit, _ = strconv.Atoi(limitStr)
	}

	offsetStr := c.Params("offset")
	if offsetStr == "" {
		offset = 0
	} else {
		offset, _ = strconv.Atoi(offsetStr)
	}

	usersLasted, err := srvAuth.SrvUser.GetAllUsersLasted(email, limit, offset)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el listado de los últimos usuarios, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	for _, user := range usersLasted {
		validation, _, err := srvWf.SrvStatusReq.GetStatusRequestByUserID(user.ID)
		if err != nil {
			logger.Error.Printf("No se pudo obtener el estado del usuario, error: %s", err.Error())
			continue
		}

		if validation == nil {
			continue
		}
		res.Data = append(res.Data, &UserStatus{
			ID:            user.ID,
			Email:         user.Email,
			FirstName:     user.FirstName,
			SecondName:    user.SecondName,
			SecondSurname: user.SecondSurname,
			FirstSurname:  user.FirstSurname,
			Status:        validation.Status,
			UpdatedAt:     user.UpdatedAt,
		})
	}

	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// getUsersDataPending godoc
// @Summary Obtiene la cantidad de usuarios que no cargaron información requerida
// @Description Método para el obtener la cantidad de usuarios que no han cargado la información básica como la selfie, el documento de identidad y la información básica
// @tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Success 200 {object} resGetUsersDataPending
// @Router /api/v1/user/data-pending [get]
func (h *handlerUser) getUsersDataPending(c *fiber.Ctx) error {
	res := resGetUsersDataPending{}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	selfieData, err := srvAuth.SrvUser.GetAllNotUploadFile(1)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el listado de usuarios que no cargaron la selfie, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	documentsData, err := srvAuth.SrvUser.GetAllNotUploadFile(2)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el listado de usuarios que no cargaron el documento, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	basicInformation, err := srvAuth.SrvUser.GetAllNotStarted()
	if err != nil {
		logger.Error.Printf("No se pudo obtener el listado de usuarios que no cargaron la información básica, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = DataPending{
		Selfie:           len(selfieData),
		Document:         len(documentsData),
		BasicInformation: len(basicInformation),
	}
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// validateUser godoc
// @Summary Verifica si el usuario ha validado su identidad
// @Description Método para verificar si el usuario ha validado su identidad
// @tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param identity_number path string true "Número de identificación del usuario"
// @Success 200 {object} responseAnny
// @Router /api/v1/user/validate/{identity_number} [get]
func (h *handlerUser) validateUser(c *fiber.Ctx) error {
	res := responseAnny{Error: true}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	documentNumber, err := strconv.ParseInt(c.Params("identity_number"), 10, 64)
	if err != nil {
		logger.Error.Printf("el numero de identificacion es incorrecto, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	user, code, err := srvAuth.SrvUser.GetUsersByIdentityNumber(documentNumber)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el usuario por su identificacion, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if user == nil {
		res.Code, res.Type, res.Msg = 5, 1, "No existe un usuario registrado con la información proporcionada"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	validation, code, err := srvAuth.SrvValidationUsers.GetValidationUsersByUserID(user.ID)
	if err != nil {
		logger.Error.Printf("No se pudo consultar la validacion de identidad, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = true
	if validation == nil {
		res.Data = false
	}

	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// validationFace godoc
// @Summary Verifica la identidad de un usuario
// @Description Método para verificar la identidad de una persona
// @tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param ReqValidationFace body ReqValidationFace true "Datos para la verificación de identidad"
// @Success 200 {object} responseAnny
// @Router /api/v1/user/validation [post]
func (h *handlerUser) validationFace(c *fiber.Ctx) error {
	res := responseAnny{Error: true}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)
	srvCfg := cfg.NewServerCfg(h.DB, nil, h.TxID)
	req := ReqValidationFace{}
	err := c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("no se pudo parsear el cuerpo de la solicitud, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	user, code, err := srvAuth.SrvUser.GetUsersByIdentityNumber(req.DocumentNumber)
	if err != nil {
		logger.Error.Printf("no se pudo obtener el usuario a validar, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if user == nil {
		res.Code, res.Type, res.Msg = 22, 1, "No hay un usuario registrado con la información proporcionada"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	userValidation, code, err := srvAuth.SrvValidationUsers.GetValidationUsersByUserID(user.ID)
	if err != nil {
		logger.Error.Printf("no se pudo obtener el id de la validación de identidad, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if userValidation == nil {
		res.Code, res.Type, res.Msg = 403, 1, "El usuario aun no ha validado su identidad en el portal"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	fileDocFront, code, err := srvCfg.SrvFiles.GetFilesByTypeAndUserID(2, user.ID)
	if err != nil {
		logger.Error.Printf("no se pudo obtener la imagen del documento de identidad, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if fileDocFront == nil {
		res.Code, res.Type, res.Msg = 22, 1, "El usuario no ha cargado su documento de identidad aun"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	documentB64, code, err := srvCfg.SrvFilesS3.GetFileByPath(fileDocFront.Path, fileDocFront.Name)
	if err != nil {
		logger.Error.Printf("no se pudo obtener la imagen del documento de identidad, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	selfieBytes, err := base64.StdEncoding.DecodeString(req.FaceImage)
	if err != nil {
		logger.Error.Printf("couldn't decode selfie: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	documentFrontBytes, err := base64.StdEncoding.DecodeString(documentB64.Encoding)
	if err != nil {
		logger.Error.Printf("couldn't decode document front: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	resp, err := aws_ia.CompareFaces(selfieBytes, documentFrontBytes)
	if err != nil {
		logger.Error.Printf("no se pudo comparar los rostros de la persona y el documento de identidad: %v", err)
		res.Code, res.Type, res.Msg = 22, 1, "no se pudo comparar los rostros de la persona y el documento de identidad"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	reqWs := ReqWsValidation{}

	if !resp {
		res.Code, res.Type, res.Msg = 22, 1, "La persona no es la misma que la del documento de identidad"
		reqWs = ReqWsValidation{
			TransactionId:  "-",
			UserId:         user.ID,
			DocumentNumber: strconv.FormatInt(user.DocumentNumber, 10),
			ValidatedAt:    time.Now().UTC().String(),
			ValidatorId:    "",
			RequestId:      req.RequestID,
		}
	} else {
		res.Data = "Validación de identidad realizada correctamente, la persona es la misma que la del documento de identificación"
		res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
		res.Error = false
		reqWs = ReqWsValidation{
			TransactionId:  userValidation.TransactionId,
			UserId:         user.ID,
			DocumentNumber: strconv.FormatInt(user.DocumentNumber, 10),
			ValidatedAt:    time.Now().UTC().String(),
			ValidatorId:    "",
			RequestId:      req.RequestID,
		}
	}

	if req.Nit != "" {
		client, code, err := srvCfg.SrvClients.GetClientsByNit(req.Nit)
		if err != nil {
			res.Data = ""
			logger.Error.Printf("No se pudo obtener el cliente: %v", err)
			res.Code, res.Type, res.Msg = code, 1, "No se pudo obtener el cliente"
			return c.Status(http.StatusAccepted).JSON(res)
		}

		if client == nil {
			res.Data = ""
			res.Code, res.Type, res.Msg = code, 1, "No se encontró un cliente con los datos proporcionados"
			return c.Status(http.StatusAccepted).JSON(res)
		}

		if client.UrlApi != "" {
			reqBytes, _ := json.Marshal(&reqWs)
			_, code, err := ws.ConsumeWS(reqBytes, client.UrlApi, "POST", "", nil)
			if err != nil {
				logger.Error.Printf("No se pudo enviar la petición para registra la validación de identidad: %v", err)
				res.Data = ""
				res.Code, res.Type, res.Msg = 403, 1, "No se pudo enviar la petición para registra la validación de identidad"
				return c.Status(http.StatusAccepted).JSON(res)
			}

			if code != 200 {
				logger.Error.Printf("El servicio para registrar la validación de identidad respondió con un código diferente al 200, código: %d", code)
				res.Data = ""
				res.Code, res.Type, res.Msg = code, 1, fmt.Sprintf("El servicio para registrar la validación de identidad respondió con un código diferente al 200, código: %d", code)
				return c.Status(http.StatusAccepted).JSON(res)
			}
		}
	}

	return c.Status(http.StatusOK).JSON(res)
}
