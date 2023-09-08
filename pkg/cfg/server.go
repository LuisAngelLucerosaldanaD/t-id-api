package cfg

import (
	"check-id-api/internal/models"
	"check-id-api/pkg/cfg/client"
	"check-id-api/pkg/cfg/file"
	"check-id-api/pkg/cfg/files_s3"
	"check-id-api/pkg/cfg/messages"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	SrvMessage messages.PortsServerMessages
	SrvFiles   file.PortsServerFile
	SrvFilesS3 files_s3.PortsServerFile
	SrvClients client.PortsServerClient
}

func NewServerCfg(db *sqlx.DB, user *models.User, txID string) *Server {

	repoMessage := messages.FactoryStorage(db, user, txID)
	srvMessage := messages.NewMessagesService(repoMessage, user, txID)

	repoFiles := file.FactoryStorage(db, user, txID)
	srvFiles := file.NewFileService(repoFiles, user, txID)

	repoS3File := files_s3.FactoryFileDocumentRepository(user, txID)
	srvFilesS3 := files_s3.NewFileService(repoS3File, user, txID)

	repoClients := client.FactoryStorage(db, user, txID)
	srvClients := client.NewClientService(repoClients, user, txID)

	return &Server{
		SrvMessage: srvMessage,
		SrvFiles:   srvFiles,
		SrvFilesS3: srvFilesS3,
		SrvClients: srvClients,
	}
}
