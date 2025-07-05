package server

import (
	"html"
	"net/http"

	"github.com/platform-mesh/golang-commons/logger"
)

type requestLogger struct {
	log *logger.Logger
}

func (rl *requestLogger) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.log.Debug().Msgf("Request from %s %s %s", r.RemoteAddr, r.Method, html.EscapeString(r.URL.Path))
		h.ServeHTTP(w, r)
	})
}
