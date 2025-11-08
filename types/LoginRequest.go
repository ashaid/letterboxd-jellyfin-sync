package types

type LoginRequest struct {
	Username string `json:"Username"`
	Pw       string `json:"Pw"`
}
