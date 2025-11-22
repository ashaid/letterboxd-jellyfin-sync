package types

type AuthCredentials struct {
	Cookies struct {
		CSRF       string
		GA         string
		GAD3       string
		SignedInAs string
	}
	RefreshToken string
	ClientSecret string
	ClientID     string
}
