package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"roulettept/domain/models"
	"roulettept/domain/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyUsed   = errors.New("email already used")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserInactive       = errors.New("user is inactive")
	ErrSecretNotSet       = errors.New("SECRET is not set")
)

type Service struct {
	users repository.UserRepository
	ttl   time.Duration
}

func NewService(users repository.UserRepository) *Service {
	return &Service{
		users: users,
		ttl:   12 * time.Hour,
	}
}

func (s *Service) SignUp(ctx context.Context, in SignUpInput) error {
	exists, err := s.users.FindByEmail(ctx, in.Email)
	if err != nil {
		return err
	}
	if exists != nil {
		return ErrEmailAlreadyUsed
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u := &models.User{
		Email:        in.Email,
		PasswordHash: string(hash),
		Role:         models.RoleUser,
		TokenVersion: 0,
		PointBalance: 0,
		IsActive:     true,
	}

	return s.users.Create(ctx, u)
}

func (s *Service) Login(ctx context.Context, in LoginInput) (AuthOutput, error) {
	u, err := s.users.FindByEmail(ctx, in.Email)
	if err != nil {
		return AuthOutput{}, err
	}
	if u == nil {
		return AuthOutput{}, ErrInvalidCredentials
	}
	if !u.IsActive {
		return AuthOutput{}, ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)); err != nil {
		return AuthOutput{}, ErrInvalidCredentials
	}

	secret := os.Getenv("SECRET")
	if secret == "" {
		return AuthOutput{}, ErrSecretNotSet
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"user_id":       u.ID,
		"token_version": u.TokenVersion,
		"exp":           now.Add(s.ttl).Unix(),
		"iat":           now.Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString([]byte(secret))
	if err != nil {
		return AuthOutput{}, err
	}

	return AuthOutput{AccessToken: token}, nil
}

func (s *Service) LogoutAll(ctx context.Context, userID int64) error {
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if u == nil {
		return ErrUserNotFound
	}

	_, rows, err := s.users.IncrementTokenVersion(ctx, userID)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}
