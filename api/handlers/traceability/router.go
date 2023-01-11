package users

import (
	"check-id-api/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func RouterTraceability(app *fiber.App, db *sqlx.DB, txID string) {
	h := handlerTraceability{DB: db, TxID: txID}
	api := app.Group("/api")
	v1 := api.Group("/v1")
	user := v1.Group("/traceability")
	user.Get("/user-session/:userID", middleware.JWTProtected(), h.getTraceabilitySession)
}