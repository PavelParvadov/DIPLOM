package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"happyhouse/backend/internal/domain"
	"happyhouse/backend/pkg/httpx"
)

const maxUploadSize = 8 << 20

func (h *Handler) decodeCreatePostInput(r *http.Request) (domain.CreatePostInput, error) {
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		var input domain.CreatePostInput
		if err := httpx.DecodeJSON(r, &input); err != nil {
			return domain.CreatePostInput{}, fmt.Errorf("invalid request body")
		}
		return input, nil
	}

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		return domain.CreatePostInput{}, fmt.Errorf("не удалось прочитать форму")
	}

	categoryID, err := strconv.ParseInt(strings.TrimSpace(r.FormValue("categoryId")), 10, 64)
	if err != nil {
		return domain.CreatePostInput{}, fmt.Errorf("invalid categoryId")
	}

	imageURL, err := h.saveOptionalImage(r, "image")
	if err != nil {
		return domain.CreatePostInput{}, err
	}

	return domain.CreatePostInput{
		CategoryID: categoryID,
		Title:      strings.TrimSpace(r.FormValue("title")),
		Content:    strings.TrimSpace(r.FormValue("content")),
		ImageURL:   imageURL,
	}, nil
}

func (h *Handler) decodeCreateChatMessageInput(r *http.Request) (domain.CreateChatMessageInput, error) {
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		var input domain.CreateChatMessageInput
		if err := httpx.DecodeJSON(r, &input); err != nil {
			return domain.CreateChatMessageInput{}, fmt.Errorf("invalid request body")
		}
		return input, nil
	}

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		return domain.CreateChatMessageInput{}, fmt.Errorf("не удалось прочитать форму")
	}

	imageURL, err := h.saveOptionalImage(r, "image")
	if err != nil {
		return domain.CreateChatMessageInput{}, err
	}

	return domain.CreateChatMessageInput{
		Content:  strings.TrimSpace(r.FormValue("content")),
		ImageURL: imageURL,
	}, nil
}

func (h *Handler) saveOptionalImage(r *http.Request, fieldName string) (string, error) {
	file, header, err := r.FormFile(fieldName)
	if err != nil {
		if err == http.ErrMissingFile {
			return "", nil
		}
		return "", fmt.Errorf("не удалось прочитать изображение")
	}
	defer file.Close()

	return h.saveUploadedImage(r.Context(), file, header)
}

func (h *Handler) saveUploadedImage(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	data, err := io.ReadAll(io.LimitReader(file, maxUploadSize+1))
	if err != nil {
		return "", fmt.Errorf("не удалось сохранить изображение")
	}
	if len(data) == 0 {
		return "", nil
	}
	if len(data) > maxUploadSize {
		return "", fmt.Errorf("изображение слишком большое")
	}

	contentType := http.DetectContentType(data)
	extension, ok := allowedImageExtensions[contentType]
	if !ok {
		return "", fmt.Errorf("поддерживаются только JPG, PNG, GIF и WEBP")
	}
	if fromName := strings.ToLower(filepath.Ext(header.Filename)); fromName != "" {
		switch fromName {
		case ".jpg", ".jpeg", ".png", ".gif", ".webp":
			extension = fromName
		}
	}

	fileName, err := randomFileName(extension)
	if err != nil {
		return "", fmt.Errorf("не удалось сохранить изображение")
	}

	if _, err := h.media.Create(ctx, fileName, contentType, data); err == nil {
		return "/uploads/" + fileName, nil
	}

	targetPath := filepath.Join(h.uploadDir, fileName)
	if err := os.WriteFile(targetPath, data, 0o644); err != nil {
		return "", fmt.Errorf("не удалось сохранить изображение")
	}

	return "/uploads/" + fileName, nil
}

func randomFileName(extension string) (string, error) {
	buffer := make([]byte, 16)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return hex.EncodeToString(buffer) + extension, nil
}

var allowedImageExtensions = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}
