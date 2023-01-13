package work

type resAllWork struct {
	Error bool   `json:"error"`
	Data  Status `json:"data"`
	Code  int    `json:"code"`
	Type  int    `json:"type"`
	Msg   string `json:"msg"`
}

type Status struct {
	Valid     int `json:"valid"`
	Pending   int `json:"pending"`
	Refused   int `json:"refused"`
	Total     int `json:"total"`
	Expired   int `json:"expired"`
	NotStated int `json:"not_stated"`
}

type resAnny struct {
	Error bool   `json:"error"`
	Data  string `json:"data"`
	Code  int    `json:"code"`
	Type  int    `json:"type"`
	Msg   string `json:"msg"`
}

type ReqAccept struct {
	UserID string `json:"user_id"`
}
