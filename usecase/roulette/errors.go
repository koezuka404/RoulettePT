package roulette

import "errors"

var (
	ErrInvalidKey      = errors.New("invalid idempotency key")
	ErrInvalidPage     = errors.New("invalid page")
	ErrInvalidLimit    = errors.New("invalid limit")
	ErrBalanceOverflow = errors.New("balance overflow")
)
