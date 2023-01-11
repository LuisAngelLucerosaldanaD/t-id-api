package users

import (
	"check-id-api/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func RouterUser(app *fiber.App, db *sqlx.DB, txID string) {
	h := handlerUser{DB: db, TxID: txID}
	api := app.Group("/api")
	v1 := api.Group("/v1")
	user := v1.Group("/user")
	user.Post("/validate-identity", middleware.JWTProtected(), h.validateIdentity)
	user.Get("/user-session/:email", middleware.JWTProtected(), h.getUserSession)
}
