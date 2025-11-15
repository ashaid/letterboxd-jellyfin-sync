package types

type UnscopedTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type TokenResponse struct {
	UnscopedTokenResponse
	RefreshToken string `json:"refresh_token"`
}
