package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/raymonstah/ardanlabs-go-service/internal/mid"
	"github.com/raymonstah/ardanlabs-go-service/internal/platform/web"
)

// API ...
func API(build string, shutdown chan os.Signal, log *log.Logger) *web.App {
	app := web.NewApp(shutdown, mid.Logger(log))

	app.Handle(http.MethodGet, "/test", health)

	return app
}
