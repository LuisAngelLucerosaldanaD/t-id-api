package users

import (
	"check-id-api/pkg/auth/users"
	"time"
)

type responseAnny struct {
	Error bool        `json:"error"`
	Data  interface{} `json:"data"`
	Code  int         `json:"code"`
	Type  int         `json:"type"`
	Msg   string      `json:"msg"`
}

type resCreateUser struct {
	Error bool         `json:"error"`
	Data  *users.Users `json:"data"`
	Code  int          `json:"code"`
	Type  int          `json:"type"`
	Msg   string       `json:"msg"`
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

type requestValidateIdentity struct {
	Id             string     `json:"id"`
	TypeDocument   string     `json:"type_document"`
	DocumentNumber int64      `json:"document_number"`
	ExpeditionDate *time.Time `json:"expedition_date"`
	Email          string     `json:"email"`
	FirstName      string     `json:"first_name"`
	SecondName     string     `json:"second_name"`
	SecondSurname  string     `json:"second_surname"`
	Age            int32      `json:"age"`
	Gender         string     `json:"gender"`
	Nationality    string     `json:"nationality"`
	CivilStatus    string     `json:"civil_status"`
	FirstSurname   string     `json:"first_surname"`
	BirthDate      *time.Time `json:"birth_date"`
	Country        string     `json:"country"`
	Department     string     `json:"department"`
	City           string     `json:"city"`
}

type resGetUserSession struct {
	Error bool            `json:"error"`
	Data  *UserValidation `json:"data"`
	Code  int             `json:"code"`
	Type  int             `json:"type"`
	Msg   string          `json:"msg"`
}

type UserValidation struct {
	ID               string     `json:"id"`
	TypeDocument     string     `json:"type_document"`
	DocumentNumber   int64      `json:"document_number"`
	ExpeditionDate   *time.Time `json:"expedition_date"`
	Email            string     `json:"email"`
	FirstName        string     `json:"first_name"`
	SecondName       string     `json:"second_name"`
	SecondSurname    string     `json:"second_surname"`
	Age              int32      `json:"age"`
	SelfieImg        string     `json:"selfie_img"`
	BackDocumentImg  string     `json:"back_document_img"`
	FrontDocumentImg string     `json:"front_document_img"`
	Gender           string     `json:"gender"`
	Nationality      string     `json:"nationality"`
	CivilStatus      string     `json:"civil_status"`
	FirstSurname     string     `json:"first_surname"`
	BirthDate        *time.Time `json:"birth_date"`
	Country          string     `json:"country"`
	TransactionId    string     `json:"transaction_id"`
	Department       string     `json:"department"`
	City             string     `json:"city"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
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
