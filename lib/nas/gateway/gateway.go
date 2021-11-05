package gateway

type Gateway interface {
	Authorize(redirect string, state string) (string, error)
	AccessToken(callbackData map[string]string, redirect string) (*AccessToken, error)
	FreshToken(refreshToken string) (*AccessToken, error)
	User(accessToken string) (*User, error)
}
