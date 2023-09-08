package models

import (
	"time"
)

type User struct {
	ID             string     `json:"id" db:"id" valid:"required,uuid"`
	TypeDocument   string     `json:"type_document" db:"type_document" valid:"-"`
	DocumentNumber string     `json:"document_number" db:"document_number" valid:"-"`
	ExpeditionDate *time.Time `json:"expedition_date" db:"expedition_date" valid:"-"`
	Email          string     `json:"email" db:"email" valid:"required"`
	FirstName      string     `json:"first_name" db:"first_name" valid:"-"`
	SecondName     string     `json:"second_name" db:"second_name" valid:"-"`
	SecondSurname  string     `json:"second_surname" db:"second_surname" valid:"-"`
	Age            int32      `json:"age" db:"age" valid:"-"`
	Gender         string     `json:"gender" db:"gender" valid:"-"`
	Nationality    string     `json:"nationality" db:"nationality" valid:"-"`
	CivilStatus    string     `json:"civil_status" db:"civil_status" valid:"-"`
	FirstSurname   string     `json:"first_surname" db:"first_surname" valid:"-"`
	BirthDate      *time.Time `json:"birth_date" db:"birth_date" valid:"-"`
	Country        string     `json:"country" db:"country" valid:"-"`
	Department     string     `json:"department" db:"department" valid:"-"`
	City           string     `json:"city" db:"city" valid:"-"`
	RealIp         string     `json:"real_ip" db:"real_ip" valid:"required"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

type UserToken struct {
	ID                 string    `json:"id" db:"id" valid:"required,uuid"`
	Nickname           string    `json:"nickname" db:"nickname" valid:"required"`
	Email              string    `json:"email" db:"email" valid:"required"`
	Password           string    `json:"password,omitempty" db:"password" valid:"required"`
	Name               string    `json:"name" db:"name" valid:"required"`
	Lastname           string    `json:"lastname" db:"lastname" valid:"required"`
	IdType             int       `json:"id_type" db:"id_type" valid:"required"`
	IdNumber           string    `json:"id_number" db:"id_number" valid:"required"`
	Cellphone          string    `json:"cellphone" db:"cellphone" valid:"required"`
	StatusId           int       `json:"status_id" db:"status_id" valid:"required"`
	FailedAttempts     int       `json:"failed_attempts,omitempty" db:"failed_attempts"`
	BlockDate          time.Time `json:"block_date,omitempty" db:"block_date"`
	DisabledDate       time.Time `json:"disabled_date,omitempty" db:"disabled_date"`
	LastLogin          time.Time `json:"last_login,omitempty" db:"last_login" `
	LastChangePassword time.Time `json:"last_change_password,omitempty" db:"last_change_password"`
	BirthDate          time.Time `json:"birth_date,omitempty" db:"birth_date"`
	VerifiedCode       string    `json:"verified_code,omitempty" db:"verified_code"`
	VerifiedAt         time.Time `json:"verified_at,omitempty" db:"verified_at"`
	IsDeleted          bool      `json:"is_deleted,omitempty" db:"is_deleted"`
	IdUser             string    `json:"id_user,omitempty" db:"id_user"`
	IdRole             int       `json:"id_role" db:"id_role" valid:"required"`
	FullPathPhoto      string    `json:"full_path_photo,omitempty" db:"full_path_photo"`
	RecoveryAccountAt  time.Time `json:"recovery_account_at,omitempty" db:"recovery_account_at"`
	RealIP             string    `json:"real_ip,omitempty"`
	DeletedAt          time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	CreatedAt          time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at,omitempty" db:"updated_at"`
}
