package server

import (
	"net/http"

	"github.com/meh-hackathon/meh/logger"
	"github.com/oapi-codegen/runtime/types"
)

func (*API) CreateQrCode(w http.ResponseWriter, r *http.Request) {
	logger.Info("CreateQrCode")
}

func (*API) GetQrCodeById(q http.ResponseWriter, r *http.Request, uuid types.UUID) {
	logger.Info("GetQrCode")
}

func (*API) GetQrCodeBySlug(q http.ResponseWriter, r *http.Request, slug string) {
	logger.Info("GetQrCode")
}
