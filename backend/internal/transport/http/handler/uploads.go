package handler

import (
	"bytes"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"

	"happyhouse/backend/internal/domain"
)

func (h *Handler) handleUploadedAsset(w http.ResponseWriter, r *http.Request) {
	publicID := strings.TrimSpace(chi.URLParam(r, "*"))
	if publicID == "" {
		http.NotFound(w, r)
		return
	}

	asset, err := h.media.GetByPublicID(r.Context(), publicID)
	if err == nil {
		if asset.ContentType != "" {
			w.Header().Set("Content-Type", asset.ContentType)
		}
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		http.ServeContent(w, r, publicID, asset.CreatedAt, bytes.NewReader(asset.Data))
		return
	}
	if !errors.Is(err, domain.ErrNotFound) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	targetPath := filepath.Join(h.uploadDir, filepath.Base(publicID))
	if _, statErr := os.Stat(targetPath); statErr == nil {
		http.ServeFile(w, r, targetPath)
		return
	}

	http.NotFound(w, r)
}
