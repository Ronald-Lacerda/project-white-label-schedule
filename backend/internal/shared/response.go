package shared

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Meta struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
	Total   int `json:"total"`
}

type successResponse struct {
	Data any   `json:"data"`
	Meta *Meta `json:"meta,omitempty"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Error errorDetail `json:"error"`
}

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(successResponse{Data: data})
}

func JSONList(w http.ResponseWriter, status int, data any, meta Meta) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(successResponse{Data: data, Meta: &meta})
}

func JSONError(w http.ResponseWriter, err error) {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(domainErr.Status)
		json.NewEncoder(w).Encode(errorResponse{Error: errorDetail{
			Code:    domainErr.Code,
			Message: domainErr.Message,
		}})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(errorResponse{Error: errorDetail{
		Code:    "INTERNAL_ERROR",
		Message: "Erro interno do servidor.",
	}})
}

func Decode(r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return ErrInvalidInput
	}
	return nil
}
