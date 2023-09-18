package blockchain

import (
	"check-id-api/internal/ciphers"
	"check-id-api/internal/env"
	"check-id-api/internal/grpc/accounting_proto"
	"check-id-api/internal/grpc/wallet_proto"
	"check-id-api/internal/logger"
	"check-id-api/internal/ws"
	"check-id-api/pkg/auth/user"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	grpcMetadata "google.golang.org/grpc/metadata"
	"os"
	"strconv"
	"strings"
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

func CreateTransaction(user *user.User, nameTransaction, descriptionTransaction, to string) (string, error) {

	e := env.NewConfiguration()

	identifier := []Identifier{
		{
			Name: "Información básica",
			Attributes: []Attribute{
				{
					Id:    1,
					Name:  "Primer Nombre",
					Value: strings.TrimSpace(*user.FirstName),
				},
				{
					Id:    2,
					Name:  "Segundo Nombre",
					Value: strings.TrimSpace(*user.SecondName),
				},
				{
					Id:    3,
					Name:  "Primer Apellido",
					Value: strings.TrimSpace(*user.FirstSurname),
				},
				{
					Id:    4,
					Name:  "Segundo Apellido",
					Value: strings.TrimSpace(*user.SecondSurname),
				},
				{
					Id:    6,
					Name:  "Número de Documento",
					Value: user.DocumentNumber,
				},
				{
					Id:    7,
					Name:  "Correo Electrónico",
					Value: user.Email,
				},
				{
					Id:    8,
					Name:  "Edad",
					Value: strconv.Itoa(int(*user.Age)),
				},
				{
					Id:    9,
					Name:  "Sexo",
					Value: *user.Gender,
				},
				{
					Id:    10,
					Name:  "Fecha de Nacimiento",
					Value: user.BirthDate.String(),
				},
				{
					Id:    12,
					Name:  "IP de Dispositivo",
					Value: user.RealIp,
				},
				{
					Id:    13,
					Name:  "Nacionalidad",
					Value: *user.Nationality,
				},
				{
					Id:    14,
					Name:  "Fecha de Actualización",
					Value: time.Now().UTC().String(),
				},
			},
		},
	}
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

	_, publicKeyPem, err := GetPublicKeyUser(user.DocumentNumber)
	if err != nil {
		logger.Error.Println("couldn't get public key of client", err)
		return "", err
	}

	publicKey, err := ciphers.DecodePublic(publicKeyPem)
	if err != nil {
		logger.Error.Println("couldn't decode client public key", err)
		return "", err
	}

	cryptoMessage, err := ciphers.CipherDH(*privateKey, *publicKey, dataBytes)

	transactionRq := Transaction{
		From:   e.Blockchain.Wallet,
		To:     to,
		TypeId: 18,
		Amount: 1,
		Files:  []*File{},
		Data:   cryptoMessage,
	}

	token := GetToken(e.Blockchain.UrlAuth, e.Blockchain.Email, e.Blockchain.Password)

	bodyRq, err := json.Marshal(transactionRq)
	if err != nil {
		logger.Error.Println("couldn't bind request", err)
		return "", err
	}

	hash := ciphers.StringToHashSha256(string(bodyRq))
	signValue, err := ciphers.SignWithEcdsa([]byte(hash), *privateKey)
	if err != nil {
		return "", err
	}

	headers := map[string]string{
		"sign":            signValue,
		"identity_number": user.DocumentNumber,
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

func CreateWallet(user *user.User) (*WalletInfo, error) {
	e := env.NewConfiguration()

	connAuth, err := grpc.Dial(e.AuthService.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error.Printf("error conectando con el servicio trx de blockchain: %s", err)
		return nil, err
	}
	defer connAuth.Close()

	clientWallet := wallet_proto.NewWalletServicesWalletClient(connAuth)
	clientAccount := accounting_proto.NewAccountingServicesAccountingClient(connAuth)

	token := GetToken(e.Blockchain.UrlAuth, e.Blockchain.Email, e.Blockchain.Password)

	ctx := grpcMetadata.AppendToOutgoingContext(context.Background(), "authorization", token)

	wallet, err := clientWallet.CreateWallet(ctx, &wallet_proto.RequestCreateWallet{
		IdentityNumber: user.DocumentNumber,
	})
	if err != nil {
		logger.Error.Printf("couldn't create wallet: %v", err)
		return nil, err
	}

	if wallet == nil {
		logger.Error.Printf("couldn't create wallet")
		return nil, err
	}

	if wallet.Error {
		logger.Error.Printf(wallet.Msg)
		return nil, err
	}

	resAccountTo, err := clientAccount.CreateAccounting(ctx, &accounting_proto.RequestCreateAccounting{
		Id:       uuid.New().String(),
		IdWallet: wallet.Data.Id,
		Amount:   0,
		IdUser:   user.ID,
	})
	if err != nil {
		logger.Error.Printf("couldn't create accounting to wallet: %v", err)
		return nil, err
	}

	if resAccountTo == nil {
		logger.Error.Printf("couldn't create accounting to wallet: %v", err)
		return nil, err
	}

	if resAccountTo.Error {
		logger.Error.Printf(resAccountTo.Msg)
		return nil, err
	}

	return &WalletInfo{
		Id:       wallet.Data.Id,
		Public:   wallet.Data.Key.Public,
		Private:  wallet.Data.Key.Private,
		Mnemonic: wallet.Data.Mnemonic,
	}, nil
}

func GetPublicKeyUser(identityNumber string) (string, string, error) {
	e := env.NewConfiguration()

	connAuth, err := grpc.Dial(e.AuthService.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error.Printf("error conectando con el servicio trx de blockchain: %s", err)
		return "", "", err
	}
	defer connAuth.Close()

	clientWallet := wallet_proto.NewWalletServicesWalletClient(connAuth)

	token := GetToken(e.Blockchain.UrlAuth, e.Blockchain.Email, e.Blockchain.Password)

	ctx := grpcMetadata.AppendToOutgoingContext(context.Background(), "authorization", token)

	resWs, err := clientWallet.GetWalletByIdentityNumber(ctx, &wallet_proto.RqGetByIdentityNumber{IdentityNumber: identityNumber})
	if err != nil {
		logger.Error.Printf("error al obtener la wallet del usuario: %s", err)
		return "", "", err
	}
	if resWs == nil {
		logger.Error.Printf("error al obtener la wallet del usuario")
		return "", "", fmt.Errorf("error al obtener la wallet del usuario")
	}

	if resWs.Error {
		logger.Error.Printf(resWs.Msg)
		return "", "", fmt.Errorf(resWs.Msg)
	}

	if resWs.Data == nil {
		return "", "", fmt.Errorf("no existe una wallet asociada a ese numero de identificacion")
	}

	return resWs.Data.Id, resWs.Data.Public, nil

}
