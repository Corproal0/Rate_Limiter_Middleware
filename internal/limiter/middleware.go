package limiter

import (
	"net"
	"net/http"
)

type Middleware struct {
	store Store
}

func NewMiddleware(store Store) *Middleware {
	return &Middleware{store: store}
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		allowed, err := m.store.Allow(r.Context(), ip)
		if err != nil || !allowed {
			w.Header().Set("Retry-After", "1")
			http.Error(w, "429 Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
