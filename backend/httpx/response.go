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
	somaErr, ok := err.(*apperror.AppError)
	if !ok {
		somaErr = apperror.Define("api:internal_error", "An internal error occurred").WithCause(err)
	}

	if somaErr.StatusCode == 0 {
		logger.Error("Error without status code", "error", somaErr)
		somaErr.StatusCode = http.StatusInternalServerError
	}

	if somaErr.StatusCode >= http.StatusInternalServerError {
		logger.Error("Unhandled error", "error", somaErr)
	}

	if config.Env == config.Dev {
		cause := ""
		if somaErr.Cause != nil {
			cause = somaErr.Cause.Error()
		}

		devErr := devAppError{
			Code:       somaErr.Code,
			Message:    somaErr.Message,
			ApiMessage: somaErr.ApiMessage,
			Origin:     somaErr.Origin,
			StatusCode: somaErr.StatusCode,
			Values:     somaErr.Values,
			Cause:      cause,
		}

		WriteJSON(w, somaErr.StatusCode, devErr)
		return
	}

	WriteJSON(w, somaErr.StatusCode, somaErr)
}
