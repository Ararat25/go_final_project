package middleware

import (
	"net/http"

	"github.com/Ararat25/go_final_project/task"
)

// Middleware структура для обработчика-посредника
type Middleware struct {
	service *task.Service
}

// NewMiddleware создание объекта Middleware
func NewMiddleware(service *task.Service) *Middleware {
	return &Middleware{
		service: service,
	}
}

// Auth middleware для проверки авторизации
func (m Middleware) Auth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := m.service.Config.Password
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
