package models

import (
	"time"
)

type User struct {
	ID             string    `json:"id" db:"id" valid:"required,uuid"`
	TypeDocument   string    `json:"type_document" db:"type_document" valid:"required"`
	DocumentNumber int64     `json:"document_number" db:"document_number" valid:"required"`
	ExpeditionDate time.Time `json:"expedition_date" db:"expedition_date" valid:"required"`
	Email          string    `json:"email" db:"email" valid:"required"`
	FirstName      string    `json:"first_name" db:"first_name" valid:"required"`
	SecondName     string    `json:"second_name" db:"second_name" valid:"required"`
	SecondSurname  string    `json:"second_surname" db:"second_surname" valid:"required"`
	Age            int32     `json:"age" db:"age" valid:"required"`
	Gender         string    `json:"gender" db:"gender" valid:"required"`
	Nationality    string    `json:"nationality" db:"nationality" valid:"required"`
	CivilStatus    string    `json:"civil_status" db:"civil_status" valid:"required"`
	FirstSurname   string    `json:"first_surname" db:"first_surname" valid:"required"`
	BirthDate      time.Time `json:"birth_date" db:"birth_date" valid:"required"`
	Country        string    `json:"country" db:"country" valid:"required"`
	Department     string    `json:"department" db:"department" valid:"required"`
	City           string    `json:"city" db:"city" valid:"required"`
	RealIp         string    `json:"real_ip" db:"real_ip" valid:"required"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
