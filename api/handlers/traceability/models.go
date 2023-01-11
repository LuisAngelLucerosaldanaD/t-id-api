package users

import "check-id-api/pkg/trx/traceability"

type resTraceability struct {
	Error bool                         `json:"error"`
	Data  []*traceability.Traceability `json:"data"`
	Code  int                          `json:"code"`
	Type  int                          `json:"type"`
	Msg   string                       `json:"msg"`
}
