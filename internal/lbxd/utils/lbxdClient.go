package utils

import "net/http"

type LbxdClient struct {
	BaseURL       string
	Username      string
	Password      string
	Client_Id     string
	Client_Secret string
	Client        *http.Client
}

func NewLbxdClient(baseURL, username, password, clientID, clientSecret string) *LbxdClient {
	return &LbxdClient{
		BaseURL:       baseURL,
		Username:      username,
		Password:      password,
		Client_Id:     clientID,
		Client_Secret: clientSecret,
		Client:        &http.Client{},
	}
}
