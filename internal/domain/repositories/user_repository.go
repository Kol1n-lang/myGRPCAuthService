package repositories

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	InsertUser(ctx context.Context, email string, hashedPassword []byte) error
	CheckUserExist(ctx context.Context, email string) (bool, error)
	GetUserCredentials(ctx context.Context, email string) (uuid.UUID, []byte, error)
}
