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
	user.Post("/register", h.createUser)
	user.Post("/login", h.Login)
	user.Get("/user-session/:identifier", middleware.JWTProtected(), h.getUserSession)
	user.Get("/role", middleware.JWTProtected(), h.getUserSession)
	user.Get("/validate", middleware.JWTProtected(), h.validateUser)
	user.Get("/finish-onboarding", middleware.JWTProtected(), h.getFinishOnboarding)
	user.Get("/finish-validation", middleware.JWTProtected(), h.getFinishValidationIdentity)
	user.Get("/file/:id", middleware.JWTProtected(), h.getUserFile)
	user.Get("/:identifier", middleware.JWTProtected(), h.getUserSession)
}
