package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// LocalProvider implements password-based authentication.
type LocalProvider struct {
	store *MySQLStore
}

func NewLocalProvider(store *MySQLStore) *LocalProvider {
	return &LocalProvider{store: store}
}

func (p *LocalProvider) Name() string { return "local" }

func (p *LocalProvider) Authenticate(ctx context.Context, req map[string]any) (*User, error) {
	username, _ := req["username"].(string)
	password, _ := req["password"].(string)

	if username == "" || password == "" {
		return nil, ErrInvalidCredentials
	}

	user, hash, err := p.store.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
