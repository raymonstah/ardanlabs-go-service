package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/raymonstah/ardanlabs-go-service/internal/mid"
	"github.com/raymonstah/ardanlabs-go-service/internal/platform/web"
)

// API ...
func API(build string, shutdown chan os.Signal, log *log.Logger, db *sqlx.DB) *web.App {
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log))

	check := check{build: build, db: db}
	app.Handle(http.MethodGet, "/test", check.health)

	return app
}
