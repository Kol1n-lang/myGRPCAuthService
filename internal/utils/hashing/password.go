package hashing

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordTooLong = errors.New("password too long")
	ErrInvalidHash     = errors.New("invalid hash format")
	ErrInvalidPassword = errors.New("invalid password")
)

func HashPassword(password string) ([]byte, error) {
	if len(password) > 64 {
		return nil, ErrPasswordTooLong
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	return bytes, nil
}

func VerifyPassword(password string, hash []byte) error {
	if len(hash) == 0 {
		return ErrInvalidHash
	}

	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return ErrInvalidPassword
		case errors.Is(err, bcrypt.ErrHashTooShort):
			return ErrInvalidHash
		default:
			return fmt.Errorf("password verification failed: %w", err)
		}
	}

	return nil
}
