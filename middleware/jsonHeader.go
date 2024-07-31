package middleware

import (
	"net/http"
)

// JsonHeader устанавливает заголовок application/json в ответ
func JsonHeader(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
