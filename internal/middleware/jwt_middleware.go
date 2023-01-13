package middleware

import (
	"check-id-api/internal/env"
	"check-id-api/internal/logger"
	"crypto/rsa"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt/v4"
	"os"
)

var (
	verifyKey *rsa.PublicKey
)

// TokenMetadata struct to describe metadata in JWT.
type TokenMetadata struct {
	Expires int64
	Ip      string
}

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
