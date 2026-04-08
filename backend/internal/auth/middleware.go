package auth

import (
	"net/http"
	"strings"

	"schedule/internal/shared"
)

type Middleware struct {
	svc *Service
}

func NewMiddleware(svc *Service) *Middleware {
	return &Middleware{svc: svc}
}

// Authenticate valida o JWT e injeta user_id + establishment_id no contexto.
func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			shared.JSONError(w, shared.ErrUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := m.svc.ValidateAccessToken(tokenStr)
		if err != nil {
			shared.JSONError(w, shared.ErrUnauthorized)
			return
		}

		ctx := shared.WithUserID(r.Context(), claims.UserID)
		ctx = shared.WithEstablishmentID(ctx, claims.EstablishmentID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
