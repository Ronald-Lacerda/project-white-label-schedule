package catalog

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"schedule/internal/shared"
)

type ProfessionalHandler struct {
	svc *ProfessionalService
}

func NewProfessionalHandler(svc *ProfessionalService) *ProfessionalHandler {
	return &ProfessionalHandler{svc: svc}
}

func (h *ProfessionalHandler) List(w http.ResponseWriter, r *http.Request) {
	estID := shared.EstablishmentIDFromContext(r.Context())
	list, err := h.svc.List(r.Context(), estID)
	if err != nil {
		shared.JSONError(w, err)
		return
	}
	shared.JSON(w, http.StatusOK, list)
}

func (h *ProfessionalHandler) Create(w http.ResponseWriter, r *http.Request) {
	estID := shared.EstablishmentIDFromContext(r.Context())

	var req struct {
		Name         string  `json:"name"`
		AvatarURL    *string `json:"avatar_url"`
		Phone        *string `json:"phone"`
		DisplayOrder int     `json:"display_order"`
	}
	if err := shared.Decode(r, &req); err != nil {
		shared.JSONError(w, err)
		return
	}
	if req.Name == "" {
		shared.JSONError(w, shared.ErrInvalidInput)
		return
	}

	p, err := h.svc.Create(r.Context(), estID, ProfessionalInput{
		Name:         req.Name,
		AvatarURL:    req.AvatarURL,
		Phone:        req.Phone,
		DisplayOrder: req.DisplayOrder,
	})
	if err != nil {
		shared.JSONError(w, err)
		return
	}
	shared.JSON(w, http.StatusCreated, p)
}

func (h *ProfessionalHandler) Get(w http.ResponseWriter, r *http.Request) {
	estID := shared.EstablishmentIDFromContext(r.Context())
	id := chi.URLParam(r, "id")

	p, err := h.svc.Get(r.Context(), id, estID)
	if err != nil {
		shared.JSONError(w, err)
		return
	}
	shared.JSON(w, http.StatusOK, p)
}

func (h *ProfessionalHandler) Update(w http.ResponseWriter, r *http.Request) {
	estID := shared.EstablishmentIDFromContext(r.Context())
	id := chi.URLParam(r, "id")

	var req struct {
		Name         *string `json:"name"`
		AvatarURL    *string `json:"avatar_url"`
		Phone        *string `json:"phone"`
		DisplayOrder *int    `json:"display_order"`
		Active       *bool   `json:"active"`
	}
	if err := shared.Decode(r, &req); err != nil {
		shared.JSONError(w, err)
		return
	}

	if req.Name != nil {
		if *req.Name == "" {
			shared.JSONError(w, shared.ErrInvalidInput)
			return
		}

		displayOrder := 0
		if req.DisplayOrder != nil {
			displayOrder = *req.DisplayOrder
		}

		p, err := h.svc.Update(r.Context(), id, estID, ProfessionalInput{
			Name:         *req.Name,
			AvatarURL:    req.AvatarURL,
			Phone:        req.Phone,
			DisplayOrder: displayOrder,
		})
		if err != nil {
			shared.JSONError(w, err)
			return
		}
		shared.JSON(w, http.StatusOK, p)
		return
	}

	if req.Active == nil && req.DisplayOrder == nil && req.Phone == nil && req.AvatarURL == nil {
		shared.JSONError(w, shared.ErrInvalidInput)
		return
	}

	p, err := h.svc.Patch(r.Context(), id, estID, ProfessionalPatchInput{
		AvatarURL:         req.AvatarURL,
		AvatarURLProvided: req.AvatarURL != nil,
		Phone:             req.Phone,
		PhoneProvided:     req.Phone != nil,
		DisplayOrder:      req.DisplayOrder,
		Active:            req.Active,
	})
	if err != nil {
		shared.JSONError(w, err)
		return
	}
	shared.JSON(w, http.StatusOK, p)
}

func (h *ProfessionalHandler) Delete(w http.ResponseWriter, r *http.Request) {
	estID := shared.EstablishmentIDFromContext(r.Context())
	id := chi.URLParam(r, "id")

	if err := h.svc.Delete(r.Context(), id, estID); err != nil {
		shared.JSONError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProfessionalHandler) UpdateHours(w http.ResponseWriter, r *http.Request) {
	estID := shared.EstablishmentIDFromContext(r.Context())
	id := chi.URLParam(r, "id")

	var req struct {
		Hours []struct {
			DayOfWeek     int    `json:"day_of_week"`
			StartTime     string `json:"start_time"`
			EndTime       string `json:"end_time"`
			IsUnavailable bool   `json:"is_unavailable"`
		} `json:"hours"`
	}
	if err := shared.Decode(r, &req); err != nil {
		shared.JSONError(w, err)
		return
	}

	inputs := make([]ProfessionalHourInput, 0, len(req.Hours))
	for _, h := range req.Hours {
		inputs = append(inputs, ProfessionalHourInput{
			DayOfWeek:     h.DayOfWeek,
			StartTime:     h.StartTime,
			EndTime:       h.EndTime,
			IsUnavailable: h.IsUnavailable,
		})
	}

	hours, err := h.svc.UpdateHours(r.Context(), id, estID, inputs)
	if err != nil {
		shared.JSONError(w, err)
		return
	}
	shared.JSON(w, http.StatusOK, hours)
}

func (h *ProfessionalHandler) UpdateServices(w http.ResponseWriter, r *http.Request) {
	estID := shared.EstablishmentIDFromContext(r.Context())
	id := chi.URLParam(r, "id")

	var req struct {
		ServiceIDs []string `json:"service_ids"`
	}
	if err := shared.Decode(r, &req); err != nil {
		shared.JSONError(w, err)
		return
	}

	if err := h.svc.UpdateServices(r.Context(), id, estID, req.ServiceIDs); err != nil {
		shared.JSONError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
