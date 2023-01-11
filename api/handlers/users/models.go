package users

import (
	"time"
)

type responseValidateUser struct {
	Error bool   `json:"error"`
	Data  string `json:"data"`
	Code  int    `json:"code"`
	Type  int    `json:"type"`
	Msg   string `json:"msg"`
}

type WalletIdentity struct {
	ID         string `json:"id" db:"id"`
	Mnemonic   string `json:"mnemonic" db:"mnemonic"`
	RsaPublic  string `json:"rsa_public" db:"rsa_public"`
	RsaPrivate string `json:"rsa_private" db:"rsa_private"`
}

type requestValidateIdentity struct {
	DocumentFrontImg string    `json:"document_front_img"`
	DocumentBackImg  string    `json:"document_back_img"`
	SelfieImg        string    `json:"selfie_img"`
	TypeDocument     string    `json:"type_document"`
	DocumentNumber   int64     `json:"document_number"`
	ExpeditionDate   time.Time `json:"expedition_date"`
	Email            string    `json:"email"`
	FirstName        string    `json:"first_name"`
	SecondName       string    `json:"second_name"`
	SecondSurname    string    `json:"second_surname"`
	Age              int32     `json:"age"`
	Gender           string    `json:"gender"`
	Nationality      string    `json:"nationality"`
	CivilStatus      string    `json:"civil_status"`
	FirstSurname     string    `json:"first_surname"`
	BirthDate        time.Time `json:"birth_date"`
	Country          string    `json:"country"`
	Department       string    `json:"department"`
	City             string    `json:"city"`
}

type resGetUserSession struct {
	Error bool            `json:"error"`
	Data  *UserValidation `json:"data"`
	Code  int             `json:"code"`
	Type  int             `json:"type"`
	Msg   string          `json:"msg"`
}

type UserValidation struct {
	ID               string    `json:"id"`
	TypeDocument     string    `json:"type_document"`
	DocumentNumber   int64     `json:"document_number"`
	ExpeditionDate   time.Time `json:"expedition_date"`
	Email            string    `json:"email"`
	FirstName        string    `json:"first_name"`
	SecondName       string    `json:"second_name"`
	SecondSurname    string    `json:"second_surname"`
	Age              int32     `json:"age"`
	SelfieImg        string    `json:"selfie_img"`
	BackDocumentImg  string    `json:"back_document_img"`
	FrontDocumentImg string    `json:"front_document_img"`
	Gender           string    `json:"gender"`
	Nationality      string    `json:"nationality"`
	CivilStatus      string    `json:"civil_status"`
	FirstSurname     string    `json:"first_surname"`
	BirthDate        time.Time `json:"birth_date"`
	Country          string    `json:"country"`
	Role             string    `json:"role"`
	TransactionId    string    `json:"transaction_id"`
	Department       string    `json:"department"`
	City             string    `json:"city"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
