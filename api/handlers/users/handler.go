package users

import (
	"check-id-api/internal/logger"
	"check-id-api/internal/middleware"
	"check-id-api/internal/msg"
	"check-id-api/pkg/auth"
	user2 "check-id-api/pkg/auth/user"
	"check-id-api/pkg/cfg"
	"check-id-api/pkg/trx"
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

	err := c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("couldn't bind model create wallets: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	fileSelfie, code, err := srvCfg.SrvFiles.GetFileByTypeAndUserID(1, req.UserID)
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

	err := c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("couldn't bind model create wallets: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	fileTemp, code, err := srvCfg.SrvFiles.GetFileByTypeAndUserID(2, req.UserID)
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

	_, code, err = srvCfg.SrvFiles.CreateFile(f.Path, f.FileName, 2, req.UserID)
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

	res.Data = "Documento cargado correctamente"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// createUser godoc
// @Summary Metodo que permite la creación de un usuario
// @Description Metodo que permite la creación de un usuario
// @tags User
// @Accept json
// @Produce json
// @Param BasicInformation body RequestCreateUser true "request of validate user identity"
// @Success 200 {object} responseAnny
// @Router /api/v1/user/create [post]
func (h *handlerUser) createUser(c *fiber.Ctx) error {
	res := responseAnny{Error: true}
	req := RequestCreateUser{}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)
	srvTrx := trx.NewServerTrx(h.DB, nil, h.TxID)

	err := c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("couldn't bind model create user: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	user, code, err := srvAuth.SrvUser.CreateUser(&user2.User{
		ID:             uuid.New().String(),
		Nickname:       req.Email,
		Email:          req.Email,
		Password:       req.Password,
		DocumentNumber: req.DocumentNumber,
		Cellphone:      req.Cellphone,
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

	_, code, err = srvTrx.SrvTraceability.CreateTraceability("Registro", "info",
		"Creación de la cuenta en la plataforma", user.ID)
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

// getUserSession godoc
// @Summary Obtiene los datos registrados del usuario por su email o su id
// @Description Método para el obtener la información del usuario en sesión por su email o id
// @tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Success 200 {object} resGetUserSession
// @Router /api/v1/user/user-session [get]
func (h *handlerUser) getUserSession(c *fiber.Ctx) error {
	res := resGetUserSession{}
	userToken, err := middleware.GetUser(c)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el usuario del token, error: %s", err.Error())
		// TODO pendiente de agregar mensaje
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)
	srvCfg := cfg.NewServerCfg(h.DB, nil, h.TxID)

	userTmp, code, err := srvAuth.SrvUser.GetUserByIdentityNumber(userToken.IdNumber)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el usuario por su numero de identificacion, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if userTmp == nil {
		userTmp, code, err = srvAuth.SrvUser.GetUserByEmail(userToken.Email)
		if err != nil {
			logger.Error.Printf("No se pudo obtener el usuario por su email, error: %s", err.Error())
			res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}
	}

	if userTmp == nil {
		res.Error = false
		res.Code, res.Type, res.Msg = msg.GetByCode(95, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	user := userTmp

	onboarding, code, err := srvAuth.SrvOnboarding.GetOnboardingByUserID(user.ID)
	if err != nil {
		logger.Error.Printf("No se pudo obtener la validación del usuario, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	role, code, err := srvAuth.SrvRole.GetRoleByUserID(user.ID)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el rol del usuario, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if role == nil {
		logger.Error.Printf("El usuario no tiene un rol asigando")
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	transactionID := ""
	selfieImg := ""
	frontDocument := ""
	backDocument := ""

	if onboarding != nil {
		transactionID = onboarding.TransactionId
	}

	files, code, err := srvCfg.SrvFiles.GetFilesByUserID(user.ID)
	if err != nil {
		logger.Error.Printf("No se pudo obtener la validación del usuario, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if files != nil {
		for _, file := range files {
			fileID := strconv.FormatInt(file.ID, 10)
			switch file.Type {
			case 1:
				selfieImg = fileID
				break
			case 2:
				frontDocument = fileID
				break
			default:
				backDocument = fileID
				break
			}
		}
	}

	res.Data = &User{
		ID:               user.ID,
		TypeDocument:     user.TypeDocument,
		DocumentNumber:   user.DocumentNumber,
		Email:            user.Email,
		FirstName:        *user.FirstName,
		SecondName:       *user.SecondName,
		SecondSurname:    *user.SecondSurname,
		Age:              user.Age,
		Gender:           user.Gender,
		Nationality:      user.Nationality,
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
		Role:             role.Name,
		UpdatedAt:        user.UpdatedAt,
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
// @Success 200 {object} responseAnny
// @Router /api/v1/user/validate [get]
func (h *handlerUser) validateUser(c *fiber.Ctx) error {
	res := responseAnny{Error: true}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	userToken, err := middleware.GetUser(c)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el usuario del token, error: %s", err.Error())
		// TODO pendiente de agregar mensaje
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	user, code, err := srvAuth.SrvUser.GetUserByIdentityNumber(userToken.IdNumber)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el usuario por su identificacion, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if user == nil {
		res.Code, res.Type, res.Msg = 5, 1, "No existe un usuario registrado con la información proporcionada"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if user == nil {
		logger.Error.Printf("Usuario no encontrado, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	onboarding, code, err := srvAuth.SrvOnboarding.GetOnboardingByUserID(user.ID)
	if err != nil {
		logger.Error.Printf("No se pudo consultar la validacion de identidad, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if onboarding == nil {
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = onboarding.TransactionId != ""
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
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Success 200 {object} responseFinishOnboarding
// @Router /api/v1/user/finish-onboarding [get]
func (h *handlerUser) getFinishOnboarding(c *fiber.Ctx) error {
	res := responseFinishOnboarding{Error: true}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	userToken, err := middleware.GetUser(c)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el usuario del token, error: %s", err.Error())
		// TODO pendiente de agregar mensaje
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	user, code, err := srvAuth.SrvUser.GetUserByIdentityNumber(userToken.IdNumber)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el usuario por su identificacion, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if user == nil {
		logger.Error.Printf("Usuario no encontrado, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	onboarding, code, err := srvAuth.SrvOnboarding.GetOnboardingByUserID(user.ID)
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
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Success 200 {object} responseFinishOnboarding
// @Router /api/v1/user/finish-validation [get]
func (h *handlerUser) getFinishValidationIdentity(c *fiber.Ctx) error {
	res := responseFinishOnboarding{Error: true}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	userToken, err := middleware.GetUser(c)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el usuario del token, error: %s", err.Error())
		// TODO pendiente de agregar mensaje
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	user, code, err := srvAuth.SrvUser.GetUserByIdentityNumber(userToken.IdNumber)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el usuario por su identificacion, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if user == nil {
		logger.Error.Printf("Usuario no encontrado, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	lifeTest, code, err := srvAuth.SrvLifeTest.GetLifeTestByUserID(user.ID)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el registro de validación de identidad del usuario: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if lifeTest == nil {
		logger.Error.Printf("No se pudo obtener el registro de validación de identidad del usuario")
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = lifeTest.Status == "finished"
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
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param id path string true "Id del archivo"
// @Success 200 {object} ResponseGetUserFile
// @Router /api/v1/user/file [get]
func (h *handlerUser) getUserFile(c *fiber.Ctx) error {
	res := ResponseGetUserFile{Error: true}
	srvCfg := cfg.NewServerCfg(h.DB, nil, h.TxID)
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	fileID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		logger.Error.Printf("El id del archivo es invalido")
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	userToken, err := middleware.GetUser(c)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el usuario del token, error: %s", err.Error())
		// TODO pendiente de agregar mensaje
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	user, code, err := srvAuth.SrvUser.GetUserByIdentityNumber(userToken.IdNumber)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el usuario por su identificacion, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if user == nil {
		logger.Error.Printf("Usuario no encontrado, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	file, code, err := srvCfg.SrvFiles.GetFileByID(fileID)
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

	if file.UserId != user.ID {
		logger.Error.Printf("No esta autorizado para ver este recurso")
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
