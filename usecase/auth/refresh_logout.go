package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"roulettept/domain/repository"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrRefreshTokenInvalid  = errors.New("refresh token invalid")
	ErrRefreshTokenReused   = errors.New("refresh token reused")
	ErrTokenVersionMismatch = errors.New("token version mismatch")
)

// Refresh/Logout に必要なので Service に repo を追加して使う想定
// （DIで注入してね）
type RefreshLogoutDeps struct {
	Users      repository.UserRepository
	RT         repository.RefreshTokenRepository
	AccessTTL  time.Duration // 15分推奨（仕様：短命） :contentReference[oaicite:4]{index=4}
	RefreshTTL time.Duration // 14日 :contentReference[oaicite:5]{index=5}
}

// RefreshAccessToken: refresh cookie を検証し、ローテーションして新AccessTokenを返す
// - CSRF は middleware で弾く前提（Refresh/Logout必須） :contentReference[oaicite:6]{index=6}
func RefreshAccessToken(ctx context.Context, d RefreshLogoutDeps, refreshPlain string, currentAccessToken string, userAgent string, ip string) (newAccess string, newRefreshPlain string, err error) {
	if d.Users == nil || d.RT == nil {
		return "", "", errors.New("deps not set")
	}
	if d.AccessTTL == 0 {
		d.AccessTTL = 15 * time.Minute
	}
	if d.RefreshTTL == 0 {
		d.RefreshTTL = 14 * 24 * time.Hour
	}
	if refreshPlain == "" {
		return "", "", ErrRefreshTokenInvalid
	}

	// refresh cookie(平文) -> sha256 hash
	hash := sha256Hex(refreshPlain)

	// DB lookup
	rt, err := d.RT.FindByHash(ctx, hash)
	if err != nil {
		return "", "", ErrRefreshTokenInvalid
	}
	// 期限切れは無効 :contentReference[oaicite:7]{index=7}
	if time.Now().After(rt.ExpiresAt) {
		_ = d.RT.DeleteByHash(ctx, hash)
		return "", "", ErrRefreshTokenInvalid
	}

	// 再利用検知：used_at != NULL なら全削除して401 :contentReference[oaicite:8]{index=8}
	if rt.UsedAt != nil {
		_ = d.RT.DeleteByUserID(ctx, rt.UserID)
		return "", "", ErrRefreshTokenReused
	}

	// token_version 照合：tv ≠ users.token_version なら 401 :contentReference[oaicite:9]{index=9}
	u, err := d.Users.FindByID(ctx, rt.UserID)
	if err != nil || u == nil {
		return "", "", ErrRefreshTokenInvalid
	}
	tv, ok := readTVIgnoringExp(currentAccessToken)
	if !ok || tv != u.TokenVersion {
		return "", "", ErrTokenVersionMismatch
	}

	// 先に old refresh を used にする（競合時は RowsAffected==0 で NotFound になる想定）
	now := time.Now()
	if err := d.RT.MarkUsed(ctx, rt.ID, now); err != nil {
		// ここに来たら同時更新/再利用の可能性 → 全削除して 401 REUSED :contentReference[oaicite:10]{index=10}
		_ = d.RT.DeleteByUserID(ctx, rt.UserID)
		return "", "", ErrRefreshTokenReused
	}

	// 新 refresh を発行（ローテーション）
	newRefreshPlain, err = randomTokenHex(32)
	if err != nil {
		return "", "", err
	}
	newHash := sha256Hex(newRefreshPlain)

	// DB保存（hashのみ。平文保存禁止） :contentReference[oaicite:11]{index=11}
	rt2 := *rt
	rt2.ID = 0
	rt2.TokenHash = newHash
	rt2.UsedAt = nil
	rt2.ExpiresAt = now.Add(d.RefreshTTL)
	rt2.UserAgent = userAgent
	rt2.IP = ip
	rt2.CreatedAt = now

	if err := d.RT.Create(ctx, &rt2); err != nil {
		return "", "", err
	}

	// 新 access token 発行（tv を入れる） :contentReference[oaicite:12]{index=12}
	newAccess, err = issueAccessToken(u.ID, string(u.Role), u.TokenVersion, d.AccessTTL)
	if err != nil {
		return "", "", err
	}

	return newAccess, newRefreshPlain, nil
}

// Logout: refresh cookie を削除（DB & cookieはhandlerでMaxAge=-1） :contentReference[oaicite:13]{index=13}
func Logout(ctx context.Context, d RefreshLogoutDeps, userID int64, refreshPlain string) error {
	if d.Users == nil || d.RT == nil {
		return errors.New("deps not set")
	}
	// （JWT検証やuserID取得は middleware/handler 側の責務）
	if refreshPlain == "" {
		// cookie無しでも成功扱いでOK（仕様上は401の規定はないので安全側で） :contentReference[oaicite:14]{index=14}
		return nil
	}
	hash := sha256Hex(refreshPlain)

	// そのrefreshがそのユーザーのものか軽く確認したいなら FindByHash→UserID比較でもOK。
	// 最短シンプル：hashで削除
	return d.RT.DeleteByHash(ctx, hash)
}

func issueAccessToken(userID int64, role string, tokenVersion int64, ttl time.Duration) (string, error) {
	secret := os.Getenv("SECRET")
	if secret == "" {
		return "", ErrSecretNotSet
	}
	now := time.Now()

	// 仕様：sub, role, tv, iat, exp :contentReference[oaicite:15]{index=15}
	claims := jwt.MapClaims{
		"sub":  itoa(userID),
		"role": role,
		"tv":   tokenVersion,
		"iat":  now.Unix(),
		"exp":  now.Add(ttl).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}

func readTVIgnoringExp(accessToken string) (int64, bool) {
	if accessToken == "" {
		return 0, false
	}
	secret := os.Getenv("SECRET")
	if secret == "" {
		return 0, false
	}

	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	claims := jwt.MapClaims{}
	_, err := parser.ParseWithClaims(accessToken, claims, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return 0, false
	}
	v, ok := claims["tv"]
	if !ok {
		return 0, false
	}
	switch x := v.(type) {
	case float64:
		return int64(x), true
	case int64:
		return x, true
	case int:
		return int64(x), true
	default:
		return 0, false
	}
}

func randomTokenHex(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func sha256Hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

func itoa(n int64) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + (n % 10))
		n /= 10
	}
	return string(b[i:])
}
