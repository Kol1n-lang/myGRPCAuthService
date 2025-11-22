package service

import (
	"context"
	"log"

	"authService/internal/domain/repositories"
	"authService/internal/domain/value_objects"
	"authService/internal/utils/hashing"
	"authService/internal/utils/service_errors"
)

type UserService interface {
	Register(ctx context.Context, userRegistry *value_objects.UserRegistryVO) error
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

func (u *UserServiceImpl) Register(ctx context.Context, userRegistry *value_objects.UserRegistryVO) error {
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
