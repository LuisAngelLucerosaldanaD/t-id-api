package blockchain

import (
	"check-id-api/internal/ciphers"
	"check-id-api/internal/env"
	"check-id-api/internal/grpc/accounting_proto"
	"check-id-api/internal/grpc/auth_proto"
	"check-id-api/internal/grpc/users_proto"
	"check-id-api/internal/grpc/wallet_proto"
	"check-id-api/internal/logger"
	"check-id-api/internal/models"
	"check-id-api/internal/ws"
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

func CreateAccountAndWallet(user models.User) (*WalletInfo, error) {
	e := env.NewConfiguration()

	connAuth, err := grpc.Dial(e.AuthService.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error.Printf("error conectando con el servicio trx de blockchain: %s", err)
		return nil, err
	}
	defer connAuth.Close()

	clientWallet := wallet_proto.NewWalletServicesWalletClient(connAuth)
	clientUser := users_proto.NewAuthServicesUsersClient(connAuth)
	clientAuth := auth_proto.NewAuthServicesUsersClient(connAuth)
	clientAccount := accounting_proto.NewAccountingServicesAccountingClient(connAuth)

	resAuth, err := clientAuth.Login(context.Background(), &auth_proto.LoginRequest{
		Email:    &e.Blockchain.Email,
		Nickname: nil,
		Password: e.Blockchain.Password,
	})
	if err != nil {
		logger.Error.Printf("error al obtener el token de autorización: %s", err.Error())
		return nil, err
	}

	if resAuth == nil {
		logger.Error.Printf("error al obtener el token de autorización")
		return nil, fmt.Errorf("error al obtener el token de autorización")
	}

	if resAuth.Error {
		logger.Error.Printf(resAuth.Msg)
		return nil, fmt.Errorf(resAuth.Msg)
	}

	ctx := grpcMetadata.AppendToOutgoingContext(context.Background(), "authorization", resAuth.Data.AccessToken)

	resUser, err := clientUser.CreateUserBySystem(ctx, &users_proto.RequestCreateUserBySystem{
		Nickname:      user.FirstName + user.FirstSurname,
		Email:         user.Email,
		Password:      user.FirstName + strconv.FormatInt(user.DocumentNumber, 10),
		FullPathPhoto: "",
		Name:          strings.TrimSpace(user.FirstName + " " + user.SecondName),
		Lastname:      strings.TrimSpace(user.FirstSurname + " " + user.SecondSurname),
		IdType:        8,
		IdNumber:      strconv.FormatInt(user.DocumentNumber, 10),
		Cellphone:     "",
		BirthDate:     user.BirthDate.String(),
	})
	if err != nil {
		logger.Error.Printf("error al crear el usuario: %s", err.Error())
		return nil, err
	}

	if resUser == nil {
		logger.Error.Printf("error al crear el usuario")
		return nil, fmt.Errorf("error al crear el usuario")
	}

	if resUser.Error {
		logger.Error.Printf(resUser.Msg)
		return nil, fmt.Errorf(resUser.Msg)
	}

	wallet, err := clientWallet.CreateWalletBySystem(ctx, &wallet_proto.RqCreateWalletBySystem{
		IdentityNumber: resUser.Data.IdNumber,
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

	resUserWallet, err := clientUser.CreateUserWallet(ctx, &users_proto.RqCreateUserWallet{
		UserId:   resUser.Data.Id,
		WalletId: wallet.Data.Id,
	})
	if err != nil {
		logger.Error.Printf("couldn't create user wallet: %v", err)
		return nil, err
	}

	if resUserWallet == nil {
		logger.Error.Printf("couldn't create user wallet: %v", err)
		return nil, err
	}

	if resUserWallet.Error {
		logger.Error.Printf(resUserWallet.Msg)
		return nil, err
	}

	resAccountTo, err := clientAccount.CreateAccounting(ctx, &accounting_proto.RequestCreateAccounting{
		Id:       uuid.New().String(),
		IdWallet: wallet.Data.Id,
		Amount:   0,
		IdUser:   resUser.Data.Id,
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
