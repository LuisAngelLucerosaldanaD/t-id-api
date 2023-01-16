package main

import (
	"check-id-api/api"
	"check-id-api/internal/env"
	"os"
)

// @title Check ID OnBoarding
// @version 1.4
// @description Api para OnBoarding y validación de identidad
// @termsOfService https://www.bjungle.net/terms/
// @contact.name API Support
// @contact.email luis.lucero@bjungle.net
// @license.name Software Owner
// @license.url https://www.bjungle.net/terms/licenses
// @host http://172.147.77.149:50050
// @tag.name User
// @tag.description Métodos referentes al usuario
// @tag.name Traceability
// @tag.description Métodos referentes a la trazabilidad
// @tag.name Work
// @tag.description Métodos referentes al trabajo registrado
// @BasePath /
func main() {
	c := env.NewConfiguration()
	_ = os.Setenv("AWS_ACCESS_KEY_ID", c.Aws.AWSACCESSKEYID)
	_ = os.Setenv("AWS_SECRET_ACCESS_KEY", c.Aws.AWSSECRETACCESSKEY)
	_ = os.Setenv("AWS_DEFAULT_REGION", c.Aws.AWSDEFAULTREGION)

	api.Start(c.App.Port, c.App.ServiceName, c.App.LoggerHttp, c.App.AllowedDomains)
}
