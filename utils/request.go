package utils

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

type HTTPRequest struct {
	client   *http.Client
	URL      url.URL
	Method   string
	Body     io.Reader
	Error    error
	Response *http.Response
	Request  *http.Request
}

func NewRequest(URL, method string, body io.Reader) *HTTPRequest {
	return &HTTPRequest{
		URL:    buildURLParts(URL),
		Method: method,
		Body:   body,
		client: &http.Client{},
	}
}

func buildURLParts(URL string) url.URL {
	return url.URL{
		Scheme: strings.Split(URL, "://")[0],
		Host:   strings.Split(strings.Join(strings.Split(URL, "//")[1:], "/"), "/")[0],
		Path:   strings.Join(strings.Split(strings.Join(strings.Split(URL, "//")[1:], "/"), "/")[1:], "/"),
	}
}

func (r *HTTPRequest) Do() error {
	r.Request, r.Error = r.request()

	if r.Error != nil {
		return r.Error
	}

	r.Response, r.Error = r.client.Do(r.Request)

	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (r *HTTPRequest) DoConcurrently(checker chan bool) {
	r.Request, r.Error = r.request()

	if r.Error != nil {
		checker <- false
		return
	}

	r.Response, r.Error = r.client.Do(r.Request)

	if r.Error != nil {
		checker <- false
		return
	}

	checker <- true
	return
}

func (r *HTTPRequest) request() (*http.Request, error) {
	return http.NewRequest(r.Method, r.URL.String(), r.Body)
}
