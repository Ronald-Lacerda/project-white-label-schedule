package calendar

import (
	"log"
	"net/http"
	"os"

	"schedule/internal/shared"
)

// Handler expõe os endpoints REST da integração Google Agenda.
type Handler struct {
	svc *Service
}

// NewHandler cria um novo Handler com o Service injetado.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// frontendURL retorna a URL base do frontend lida do ambiente.
func frontendURL() string {
	u := os.Getenv("FRONTEND_URL")
	if u == "" {
		u = "http://localhost:3000"
	}
	return u
}

// GetAuthURL gera e retorna a URL de autorização OAuth2 do Google.
//
// GET /api/v1/google/auth-url
// Resposta: { "data": { "url": "https://accounts.google.com/..." } }
func (h *Handler) GetAuthURL(w http.ResponseWriter, r *http.Request) {
	establishmentID := shared.EstablishmentIDFromContext(r.Context())

	url, err := h.svc.GenerateAuthURL(r.Context(), establishmentID)
	if err != nil {
		shared.JSONError(w, err)
		return
	}

	shared.JSON(w, http.StatusOK, map[string]string{"url": url})
}

// Callback processa o redirecionamento do Google após a autorização OAuth2.
// Após processar, redireciona o browser para o frontend.
//
// GET /api/v1/google/callback  (público — sem autenticação JWT)
func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	errParam := r.URL.Query().Get("error")

	base := frontendURL() + "/dashboard/settings/google"

	if errParam != "" || code == "" || state == "" {
		log.Printf("calendar: oauth callback error: %q code=%q state=%q", errParam, code, state)
		http.Redirect(w, r, base+"?status=error", http.StatusFound)
		return
	}

	if err := h.svc.HandleCallback(r.Context(), code, state); err != nil {
		log.Printf("calendar: HandleCallback error: %v", err)
		http.Redirect(w, r, base+"?status=error", http.StatusFound)
		return
	}

	http.Redirect(w, r, base+"?status=success", http.StatusFound)
}

// Disconnect desconecta a integração Google Agenda do estabelecimento autenticado.
//
// DELETE /api/v1/google/disconnect
// Resposta: 204 No Content
func (h *Handler) Disconnect(w http.ResponseWriter, r *http.Request) {
	establishmentID := shared.EstablishmentIDFromContext(r.Context())

	if err := h.svc.Disconnect(r.Context(), establishmentID); err != nil {
		shared.JSONError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetStatus retorna o status da integração e a lista de profissionais com seus calendários.
//
// GET /api/v1/google/status
// Resposta: { "data": { "connected": true, "professionals": [...] } }
func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	establishmentID := shared.EstablishmentIDFromContext(r.Context())

	status, err := h.svc.GetStatus(r.Context(), establishmentID)
	if err != nil {
		shared.JSONError(w, err)
		return
	}

	shared.JSON(w, http.StatusOK, status)
}
