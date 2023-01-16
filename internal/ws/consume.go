package ws

import (
	"bytes"
	"check-id-api/internal/logger"
	"io"
	"io/ioutil"
	"net/http"
)

func ConsumeWS(jsonBytes []byte, url, method, token string, headers *map[string]string) ([]byte, int, error) {
	var req http.Request

	if method == "POST" {
		resp, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBytes))
		if err != nil {
			logger.Error.Printf("no se  puedo obtener respuesta: %v  -- log: ", err)
			return nil, 1, err
		}
		req = *resp
	}

	if method == "GET" {
		resp, err := http.NewRequest(method, url, nil)
		if err != nil {
			logger.Error.Printf("no se  puedo obtener respuesta: %v  -- log: ", err)
			return nil, 1, err
		}
		req = *resp
	}

	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", token)
	}

	if headers != nil && len(*headers) > 0 {
		for key, value := range *headers {
			req.Header.Set(key, value)
		}
	}

	client := &http.Client{}

	resp, err := client.Do(&req)
	if err != nil {
		logger.Error.Printf("no se  puedo enviar la petición: %v  -- log: ", err)
		return nil, 500, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error.Printf("no se pudo ejecutar defer body close: %v  -- log: ", err)
		}
	}(resp.Body)

	rsBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Printf("no se  puedo obtener respuesta: %v  -- log: ", err)
		return rsBody, resp.StatusCode, err
	}

	return rsBody, resp.StatusCode, nil
}
