package clients

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func RouterClients(app *fiber.App, db *sqlx.DB, txID string) {
	h := handlerWork{DB: db, TxID: txID}
	api := app.Group("/api")
	v1 := api.Group("/v1")
	user := v1.Group("/client")
	user.Get("/:nit", h.getDataClient)
	user.Post("/", h.CreateClient)
}
