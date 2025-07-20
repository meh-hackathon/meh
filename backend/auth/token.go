package auth

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/meh-hackathon/meh/constants"
	"github.com/meh-hackathon/meh/db"
)

var (
	secretKey            []byte
	tokenType            = "Bearer"
	accessTokenLifetime  = 15 * time.Minute   // Default access token lifetime
	refreshTokenLifetime = 7 * 24 * time.Hour // Default refresh token lifetime
)

type jwtClaims struct {
	UserId   uuid.UUID      `json:"user_id"`
	Username string         `json:"username"`
	Email    string         `json:"email,omitempty"`
	Roles    []Role         `json:"roles,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
	jwt.RegisteredClaims
}

type TokenOptions struct {
	ExpiresAt *time.Time
	IssuedAt  *time.Time
	NotBefore *time.Time
	Issuer    string
	Subject   string
	Audience  []string
	ID        string
}

func (handler *oAuthHandler) generateTokenPair(ctx context.Context, user *User) (tokenPair tokenResponse, err error) {
	refreshToken, err := generateRefreshToken(32)
	if err != nil {
		return tokenResponse{}, err
	}

	id, _ := uuid.NewV7()

	accessTokenExpiresAt := time.Now().Add(accessTokenLifetime)
	accessToken, err := generateAccessToken(user, &TokenOptions{ExpiresAt: &accessTokenExpiresAt, ID: id.String()})
	if err != nil {
		return tokenResponse{}, err
	}

	id, _ = uuid.NewV7()
	_, err = db.ExecContext(
		ctx,
		"INSERT INTO auth_token (id, user_id, access_token, refresh_token, access_token_expires_at, refresh_token_expires_at) VALUES ($1, $2, $3, $4, $5, $6)",
		id.String(),
		user.Id,
		accessToken,
		refreshToken,
		accessTokenExpiresAt,
		time.Now().Add(refreshTokenLifetime),
	)
	if err != nil {
		return tokenResponse{}, ErrInternal.WithMessage("failed to insert auth token").WithOrigin().WithCause(err)
	}

	_, err = db.ExecContext(
		ctx,
		`UPDATE "user" SET last_login_at = NOW() WHERE id = $1`,
		user.Id,
	)
	if err != nil {
		return tokenResponse{}, ErrInternal.WithMessage("failed to update user last login").WithOrigin().WithCause(err)
	}

	return tokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(accessTokenLifetime.Seconds()),
		TokenType:    tokenType,
	}, nil
}

func generateAccessToken(user *User, options *TokenOptions) (string, error) {
	if user == nil {
		return "", ErrInternal.WithMessage("User cannot be nil").WithOrigin()
	}
	if len(secretKey) == 0 {
		return "", ErrInternal.WithMessage("Secret key must be set before generating tokens").WithOrigin()
	}

	claims := &jwtClaims{
		UserId:   user.Id,
		Username: user.Username,
		Email:    user.Email,
		Roles:    user.Roles,
		Metadata: user.Metadata,
	}
	claims.Issuer = constants.SYSTEM_NAME
	claims.Audience = []string{constants.SYSTEM_NAME}
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(refreshTokenLifetime))

	if options != nil {
		if options.ExpiresAt != nil {
			claims.ExpiresAt = jwt.NewNumericDate(*options.ExpiresAt)
		}
		if options.IssuedAt != nil {
			claims.IssuedAt = jwt.NewNumericDate(*options.IssuedAt)
		}
		if options.NotBefore != nil {
			claims.NotBefore = jwt.NewNumericDate(*options.NotBefore)
		}
		if options.Issuer != "" {
			claims.Issuer = options.Issuer
		}
		if options.Subject != "" {
			claims.Subject = options.Subject
		}
		if len(options.Audience) > 0 {
			claims.Audience = options.Audience
		}
		if options.ID != "" {
			claims.ID = options.ID
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", ErrInternal.WithMessage("Failed to sign token").WithCause(err).WithOrigin()
	}
	return signedToken, nil
}

func parseAccessToken(tokenString string) (*User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken.WithMessage("Unexpected signing method")
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, ErrInvalidToken.WithMessage("Failed to parse token").WithCause(err).WithOrigin()
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken.WithMessage("Invalid token claims").WithOrigin()
	}

	user := &User{
		Id:       claims.UserId,
		Username: claims.Username,
		Email:    claims.Email,
		Roles:    claims.Roles,
		Metadata: claims.Metadata,
	}

	return user, nil
}

func generateRefreshToken(size int) (string, error) {
	random := make([]byte, size)
	if _, err := rand.Read(random); err != nil {
		return "", ErrInternal.WithMessage("failed to generate random token").WithOrigin().WithCause(err)
	}

	mac := hmac.New(sha256.New, secretKey)
	mac.Write(random)
	signature := mac.Sum(nil)

	token := append(random, signature...)
	return base64.RawURLEncoding.EncodeToString(token), nil
}
