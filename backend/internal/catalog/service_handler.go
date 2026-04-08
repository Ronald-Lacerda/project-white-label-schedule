package catalog

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"schedule/internal/shared"
)

type SvcHandler struct {
	svc *SvcService
}

func NewSvcHandler(svc *SvcService) *SvcHandler {
	return &SvcHandler{svc: svc}
}

func (h *SvcHandler) List(w http.ResponseWriter, r *http.Request) {
	estID := shared.EstablishmentIDFromContext(r.Context())
	list, err := h.svc.List(r.Context(), estID)
	if err != nil {
		shared.JSONError(w, err)
		return
	}
	shared.JSON(w, http.StatusOK, list)
}

func (h *SvcHandler) Create(w http.ResponseWriter, r *http.Request) {
	estID := shared.EstablishmentIDFromContext(r.Context())

	var req struct {
		Name            string  `json:"name"`
		Description     *string `json:"description"`
		DurationMinutes int     `json:"duration_minutes"`
		PriceCents      *int    `json:"price_cents"`
		DisplayOrder    int     `json:"display_order"`
	}
	if err := shared.Decode(r, &req); err != nil {
		shared.JSONError(w, err)
		return
	}
	if req.Name == "" || req.DurationMinutes <= 0 {
		shared.JSONError(w, shared.ErrInvalidInput)
		return
	}

	svc, err := h.svc.Create(r.Context(), estID, SvcInput{
		Name:            req.Name,
		Description:     req.Description,
		DurationMinutes: req.DurationMinutes,
		PriceCents:      req.PriceCents,
		DisplayOrder:    req.DisplayOrder,
	})
	if err != nil {
		shared.JSONError(w, err)
		return
	}
	shared.JSON(w, http.StatusCreated, svc)
}

func (h *SvcHandler) Update(w http.ResponseWriter, r *http.Request) {
	estID := shared.EstablishmentIDFromContext(r.Context())
	id := chi.URLParam(r, "id")

	var req struct {
		Name            string  `json:"name"`
		Description     *string `json:"description"`
		DurationMinutes int     `json:"duration_minutes"`
		PriceCents      *int    `json:"price_cents"`
		DisplayOrder    int     `json:"display_order"`
	}
	if err := shared.Decode(r, &req); err != nil {
		shared.JSONError(w, err)
		return
	}
	if req.Name == "" || req.DurationMinutes <= 0 {
		shared.JSONError(w, shared.ErrInvalidInput)
		return
	}

	svc, err := h.svc.Update(r.Context(), id, estID, SvcInput{
		Name:            req.Name,
		Description:     req.Description,
		DurationMinutes: req.DurationMinutes,
		PriceCents:      req.PriceCents,
		DisplayOrder:    req.DisplayOrder,
	})
	if err != nil {
		shared.JSONError(w, err)
		return
	}
	shared.JSON(w, http.StatusOK, svc)
}

func (h *SvcHandler) Delete(w http.ResponseWriter, r *http.Request) {
	estID := shared.EstablishmentIDFromContext(r.Context())
	id := chi.URLParam(r, "id")

	if err := h.svc.Delete(r.Context(), id, estID); err != nil {
		shared.JSONError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
