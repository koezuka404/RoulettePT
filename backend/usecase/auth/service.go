package auth

import (
	"context"
	"os"
	"strings"
	"time"

	user "roulettept/domain/user/model"
	userrepo "roulettept/domain/user/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	users userrepo.UserRepository
	rt    userrepo.RefreshTokenRepository

	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewService(users userrepo.UserRepository, rt userrepo.RefreshTokenRepository) *Service {
	return &Service{
		users:      users,
		rt:         rt,
		accessTTL:  15 * time.Minute,
		refreshTTL: 14 * 24 * time.Hour,
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

	u := &user.User{
		Email:        in.Email,
		PasswordHash: string(hash),
		Role:         user.RoleUser,
		TokenVersion: 0,
		PointBalance: 0,
		IsActive:     true,
	}
	return s.users.Create(ctx, u)
}

func (s *Service) Login(ctx context.Context, in LoginInput) (LoginOutput, error) {
	u, err := s.users.FindByEmail(ctx, in.Email)
	if err != nil {
		return LoginOutput{}, err
	}
	if u == nil {
		return LoginOutput{}, ErrInvalidCredentials
	}
	if !u.IsActive {
		return LoginOutput{}, ErrUserInactive
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)); err != nil {
		return LoginOutput{}, ErrInvalidCredentials
	}

	access, err := s.issueAccessToken(u)
	if err != nil {
		return LoginOutput{}, err
	}

	// refresh token（平文cookie、DBにはhashのみ）
	refreshPlain, err := randomTokenHex(32)
	if err != nil {
		return LoginOutput{}, err
	}
	refreshHash := sha256Hex(refreshPlain)

	now := time.Now()
	if err := s.rt.Create(ctx, &user.RefreshToken{
		UserID:    u.ID,
		TokenHash: refreshHash,
		ExpiresAt: now.Add(s.refreshTTL),
		UsedAt:    nil,
		UserAgent: in.UserAgent,
		IP:        in.IP,
		CreatedAt: now,
	}); err != nil {
		return LoginOutput{}, err
	}

	csrf, err := randomTokenHex(32)
	if err != nil {
		return LoginOutput{}, err
	}

	return LoginOutput{
		AccessToken:  access,
		RefreshToken: refreshPlain,
		CSRFToken:    csrf,
	}, nil
}

func (s *Service) Refresh(ctx context.Context, in RefreshInput) (RefreshOutput, error) {
	if strings.TrimSpace(in.RefreshToken) == "" {
		return RefreshOutput{}, ErrRefreshTokenInvalid
	}

	hash := sha256Hex(in.RefreshToken)

	rt, err := s.rt.FindByHash(ctx, hash)
	if err != nil || rt == nil {
		return RefreshOutput{}, ErrRefreshTokenInvalid
	}

	now := time.Now()
	if now.After(rt.ExpiresAt) {
		_ = s.rt.DeleteByHash(ctx, hash)
		return RefreshOutput{}, ErrRefreshTokenInvalid
	}

	if rt.UsedAt != nil {
		_ = s.rt.DeleteByUserID(ctx, rt.UserID)
		return RefreshOutput{}, ErrRefreshTokenReused
	}

	u, err := s.users.FindByID(ctx, rt.UserID)
	if err != nil || u == nil {
		return RefreshOutput{}, ErrRefreshTokenInvalid
	}
	if !u.IsActive {
		return RefreshOutput{}, ErrUserInactive
	}

	tv, ok := readTVIgnoringExp(in.AccessToken)
	if !ok || tv != u.TokenVersion {
		return RefreshOutput{}, ErrTokenVersionMismatch
	}

	// old を used に（二重実行はここで弾く）
	if err := s.rt.MarkUsed(ctx, rt.ID, now); err != nil {
		_ = s.rt.DeleteByUserID(ctx, rt.UserID)
		return RefreshOutput{}, ErrRefreshTokenReused
	}

	newPlain, err := randomTokenHex(32)
	if err != nil {
		return RefreshOutput{}, err
	}
	newHash := sha256Hex(newPlain)

	if err := s.rt.Create(ctx, &user.RefreshToken{
		UserID:    rt.UserID,
		TokenHash: newHash,
		ExpiresAt: now.Add(s.refreshTTL),
		UsedAt:    nil,
		UserAgent: in.UserAgent,
		IP:        in.IP,
		CreatedAt: now,
	}); err != nil {
		return RefreshOutput{}, err
	}

	access, err := s.issueAccessToken(u)
	if err != nil {
		return RefreshOutput{}, err
	}
	csrf, err := randomTokenHex(32)
	if err != nil {
		return RefreshOutput{}, err
	}

	return RefreshOutput{
		AccessToken:  access,
		RefreshToken: newPlain,
		CSRFToken:    csrf,
	}, nil
}

func (s *Service) Logout(ctx context.Context, in LogoutInput) error {
	if in.UserID == 0 {
		return ErrUserNotFound
	}
	if strings.TrimSpace(in.RefreshToken) == "" {
		return nil
	}
	hash := sha256Hex(in.RefreshToken)
	return s.rt.DeleteByHash(ctx, hash)
}

func (s *Service) LogoutAll(ctx context.Context, userID int64) error {
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if u == nil {
		return ErrUserNotFound
	}

	return s.users.IncrementTokenVersion(ctx, userID)
}

func (s *Service) issueAccessToken(u *user.User) (string, error) {
	secret := os.Getenv("SECRET")
	if secret == "" {
		return "", ErrSecretNotSet
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub":  itoa(u.ID),
		"role": string(u.Role),
		"tv":   u.TokenVersion,
		"iat":  now.Unix(),
		"exp":  now.Add(s.accessTTL).Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}
