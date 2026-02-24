package auth

import "errors"

var (
	ErrEmailAlreadyUsed   = errors.New("email already used")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserInactive       = errors.New("user is inactive")
	ErrSecretNotSet       = errors.New("SECRET is not set")

	ErrRefreshTokenInvalid  = errors.New("refresh token invalid")
	ErrRefreshTokenReused   = errors.New("refresh token reused")
	ErrTokenVersionMismatch = errors.New("token version mismatch")
)
