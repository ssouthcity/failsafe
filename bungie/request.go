package bungie

import "net/http"

type RequestOption func(req *http.Request)

func WithAPIKey(key string) RequestOption {
	return func(req *http.Request) {
		req.Header.Add("X-API-Key", key)
	}
}
