package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"authService/internal/domain/repositories"
	"github.com/Masterminds/squirrel"
)

var Psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepositoryImpl(db *sql.DB) repositories.UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (r *UserRepositoryImpl) InsertUser(ctx context.Context, email string, hashedPassword []byte) error {
	query, args, err := Psql.
		Insert("users").
		Columns("email", "password", "is_active").
		Values(email, hashedPassword, true).
		ToSql()

	if err != nil {
		log.Printf("Failed to build insert user query: %v", err)
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) CheckUserExist(ctx context.Context, email string) (bool, error) {
	query, args, err := Psql.
		Select("1").
		From("users").
		Where(squirrel.Eq{"email": email}).
		ToSql()

	if err != nil {
		log.Printf("Failed to build check user query: %v", err)
		return false, err
	}

	var exists int
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&exists)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		log.Printf("Failed to check user existence: %v", err)
		return false, err
	}

	return true, nil
}
