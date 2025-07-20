package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/meh-hackathon/meh/db"
	"github.com/meh-hackathon/meh/httpx"
)

// Based on the grant type the client will send username + password or refresh token
type tokenRequest struct {
	GrantType    string `json:"grant_type"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type oAuthHandler struct {
	authenticators []Authenticator
}

var passwordHashCost = 12

type Authenticator interface {
	authenticate(ctx context.Context, username, password string) (*User, error)
}

type OauthOption func(*oAuthHandler) error

func OAuthHandler(router *http.ServeMux, opts ...OauthOption) error {
	handler := &oAuthHandler{}
	for _, opt := range opts {
		if err := opt(handler); err != nil {
			return fmt.Errorf("failed to apply OAuth option: %w", err)
		}
	}
	if len(handler.authenticators) == 0 {
		panic("auth: at least one authenticator must be set for OAuth handler")
	}

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		req, err := httpx.ParseReqBody[tokenRequest](r)
		if err != nil {
			httpx.WriteError(w, err)
			return
		}

		switch req.GrantType {
		case "password":
			handler.handlePasswordGrant(w, r, req)
		case "refresh_token":
			handler.handleRefreshGrant(w, r, req)
		default:
			httpx.WriteError(w, ErrInvalidToken.WithApiMessagef("Unsupported grant type: %s. Supported types are: password, refresh_token", req.GrantType))
		}
	}

	router.HandleFunc("POST /oauth/token", handlerFunc)
	router.HandleFunc("GET /oauth/user", UserHandler)

	return nil
}

func (handler *oAuthHandler) handlePasswordGrant(w http.ResponseWriter, r *http.Request, req tokenRequest) {
	// check all authenticators concurrently and check for the first successful authentication
	type result struct {
		user *User
		err  error
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var wg sync.WaitGroup
	resultChannel := make(chan result, len(handler.authenticators))
	for _, authenticator := range handler.authenticators {
		wg.Add(1)
		go func(auth Authenticator) {
			defer wg.Done()
			user, err := auth.authenticate(ctx, req.Username, req.Password)
			resultChannel <- result{user: user, err: err}
		}(authenticator)
	}

	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	var user *User
	for res := range resultChannel {
		if res.err != nil && !errors.Is(res.err, ErrUnauthorized) {
			httpx.WriteError(w, res.err)
			return
		}
		if res.err == nil && res.user != nil {
			// If we found a user, cancel the context to stop other goroutines
			cancel()
			user = res.user
			break
		}
	}
	if user == nil {
		httpx.WriteError(w, ErrUnauthorized.WithApiMessage("invalid credentials"))
		return
	}

	tokenPair, err := handler.generateTokenPair(r.Context(), user)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, tokenPair)
}

func (handler *oAuthHandler) handleRefreshGrant(w http.ResponseWriter, r *http.Request, req tokenRequest) {
	type dbRecord struct {
		ID                    uuid.UUID `db:"id"`
		Username              string    `db:"username"`
		Email                 string    `db:"email"`
		FirstName             *string   `db:"firstname"`
		LastName              *string   `db:"lastname"`
		RefreshTokenExpiresAt time.Time `db:"refresh_token_expires_at"`
	}

	query := `
	SELECT
		u.id,
		u.username,
		u.email,
		u.firstname,
		u.lastname,
		at.refresh_token_expires_at
	FROM "user" u
	INNER JOIN auth_token at ON u.id = at.user_id
	WHERE at.refresh_token = $1 and u.deleted_at IS NULL`
	var rec dbRecord
	err := db.GetContext(r.Context(), &rec, query, req.RefreshToken)
	if errors.Is(err, sql.ErrNoRows) {
		httpx.WriteError(w, ErrUnauthorized)
		return
	}
	if err != nil {
		httpx.WriteError(w, ErrInternal.WithOrigin().WithCause(err).WithMessage("failed to query user"))
		return
	}

	if time.Now().After(rec.RefreshTokenExpiresAt) {
		db.ExecContext(r.Context(), "DELETE FROM auth_token WHERE access_token = $1", req.RefreshToken)
		httpx.WriteError(w, ErrTokenExpired.WithApiMessagef("Refresh token for user '%s' has expired at %s", rec.Username, rec.RefreshTokenExpiresAt))
		return
	}

	user := User{
		Id:       rec.ID,
		Username: rec.Username,
		Email:    rec.Email,
		Roles:    []Role{},
		Metadata: map[string]any{
			"firstname": rec.FirstName,
			"lastname":  rec.LastName,
		},
	}

	tokenPair, err := handler.generateTokenPair(r.Context(), &user)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, tokenPair)
}

func CleanupTokensCron() func(context.Context) error {
	return func(ctx context.Context) error {
		timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		return cleanupExpiredTokens(timeoutCtx)
	}
}

func cleanupExpiredTokens(ctx context.Context) error {
	_, err := db.ExecContext(ctx, `
		DELETE FROM auth_token
		WHERE refresh_token_expires_at < NOW()
	`)
	if err != nil {
		return ErrInternal.WithMessage("failed to cleanup expired tokens").WithOrigin().WithCause(err)
	}
	return nil
}
