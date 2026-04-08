package tenancy

import (
	"net/http"

	"schedule/internal/shared"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	establishmentID := shared.EstablishmentIDFromContext(r.Context())
	e, err := h.svc.GetByID(r.Context(), establishmentID)
	if err != nil {
		shared.JSONError(w, err)
		return
	}
	shared.JSON(w, http.StatusOK, e)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	establishmentID := shared.EstablishmentIDFromContext(r.Context())

	var req struct {
		Name                  string  `json:"name"`
		Slug                  string  `json:"slug"`
		Timezone              string  `json:"timezone"`
		ContactEmail          *string `json:"contact_email"`
		ContactPhone          *string `json:"contact_phone"`
		MinAdvanceCancelHours int     `json:"min_advance_cancel_hours"`
	}
	if err := shared.Decode(r, &req); err != nil {
		shared.JSONError(w, err)
		return
	}

	if req.Name == "" || req.Slug == "" {
		shared.JSONError(w, shared.ErrInvalidInput)
		return
	}

	e, err := h.svc.Update(r.Context(), establishmentID, UpdateInput{
		Name:                  req.Name,
		Slug:                  req.Slug,
		Timezone:              req.Timezone,
		ContactEmail:          req.ContactEmail,
		ContactPhone:          req.ContactPhone,
		MinAdvanceCancelHours: req.MinAdvanceCancelHours,
	})
	if err != nil {
		shared.JSONError(w, err)
		return
	}
	shared.JSON(w, http.StatusOK, e)
}

func (h *Handler) GetBusinessHours(w http.ResponseWriter, r *http.Request) {
	establishmentID := shared.EstablishmentIDFromContext(r.Context())
	hours, err := h.svc.GetBusinessHours(r.Context(), establishmentID)
	if err != nil {
		shared.JSONError(w, err)
		return
	}
	shared.JSON(w, http.StatusOK, hours)
}

func (h *Handler) UpdateBusinessHours(w http.ResponseWriter, r *http.Request) {
	establishmentID := shared.EstablishmentIDFromContext(r.Context())

	var req struct {
		Hours []struct {
			DayOfWeek int    `json:"day_of_week"`
			OpenTime  string `json:"open_time"`
			CloseTime string `json:"close_time"`
			IsClosed  bool   `json:"is_closed"`
		} `json:"hours"`
	}
	if err := shared.Decode(r, &req); err != nil {
		shared.JSONError(w, err)
		return
	}

	inputs := make([]BusinessHourInput, 0, len(req.Hours))
	for _, h := range req.Hours {
		inputs = append(inputs, BusinessHourInput{
			DayOfWeek: h.DayOfWeek,
			OpenTime:  h.OpenTime,
			CloseTime: h.CloseTime,
			IsClosed:  h.IsClosed,
		})
	}

	hours, err := h.svc.UpdateBusinessHours(r.Context(), establishmentID, inputs)
	if err != nil {
		shared.JSONError(w, err)
		return
	}
	shared.JSON(w, http.StatusOK, hours)
}
