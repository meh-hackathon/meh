package server

import (
	"crypto/rand"
	"database/sql"
	"math/big"
	math_rand "math/rand"
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
	err = db.Get(&qrCode,
		"INSERT INTO qr_codes (owner_id, slug, name) VALUES ($1, $2, $3) RETURNING id, owner_id, created_at, updated_at, slug",
		user.ID,
		slug,
		body.Name,
	)

	if err != nil {
		logger.Error("Failed to insert QR code", "error", err)
		httpx.WriteError(w, ErrInternalServerError.WithCause(err))
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, qrCode)
}

func (*Server) GetQrCodeById(w http.ResponseWriter, r *http.Request, id types.UUID) {
	user, err := auth.GetUser(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	var qrCode db.QrCode
	err = db.Get(
		&qrCode,
		"SELECT id, owner_id, created_at, updated_at, slug, name FROM qr_codes WHERE id = $1",
		id,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			httpx.WriteError(w, ErrNotFound.WithMessage("QR code not found"))
		} else {
			logger.Error("Failed to get QR code by ID", "error", err)
			httpx.WriteError(w, ErrInternalServerError.WithCause(err))
		}
		return
	}

	if qrCode.OwnerID.String() != user.ID.String() {
		httpx.WriteError(w, auth.ErrForbidden.WithMessage("You don't have permission to access this QR code"))
		return
	}

	httpx.WriteJSON(w, http.StatusOK, qrCode)
}

func (*Server) GetQrCodeBySlug(w http.ResponseWriter, r *http.Request, slug string) {
	var qrCode db.QrCode
	err := db.Get(&qrCode,
		"SELECT id, owner_id, created_at, updated_at, slug, name FROM qr_codes WHERE slug = $1 AND deleted_at IS NULL",
		slug,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			httpx.WriteError(w, ErrNotFound.WithOrigin().WithMessage("QR code not found"))
		} else {
			logger.Error("Failed to get QR code by slug", "error", err)
			httpx.WriteError(w, ErrInternalServerError.WithOrigin().WithCause(err))
		}
		return
	}

	type ContentWithWeight struct {
		db.Content
		Weight int `db:"weight"`
	}

	var contentItems []ContentWithWeight
	err = db.Select(
		&contentItems,
		`SELECT c.id, c.owner_id, c.created_at, c.updated_at, c.type, c.data, r.weight
		FROM contents c
		JOIN qr_code_to_content_relation r ON c.id = r.content_id
		WHERE r.qr_code_id = $1 AND c.deleted_at IS NULL`,
		qrCode.ID,
	)

	if err != nil {
		logger.Error("Failed to get QR code content", "error", err)
		httpx.WriteError(w, ErrInternalServerError.WithOrigin().WithCause(err))
		return
	}

	type QrCodeWithContent_ struct {
		QrCode  *db.QrCode
		Content *db.Content
	}
	response := QrCodeWithContent_{
		QrCode: &qrCode,
	}

	if len(contentItems) > 0 {
		totalWeight := 0
		for _, item := range contentItems {
			totalWeight += item.Weight
		}

		if totalWeight > 0 {
			randomValue := math_rand.Intn(totalWeight)

			currentWeight := 0
			for _, item := range contentItems {
				currentWeight += item.Weight
				if randomValue < currentWeight {
					response.Content = &item.Content
					break
				}
			}
		} else {
			response.Content = &contentItems[0].Content
		}
	}

	httpx.WriteJSON(w, http.StatusOK, response)
}

func (*Server) GetQrCodes(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUser(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	qrCodes := []db.QrCode{}
	err = db.Select(
		&qrCodes,
		"SELECT id, owner_id, created_at, updated_at, slug, name FROM qr_codes WHERE owner_id = $1 ORDER BY created_at DESC",
		user.ID,
	)

	if err != nil {
		logger.Error("Failed to get user's QR codes", "error", err)
		httpx.WriteError(w, ErrInternalServerError.WithCause(err))
		return
	}

	httpx.WriteJSON(w, http.StatusOK, qrCodes)
}
