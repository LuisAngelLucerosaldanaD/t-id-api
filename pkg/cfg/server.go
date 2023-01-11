package cfg

import (
	"check-id-api/internal/models"
	"check-id-api/pkg/cfg/files"
	"check-id-api/pkg/cfg/files_s3"
	"check-id-api/pkg/cfg/messages"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	SrvMessage messages.PortsServerMessages
	SrvFiles   files.PortsServerFiles
	SrvFilesS3 files_s3.PortsServerFile
}

func NewServerCfg(db *sqlx.DB, user *models.User, txID string) *Server {

	repoMessage := messages.FactoryStorage(db, user, txID)
	srvMessage := messages.NewMessagesService(repoMessage, user, txID)

	repoFiles := files.FactoryStorage(db, user, txID)
	srvFiles := files.NewFilesService(repoFiles, user, txID)

	repoS3File := files_s3.FactoryFileDocumentRepository(user, txID)
	srvFilesS3 := files_s3.NewFileService(repoS3File, user, txID)

	return &Server{
		SrvMessage: srvMessage,
		SrvFiles:   srvFiles,
		SrvFilesS3: srvFilesS3,
	}
}
