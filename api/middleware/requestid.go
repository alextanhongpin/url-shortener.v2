package middleware

import (
	"net/http"

	"github.com/alextanhongpin/pkg/requestid"

	"github.com/rs/xid"
)

// HTTP middleware setting a value on the request context
func RequestID(next http.Handler) http.Handler {
	provider := requestid.RequestID(func() (string, error) {
		return xid.New().String(), nil
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := provider(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}
