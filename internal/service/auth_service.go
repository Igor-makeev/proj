package service

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"proj/config"
	"proj/internal/entities/models"
	"proj/internal/repository"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

type tokenClaims struct {
	jwt.StandardClaims

	UserID int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
	cfg  config.AuthConfig
}

func NewAuthService(repo repository.Authorization, cfg *config.Config) *AuthService {
	return &AuthService{
		repo: repo,
		cfg:  cfg.AuthConfig,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, user models.User) error {
	user.Password = s.generatePasswordHash(user.Password)
	return s.repo.CreateUser(ctx, user)
}

func (s *AuthService) GenerateToken(ctx context.Context, login, password string) (string, error) {
	userid, err := s.repo.GetUser(ctx, login, s.generatePasswordHash(password))
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(6 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userid,
	})
	return token.SignedString([]byte(s.cfg.SigningKey))
}

func (s *AuthService) ParseToken(accesstoken string) (int, error) {
	tocken, err := jwt.ParseWithClaims(accesstoken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.cfg.SigningKey), nil
	})

	if err != nil {
		logrus.Print(err)
		return 0, nil
	}
	claims, ok := tocken.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tockenClaims")
	}

	return claims.UserID, nil
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(s.cfg.Salt)))
}
