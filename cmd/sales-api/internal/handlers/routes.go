package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/dimfeld/httptreemux/v5"
)

// API ...
func API(build string, shutdown chan os.Signal, log *log.Logger) http.Handler {
	httptreemux.New()
	return nil
}
