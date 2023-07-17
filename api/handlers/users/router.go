package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func RouterUser(app *fiber.App, db *sqlx.DB, txID string) {
	h := handlerUser{DB: db, TxID: txID}
	api := app.Group("/api")
	v1 := api.Group("/v1")
	user := v1.Group("/user")
	user.Post("/register", h.createUser)
	user.Post("/upload-selfie", h.uploadSelfie)
	user.Post("/upload-documents", h.uploadDocuments)
	user.Post("/basic-information", h.registerBasicInformation)
	user.Get("/user-session/:identifier", h.getUserSession)
	user.Get("/users-lasted/:email/:limit/:offset", h.getLastedUsers)
	user.Get("/data-pending", h.getUsersDataPending)
	user.Get("/validate/:identity_number", h.validateUser)
	// user.Post("/validation", h.validationFace)
}
