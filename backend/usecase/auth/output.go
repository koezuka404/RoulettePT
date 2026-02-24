package auth

type LoginOutput struct {
	AccessToken  string
	RefreshToken string
	CSRFToken    string
}

type RefreshOutput struct {
	AccessToken  string
	RefreshToken string
	CSRFToken    string
}
