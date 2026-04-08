package shared

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// TenantResolver é implementado pelo tenancy.Service.
type TenantResolver interface {
	GetBySlug(ctx context.Context, slug string) (interface{ GetID() string }, error)
}

// SlugTenantMiddleware resolve o establishment_id a partir do slug na URL
// e o injeta no contexto para as rotas públicas /pub/{slug}/...
func SlugTenantMiddleware(resolver func(ctx context.Context, slug string) (string, error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slug := chi.URLParam(r, "slug")
			if slug == "" {
				JSONError(w, ErrNotFound)
				return
			}

			establishmentID, err := resolver(r.Context(), slug)
			if err != nil {
				JSONError(w, ErrNotFound)
				return
			}

			ctx := WithEstablishmentID(r.Context(), establishmentID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
