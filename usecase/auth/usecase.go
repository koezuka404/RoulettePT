package auth

import "context"

type SignUpInput struct {
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthOutput struct {
	AccessToken string
}

type Usecase interface {
	SignUp(ctx context.Context, in SignUpInput) error
	Login(ctx context.Context, in LoginInput) (AuthOutput, error)
	LogoutAll(ctx context.Context, userID int64) error
}
