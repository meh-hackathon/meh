package server

import (
	"crypto/rand"
	"math/big"
	"net/http"

	"github.com/meh-hackathon/meh/auth"
	"github.com/meh-hackathon/meh/db"
	"github.com/meh-hackathon/meh/httpx"
	"github.com/meh-hackathon/meh/logger"
	"github.com/oapi-codegen/runtime/types"
)

func generateSlug() string {
	// generate a random 10 character string
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 10

	result := make([]byte, length)

	// Using crypto/rand for secure random generation
	for i := range result {
		// Generate random index within charset
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// Log the error
			logger.Error("Failed to generate random number using crypto/rand", "error", err)
			// In case of error, use a simple modulo operation as fallback
			// This is less secure but ensures we still get a slug
			result[i] = charset[i%len(charset)]
			continue
		}

		result[i] = charset[randomIndex.Int64()]
	}

	return string(result)
}

func (*Server) CreateQrCode(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUser(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	body, err := httpx.ParseReqBody[QrCodeCreateRequest](r)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	var slug string
	if body.Slug == nil {
		slug = generateSlug()
	} else {
		slug = *body.Slug
	}

	var qrCode db.QrCode
	err = db.DB.QueryRowx(
		"INSERT INTO qr_codes (owner_id, slug) VALUES ($1, $2) RETURNING id, owner_id, created_at, updated_at, slug",
		user.ID,
		slug,
	).StructScan(&qrCode)

	if err != nil {
		logger.Error("Failed to insert QR code", "error", err)
		httpx.WriteError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, qrCode)
}

func (*Server) GetQrCodeById(w http.ResponseWriter, r *http.Request, uuid types.UUID) {
	logger.Info("GetQrCode")
}

func (*Server) GetQrCodeBySlug(w http.ResponseWriter, r *http.Request, slug string) {
	logger.Info("GetQrCode")
}

func (*Server) GetQrCodes(w http.ResponseWriter, r *http.Request) {
	logger.Info("GetQrCodes")
}
