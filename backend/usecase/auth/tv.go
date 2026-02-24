package auth

import (
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// expは無視してtvだけ読みたい（refresh用）
func readTVIgnoringExp(accessToken string) (int64, bool) {
	accessToken = strings.TrimSpace(accessToken)
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
