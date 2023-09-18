package middleware

import (
	"check-id-api/internal/env"
	"check-id-api/internal/logger"
	"check-id-api/internal/models"
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"os"
)

var (
	verifyKey *rsa.PublicKey
)

// init lee los archivos de firma y validación RSA
func init() {
	c := env.NewConfiguration()

	verifyBytes, err := os.ReadFile(c.App.RSAPublicKey)
	if err != nil {
		logger.Error.Printf("leyendo el archivo público de confirmación: %s", err)
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		logger.Error.Printf("realizando el parse en jwt RSA public: %s", err)
	}
}

type jwtCustomClaims struct {
	User      *models.User `json:"user"`
	IPAddress string       `json:"ip_address"`
	jwt.RegisteredClaims
}

func JWTProtected() fiber.Handler {
	config := jwtware.Config{
		ErrorHandler:  jwtError,
		SigningKey:    verifyKey,
		SigningMethod: "RS256",
	}
	return jwtware.New(config)
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}

func GetUser(c *fiber.Ctx) (*models.User, error) {
	bearer := c.Get("Authorization")
	tkn := bearer[7:]

	var u *models.User
	verifyFunction := func(tkn *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	}

	token, err := jwt.ParseWithClaims(tkn, &jwtCustomClaims{}, verifyFunction)
	if err != nil {
		var validationError *jwt.ValidationError
		switch {
		case errors.As(err, &validationError):
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				logger.Warning.Printf("token expirado: %v", err)
				return u, err
			default:
				logger.Warning.Printf("Error de validacion del token: %v", err)
				return u, err
			}
		default:
			logger.Warning.Printf("Error al procesar el token: %v", err)
			return u, err
		}
	}
	u = token.Claims.(*jwtCustomClaims).User
	if !token.Valid {
		logger.Warning.Printf("Token no Valido: %v", err)
		return u, fmt.Errorf("token no Valido")
	}
	if c.IP() != u.RealIp {
		logger.Warning.Printf("token creado en un origen diferente : %v", err)
		return u, fmt.Errorf("token creado en un origen diferente")
	}
	return u, nil
}
