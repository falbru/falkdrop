package middleware

import (
	"errors"
	"github.com/falbru/falkdrop/internal/api/errors"
	"log/slog"
	"net/http"
)

type HTTPHandlerWithErr func(w http.ResponseWriter, r *http.Request) error

func ErrorHandler(fn HTTPHandlerWithErr) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("incoming request", "method", r.Method, "path", r.URL.Path, "client_ip", r.RemoteAddr)

		if err := fn(w, r); err != nil {
			var httpErr httperror.HTTPError
			if errors.As(err, &httpErr) {
				http.Error(w, err.Error(), httpErr.Code)
				if httpErr.Code >= 500 {
					slog.Error("server error", "code", httpErr.Code, "error", err.Error())
				} else {
					slog.Warn("client error", "code", httpErr.Code, "error", err.Error())
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				slog.Error("internal server error", "error", err.Error())
			}
		}
	}
}
