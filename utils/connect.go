package utils

import "net/url"

func ConnectURL(path, query string) url.URL {
	return url.URL{
		Scheme:   "https",
		Host:     Getenv("CONNECT_URL", "auth-dev.vatsim.net"),
		Path:     path,
		RawQuery: query,
	}
}
