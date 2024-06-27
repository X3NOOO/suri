package middleware

import (
	"mime"
	"net/http"
)

func JsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")

		content_type := r.Header.Get("Content-Type")

		mt, _, err := mime.ParseMediaType(content_type)
		if err != nil {
			http.Error(w, "Malformed Content-Type header", http.StatusBadRequest)
			return
		}

		if mt != "application/json" {
			http.Error(w, "Expected application/json Content-Type", http.StatusUnsupportedMediaType)
			return
		}

		next.ServeHTTP(w, r)
	})
}
