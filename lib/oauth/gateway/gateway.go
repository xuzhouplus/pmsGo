package gateway

type Gateway interface {
	Scope() string
	GrantType() string
	AuthorizeUrl(scope string, redirect string, state string) (string, string, error)
	AccessToken(callbackData map[string]string, redirect string) (map[string]string, error)
	User(accessToken string) (map[string]string, error)
}
