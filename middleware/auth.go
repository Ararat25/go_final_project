package middleware

import (
	"github.com/Ararat25/go_final_project/model"
	"net/http"
	"os"
)

type Middleware struct {
	service *model.Service
}

func NewMiddleware(service *model.Service) *Middleware {
	return &Middleware{
		service: service,
	}
}

func (m Middleware) Auth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			token, err := r.Cookie("token")
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			var valid bool

			err = m.service.CheckToken(token.Value, pass)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			valid = true

			if !valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
