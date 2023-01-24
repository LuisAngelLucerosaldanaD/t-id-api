package env

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

var (
	once   sync.Once
	config = &configuration{}
)

type configuration struct {
	App         App         `json:"app"`
	DB          DB          `json:"db"`
	Template    Template    `json:"template"`
	SendGrid    SendGrid    `json:"send_grid"`
	Files       Files       `json:"files"`
	Aws         Aws         `json:"aws"`
	Blockchain  Blockchain  `json:"blockchain"`
	AuthService AuthService `json:"auth_service"`
}

type App struct {
	ServiceName       string `json:"service_name"`
	Port              int    `json:"port"`
	AllowedDomains    string `json:"allowed_domains"`
	PathLog           string `json:"path_log"`
	LogReviewInterval int    `json:"log_review_interval"`
	EcdsaPrivateKey   string `json:"ecdsa_private_key"`
	RegisterLog       bool   `json:"register_log"`
	RSAPublicKey      string `json:"rsa_public_key"`
	LoggerHttp        bool   `json:"logger_http"`
	Language          string `json:"language"`
}

type Template struct {
	EmailCode        string `json:"email_code"`
	EmailToken       string `json:"email_token"`
	EmailWalletToken string `json:"email_wallet_token"`
}

type DB struct {
	Engine   string `json:"engine"`
	Server   string `json:"server"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	Instance string `json:"instance"`
	IsSecure bool   `json:"is_secure"`
	SSLMode  string `json:"ssl_mode"`
}

type SendGrid struct {
	Key      string `json:"key"`
	FromMail string `json:"from_mail"`
	FromName string `json:"from_name"`
}

type Files struct {
	Repo string `json:"repo"`
	S3   struct {
		Bucket     string `json:"bucket"`
		BucketSign string `json:"bucket_sign"`
		Region     string `json:"region"`
	} `json:"s3"`
}

type Blockchain struct {
	UrlAuth  string `json:"url_auth"`
	UrlApi   string `json:"url_api"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Wallet   string `json:"wallet"`
}

type Aws struct {
	AWSACCESSKEYID     string `json:"AWS_ACCESS_KEY_ID"`
	AWSSECRETACCESSKEY string `json:"AWS_SECRET_ACCESS_KEY"`
	AWSDEFAULTREGION   string `json:"AWS_DEFAULT_REGION"`
}

type AuthService struct {
	Port string `json:"port"`
}

func NewConfiguration() *configuration {
	fromFile()
	return config
}

// LoadConfiguration lee el archivo configuration.json
// y lo carga en un objeto de la estructura Configuration
func fromFile() {
	once.Do(func() {
		b, err := os.ReadFile("config.json")
		if err != nil {
			log.Fatalf("no se pudo leer el archivo de configuración: %s", err.Error())
		}

		err = json.Unmarshal(b, config)
		if err != nil {
			log.Fatalf("no se pudo parsear el archivo de configuración: %s", err.Error())
		}

		if config.DB.Engine == "" {
			log.Fatal("no se ha cargado la información de configuración")
		}
	})
}
