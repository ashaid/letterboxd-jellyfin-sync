package utils

import (
	"net/http"
	"net/url"

	"github.com/ashaid/letterboxd-jellyfin-sync/internal/lbxd/types"
)

func SetAPIHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Letterboxd/6553 CFNetwork/3826.600.41 Darwin/24.6.0")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
}

func SetWebHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
}

func SetAuthCredentials(req *http.Request, creds *types.AuthCredentials) {
	req.AddCookie(&http.Cookie{Name: "com.xk72.webparts.csrf", Value: creds.Cookies.CSRF})
	req.AddCookie(&http.Cookie{Name: "_ga", Value: creds.Cookies.GA})
	req.AddCookie(&http.Cookie{Name: "_ga_D3**", Value: creds.Cookies.GAD3})
	req.AddCookie(&http.Cookie{Name: "letterboxd.signed.in.as", Value: creds.Cookies.SignedInAs})
}

func SetAuthFormData(formData url.Values, creds *types.AuthCredentials) {
	formData.Set("refresh_token", creds.RefreshToken)
	formData.Set("client_secret", creds.ClientSecret)
	formData.Set("client_id", creds.ClientID)
}
