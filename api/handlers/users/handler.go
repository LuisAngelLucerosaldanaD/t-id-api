package users

import (
	"check-id-api/api/handlers/onboarding"
	"check-id-api/internal/env"
	"check-id-api/internal/jwt"
	"check-id-api/internal/logger"
	"check-id-api/internal/middleware"
	"check-id-api/internal/models"
	"check-id-api/internal/msg"
	"check-id-api/internal/password"
	"check-id-api/pkg/auth"
	"check-id-api/pkg/auth/user"
	"check-id-api/pkg/cfg"
	"check-id-api/pkg/trx"
	"fmt"
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

// createUser godoc
// @Summary Método que permite la creación de un usuario con los datos básicos
// @Description Método que permite la creación de un usuario con los datos básicos y permite iniciar el proceso de validación de identidad usando checkid como cliente para la solicitud
// @tags User
// @Accept json
// @Produce json
// @Param RequestCreateUser body RequestCreateUser true "Datos para la creación del usuario"
// @Success 200 {object} responseCreateUser
// @Router /api/v1/user/create [post]
func (h *handlerUser) createUser(c *fiber.Ctx) error {
	res := responseCreateUser{Error: true}
	req := RequestCreateUser{}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)
	srvTrx := trx.NewServerTrx(h.DB, nil, h.TxID)
	e := env.NewConfiguration()

	err := c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("couldn't bind model create user: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	userFound, code, err := srvAuth.SrvUser.GetUserByEmail(req.Email)
	if err != nil {
		logger.Error.Printf("couldn't bind get user by email: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if userFound != nil {
		res.Code, res.Type, res.Msg = 102, 1, "Ya existe un usuario con el correo proporcionado"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	newUser, code, err := srvAuth.SrvUser.CreateUser(&user.User{
		ID:             uuid.New().String(),
		Nickname:       req.Email,
		Email:          req.Email,
		Password:       password.Encrypt(req.Password),
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

	_, code, err = srvAuth.SrvUserRole.CreateUseRole(uuid.New().String(), newUser.ID, "14cbf8d2-485a-4fbe-baa6-16273c765f14")
	if err != nil {
		logger.Error.Printf("couldn't create user role, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvTrx.SrvTraceability.CreateTraceability("Registro", "info", "Registro de información básica", newUser.ID)
	if err != nil {
		logger.Error.Printf("couldn't create traceability, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	newRequestOnboarding, code, err := srvAuth.SrvOnboardingCheckId.CreateOnboardingCheckId(newUser.ID, c.IP())
	if err != nil {
		logger.Error.Printf("couldn't create onboarding check id, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	newOnboarding, code, err := srvAuth.SrvOnboarding.CreateOnboarding(uuid.New().String(), 5, strconv.FormatInt(newRequestOnboarding.ID, 10), newUser.ID, "started", "")
	if err != nil {
		logger.Error.Printf("couldn't create onboarding, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = &onboarding.Onboarding{
		Url:    e.OnlyOne.Url + e.OnlyOne.Onboarding + fmt.Sprintf("?process=enroll&onboarding_id=%s&user_id=%s&email=%s", newOnboarding.ID, newUser.ID, req.Email),
		Method: "onboarding",
	}
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// Login godoc
// @Summary Método que permite autenticar al usuario en el sistema
// @Description Método que permite autenticar al usuario en el sistema
// @tags User
// @Accept json
// @Produce json
// @Param requestLogin body requestLogin true "Datos para la autenticación"
// @Success 200 {object} ResponseLogin
// @Router /api/v1/user/login [post]
func (h *handlerUser) Login(c *fiber.Ctx) error {
	res := ResponseLogin{Error: true}
	req := requestLogin{}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	err := c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("couldn't bind model create wallets: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	userFound, code, err := srvAuth.SrvUser.GetUserByEmail(req.Email)
	if err != nil {
		logger.Error.Printf("couldn't get userFound by email: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if userFound == nil {
		res.Code, res.Type, res.Msg = msg.GetByCode(10, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if !password.Compare(userFound.ID, userFound.Password, req.Password) {
		res.Code, res.Type, res.Msg = msg.GetByCode(10, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	userRole, code, err := srvAuth.SrvUserRole.GetUseRoleByUserID(userFound.ID)
	if err != nil {
		logger.Error.Printf("couldn't get userFound role: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if userRole == nil {
		res.Code, res.Type, res.Msg = msg.GetByCode(95, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	role, code, err := srvAuth.SrvRole.GetRoleByID(userRole.RoleId)
	if err != nil {
		logger.Error.Printf("couldn't get role: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if role == nil {
		res.Code, res.Type, res.Msg = msg.GetByCode(72, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	userFound.RealIp = c.IP()
	userFound.Password = ""

	token, code, err := jwt.GenerateJWT((*models.User)(userFound), role.Name)
	if err != nil {
		logger.Error.Printf("couldn't generate userFound: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(10, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = Token{
		AccessToken:  token,
		RefreshToken: token,
	}
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
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

// getUserById godoc
// @Summary Obtiene los datos registrados del usuario
// @Description Método para obtener los datos registrados del usuario
// @tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Success 200 {object} resGetUser
// @Router /api/v1/user [get]
func (h *handlerUser) getUserById(c *fiber.Ctx) error {
	res := resGetUser{}
	userToken, err := middleware.GetUser(c)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el usuario del token, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(95, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)
	srvCfg := cfg.NewServerCfg(h.DB, nil, h.TxID)
	e := env.NewConfiguration()

	currentUser, code, err := srvAuth.SrvUser.GetUserByIdentityNumber(userToken.DocumentNumber)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el usuario por su numero de identificacion, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if currentUser == nil {
		currentUser, code, err = srvAuth.SrvUser.GetUserByEmail(userToken.Email)
		if err != nil {
			logger.Error.Printf("No se pudo obtener el usuario por su email, error: %s", err.Error())
			res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}
	}

	if currentUser == nil {
		res.Error = false
		res.Code, res.Type, res.Msg = msg.GetByCode(95, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	role, code, err := srvAuth.SrvRole.GetRoleByUserID(currentUser.ID)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el rol del usuario, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if role == nil {
		logger.Error.Printf("El usuario no tiene un rol asignado")
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	files, code, err := srvCfg.SrvFiles.GetFilesByUserID(currentUser.ID)
	if err != nil {
		logger.Error.Printf("No se pudo obtener la validación del usuario, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	userResponse := &User{
		ID:                 currentUser.ID,
		Nickname:           currentUser.Nickname,
		Email:              currentUser.Email,
		FirstName:          currentUser.FirstName,
		SecondName:         currentUser.SecondName,
		FirstSurname:       currentUser.FirstSurname,
		SecondSurname:      currentUser.SecondSurname,
		Age:                currentUser.Age,
		TypeDocument:       currentUser.TypeDocument,
		DocumentNumber:     currentUser.DocumentNumber,
		Cellphone:          currentUser.Cellphone,
		Gender:             currentUser.Gender,
		Nationality:        currentUser.Nationality,
		Country:            currentUser.Country,
		Department:         currentUser.Department,
		City:               currentUser.City,
		RealIp:             currentUser.RealIp,
		StatusId:           currentUser.StatusId,
		FailedAttempts:     currentUser.FailedAttempts,
		BlockDate:          currentUser.BlockDate,
		DisabledDate:       currentUser.DisabledDate,
		LastLogin:          currentUser.LastLogin,
		LastChangePassword: currentUser.LastChangePassword,
		BirthDate:          currentUser.BirthDate,
		VerifiedCode:       currentUser.VerifiedCode,
		IsDeleted:          currentUser.IsDeleted,
		DeletedAt:          currentUser.DeletedAt,
		CreatedAt:          currentUser.CreatedAt,
		UpdatedAt:          currentUser.UpdatedAt,
	}

	if files != nil {
		for _, file := range files {
			fileID := strconv.FormatInt(file.ID, 10)
			switch file.Type {
			case 1:
				userResponse.SelfieImg = fileID
				break
			case 2:
				userResponse.FrontDocumentImg = fileID
				break
			default:
				userResponse.BackDocumentImg = fileID
				break
			}
		}
	}

	lifeTestUser, code, err := srvAuth.SrvLifeTest.GetLifeTestByUserID(currentUser.ID)
	if err != nil {
		logger.Error.Printf("No se pudo obtener la prueba de vida del usuario, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if lifeTestUser != nil {
		userResponse.ProcessURL = e.OnlyOne.Url + e.OnlyOne.Onboarding + fmt.Sprintf("?process=validation&validation_id=%d&user_id=%s&email=%s", lifeTestUser.ID, currentUser.ID, currentUser.Email)
		res.Data = userResponse
		res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
		res.Error = false
		return c.Status(http.StatusOK).JSON(res)
	}

	onboardingUser, code, err := srvAuth.SrvOnboarding.GetOnboardingByUserID(currentUser.ID)
	if err != nil {
		logger.Error.Printf("No se pudo obtener la validación del usuario, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if onboardingUser == nil {
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	userResponse.TransactionId = onboardingUser.TransactionId

	if onboardingUser.Status == "started" {
		userResponse.ProcessURL = e.OnlyOne.Url + e.OnlyOne.Onboarding + fmt.Sprintf("?process=enroll&onboarding_id=%s&user_id=%s&email=%s", onboardingUser.ID, currentUser.ID, currentUser.Email)
	}

	res.Data = userResponse
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
		res.Code, res.Type, res.Msg = msg.GetByCode(95, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	userFound, code, err := srvAuth.SrvUser.GetUserByIdentityNumber(userToken.DocumentNumber)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el usuario por su identificacion, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if userFound == nil {
		res.Code, res.Type, res.Msg = 5, 1, "No existe un usuario registrado con la información proporcionada"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if userFound == nil {
		logger.Error.Printf("Usuario no encontrado, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	onboardingUser, code, err := srvAuth.SrvOnboarding.GetOnboardingByUserID(userFound.ID)
	if err != nil {
		logger.Error.Printf("No se pudo consultar la validacion de identidad, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if onboardingUser == nil {
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = onboardingUser.TransactionId != ""
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
		res.Code, res.Type, res.Msg = msg.GetByCode(95, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	userFound, code, err := srvAuth.SrvUser.GetUserByIdentityNumber(userToken.DocumentNumber)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el usuario por su identificacion, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if userFound == nil {
		logger.Error.Printf("Usuario no encontrado, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	onboardingUser, code, err := srvAuth.SrvOnboarding.GetOnboardingByUserID(userFound.ID)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el registro de enrolamiento del usuario: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if onboardingUser == nil {
		logger.Error.Printf("No se pudo obtener el registro de enrolamiento del usuario")
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = onboardingUser.Status == "finished"
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
		res.Code, res.Type, res.Msg = msg.GetByCode(95, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	userFound, code, err := srvAuth.SrvUser.GetUserByIdentityNumber(userToken.DocumentNumber)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el usuario por su identificacion, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if userFound == nil {
		logger.Error.Printf("Usuario no encontrado, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	lifeTest, code, err := srvAuth.SrvLifeTest.GetLifeTestByUserID(userFound.ID)
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
// @Router /api/v1/user/file/{id} [get]
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
		res.Code, res.Type, res.Msg = msg.GetByCode(95, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	userFound, code, err := srvAuth.SrvUser.GetUserByIdentityNumber(userToken.DocumentNumber)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el usuario por su identificacion, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if userFound == nil {
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

	if file.UserId != userFound.ID {
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
