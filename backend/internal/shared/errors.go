package shared

import "net/http"

// DomainError representa um erro de negócio com código e status HTTP.
type DomainError struct {
	Code    string
	Message string
	Status  int
}

func (e *DomainError) Error() string {
	return e.Message
}

var (
	ErrSlotUnavailable         = &DomainError{Code: "SLOT_UNAVAILABLE", Message: "O horário selecionado não está mais disponível.", Status: http.StatusConflict}
	ErrOutsideBusinessHours    = &DomainError{Code: "OUTSIDE_BUSINESS_HOURS", Message: "Horário fora do expediente do estabelecimento.", Status: http.StatusUnprocessableEntity}
	ErrProfessionalInactive    = &DomainError{Code: "PROFESSIONAL_INACTIVE", Message: "Profissional indisponível.", Status: http.StatusUnprocessableEntity}
	ErrServiceDurationConflict = &DomainError{Code: "SERVICE_DURATION_CONFLICT", Message: "O serviço não cabe no horário selecionado.", Status: http.StatusUnprocessableEntity}
	ErrNotFound                = &DomainError{Code: "NOT_FOUND", Message: "Recurso não encontrado.", Status: http.StatusNotFound}
	ErrUnauthorized            = &DomainError{Code: "UNAUTHORIZED", Message: "Não autenticado.", Status: http.StatusUnauthorized}
	ErrForbidden               = &DomainError{Code: "FORBIDDEN", Message: "Sem permissão para esta operação.", Status: http.StatusForbidden}
	ErrEmailConflict           = &DomainError{Code: "EMAIL_CONFLICT", Message: "Este e-mail já está cadastrado.", Status: http.StatusConflict}
	ErrSlugConflict            = &DomainError{Code: "SLUG_CONFLICT", Message: "Este identificador já está em uso.", Status: http.StatusConflict}
	ErrInvalidInput            = &DomainError{Code: "INVALID_INPUT", Message: "Dados inválidos na requisição.", Status: http.StatusBadRequest}
	ErrCancelTooLate           = &DomainError{Code: "CANCEL_TOO_LATE", Message: "Cancelamento fora do prazo mínimo.", Status: http.StatusUnprocessableEntity}
	ErrIntegrationNotConfigured = &DomainError{Code: "INTEGRATION_NOT_CONFIGURED", Message: "A integração ainda não foi configurada no servidor.", Status: http.StatusPreconditionFailed}
)
