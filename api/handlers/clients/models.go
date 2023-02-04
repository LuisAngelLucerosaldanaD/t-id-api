package clients

type ResClient struct {
	Error bool   `json:"error"`
	Data  Client `json:"data"`
	Code  int    `json:"code"`
	Type  int    `json:"type"`
	Msg   string `json:"msg"`
}

type Client struct {
	ID          int64  `json:"id"`
	FullName    string `json:"full_name"`
	Nit         string `json:"nit"`
	Banner      string `json:"banner"`
	LogoSmall   string `json:"logo_small"`
	MainColor   string `json:"main_color"`
	SecondColor string `json:"second_color"`
	UrlRedirect string `json:"url_redirect"`
	UrlApi      string `json:"url_api"`
}
