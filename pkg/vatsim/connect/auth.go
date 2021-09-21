package connect

import (
	"auth/utils"
	"net/http"
	"net/url"
)

func Login(w http.ResponseWriter, r *http.Request) {
	value := url.Values{}
	value.Set("client_id", utils.Getenv("CONNECT_CLIENT_ID", ""))
	value.Set("redirect_uri", utils.Getenv("CONNECT_REDIRECT", ""))
	value.Set("response_type", "code")
	value.Set("scope", utils.Getenv("CONNECT_SCOPES", ""))
	value.Set("required_scopes", utils.Getenv("CONNECT_SCOPES", ""))
	value.Set("state", utils.RandomString(40))

	redirectURL := url.URL{
		Scheme:      "https",
		Host:        utils.Getenv("CONNECT_URL", "auth-dev.vatsim.net"),
		Path:        "oauth/authorize",
		RawQuery: value.Encode(),
	}

	http.Redirect(w, r, redirectURL.String(), http.StatusTemporaryRedirect)
}
