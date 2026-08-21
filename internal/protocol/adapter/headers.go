package adapter

import (
	"net/http"
	"strings"
)

func SetProtocolHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Protocol-Version", "1")
	w.Header().Set("Cache-Control", "no-store")
}
func IsParticipantRequest(r *http.Request) bool {
	return strings.HasPrefix(r.URL.Path, "/v1/participant") || strings.HasPrefix(r.URL.Path, "/v1/rounds")
}
func RequireRequestID(r *http.Request) string {
	if v := r.Header.Get("X-Request-ID"); v != "" {
		return v
	}
	return "missing"
}
