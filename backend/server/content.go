package server

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/meh-hackathon/meh/auth"
	"github.com/meh-hackathon/meh/db"
	"github.com/meh-hackathon/meh/httpx"
	"github.com/meh-hackathon/meh/logger"
	"github.com/oapi-codegen/runtime/types"
)

// CreateContent creates a new content item
func (s *Server) CreateContent(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUser(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	body, err := httpx.ParseReqBody[ContentCreateRequest](r)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	if body.Type == "" {
		httpx.WriteError(w, ErrBadRequest.WithMessage("Content type is required"))
		return
	}

	if body.Data == "" {
		httpx.WriteError(w, ErrBadRequest.WithMessage("Content data is required"))
		return
	}

	var content db.Content
	err = db.Get(&content,
		`INSERT INTO contents (owner_id, type, data)
		VALUES ($1, $2, $3)
		RETURNING id, owner_id, created_at, updated_at, type, data`,
		user.ID,
		body.Type,
		body.Data,
	)

	if err != nil {
		logger.Error("Failed to insert content", "error", err)
		httpx.WriteError(w, ErrInternalServerError.WithCause(err))
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, content)
}

// GetContentById retrieves a content item by its ID
func (s *Server) GetContentById(w http.ResponseWriter, r *http.Request, id types.UUID) {
	user, err := auth.GetUser(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	var content db.Content
	err = db.Get(&content,
		`SELECT id, owner_id, created_at, updated_at, type, data
		FROM contents
		WHERE id = $1
		AND deleted_at IS NULL`,
		id,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			httpx.WriteError(w, ErrNotFound.WithOrigin().WithMessage("Content not found"))
		} else {
			logger.Error("Failed to get content by ID", "error", err)
			httpx.WriteError(w, ErrInternalServerError.WithOrigin().WithCause(err))
		}
		return
	}

	if content.OwnerID.String() != user.ID.String() {
		httpx.WriteError(w, auth.ErrForbidden.WithMessage("You don't have permission to access this content"))
		return
	}

	httpx.WriteJSON(w, http.StatusOK, content)
}

// GetContentItems retrieves all content items owned by the authenticated user
func (s *Server) GetContentItems(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUser(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	contentItems := []db.Content{}
	err = db.Select(
		&contentItems,
		`SELECT id, owner_id, created_at, updated_at, type, data
		FROM contents
		WHERE owner_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC`,
		user.ID,
	)

	if err != nil {
		logger.Error("Failed to get user's content", "error", err)
		httpx.WriteError(w, ErrInternalServerError.WithOrigin().WithCause(err))
		return
	}

	httpx.WriteJSON(w, http.StatusOK, contentItems)
}

// UpdateContent updates an existing content item
func (s *Server) UpdateContent(w http.ResponseWriter, r *http.Request, id types.UUID) {
	user, err := auth.GetUser(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	var existingContent db.Content
	err = db.Get(&existingContent,
		`SELECT id, owner_id FROM contents WHERE id = $1 AND deleted_at IS NULL`,
		id,
	)

	if errors.Is(err, sql.ErrNoRows) {
		httpx.WriteError(w, ErrNotFound.WithOrigin().WithMessage("Content not found"))
		return
	}

	if err != nil {
		logger.Error("Failed to check content ownership", "error", err)
		httpx.WriteError(w, ErrInternalServerError.WithCause(err))
		return
	}

	if existingContent.OwnerID.String() != user.ID.String() {
		httpx.WriteError(w, auth.ErrForbidden.WithOrigin().WithMessage("You don't have permission to update this content"))
		return
	}

	body, err := httpx.ParseReqBody[ContentUpdateRequest](r)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	query := "UPDATE contents SET updated_at = NOW()"
	params := []interface{}{}
	paramIndex := 1

	if body.Type != nil && *body.Type != "" {
		query += fmt.Sprintf(", type = $%d", paramIndex)
		params = append(params, body.Type)
		paramIndex++
	}

	if body.Data != nil && *body.Data != "" {
		query += fmt.Sprintf(", data = $%d", paramIndex)
		params = append(params, body.Data)
		paramIndex++
	}

	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, owner_id, created_at, updated_at, type, data", paramIndex)
	params = append(params, id)

	var updatedContent db.Content
	err = db.Get(&updatedContent, query, params...)

	if err != nil {
		logger.Error("Failed to update content", "error", err)
		httpx.WriteError(w, ErrInternalServerError.WithOrigin().WithCause(err))
		return
	}

	httpx.WriteJSON(w, http.StatusOK, updatedContent)
}

// DeleteContent soft-deletes a content item
func (s *Server) DeleteContent(w http.ResponseWriter, r *http.Request, id types.UUID) {
	user, err := auth.GetUser(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	var existingContent db.Content
	err = db.Get(&existingContent,
		`SELECT id, owner_id FROM contents WHERE id = $1 AND deleted_at IS NULL`,
		id,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			httpx.WriteError(w, ErrNotFound.WithOrigin().WithMessage("Content not found"))
		} else {
			logger.Error("Failed to check content ownership", "error", err)
			httpx.WriteError(w, ErrInternalServerError.WithOrigin().WithCause(err))
		}
		return
	}

	if existingContent.OwnerID.String() != user.ID.String() {
		httpx.WriteError(w, auth.ErrForbidden.WithMessage("You don't have permission to delete this content"))
		return
	}

	_, err = db.Exec(
		`UPDATE contents SET deleted_at = NOW() WHERE id = $1`,
		id,
	)

	if err != nil {
		logger.Error("Failed to delete content", "error", err)
		httpx.WriteError(w, ErrInternalServerError.WithCause(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// LinkContentToQrCode links a content item to a QR code
func (s *Server) LinkContentToQrCode(w http.ResponseWriter, r *http.Request, qrCodeId types.UUID) {
	user, err := auth.GetUser(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	body, err := httpx.ParseReqBody[LinkContentRequest](r)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	contentId := body.ContentId

	var content db.Content
	err = db.Get(&content,
		`SELECT id, owner_id FROM contents WHERE id = $1 AND deleted_at IS NULL`,
		contentId,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			httpx.WriteError(w, ErrNotFound.WithOrigin().WithMessage("Content not found"))
		} else {
			logger.Error("Failed to get content", "error", err)
			httpx.WriteError(w, ErrInternalServerError.WithOrigin().WithCause(err))
		}
		return
	}

	if content.OwnerID.String() != user.ID.String() {
		httpx.WriteError(w, auth.ErrForbidden.WithMessage("You don't have permission to access this content"))
		return
	}

	var qrCode db.QrCode
	err = db.Get(&qrCode,
		`SELECT id, owner_id FROM qr_codes WHERE id = $1 AND deleted_at IS NULL`,
		qrCodeId,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			httpx.WriteError(w, ErrNotFound.WithOrigin().WithMessage("QR code not found"))
		} else {
			logger.Error("Failed to get QR code", "error", err)
			httpx.WriteError(w, ErrInternalServerError.WithOrigin().WithCause(err))
		}
		return
	}

	if qrCode.OwnerID.String() != user.ID.String() {
		httpx.WriteError(w, auth.ErrForbidden.WithMessage("You don't have permission to access this QR code"))
		return
	}

	_, err = db.Exec(
		`INSERT INTO qr_code_to_content_relation (qr_code_id, content_id, weight)
		VALUES ($1, $2, 1)
		ON CONFLICT (qr_code_id, content_id) DO UPDATE SET weight = EXCLUDED.weight`,
		qrCodeId,
		contentId,
	)

	if err != nil {
		logger.Error("Failed to link content to QR code", "error", err)
		httpx.WriteError(w, ErrInternalServerError.WithOrigin().WithCause(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// UnlinkContentFromQrCode removes the link between a content item and QR code
func (s *Server) UnlinkContentFromQrCode(w http.ResponseWriter, r *http.Request, qrCodeId types.UUID) {
	user, err := auth.GetUser(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	body, err := httpx.ParseReqBody[LinkContentRequest](r)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	contentId := body.ContentId

	var content db.Content
	err = db.Get(&content,
		`SELECT id, owner_id FROM contents WHERE id = $1 AND deleted_at IS NULL`,
		contentId,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			httpx.WriteError(w, ErrNotFound.WithOrigin().WithMessage("Content not found"))
		} else {
			logger.Error("Failed to get content", "error", err)
			httpx.WriteError(w, ErrInternalServerError.WithOrigin().WithCause(err))
		}
		return
	}

	if content.OwnerID.String() != user.ID.String() {
		httpx.WriteError(w, auth.ErrForbidden.WithMessage("You don't have permission to access this content"))
		return
	}

	var qrCode db.QrCode
	err = db.Get(&qrCode,
		`SELECT id, owner_id FROM qr_codes WHERE id = $1 AND deleted_at IS NULL`,
		qrCodeId,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			httpx.WriteError(w, ErrNotFound.WithOrigin().WithMessage("QR code not found"))
		} else {
			logger.Error("Failed to get QR code", "error", err)
			httpx.WriteError(w, ErrInternalServerError.WithOrigin().WithCause(err))
		}
		return
	}

	if qrCode.OwnerID.String() != user.ID.String() {
		httpx.WriteError(w, auth.ErrForbidden.WithMessage("You don't have permission to access this QR code"))
		return
	}

	_, err = db.Exec(
		`DELETE FROM qr_code_to_content_relation WHERE qr_code_id = $1 AND content_id = $2`,
		qrCodeId,
		contentId,
	)

	if err != nil {
		logger.Error("Failed to unlink content from QR code", "error", err)
		httpx.WriteError(w, ErrInternalServerError.WithOrigin().WithCause(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetQrCodeContent retrieves all content linked to a specific QR code
func (s *Server) GetQrCodeContent(w http.ResponseWriter, r *http.Request, qrCodeId types.UUID) {
	user, err := auth.GetUser(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	var qrCode db.QrCode
	err = db.Get(&qrCode,
		`SELECT id, owner_id FROM qr_codes WHERE id = $1 AND deleted_at IS NULL`,
		qrCodeId,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			httpx.WriteError(w, ErrNotFound.WithOrigin().WithMessage("QR code not found"))
		} else {
			logger.Error("Failed to get QR code", "error", err)
			httpx.WriteError(w, ErrInternalServerError.WithOrigin().WithCause(err))
		}
		return
	}

	if qrCode.OwnerID.String() != user.ID.String() {
		httpx.WriteError(w, auth.ErrForbidden.WithMessage("You don't have permission to access this QR code"))
		return
	}

	contentItems := []db.Content{}
	err = db.Select(
		&contentItems,
		`SELECT c.id, c.owner_id, c.created_at, c.updated_at, c.type, c.data
		FROM contents c
		JOIN qr_code_to_content_relation r ON c.id = r.content_id
		WHERE r.qr_code_id = $1 AND c.deleted_at IS NULL
		ORDER BY r.weight ASC, c.created_at DESC`,
		qrCodeId,
	)

	if err != nil {
		logger.Error("Failed to get QR code content", "error", err)
		httpx.WriteError(w, ErrInternalServerError.WithOrigin().WithCause(err))
		return
	}

	httpx.WriteJSON(w, http.StatusOK, contentItems)
}
