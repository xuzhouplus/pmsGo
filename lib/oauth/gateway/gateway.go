package gateway

type Gateway interface {
	Scope() string
	GrantType() string
	AuthorizeUrl(scope string, redirect string, state string) string
	AccessToken(code string, redirect string, state string) (string, error)
	User(accessToken string) (map[string]string, error)
}
