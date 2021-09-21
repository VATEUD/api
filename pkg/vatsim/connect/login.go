package connect

import (
	"auth/utils"
	"net/http"
	"net/url"
)

func Login(w http.ResponseWriter, r *http.Request) {
	values := url.Values{}
	values.Set("client_id", utils.Getenv("CONNECT_CLIENT_ID", ""))
	values.Set("redirect_uri", utils.Getenv("CONNECT_REDIRECT", ""))
	values.Set("response_type", "code")
	values.Set("scope", utils.Getenv("CONNECT_SCOPES", ""))
	values.Set("required_scopes", utils.Getenv("CONNECT_SCOPES", ""))
	values.Set("state", utils.RandomString(40))

	redirectURL := utils.ConnectURL("oauth/authorize", values.Encode())

	http.Redirect(w, r, redirectURL.String(), http.StatusTemporaryRedirect)
}
