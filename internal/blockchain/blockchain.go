package blockchain

import (
	"check-id-api/internal/ciphers"
	"check-id-api/internal/env"
	"check-id-api/internal/logger"
	"check-id-api/internal/ws"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"os"
	"time"
)

var privateKey *ecdsa.PrivateKey

func init() {
	c := env.NewConfiguration()
	private := c.App.EcdsaPrivateKey
	signBytes, err := os.ReadFile(private)
	if err != nil {
		logger.Error.Printf("leyendo el archivo privado de firma: %s", err)
	}

	privateKey, err = ciphers.DecodePrivate(string(signBytes))
	if err != nil {
		logger.Error.Printf("realizando el parse de la llave ecdsa: %s", err)
	}
}

func CreateTransaction(identifier []Identifier, nameTransaction, descriptionTransaction, to, identityNumber string) (string, error) {
	e := env.NewConfiguration()
	res := ResponseCreateTransaction{}
	resData := DataResponseCreateTransaction{}
	dataTrx := DataCreateTransaction{
		Category:    "2e59a864-b7ff-45d9-be8c-7d1b9513f7c5",
		Name:        nameTransaction,
		Description: descriptionTransaction,
		Identifiers: identifier,
		Type:        1,
		Id:          uuid.New().String(),
		Status:      "active",
		CreatedAt:   time.Now().String(),
	}

	dataBytes, _ := json.Marshal(dataTrx)

	transactionRq := Transaction{
		From:           e.Blockchain.Wallet,
		To:             to,
		TypeId:         18,
		Amount:         1,
		IdentityNumber: identityNumber,
		Files:          []*File{},
		Data:           string(dataBytes),
	}

	token := GetToken(e.Blockchain.UrlAuth, e.Blockchain.Email, e.Blockchain.Password)

	bodyRq, err := json.Marshal(transactionRq)
	if err != nil {
		logger.Error.Println("couldn't bind request", err)
		return "", err
	}

	hash := ciphers.StringToHashSha256(string(dataBytes))
	signValue, err := ciphers.SignWithEcdsa([]byte(hash), *privateKey)
	if err != nil {
		return "", err
	}

	headers := map[string]string{
		"sign": signValue,
	}

	rs, codeHTTP, err := ws.ConsumeWS(bodyRq, e.Blockchain.UrlApi, "POST", token, &headers)
	if err := json.Unmarshal(rs, &res); err != nil {
		logger.Error.Println("don't bind response in struct", err)
		return "", err
	}
	if codeHTTP != 200 {
		err = errors.New(fmt.Sprintf("respuesta diferente a http 200, %d", codeHTTP))
		return "", err
	}
	if res.Error {
		err = errors.New(fmt.Sprintf("respuesta con error, %d", res.Code))
		return "", err
	}
	if res.Data == nil {
		return "", err
	}

	byteData, _ := json.Marshal(res.Data)
	err = json.Unmarshal(byteData, &resData)
	if err != nil {
		logger.Error.Println("couldn't bind response un Unmarshal", err)
		return "", err
	}
	return resData.Id, nil

}
