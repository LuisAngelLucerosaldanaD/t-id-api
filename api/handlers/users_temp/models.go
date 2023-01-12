package user_temp

type resAnny struct {
	Error bool   `json:"error"`
	Data  string `json:"data"`
	Code  int    `json:"code"`
	Type  int    `json:"type"`
	Msg   string `json:"msg"`
}

type RqCreateUserTemp struct {
	FullName string `json:"full_name"`
	Surname  string `json:"surname"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	Email    string `json:"email"`
	Domain   string `json:"domain"`
}
