package traceability

import (
	"check-id-api/pkg/trx/traceability"
	"time"
)

type resTraceability struct {
	Error bool                         `json:"error"`
	Data  []*traceability.Traceability `json:"data"`
	Code  int                          `json:"code"`
	Type  int                          `json:"type"`
	Msg   string                       `json:"msg"`
}

type ResTrackingValidation struct {
	Error bool        `json:"error"`
	Data  []*Tracking `json:"data"`
	Code  int         `json:"code"`
	Type  int         `json:"type"`
	Msg   string      `json:"msg"`
}

type Tracking struct {
	ID               int64     `json:"id"`
	ClientId         int64     `json:"client_id"`
	MaxNumValidation int       `json:"max_num_validation"`
	RequestId        string    `json:"request_id"`
	ExpiredAt        time.Time `json:"expired_at"`
	UserID           string    `json:"user_id"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
