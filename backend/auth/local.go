package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/meh-hackathon/meh/db"
)

var _ Authenticator = (*LocalAuthenticator)(nil)

type LocalAuthenticator struct{}

func NewLocalAuthenticator() *LocalAuthenticator { return &LocalAuthenticator{} }

var WithLocalAuthenticator OauthOption = func(o *OAuthHandler) error {
	o.authenticators = append(o.authenticators, &LocalAuthenticator{})
	return nil
}

func (a *LocalAuthenticator) authenticate(ctx context.Context, username, password string) (*User, error) {
	type dbRecord struct {
		ID           uuid.UUID `db:"id"`
		Username     string    `db:"username"`
		Email        string    `db:"email"`
		FirstName    *string   `db:"firstname"`
		LastName     *string   `db:"lastname"`
		PasswordHash string    `db:"password_hash"`
	}

	var rec dbRecord
	query := `SELECT u.id, u.username, u.email, u.firstname, u.lastname, p.password_hash
	FROM "user" u INNER JOIN auth_password p ON u.id = p.user_id
	WHERE username = $1 AND deleted_at IS NULL;`
	err := db.GetContext(ctx, &rec, query, username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUnauthorized
	}
	if err != nil {
		return nil, ErrInternal.WithOrigin().WithCause(err).WithMessage("failed to query user")
	}

	if !CheckPasswordHash(password, rec.PasswordHash) {
		return nil, ErrUnauthorized
	}

	user := User{
		ID:       rec.ID,
		Username: rec.Username,
		Email:    rec.Email,
		Roles:    []Role{},
		Metadata: map[string]any{
			"firstname": rec.FirstName,
			"lastname":  rec.LastName,
		},
	}

	return &user, nil
}
