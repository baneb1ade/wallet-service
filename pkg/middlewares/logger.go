package middlewares

import (
	"github.com/julienschmidt/httprouter"
	"log/slog"
	"net/http"
)

func LoggingMiddleware(logger *slog.Logger, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		logger.Info("From " + r.RemoteAddr + " " + r.Method + " " + r.URL.Path)
		next(w, r, p)
	}
}
