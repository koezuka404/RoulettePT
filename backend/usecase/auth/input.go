package auth

type SignUpInput struct {
	Email    string
	Password string
}

type LoginInput struct {
	Email     string
	Password  string
	UserAgent string
	IP        string
}

type RefreshInput struct {
	RefreshToken string
	AccessToken  string
	UserAgent    string
	IP           string
}

type LogoutInput struct {
	UserID       int64
	RefreshToken string
}
