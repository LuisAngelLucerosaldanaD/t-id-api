package traceability

import (
	"check-id-api/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func RouterTraceability(app *fiber.App, db *sqlx.DB, txID string) {
	h := handlerTraceability{DB: db, TxID: txID}
	api := app.Group("/api")
	v1 := api.Group("/v1")
	tracking := v1.Group("/traceability")
	tracking.Get("/", middleware.JWTProtected(), h.getTraceability)
	tracking.Get("/life-test", middleware.JWTProtected(), h.getTraceabilityLifeTest)
}
