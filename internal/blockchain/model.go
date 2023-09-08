package blockchain

import "time"

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Error bool     `json:"error"`
	Data  dataAuth `json:"data"`
	Code  int      `json:"code"`
	Msg   string   `json:"msg"`
}

type dataAuth struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Transaction struct {
	From   string  `json:"from,omitempty"`
	To     string  `json:"to,omitempty"`
	Amount float64 `json:"amount,omitempty"`
	TypeId int     `json:"type_id,omitempty"`
	Data   string  `json:"data,omitempty"`
	Files  []*File `json:"file,omitempty"`
}
type DataCreateTransaction struct {
	Category    string       `json:"category"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Identifiers []Identifier `json:"identifiers"`
	ExpiresAt   *time.Time   `json:"expires_at"`
	Type        int32        `json:"type"`
	Id          string       `json:"id"`         // id de la credencial
	Status      string       `json:"status"`     // estado de la credencial
	CreatedAt   string       `json:"created_at"` // fecha de creación de la credencial
}

type Identifier struct {
	Name       string      `json:"name"`
	Attributes []Attribute `json:"attributes"`
}

type Attribute struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type File struct {
	FileID     int    `json:"id_file"`
	Name       string `json:"name"`
	FileEncode string `json:"file_encode"`
}

type ResponseCreateTransaction struct {
	Error bool                           `json:"error"`
	Data  *DataResponseCreateTransaction `json:"data"`
	Code  int                            `json:"code"`
	Type  int                            `json:"type"`
	Msg   string                         `json:"msg"`
}

type DataResponseCreateTransaction struct {
	Id        string  `json:"id,omitempty"`
	From      string  `json:"from,omitempty"`
	To        string  `json:"to,omitempty"`
	Amount    float64 `json:"amount,omitempty"`
	TypeId    int32   `json:"type_id,omitempty"`
	Data      string  `json:"data,omitempty"`
	Block     int64   `json:"block,omitempty"`
	Files     string  `json:"file,omitempty"`
	CreatedAt string  `json:"created_at,omitempty"`
	UpdatedAt string  `json:"updated_at,omitempty"`
}

type WalletInfo struct {
	Id       string `json:"id"`
	Public   string `json:"public"`
	Private  string `json:"private"`
	Mnemonic string `json:"mnemonic"`
}
