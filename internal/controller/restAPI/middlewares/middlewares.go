package middlewares

import (
	"fmt"
	"net/http"
	"users-app/pkg/logger"
)

type Middleware struct {
	log logger.Logger
}

func New(log logger.Logger) *Middleware {
	return &Middleware{
		log: log,
	}
}

func (m *Middleware) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var headers string

		for k, v := range r.Header {
			if k == "authorization" || k == "Cookie" {
				continue
			}

			headers += fmt.Sprintf("%s: %s,\n", k, v)
		}

		m.log.InfoF("incoming request: method = %s, url = %s headers = %s, user_ip = %s",
			r.Method, r.URL.String(), headers, r.RemoteAddr)

		next.ServeHTTP(w, r.WithContext(r.Context()))
	})
}
