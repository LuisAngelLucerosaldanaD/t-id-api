package users

import (
	"check-id-api/internal/logger"
	"check-id-api/internal/msg"
	"check-id-api/pkg/auth"
	"check-id-api/pkg/cfg"
	"check-id-api/pkg/trx"
	"check-id-api/pkg/wf"
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

// validateIdentity godoc
// @Summary Registro de validación de identidad
// @Description Método para el registro de los datos para la validación de identidad de una persona
// @tags User
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param validateIdentity body requestValidateIdentity true "request of validate user identity"
// @Success 200 {object} responseValidateUser
// @Router /api/v1/user/validate-identity [post]
func (h *handlerUser) validateIdentity(c *fiber.Ctx) error {
	res := responseValidateUser{Error: true}
	req := requestValidateIdentity{}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)
	srvCfg := cfg.NewServerCfg(h.DB, nil, h.TxID)
	srvTrx := trx.NewServerTrx(h.DB, nil, h.TxID)
	srvWf := wf.NewServerWf(h.DB, nil, h.TxID)

	err := c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("couldn't bind model create wallets: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	resUsr, code, err := srvAuth.SrvUser.GetUsersByEmail(req.Email)
	if err != nil {
		logger.Error.Printf("couldn't create user, error validation user, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resUsr != nil {
		res.Code, res.Type, res.Msg = msg.GetByCode(5, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	user, code, err := srvAuth.SrvUser.CreateUsers(uuid.New().String(), req.TypeDocument, req.DocumentNumber, req.ExpeditionDate, req.Email, req.FirstName, req.SecondName, req.SecondSurname, req.Age, req.Gender, req.Nationality, req.CivilStatus, req.FirstSurname, req.BirthDate, req.Country, req.Department, req.City, c.IP())
	if err != nil {
		logger.Error.Printf("couldn't create user, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvTrx.SrvTraceability.CreateTraceability("Registro", "info", "Inicio de validación de identidad", user.ID)
	if err != nil {
		logger.Error.Printf("couldn't create traceability, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	f, err := srvCfg.SrvFilesS3.UploadFile(req.DocumentNumber, strconv.FormatInt(req.DocumentNumber, 10)+"_selfie.jpg", req.SelfieImg)
	if err != nil {
		logger.Error.Printf("couldn't upload file s3: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvCfg.SrvFiles.CreateFiles(f.Path, f.FileName, 1, user.ID)
	if err != nil {
		logger.Error.Printf("couldn't create selfie image, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		res.Msg = err.Error()
		return c.Status(http.StatusAccepted).JSON(res)
	}

	f, err = srvCfg.SrvFilesS3.UploadFile(req.DocumentNumber, strconv.FormatInt(req.DocumentNumber, 10)+"_doc_front.jpg", req.DocumentFrontImg)
	if err != nil {
		logger.Error.Printf("couldn't upload file document front to s3: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvCfg.SrvFiles.CreateFiles(f.Path, f.FileName, 2, user.ID)
	if err != nil {
		logger.Error.Printf("couldn't create file document front, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		res.Msg = err.Error()
		return c.Status(http.StatusAccepted).JSON(res)
	}

	f, err = srvCfg.SrvFilesS3.UploadFile(req.DocumentNumber, strconv.FormatInt(req.DocumentNumber, 10)+"_doc_back.jpg", req.DocumentBackImg)
	if err != nil {
		logger.Error.Printf("couldn't upload file document back to s3: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvCfg.SrvFiles.CreateFiles(f.Path, f.FileName, 3, user.ID)
	if err != nil {
		logger.Error.Printf("couldn't create file document back, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		res.Msg = err.Error()
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvWf.SrvWork.CreateWorkValidation("start", user.ID)
	if err != nil {
		_, _ = srvAuth.SrvUser.DeleteUsers(user.ID)
		_, _ = srvTrx.SrvTraceability.DeleteTraceabilityByUserID(user.ID)
		_, _ = srvCfg.SrvFiles.DeleteFilesByUserID(user.ID)
		logger.Error.Printf("couldn't start work, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		res.Msg = err.Error()
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, _ = srvAuth.SrvUserTemp.DeleteUserTempByEmail(user.Email)

	res.Data = "Datos registrados correctamente"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// getUserSession godoc
// @Summary Obtiene los datos registrados del usuario
// @Description Método para el obtener la información del usuario en sesión
// @tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param email path string true "Correo electrónico del usuario"
// @Success 200 {object} resGetUserSession
// @Router /api/v1/user/user-session/{email} [get]
func (h *handlerUser) getUserSession(c *fiber.Ctx) error {
	res := resGetUserSession{}
	email := c.Params("email")
	if email == "" {
		logger.Error.Printf("el email del usuario es requerido")
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)
	srvCfg := cfg.NewServerCfg(h.DB, nil, h.TxID)

	user, code, err := srvAuth.SrvUser.GetUsersByEmail(email)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el usuario, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
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
