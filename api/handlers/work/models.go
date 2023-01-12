package work

import (
	"check-id-api/pkg/wf/work_validation"
)

type resAllWork struct {
	Error bool                              `json:"error"`
	Data  []*work_validation.WorkValidation `json:"data"`
	Code  int                               `json:"code"`
	Type  int                               `json:"type"`
	Msg   string                            `json:"msg"`
}
