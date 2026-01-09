package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ekosachev/go-backend-template/internal/models"
	"github.com/ekosachev/go-backend-template/internal/repository"
	"github.com/ekosachev/go-backend-template/pkg/password"
	"github.com/ekosachev/go-backend-template/pkg/timezone"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct {
	usersRepo repository.Repository[models.User]
	jwtSecret string
}

func NewAuthService(repo repository.Repository[models.User], jwtSecret string) *AuthService {
	return &AuthService{
		usersRepo: repo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(ctx context.Context, email, plainPassword string) (*models.User, error) {
	hash, err := password.Hash(plainPassword)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Email:        email,
		PasswordHash: hash,
	}
	if err := s.usersRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, plainPassword string) (string, *models.User, error) {
	user, err := s.usersRepo.FindOne(ctx, map[string]any{"email": email})
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}
	if err := password.Compare(user.PasswordHash, plainPassword); err != nil {
		return "", nil, ErrInvalidCredentials
	}
	token, err := s.generateJWT(user)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func (s *AuthService) generateJWT(user *models.User) (string, error) {
	// Standard claims: subject = user ID, expires in 24h
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"email": user.Email,
		"iat":   time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) UpdateTimezone(ctx context.Context, userID uuid.UUID, timezoneStr string) error {
	if !timezone.IsValidTimezone(timezoneStr) {
		return fmt.Errorf("invalid timezone: %s", timezoneStr)
	}

	user, err := s.usersRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Timezone = timezoneStr
	return s.usersRepo.Update(ctx, user)
}

func (s *AuthService) GetSupportedTimezones(ctx context.Context) []string {
	return timezone.SupportedTimezones
}
