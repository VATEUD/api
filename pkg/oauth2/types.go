package oauth2

import (
	"api/internal/pkg/database"
	"api/pkg/models"
	"api/utils"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

const (
	cookieName  = "token"
	contentType = "application/x-www-form-urlencoded"
)

var scopes = []string{"full_name", "email", "vatsim_details", "country"}

type authorizationRequest struct {
	ResponseType, ClientID, RedirectURI, State, Scopes string
}

type accessTokenError struct {
	Err           string `json:"error"`
	Description   string `json:"error_description"`
	Code          int    `json:"-"`
	internalError error
}

func (err *accessTokenError) Json() ([]byte, error) {
	return json.Marshal(err)
}

type requestError struct {
	Response url.Values
}

func (err requestError) Error() string {
	return err.Response.Encode()
}

func newRequest(r *http.Request) (*authorizationRequest, error) {
	if r == nil {
		return nil, errors.New("please provide a valid request")
	}

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	return &authorizationRequest{
		ResponseType: r.Form.Get("response_type"),
		ClientID:     r.Form.Get("client_id"),
		RedirectURI:  r.Form.Get("redirect_uri"),
		Scopes:       r.Form.Get("scope"),
		State:        r.Form.Get("state"),
	}, nil
}

func (request *authorizationRequest) Validate() (*models.OauthClient, error) {
	if len(request.ClientID) < 1 || len(request.ResponseType) < 1 || len(request.RedirectURI) < 1 {
		return nil, requestError{
			Response: request.formatURL("invalid_request"),
		}
	}

	if request.ResponseType != "code" {
		return nil, requestError{Response: request.formatURL("unsupported_response_type")}
	}

	s := strings.Split(request.Scopes, " ")

	// also ensure the first element is an empty string because strings.Split returns an empty string
	if len(s) > 0 && s[0] != "" {
		for _, scope := range s {
			if !request.isValidScope(scope) {
				return nil, requestError{Response: request.formatURL("invalid_scope")}
			}
		}
	} else {
		request.Scopes = strings.Join(scopes, " ")
	}

	client := models.OauthClient{}
	if err := database.DB.Where("id = ?", request.ClientID).First(&client).Error; err != nil {
		return nil, requestError{Response: request.formatURL("unauthorized_client")}
	}

	if client.Revoked {
		return nil, requestError{Response: request.formatURL("unauthorized_client")}
	}

	if !client.IsValidRedirectURI(request.RedirectURI) {
		return nil, requestError{Response: request.formatURL("unauthorized_client")}
	}

	return &client, nil
}

func (request authorizationRequest) formatURL(err string) url.Values {
	val := url.Values{}
	val.Set("error", err)
	val.Set("error_description", getError(err))

	if len(request.State) > 0 {
		val.Set("state", request.State)
	}

	return val
}

func getError(err string) string {
	authErrors := map[string]string{
		"invalid_request":           "The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed.",
		"unauthorized_client":       "The client is not authorized to request an authorization code using this method.",
		"access_denied":             "The resource owner or authorization server denied the request.",
		"unsupported_response_type": "The authorization server does not support obtaining an authorization code using this method.",
		"invalid_scope":             "The requested scope is invalid, unknown, or malformed.",
		"server_error":              "The authorization server encountered an unexpected condition that prevented it from fulfilling the request.",
		"temporarily_unavailable":   "The authorization server is currently unable to handle the request due to a temporary overloading or maintenance of the server.",
		"invalid_grant":             "The provided authorization grant (e.g., authorization code, resource owner credentials) or refresh token is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client.",
		"invalid_client":            "Client authentication failed.",
		"unsupported_grant_type":    "The authorization grant type is not supported by the authorization server.",
	}

	e, ok := authErrors[err]

	if !ok {
		return "Error occurred."
	}

	return e
}

func (request authorizationRequest) isValidScope(scope string) bool {
	for _, availableScope := range scopes {
		if availableScope == scope {
			return true
		}
	}

	return false
}

type accessTokenRequest struct {
	ContentType, GrantType, Code, RedirectURI, ClientID, ClientSecret, State string
}

func newAccessTokenRequest(r *http.Request) (*accessTokenRequest, error) {
	if r == nil {
		return nil, errors.New("please provide a valid request")
	}

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	return &accessTokenRequest{
		ContentType:  r.Header.Get("Content-type"),
		GrantType:    r.PostForm.Get("grant_type"),
		Code:         r.PostForm.Get("code"),
		RedirectURI:  r.PostForm.Get("redirect_uri"),
		ClientID:     r.PostForm.Get("client_id"),
		ClientSecret: r.PostForm.Get("client_secret"),
		State:        r.PostForm.Get("state"),
	}, nil
}

func (r accessTokenRequest) Validate() (*models.OauthClient, *accessTokenError) {
	if !r.hasRequiredAttributes() {
		return nil, &accessTokenError{
			Err:           "invalid_request",
			Code:          http.StatusBadRequest,
			Description:   getError("invalid_request"),
			internalError: errors.New("required parameters weren't provided"),
		}
	}

	if r.ContentType != contentType {
		return nil, &accessTokenError{
			Err:           "invalid_request",
			Code:          http.StatusBadRequest,
			Description:   "Invalid content type provided",
			internalError: errors.New("invalid content type provided"),
		}
	}

	if r.GrantType != "authorization_code" {
		return nil, &accessTokenError{
			Err:           "invalid_grant",
			Code:          http.StatusBadRequest,
			Description:   getError("invalid_grant"),
			internalError: errors.New("invalid grant type provided"),
		}
	}

	code := models.OauthAuthCode{}
	if err := database.DB.Where("id = ?", r.Code).First(&code).Error; err != nil {
		return nil, &accessTokenError{
			Err:           "invalid_client",
			Code:          http.StatusUnauthorized,
			Description:   getError("invalid_client"),
			internalError: errors.New("auth code not found"),
		}
	}

	client := models.OauthClient{}
	if err := database.DB.Where("id = ? AND secret = ?", r.ClientID, r.ClientSecret).First(&client).Error; err != nil {
		return nil, &accessTokenError{
			Err:           "invalid_client",
			Code:          http.StatusUnauthorized,
			Description:   getError("invalid_client"),
			internalError: errors.New("client not found"),
		}
	}

	if client.Revoked {
		return nil, &accessTokenError{
			Err:           "invalid_client",
			Code:          http.StatusUnauthorized,
			Description:   getError("invalid_client"),
			internalError: errors.New("client is revoked"),
		}
	}

	if client.Secret != r.ClientSecret {
		return nil, &accessTokenError{
			Err:           "invalid_client",
			Code:          http.StatusUnauthorized,
			Description:   getError("invalid_client"),
			internalError: errors.New("client secret does not match"),
		}
	}

	if client.ID != code.ClientID {
		return nil, &accessTokenError{
			Err:           "invalid_client",
			Code:          http.StatusUnauthorized,
			Description:   getError("invalid_client"),
			internalError: errors.New("auth code client and provided client don't match not found"),
		}
	}

	if !client.IsValidRedirectURI(r.RedirectURI) {
		return nil, &accessTokenError{
			Err:           "invalid_client",
			Code:          http.StatusUnauthorized,
			Description:   getError("invalid_client"),
			internalError: errors.New("invalid redirect URI"),
		}
	}

	return &client, nil
}

func (r accessTokenRequest) hasRequiredAttributes() bool {
	if utils.IsEmptyString(r.GrantType) {
		return false
	}

	if utils.IsEmptyString(r.Code) {
		return false
	}

	if utils.IsEmptyString(r.RedirectURI) {
		return false
	}

	if utils.IsEmptyString(r.ClientSecret) {
		return false
	}

	if utils.IsEmptyString(r.ClientID) {
		return false
	}

	return true
}
