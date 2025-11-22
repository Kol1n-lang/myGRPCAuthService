package repositories

import "context"

type UserRepository interface {
	InsertUser(ctx context.Context, email string, hashedPassword []byte) error
	CheckUserExist(ctx context.Context, email string) (bool, error)
}
