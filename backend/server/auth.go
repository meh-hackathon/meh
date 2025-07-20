package server

import (
	"net/http"

	"github.com/meh-hackathon/meh/logger"
)

func (*API) LoginUser(w http.ResponseWriter, r *http.Request) {
	logger.Info("LoginUser")
}

func (*API) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	logger.Info("GetCurrentUser")
}
