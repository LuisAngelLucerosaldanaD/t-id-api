package onboarding

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func RouterOnboarding(app *fiber.App, db *sqlx.DB, txID string) {
	h := handlerOnboarding{DB: db, TxID: txID}
	api := app.Group("/api")
	v1 := api.Group("/v1")
	user := v1.Group("/onboarding")
	user.Post("/", h.Onboarding)
	user.Post("/process", h.FinishOnboarding)
	user.Post("/validate-identity", h.ValidateIdentity)
}
