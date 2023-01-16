package blockchain

import (
	"bytes"
	"check-id-api/internal/logger"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetToken(url, email, password string) string {
	var token string
	request := AuthRequest{Email: email, Password: password}
	var response AuthResponse

	reqByte, err := json.Marshal(request)
	if err != nil {
		logger.Error.Printf("Error request en Timer: %v", err)
		return token
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqByte))
	if err != nil {
		logger.Error.Printf("Error http request en Timer: %v", err)
		return token
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error.Printf("Error enviando peticion authentication en Timer: %v", err)
		return token
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Printf("Error obteniendo respuesta de peticion authentication en Timer: %v", err)
		return token
	}

	err = json.Unmarshal(b, &response)
	if err != nil {
		logger.Error.Printf("Error decodificando respuesta de peticion authentication en Timer: %v", err)
		return token
	}
	token = fmt.Sprintf("Bearer %s", response.Data.AccessToken)
	return token

}
