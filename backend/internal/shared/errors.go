package shared

import "net/http"

// DomainError representa um erro de negocio com codigo e status HTTP.
type DomainError struct {
	Code    string
	Message string
	Status  int
}

func (e *DomainError) Error() string {
	return e.Message
}

var (
	ErrSlotUnavailable         = &DomainError{Code: "SLOT_UNAVAILABLE", Message: "O horario selecionado nao esta mais disponivel.", Status: http.StatusConflict}
	ErrOutsideBusinessHours    = &DomainError{Code: "OUTSIDE_BUSINESS_HOURS", Message: "Horario fora do expediente do estabelecimento.", Status: http.StatusUnprocessableEntity}
	ErrProfessionalInactive    = &DomainError{Code: "PROFESSIONAL_INACTIVE", Message: "Profissional indisponivel.", Status: http.StatusUnprocessableEntity}
	ErrServiceDurationConflict = &DomainError{Code: "SERVICE_DURATION_CONFLICT", Message: "O servico nao cabe no horario selecionado.", Status: http.StatusUnprocessableEntity}
	ErrNotFound                = &DomainError{Code: "NOT_FOUND", Message: "Recurso nao encontrado.", Status: http.StatusNotFound}
	ErrUnauthorized            = &DomainError{Code: "UNAUTHORIZED", Message: "Nao autenticado.", Status: http.StatusUnauthorized}
	ErrForbidden               = &DomainError{Code: "FORBIDDEN", Message: "Sem permissao para esta operacao.", Status: http.StatusForbidden}
	ErrEmailConflict           = &DomainError{Code: "EMAIL_CONFLICT", Message: "Este e-mail ja esta cadastrado.", Status: http.StatusConflict}
	ErrSlugConflict            = &DomainError{Code: "SLUG_CONFLICT", Message: "Este identificador ja esta em uso.", Status: http.StatusConflict}
	ErrInvalidInput            = &DomainError{Code: "INVALID_INPUT", Message: "Dados invalidos na requisicao.", Status: http.StatusBadRequest}
	ErrCancelTooLate           = &DomainError{Code: "CANCEL_TOO_LATE", Message: "Cancelamento fora do prazo minimo.", Status: http.StatusUnprocessableEntity}
	ErrIntegrationNotConfigured = &DomainError{Code: "INTEGRATION_NOT_CONFIGURED", Message: "A integracao ainda nao foi configurada no servidor.", Status: http.StatusPreconditionFailed}
)
