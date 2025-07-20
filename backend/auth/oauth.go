package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/meh-hackathon/meh/db"
)

// Based on the grant type the client will send username + password or refresh token
type TokenRequest struct {
	GrantType    string `json:"grant_type"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type OAuthHandler struct {
	authenticators []Authenticator
}

var passwordHashCost = 12

type Authenticator interface {
	authenticate(ctx context.Context, username, password string) (*User, error)
}

type OauthOption func(*OAuthHandler) error

func NewOAuthHandler(opts ...OauthOption) (*OAuthHandler, error) {
	handler := &OAuthHandler{}
	for _, opt := range opts {
		if err := opt(handler); err != nil {
			return nil, fmt.Errorf("failed to apply OAuth option: %w", err)
		}
	}
	if len(handler.authenticators) == 0 {
		return nil, fmt.Errorf("auth: at least one authenticator must be set for OAuth handler")
	}

	return handler, nil
}

func (handler *OAuthHandler) HandlePasswordGrant(ctx context.Context, req TokenRequest) (TokenResponse, error) {
	// check all authenticators concurrently and check for the first successful authentication
	type result struct {
		user *User
		err  error
	}

	ctx, cancel := context.WithCancel(ctx)
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
			return TokenResponse{}, res.err
		}
		if res.err == nil && res.user != nil {
			// If we found a user, cancel the context to stop other goroutines
			cancel()
			user = res.user
			break
		}
	}
	if user == nil {
		return TokenResponse{}, ErrUnauthorized.WithApiMessage("invalid credentials")
	}

	tokenPair, err := handler.generateTokenPair(ctx, user)
	if err != nil {
		return TokenResponse{}, err
	}
	return tokenPair, nil
}

func (handler *OAuthHandler) HandleRefreshGrant(ctx context.Context, req TokenRequest) (TokenResponse, error) {
	type dbRecord struct {
		ID                    uuid.UUID `db:"id"`
		Username              string    `db:"username"`
		Email                 string    `db:"email"`
		RefreshTokenExpiresAt time.Time `db:"refresh_token_expires_at"`
	}

	query := `
	SELECT
		u.id,
		u.username,
		u.email,
		at.refresh_token_expires_at
	FROM "users" u
	INNER JOIN auth_tokens at ON u.id = at.user_id
	WHERE at.refresh_token = $1 and u.deleted_at IS NULL`
	var rec dbRecord
	err := db.GetContext(ctx, &rec, query, req.RefreshToken)
	if errors.Is(err, sql.ErrNoRows) {
		return TokenResponse{}, ErrInvalidToken.WithApiMessagef("Refresh token '%s' not found", req.RefreshToken)
	}
	if err != nil {
		return TokenResponse{}, ErrInternal.WithMessage("failed to get user by refresh token").WithOrigin().WithCause(err)
	}

	if time.Now().After(rec.RefreshTokenExpiresAt) {
		db.ExecContext(ctx, "DELETE FROM auth_tokens WHERE access_token = $1", req.RefreshToken)
		return TokenResponse{}, ErrTokenExpired.WithApiMessagef("Refresh token for user '%s' has expired at %s", rec.Username, rec.RefreshTokenExpiresAt)
	}

	user := User{
		Id:       rec.ID,
		Username: rec.Username,
		Email:    rec.Email,
		Roles:    []Role{},
		Metadata: map[string]any{},
	}

	tokenPair, err := handler.generateTokenPair(ctx, &user)
	if err != nil {
		return TokenResponse{}, err
	}
	return tokenPair, nil
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
		DELETE FROM auth_tokens
		WHERE refresh_token_expires_at < NOW()
	`)
	if err != nil {
		return ErrInternal.WithMessage("failed to cleanup expired tokens").WithOrigin().WithCause(err)
	}
	return nil
}
