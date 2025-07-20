package server

import (
	"net/http"

	"github.com/meh-hackathon/meh/logger"
	"github.com/oapi-codegen/runtime/types"
)

func (*Server) CreateQrCode(w http.ResponseWriter, r *http.Request) {
	logger.Info("CreateQrCode")
}

func (*Server) GetQrCodeById(w http.ResponseWriter, r *http.Request, uuid types.UUID) {
	logger.Info("GetQrCode")
}

func (*Server) GetQrCodeBySlug(w http.ResponseWriter, r *http.Request, slug string) {
	logger.Info("GetQrCode")
}

func (*Server) GetQrCodes(w http.ResponseWriter, r *http.Request) {}
