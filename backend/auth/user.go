package auth

import (
	"crypto/sha256"
	"net/http"

	"slices"

	"github.com/google/uuid"
	"github.com/meh-hackathon/meh/httpx"
	"golang.org/x/crypto/bcrypt"
)

type Role string

const (
	RoleAdmin Role = "admin"
)

type User struct {
	ID       uuid.UUID      `json:"id"`
	Username string         `json:"username"`
	Email    string         `json:"email,omitempty"`
	Roles    []Role         `json:"roles,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

func (u *User) HasRole(role Role, roles ...Role) bool {
	roles = append(roles, role)
	for _, r := range u.Roles {
		if slices.Contains(roles, r) {
			return true
		}
	}
	return false
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	u, err := GetUser(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, u)
}

func HashPassword(password string) (string, error) {
	shaHash := sha256.Sum256([]byte(password))
	bytes, err := bcrypt.GenerateFromPassword(shaHash[:], passwordHashCost)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	shaHash := sha256.Sum256([]byte(password))
	err := bcrypt.CompareHashAndPassword([]byte(hash), shaHash[:])
	return err == nil
}
