package user_temp

import (
	"check-id-api/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func RouterUserTemp(app *fiber.App, db *sqlx.DB, txID string) {
	h := handlerUserTemp{DB: db, TxID: txID}
	api := app.Group("/api")
	v1 := api.Group("/v1")
	user := v1.Group("/user_temp")
	user.Post("/create", middleware.JWTProtected(), h.createUserTemp)
}
