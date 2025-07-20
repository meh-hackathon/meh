package auth

import (
	"context"
	"slices"

	"github.com/google/uuid"
)

// GetUser retrieves the authenticated user from the context.
// If the user is not authenticated, it returns an error.
func GetUser(ctx context.Context) (*User, error) {
	unparsedUser := ctx.Value(userContextKey)
	if unparsedUser == nil {
		return nil, ErrUnauthorized
	}
	user, ok := unparsedUser.(*User)
	if !ok {
		return nil, ErrInternal.WithMessagef("Invalid user data in context: %T", unparsedUser).WithOrigin()
	}
	return user, nil
}

func GetUserByAccessToken(accessToken string) (*User, error) {
	return parseAccessToken(accessToken)
}

func MustBeAuthenticated(ctx context.Context) error {
	_, err := GetUser(ctx)
	return err
}

// MustHaveRole checks if the authenticated user has at least one of the specified roles.
func MustHaveRole(ctx context.Context, role ...Role) error {
	user, err := GetUser(ctx)
	if err != nil {
		return err
	}
	for _, r := range user.Roles {
		if slices.Contains(role, r) {
			return nil
		}
	}

	return ErrForbidden.WithMessagef("User '%s' does not have required role(s): %v", user.Username, role).WithOrigin()
}

// MustBeUser checks if the authenticated user matches the specified user ID.
func MustBeUser(ctx context.Context, id uuid.UUID) error {
	user, err := GetUser(ctx)
	if err != nil {
		return err
	}
	if user.Id != id {
		return ErrForbidden.WithMessagef("User '%s' does not match required user ID '%s'", user.Id, id).WithOrigin()
	}
	return nil
}
