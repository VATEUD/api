package oauth2

import (
	"api/internal/pkg/database"
	"api/pkg/models"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

var scopes = []string{"full_name", "email", "vatsim_details", "country"}

type authorizationRequest struct {
	ResponseType, ClientID, RedirectURI, State string
	Scopes []string
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
		Scopes:       strings.Split(r.Form.Get("scope"), " "),
		State:        r.Form.Get("state"),
	}, nil
}

func (request *authorizationRequest) Validate() error {
	if len(request.ClientID) < 1 || len(request.ResponseType) < 1 || len(request.RedirectURI) < 1 {
		return requestError{
			Response: request.formatURL("invalid_request"),
		}
	}

	if request.ResponseType != "code" {
		return requestError{Response: request.formatURL("unsupported_response_type")}
	}

	// also ensure the first element is an empty string because strings.Split returns an empty string
	if len(request.Scopes) > 0 && request.Scopes[0] != "" {
		for _, scope := range request.Scopes {
			if !request.isValidScope(scope) {
				return requestError{Response: request.formatURL("invalid_scope")}
			}
		}
	} else {
		request.Scopes = scopes
	}

	client := models.OauthClient{}
	if err := database.DB.Where("id = ?", request.ClientID).First(&client).Error; err != nil {
		return requestError{Response: request.formatURL("unauthorized_client")}
	}

	if client.Revoked {
		return requestError{Response: request.formatURL("unauthorized_client")}
	}

	if !client.IsValidRedirectURI(request.RedirectURI) {
		return requestError{Response: request.formatURL("unauthorized_client")}
	}

	return nil
}

func (request *authorizationRequest) formatURL(err string) url.Values {
	val := url.Values{}
	val.Set("error", err)
	val.Set("error_description", request.getError(err))

	if len(request.State) > 0 {
		val.Set("state", request.State)
	}

	return val
}

func (request authorizationRequest) getError(err string) string {
	authErrors := map[string]string{
		"invalid_request":           "The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed.",
		"unauthorized_client":       "The client is not authorized to request an authorization code using this method.",
		"access_denied":             "The resource owner or authorization server denied the request.",
		"unsupported_response_type": "The authorization server does not support obtaining an authorization code using this method.",
		"invalid_scope":             "The requested scope is invalid, unknown, or malformed.",
		"server_error":              "The authorization server encountered an unexpected condition that prevented it from fulfilling the request.",
		"temporarily_unavailable":   "The authorization server is currently unable to handle the request due to a temporary overloading or maintenance of the server.",
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
