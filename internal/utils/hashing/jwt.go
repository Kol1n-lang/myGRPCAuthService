package hashing

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	AlgorithmNotAllowed = errors.New("algorithm Not Allowed")
)

func CreateAccessRefreshTokens(id uuid.UUID, accessMin int, refreshDays int, secretKey string, algorithm string) (map[string]string, error) {
	accessJWTClaims := jwt.MapClaims{
		"sub":  id,
		"exp":  time.Now().Add(time.Minute * time.Duration(accessMin)).Unix(),
		"iat":  time.Now().Unix(),
		"type": "access",
	}

	refreshJWTClaims := jwt.MapClaims{
		"sub":  id,
		"exp":  time.Now().Add(time.Hour * 24 * time.Duration(refreshDays)).Unix(),
		"iat":  time.Now().Unix(),
		"type": "refresh",
	}

	if algorithm == "HS256" {
		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessJWTClaims)
		accessString, err := accessToken.SignedString([]byte(secretKey))
		if err != nil {
			return nil, err
		}

		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshJWTClaims)
		refreshString, err := refreshToken.SignedString([]byte(secretKey))
		if err != nil {
			return nil, err
		}

		return map[string]string{
			"access":  accessString,
			"refresh": refreshString,
		}, nil
	}

	return nil, AlgorithmNotAllowed
}
