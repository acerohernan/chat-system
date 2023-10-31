package auth

import (
	"net/http"
	"strings"
)

func ExtractTokenFromRequest(r *http.Request) string {
	header := r.Header.Get("Authorization")

	if header == "" {
		return ""
	}

	rawJWT, prefixFound := strings.CutPrefix(header, "Bearer ")

	if !prefixFound {
		return ""
	}

	return rawJWT
}
