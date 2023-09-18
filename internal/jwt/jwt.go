package jwt

import (
	"check-id-api/internal/env"
	"check-id-api/internal/logger"
	"check-id-api/internal/models"
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v4"
	"io/ioutil"
	"time"
)

var (
	signKey    *rsa.PrivateKey
	privateKey string
)

func init() {
	c := env.NewConfiguration()
	privateKey = c.App.RSAPrivateKey
	signBytes, err := ioutil.ReadFile(privateKey)
	if err != nil {
		logger.Error.Printf("leyendo el archivo privado de firma: %s", err)
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		logger.Error.Printf("realizando el parse en auth RSA private: %s", err)
	}
}

// GenerateJWT Genera el token
func GenerateJWT(u *models.User, role string) (string, int, error) {
	tk := jwt.New(jwt.SigningMethodRS256)
	claims := tk.Claims.(jwt.MapClaims)
	claims["user"] = u
	claims["rol"] = role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	token, err := tk.SignedString(signKey)
	if err != nil {
		logger.Error.Printf("Error en el firmado del token: %v", err)
		return "", 33, err
	}
	return token, 29, nil
}
