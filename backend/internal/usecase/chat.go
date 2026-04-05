package usecase

import (
	"context"
	"strings"

	"happyhouse/backend/internal/domain"
)

const chatMessagesLimit = 150

type ChatUseCase struct {
	houses domain.HouseRepository
	chats  domain.ChatRepository
}

func NewChatUseCase(houses domain.HouseRepository, chats domain.ChatRepository) *ChatUseCase {
	return &ChatUseCase{houses: houses, chats: chats}
}

func (uc *ChatUseCase) List(ctx context.Context, userID, houseID int64) ([]domain.ChatMessage, error) {
	if _, err := ensureMember(ctx, uc.houses, userID, houseID); err != nil {
		return nil, err
	}
	return uc.chats.ListByHouse(ctx, houseID, chatMessagesLimit)
}

func (uc *ChatUseCase) Create(ctx context.Context, userID, houseID int64, input domain.CreateChatMessageInput) (*domain.ChatMessage, error) {
	input.Content = strings.TrimSpace(input.Content)
	if !validateChatMessage(input.Content, input.ImageURL) {
		return nil, domain.ErrValidation
	}
	if _, err := ensureMember(ctx, uc.houses, userID, houseID); err != nil {
		return nil, err
	}
	return uc.chats.Create(ctx, houseID, userID, input)
}
