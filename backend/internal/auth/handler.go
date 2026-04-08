package auth

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

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerRequest struct {
	OwnerName         string  `json:"owner_name"`
	EstablishmentName string  `json:"establishment_name"`
	Email             string  `json:"email"`
	Password          string  `json:"password"`
	Slug              string  `json:"slug"`
	ContactPhone      *string `json:"contact_phone"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type logoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := shared.Decode(r, &req); err != nil {
		shared.JSONError(w, err)
		return
	}

	if req.Email == "" || req.Password == "" {
		shared.JSONError(w, shared.ErrInvalidInput)
		return
	}

	pair, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		shared.JSONError(w, err)
		return
	}

	shared.JSON(w, http.StatusOK, pair)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := shared.Decode(r, &req); err != nil {
		shared.JSONError(w, err)
		return
	}

	pair, err := h.svc.Register(r.Context(), RegisterInput{
		OwnerName:         req.OwnerName,
		EstablishmentName: req.EstablishmentName,
		Email:             req.Email,
		Password:          req.Password,
		Slug:              req.Slug,
		ContactPhone:      req.ContactPhone,
	})
	if err != nil {
		shared.JSONError(w, err)
		return
	}

	shared.JSON(w, http.StatusCreated, pair)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := shared.Decode(r, &req); err != nil {
		shared.JSONError(w, err)
		return
	}

	if req.RefreshToken == "" {
		shared.JSONError(w, shared.ErrUnauthorized)
		return
	}

	pair, err := h.svc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		shared.JSONError(w, err)
		return
	}

	shared.JSON(w, http.StatusOK, pair)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var req logoutRequest
	if err := shared.Decode(r, &req); err != nil {
		shared.JSONError(w, err)
		return
	}

	if req.RefreshToken != "" {
		h.svc.Logout(r.Context(), req.RefreshToken)
	}

	w.WriteHeader(http.StatusNoContent)
}
