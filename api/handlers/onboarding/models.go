package onboarding

type resCreateOnboarding struct {
	Error bool   `json:"error"`
	Data  string `json:"data"`
	Code  int    `json:"code"`
	Type  int    `json:"type"`
	Msg   string `json:"msg"`
}

type requestCreateOnboarding struct {
	DocumentNumber string  `json:"document_number"`
	Email          string  `json:"email"`
	FirstName      *string `json:"first_name"`
	SecondName     *string `json:"second_name"`
	FirstSurname   *string `json:"first_surname"`
	SecondSurname  *string `json:"second_surname"`
	Nationality    *string `json:"nationality"`
	Cellphone      string  `json:"cellphone"`
	ClientId       int64   `json:"client_id"`
	RequestId      string  `json:"request_id"`
}

type RequestProcessOnboarding struct {
	UserID        string `json:"user_id"`
	Selfie        string `json:"selfie"`
	DocumentFront string `json:"document_front"`
	DocumentBack  string `json:"document_back"`
	Onboarding    string `json:"onboarding"`
}

type ResProcessOnboarding struct {
	Error bool        `json:"error"`
	Data  interface{} `json:"data"`
	Code  int         `json:"code"`
	Type  int         `json:"type"`
	Msg   string      `json:"msg"`
}

type RequestValidationIdentity struct {
	FaceImage    string `json:"face_image"`
	UserID       string `json:"user_id"`
	ValidationId int64  `json:"validation_id"`
}
