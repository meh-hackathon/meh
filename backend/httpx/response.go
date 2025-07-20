package httpx

import (
	"encoding/json"
	"net/http"

	"github.com/meh-hackathon/meh/apperror"
	"github.com/meh-hackathon/meh/config"
	"github.com/meh-hackathon/meh/logger"
)

type devAppError struct {
	Code       string         `json:"code"`
	Message    string         `json:"message"`
	ApiMessage string         `json:"api_message,omitempty"`
	Origin     string         `json:"origin,omitempty"`
	StatusCode int            `json:"status_code"`
	Values     map[string]any `json:"values,omitempty"`
	Cause      string         `json:"cause,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, err error) {
	appErr, ok := err.(*apperror.AppError)
	if !ok {
		appErr = apperror.Define("api:internal_error", "An internal error occurred").WithCause(err)
	}

	if appErr.StatusCode == 0 {
		logger.Error("Error without status code", "error", appErr)
		appErr.StatusCode = http.StatusInternalServerError
	}

	if appErr.StatusCode >= http.StatusInternalServerError {
		logger.Error("Unhandled error", "error", appErr)
	}

	if config.Env == config.Dev {
		cause := ""
		if appErr.Cause != nil {
			cause = appErr.Cause.Error()
		}

		devErr := devAppError{
			Code:       appErr.Code,
			Message:    appErr.Message,
			ApiMessage: appErr.ApiMessage,
			Origin:     appErr.Origin,
			StatusCode: appErr.StatusCode,
			Values:     appErr.Values,
			Cause:      cause,
		}

		WriteJSON(w, appErr.StatusCode, devErr)
		return
	}

	WriteJSON(w, appErr.StatusCode, appErr)
}
