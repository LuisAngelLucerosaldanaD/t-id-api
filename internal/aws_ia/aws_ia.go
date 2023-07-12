package aws_ia

import (
	"bytes"
	"check-id-api/internal/env"
	"check-id-api/internal/logger"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

const (
	AWS_REGION = "us-east-1"
)

func CompareFacesV2(face1, face2 []byte) (bool, error) {
	e := env.NewConfiguration()
	face1Reader := bytes.NewReader(face1)
	face2Reader := bytes.NewReader(face2)
	body := new(bytes.Buffer)
	multipartWriter := multipart.NewWriter(body)

	face1Writer, err := multipartWriter.CreateFormFile("selfie", "selfie"+GetExtensionFromBytes(face1))
	if err != nil {
		return false, err
	}
	if _, err := io.Copy(face1Writer, face1Reader); err != nil {
		return false, err
	}

	face2Writer, err := multipartWriter.CreateFormFile("document", "document"+GetExtensionFromBytes(face2))
	if err != nil {
		return false, err
	}
	if _, err := io.Copy(face2Writer, face2Reader); err != nil {
		return false, err
	}

	multipartWriter.Close()

	request, err := http.NewRequest(http.MethodPost, e.FaceApi.CompareFace, body)
	request.Header.Add("Content-Type", multipartWriter.FormDataContentType())

	httpClient := &http.Client{}

	resp, err := httpClient.Do(request)
	if err != nil {
		logger.Error.Printf("no se  puedo enviar la petición: %v  -- log: ", err)
		return false, err
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
		return false, err
	}

	resFace := CompareFaceResponse{}

	err = json.Unmarshal(rsBody, &resFace)
	if err != nil {
		logger.Error.Printf("no se pudo parsear la respuesta del servio de comparacion de rostros: %v  -- log: ", err)
		return false, err
	}

	if resFace.Error {
		logger.Error.Printf("error al consumir el servicio ocr: " + resFace.Msg)
		return false, fmt.Errorf(resFace.Msg)
	}

	if resFace.Data.Verified == "true" {
		return true, nil
	}

	return false, nil
}

func GetExtensionFromBytes(file []byte) string {
	mimeType := http.DetectContentType(file)
	switch mimeType {
	case "image/png":
		return ".png"
	case "image/jpeg":
		return ".jpeg"
	case "image/tiff":
		return ".tif"
	}
	return ".txt"
}

func sessionAws() (client.ConfigProvider, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(AWS_REGION),
		Credentials: credentials.NewEnvCredentials(),
	})
	if err != nil {
		fmt.Println("error iniciando sesión con aws: #{err}")
		return sess, err
	}
	return sess, nil
}
