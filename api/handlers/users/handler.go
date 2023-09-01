package users

import (
	"check-id-api/internal/logger"
	"check-id-api/internal/msg"
	"check-id-api/pkg/auth"
	"check-id-api/pkg/auth/users"
	"check-id-api/pkg/cfg"
	"check-id-api/pkg/trx"
	"check-id-api/pkg/wf"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
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

	user, code, err := srvAuth.SrvUser.UpdateUsers(req.Id, &req.TypeDocument, req.DocumentNumber, req.ExpeditionDate, req.Email, &req.FirstName, &req.SecondName, &req.SecondSurname, &req.Age, &req.Gender, &req.Nationality, &req.CivilStatus, &req.FirstSurname, req.BirthDate, &req.Country, &req.Department, &req.City, c.IP(), "")
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

	user, code, err := srvAuth.SrvUser.CreateUsers(uuid.New().String(), &req.TypeDocument, req.DocumentNumber, req.ExpeditionDate, req.Email, &req.FirstName, &req.SecondName, &req.SecondSurname, &req.Age, &req.Gender, &req.Nationality, &req.CivilStatus, &req.FirstSurname, req.BirthDate, &req.Country, &req.Department, &req.City, c.IP(), "")
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
			/*fileS3, code, err := srvCfg.SrvFilesS3.GetFileByPath(file.Path, file.Name)
			if err != nil {
				logger.Error.Printf("No se pudo descargar el archivo, error: %s", err.Error())
				res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
				return c.Status(http.StatusAccepted).JSON(res)
			}*/

			switch file.Type {
			case 1:
				// selfieImg = fileS3.Encoding
				selfieImg = strconv.FormatInt(file.ID, 10)
				break
			case 2:
				frontDocument = strconv.FormatInt(file.ID, 10)
				break
			default:
				backDocument = strconv.FormatInt(file.ID, 10)
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
		FirstName:        *user.FirstName,
		SecondName:       *user.SecondName,
		SecondSurname:    *user.SecondSurname,
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
			FirstName:     *user.FirstName,
			SecondName:    *user.SecondName,
			SecondSurname: *user.SecondSurname,
			FirstSurname:  *user.FirstSurname,
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

// getFinishOnboarding godoc
// @Summary Permite validar si ha terminado el enrolamiento de un usuario
// @Description Método que permite validar si se ha finalizado el proceso de enrolamiento de un usuario
// @tags User
// @Accept json
// @Produce json
// @Param id path string true "Id del usuario"
// @Success 200 {object} responseFinishOnboarding
// @Router /api/v1/user/finish-onboarding [get]
func (h *handlerUser) getFinishOnboarding(c *fiber.Ctx) error {
	res := responseFinishOnboarding{Error: true}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	userId := c.Params("id")
	if userId == "" {
		logger.Error.Printf("El id del usuario es requerido")
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	onboarding, code, err := srvAuth.SrvOnboarding.GetOnboardingByUserID(userId)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el registro de enrolamiento del usuario: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if onboarding == nil {
		logger.Error.Printf("No se pudo obtener el registro de enrolamiento del usuario")
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = onboarding.Status == "finished"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// getFinishValidationIdentity godoc
// @Summary Permite validar si ha terminado la validación de identidad de un usuario
// @Description Método que permite validar si ha terminado la validación de identidad de un usuario
// @tags User
// @Accept json
// @Produce json
// @Param id path string true "Id del usuario"
// @Success 200 {object} responseFinishOnboarding
// @Router /api/v1/user/finish-validation [get]
func (h *handlerUser) getFinishValidationIdentity(c *fiber.Ctx) error {
	res := responseFinishOnboarding{Error: true}
	srvCfg := cfg.NewServerCfg(h.DB, nil, h.TxID)

	userId := c.Params("id")
	if userId == "" {
		logger.Error.Printf("El id del usuario es requerido")
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	validation, code, err := srvCfg.SrvValidationRequest.GetValidationRequestByUserID(userId)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el registro de validación de identidad del usuario: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if validation == nil {
		logger.Error.Printf("No se pudo obtener el registro de validación de identidad del usuario")
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = validation.Status == "finished"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// getFinishValidationIdentity godoc
// @Summary Permite validar si ha terminado la validación de identidad de un usuario
// @Description Método que permite validar si ha terminado la validación de identidad de un usuario
// @tags User
// @Accept json
// @Produce json
// @Param id path string true "Id del archivo"
// @Success 200 {object} ResponseGetUserFile
// @Router /api/v1/user/file [get]
func (h *handlerUser) getUserFile(c *fiber.Ctx) error {
	res := ResponseGetUserFile{Error: true}
	srvCfg := cfg.NewServerCfg(h.DB, nil, h.TxID)

	fileID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		logger.Error.Printf("El id del archivo es invalido")
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	file, code, err := srvCfg.SrvFiles.GetFilesByID(fileID)
	if err != nil {
		logger.Error.Printf("No se pudo obtener la validación del usuario, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if file == nil {
		logger.Error.Printf("No se pudo obtener el archivo del usuario solicitado")
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	fileS3, code, err := srvCfg.SrvFilesS3.GetFileByPath(file.Path, file.Name)
	if err != nil {
		logger.Error.Printf("No se pudo descargar el archivo, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if fileS3 == nil {
		logger.Error.Printf("No se pudo descargar el archivo solicitado")
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = fileS3.Encoding
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}
