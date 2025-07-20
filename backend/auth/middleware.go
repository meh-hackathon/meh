package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/meh-hackathon/meh/apperror"
	"github.com/meh-hackathon/meh/httpx"
)

const userContextKey = "meh:userInfo"

var (
	ErrUnauthorized = apperror.Define("auth:unauthorized", "Unauthorized access").WithStatus(http.StatusUnauthorized)
	ErrForbidden    = apperror.Define("auth:forbidden", "Forbidden access").WithStatus(http.StatusForbidden)
	ErrInternal     = apperror.Define("auth:internal", "Internal authentication error")
	ErrInvalidToken = apperror.Define("auth:invalid_token", "Invalid authentication token").WithStatus(http.StatusUnauthorized)
	ErrTokenExpired = apperror.Define("auth:token_expired", "Authentication token has expired").WithStatus(http.StatusUnauthorized)
)

func Middleware(secret []byte) func(http.Handler) http.Handler {
	if len(secret) == 0 {
		panic("auth: secret key must be set before using the middleware")
	}

	secretKey = secret

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			tokenString, err := extractBearerToken(authHeader)
			if err != nil {
				httpx.WriteError(w, err)
				return
			}

			user, err := parseAccessToken(tokenString)
			if err != nil {
				httpx.WriteError(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func extractBearerToken(authHeader string) (string, error) {
	authHeader = strings.TrimSpace(authHeader)
	if !strings.HasPrefix(authHeader, tokenType+" ") {
		return "", ErrInvalidToken.WithApiMessagef("Authorization header must start with '%s '", tokenType)
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, tokenType+" "))
	if token == "" {
		return "", ErrInvalidToken.WithApiMessagef("%s token cannot be empty", tokenType)
	}

	return token, nil
}
