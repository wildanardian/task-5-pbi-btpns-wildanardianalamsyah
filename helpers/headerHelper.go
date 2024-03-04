package helpers

import "net/http"

func GetAuthorizationHeader(req *http.Request) string {
	return req.Header.Get("Authorization")
}
