package main

import (
	"check-id-api/api"
	"check-id-api/internal/env"
	"os"
)

// @title Check ID
// @version 1.2
// @description Api para OnBoarding y validación de identidad de una persona
// @termsOfService https://www.bjungle.net/terms/
// @contact.name API Support
// @contact.email info@bjungle.net
// @license.name Software Owner
// @license.url https://www.bjungle.net/terms/licenses
// @host http://127.0.0.1:50050
// @tag.name User
// @tag.description Métodos referentes al usuario
// @tag.name Traceability
// @tag.description Métodos referentes a la trazabilidad
// @BasePath /
func main() {
	c := env.NewConfiguration()
	_ = os.Setenv("AWS_ACCESS_KEY_ID", c.Aws.AWSACCESSKEYID)
	_ = os.Setenv("AWS_SECRET_ACCESS_KEY", c.Aws.AWSSECRETACCESSKEY)
	_ = os.Setenv("AWS_DEFAULT_REGION", c.Aws.AWSDEFAULTREGION)

	api.Start(c.App.Port, c.App.ServiceName, c.App.LoggerHttp, c.App.AllowedDomains)
}
