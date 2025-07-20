//go:build !serveui

package server

import (
	"net/http"

	"github.com/meh-hackathon/meh/apperror"
	"github.com/meh-hackathon/meh/httpx"
)

var (
	ErrUiServingDisabled = apperror.Define("ui:disabled", "UI serving is disabled").WithStatus(http.StatusNotFound)
)

func GetUiHandler() (http.Handler, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteError(w, ErrUiServingDisabled.WithOrigin().WithApiMessage("This server is running as API server only. UI serving is disabled."))
	}), nil
}
