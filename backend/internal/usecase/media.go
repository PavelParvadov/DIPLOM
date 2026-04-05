package usecase

import (
	"context"

	"happyhouse/backend/internal/domain"
)

type MediaUseCase struct {
	media domain.MediaRepository
}

func NewMediaUseCase(media domain.MediaRepository) *MediaUseCase {
	return &MediaUseCase{media: media}
}

func (u *MediaUseCase) Create(ctx context.Context, publicID, contentType string, data []byte) (*domain.MediaAsset, error) {
	return u.media.Create(ctx, publicID, contentType, data)
}

func (u *MediaUseCase) GetByPublicID(ctx context.Context, publicID string) (*domain.MediaAsset, error) {
	return u.media.GetByPublicID(ctx, publicID)
}
