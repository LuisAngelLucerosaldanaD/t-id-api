package users

import (
	"check-id-api/api/handlers/onboarding"
	"time"
)

type responseAnny struct {
	Error bool        `json:"error"`
	Data  interface{} `json:"data"`
	Code  int         `json:"code"`
	Type  int         `json:"type"`
	Msg   string      `json:"msg"`
}

type responseCreateUser struct {
	Error bool                   `json:"error"`
	Data  *onboarding.Onboarding `json:"data"`
	Code  int                    `json:"code"`
	Type  int                    `json:"type"`
	Msg   string                 `json:"msg"`
}

type requestLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResponseLogin struct {
	Error bool   `json:"error"`
	Data  Token  `json:"data"`
	Code  int    `json:"code"`
	Type  int    `json:"type"`
	Msg   string `json:"msg"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type reqUploadSelfie struct {
	UserID    string `json:"user_id"`
	SelfieImg string `json:"selfie_img"`
}

type reqUploadDocument struct {
	UserID           string `json:"user_id"`
	DocumentFrontImg string `json:"document_front_img"`
	DocumentBackImg  string `json:"document_back_img"`
}

type RequestCreateUser struct {
	Id             string `json:"id"`
	DocumentNumber string `json:"document_number"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Cellphone      string `json:"cellphone"`
}

type resGetUser struct {
	Error bool   `json:"error"`
	Data  *User  `json:"data"`
	Code  int    `json:"code"`
	Type  int    `json:"type"`
	Msg   string `json:"msg"`
}

type User struct {
	ID                 string     `json:"id"`
	Nickname           string     `json:"nickname"`
	Email              string     `json:"email"`
	FirstName          *string    `json:"first_name"`
	SecondName         *string    `json:"second_name"`
	FirstSurname       *string    `json:"first_surname"`
	SecondSurname      *string    `json:"second_surname"`
	Age                *int32     `json:"age"`
	TypeDocument       *string    `json:"type_document"`
	DocumentNumber     string     `json:"document_number"`
	Cellphone          string     `json:"cellphone"`
	Gender             *string    `json:"gender"`
	Nationality        *string    `json:"nationality"`
	Country            *string    `json:"country"`
	Department         *string    `json:"department"`
	City               *string    `json:"city"`
	RealIp             string     `json:"real_ip"`
	StatusId           int32      `json:"status_id"`
	FailedAttempts     int32      `json:"failed_attempts"`
	BlockDate          *time.Time `json:"block_date"`
	DisabledDate       *time.Time `json:"disabled_date"`
	LastLogin          *time.Time `json:"last_login"`
	LastChangePassword *time.Time `json:"last_change_password"`
	BirthDate          *time.Time `json:"birth_date"`
	VerifiedCode       *string    `json:"verified_code"`
	IsDeleted          bool       `json:"is_deleted"`
	DeletedAt          *time.Time `json:"deleted_at"`
	SelfieImg          string     `json:"selfie_img"`
	FrontDocumentImg   string     `json:"front_document_img"`
	BackDocumentImg    string     `json:"back_document_img"`
	TransactionId      string     `json:"transaction_id"`
	ProcessURL         string     `json:"process_url"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type resGetUsersLasted struct {
	Error bool          `json:"error"`
	Data  []*UserStatus `json:"data"`
	Code  int           `json:"code"`
	Type  int           `json:"type"`
	Msg   string        `json:"msg"`
}

type UserStatus struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	FirstName     string    `json:"first_name"`
	SecondName    string    `json:"second_name"`
	FirstSurname  string    `json:"first_surname"`
	SecondSurname string    `json:"second_surname"`
	Status        string    `json:"status"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type resGetUsersDataPending struct {
	Error bool        `json:"error"`
	Data  DataPending `json:"data"`
	Code  int         `json:"code"`
	Type  int         `json:"type"`
	Msg   string      `json:"msg"`
}

type DataPending struct {
	Selfie           int `json:"selfie"`
	Document         int `json:"document"`
	BasicInformation int `json:"basic_information"`
}

type ReqValidationFace struct {
	FaceImage      string `json:"face_image"`
	DocumentNumber int64  `json:"document_number"`
	Nit            string `json:"nit"`
	RequestID      string `json:"request-id"`
}

type ReqWsValidation struct {
	TransactionId  string `json:"transaction_id"`
	UserId         string `json:"user_id"`
	DocumentNumber string `json:"document_number"`
	ValidatedAt    string `json:"validated_at"`
	ValidatorId    string `json:"validator_id"`
	RequestId      string `json:"request_id"`
}

type responseFinishOnboarding struct {
	Error bool   `json:"error"`
	Data  bool   `json:"data"`
	Code  int    `json:"code"`
	Type  int    `json:"type"`
	Msg   string `json:"msg"`
}

type ResponseGetUserFile struct {
	Error bool   `json:"error"`
	Data  string `json:"data"`
	Code  int    `json:"code"`
	Type  int    `json:"type"`
	Msg   string `json:"msg"`
}
