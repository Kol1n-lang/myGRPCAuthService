package http

import (
	"context"
	"errors"

	"authService/github.com/authService/api"
	"authService/internal/config"
	"authService/internal/domain/value_objects"
	"authService/internal/service"
	"authService/internal/utils/service_errors"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var validate = validator.New()

type GRPCServer struct {
	api.UnimplementedAuthServiceServer
	service service.UserService
	cfg     *config.Config
}

func NewGRPCServer(userService service.UserService, cfg *config.Config) *GRPCServer {
	return &GRPCServer{
		service: userService,
		cfg:     cfg,
	}
}

func (s *GRPCServer) Register(ctx context.Context, req *api.AuthRequest) (*api.RegisterResponse, error) {
	userRegistry := value_objects.UserVO{
		Email:    req.Email,
		Password: req.Password,
	}
	err := validate.Struct(&userRegistry)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err = s.service.Register(ctx, &userRegistry); err != nil {
		if errors.Is(err, service_errors.InternalServerError) {
			return nil, status.Error(codes.Internal, err.Error())
		}
		return nil, status.Error(codes.AlreadyExists, err.Error())
	}

	return &api.RegisterResponse{
		Success: true,
	}, nil
}

func (s *GRPCServer) Login(ctx context.Context, req *api.AuthRequest) (*api.AuthResponse, error) {
	userLogin := value_objects.UserVO{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := validate.Struct(&userLogin); err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid email or password format")
	}

	tokens, err := s.service.Login(ctx, s.cfg, &userLogin)
	if err != nil {
		switch {
		case errors.Is(err, service_errors.InvalidCredentialsError):
			return nil, status.Error(codes.Unauthenticated, "Invalid credentials or not active user")
		case errors.Is(err, service_errors.InternalServerError):
			return nil, status.Error(codes.Internal, "Internal server error")
		default:
			return nil, status.Error(codes.Internal, "Internal server error")
		}
	}

	return &api.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *GRPCServer) RefreshTokens(ctx context.Context, req *api.RefreshToken) (*api.AuthResponse, error) {
	return &api.AuthResponse{
		AccessToken:  "test_token",
		RefreshToken: "tes_token",
	}, nil

}
