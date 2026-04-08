package whitelabel

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"schedule/internal/shared"
)

type Handler struct {
	svc        *Service
	uploadsDir string
	baseURL    string
}

func NewHandler(svc *Service, uploadsDir, baseURL string) *Handler {
	return &Handler{svc: svc, uploadsDir: uploadsDir, baseURL: baseURL}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	establishmentID := shared.EstablishmentIDFromContext(r.Context())
	cfg, err := h.svc.Get(r.Context(), establishmentID)
	if err != nil {
		shared.JSONError(w, err)
		return
	}
	shared.JSON(w, http.StatusOK, cfg)
}

// UploadLogo trata POST /api/v1/whitelabel/logo (multipart, campo "logo")
func (h *Handler) UploadLogo(w http.ResponseWriter, r *http.Request) {
	establishmentID := shared.EstablishmentIDFromContext(r.Context())

	if err := r.ParseMultipartForm(5 << 20); err != nil { // 5 MB
		shared.JSONError(w, &shared.DomainError{Code: "INVALID_FORM", Message: "Erro ao processar o formulario.", Status: http.StatusBadRequest})
		return
	}

	file, header, err := r.FormFile("logo")
	if err != nil {
		shared.JSONError(w, &shared.DomainError{Code: "MISSING_FILE", Message: "Arquivo 'logo' nao encontrado.", Status: http.StatusBadRequest})
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".svg": true}
	if !allowed[ext] {
		shared.JSONError(w, &shared.DomainError{Code: "INVALID_FILE_TYPE", Message: "Tipo de arquivo nao permitido. Use jpg, png, webp ou svg.", Status: http.StatusBadRequest})
		return
	}

	if err := os.MkdirAll(h.uploadsDir, 0755); err != nil {
		shared.JSONError(w, err)
		return
	}

	filename := fmt.Sprintf("%s%s", establishmentID, ext)
	destPath := filepath.Join(h.uploadsDir, filename)

	dest, err := os.Create(destPath)
	if err != nil {
		shared.JSONError(w, err)
		return
	}
	defer dest.Close()

	if _, err := io.Copy(dest, file); err != nil {
		shared.JSONError(w, err)
		return
	}

	logoURL := fmt.Sprintf("%s/uploads/logos/%s", h.baseURL, filename)

	cfg, err := h.svc.UpdateLogoURL(r.Context(), establishmentID, logoURL)
	if err != nil {
		shared.JSONError(w, err)
		return
	}

	shared.JSON(w, http.StatusOK, cfg)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	establishmentID := shared.EstablishmentIDFromContext(r.Context())

	var req struct {
		LogoURL        *string `json:"logo_url"`
		PrimaryColor   string  `json:"primary_color"`
		SecondaryColor *string `json:"secondary_color"`
		CustomCSS      *string `json:"custom_css"`
	}
	if err := shared.Decode(r, &req); err != nil {
		shared.JSONError(w, err)
		return
	}

	cfg, err := h.svc.Update(r.Context(), establishmentID, UpdateInput{
		LogoURL:        req.LogoURL,
		PrimaryColor:   req.PrimaryColor,
		SecondaryColor: req.SecondaryColor,
		CustomCSS:      req.CustomCSS,
	})
	if err != nil {
		shared.JSONError(w, err)
		return
	}
	shared.JSON(w, http.StatusOK, cfg)
}
