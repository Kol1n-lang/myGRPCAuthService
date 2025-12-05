package service

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"authService/internal/config"
	"authService/internal/domain/repositories"
	"authService/internal/domain/value_objects"
	"authService/internal/utils/hashing"
	"authService/internal/utils/service_errors"
)

type UserService interface {
	Register(ctx context.Context, userRegistry *value_objects.UserVO) error
	Login(ctx context.Context, cfg *config.Config, userLogin *value_objects.UserVO) (value_objects.AuthResponse, error)
}

func NewUserService(repository repositories.UserRepository, brokerRepo repositories.RabbitRepository) UserService {
	return &UserServiceImpl{
		userRepo:   repository,
		brokerRepo: brokerRepo,
	}
}

type UserServiceImpl struct {
	userRepo   repositories.UserRepository
	brokerRepo repositories.RabbitRepository
}

func (u *UserServiceImpl) Register(ctx context.Context, userRegistry *value_objects.UserVO) error {
	exists, err := u.userRepo.CheckUserExist(ctx, userRegistry.Email)
	if err != nil {
		log.Printf("Error checking user existence: %v", err)
		return service_errors.InternalServerError
	}

	if exists {
		return service_errors.UserAlreadyExistsError
	}

	hashedPassword, err := hashing.HashPassword(userRegistry.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return service_errors.InternalServerError
	}

	if err = u.userRepo.InsertUser(ctx, userRegistry.Email, hashedPassword); err != nil {
		log.Printf("Error inserting user: %v", err)
		return service_errors.InternalServerError
	}

	go func() {
		if err = u.brokerRepo.CreateEmailMSG(userRegistry.Email); err != nil {
			log.Printf("Error creating email message: %v", err)
		}
	}()

	log.Printf("User registered successfully: %s", userRegistry.Email)
	return nil
}

func (u *UserServiceImpl) Login(ctx context.Context, cfg *config.Config, userLogin *value_objects.UserVO) (value_objects.AuthResponse, error) {
	userID, hashedPWD, err := u.userRepo.GetUserCredentials(ctx, userLogin.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return value_objects.AuthResponse{}, service_errors.InvalidCredentialsError
		}
		log.Printf("Error getting user credentials: %v", err)
		return value_objects.AuthResponse{}, service_errors.InternalServerError
	}

	if err := hashing.VerifyPassword(userLogin.Password, hashedPWD); err != nil {
		if errors.Is(err, hashing.ErrInvalidPassword) {
			return value_objects.AuthResponse{}, service_errors.InvalidCredentialsError
		}
		log.Printf("Password verification error: %v", err)
		return value_objects.AuthResponse{}, service_errors.InternalServerError
	}

	tokens, err := hashing.CreateAccessRefreshTokens(
		userID,
		cfg.JWT.AccessExpireMinutes,
		cfg.JWT.RefreshExpireDays,
		cfg.JWT.JWTSecret,
		cfg.JWT.Algorithm,
	)
	if err != nil {
		log.Printf("Token generation error: %v", err)
		return value_objects.AuthResponse{}, service_errors.InternalServerError
	}

	return value_objects.AuthResponse{
		AccessToken:  tokens["access"],
		RefreshToken: tokens["refresh"],
	}, nil
}
