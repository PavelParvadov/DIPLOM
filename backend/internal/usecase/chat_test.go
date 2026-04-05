package usecase

import (
	"context"
	"testing"
	"time"

	"happyhouse/backend/internal/domain"
)

type fakeChatRepo struct {
	items []domain.ChatMessage
}

func (r *fakeChatRepo) ListByHouse(_ context.Context, houseID int64, limit int) ([]domain.ChatMessage, error) {
	items := make([]domain.ChatMessage, 0)
	for _, item := range r.items {
		if item.HouseID == houseID {
			items = append(items, item)
		}
	}
	if len(items) > limit {
		items = items[len(items)-limit:]
	}
	return items, nil
}

func (r *fakeChatRepo) Create(_ context.Context, houseID, authorID int64, input domain.CreateChatMessageInput) (*domain.ChatMessage, error) {
	item := domain.ChatMessage{
		ID:         int64(len(r.items) + 1),
		HouseID:    houseID,
		AuthorID:   authorID,
		Content:    input.Content,
		ImageURL:   input.ImageURL,
		CreatedAt:  time.Now(),
		AuthorName: "Тестовый пользователь",
	}
	r.items = append(r.items, item)
	return &item, nil
}

func TestChatUseCaseCreateAllowsImageOnlyMessage(t *testing.T) {
	houseRepo := newFakeHouseRepo()
	houseRepo.memberships[[2]int64{5, 1}] = domain.Membership{
		UserID:   5,
		HouseID:  1,
		Role:     domain.RoleResident,
		JoinedAt: time.Now(),
	}
	chatRepo := &fakeChatRepo{}

	uc := NewChatUseCase(houseRepo, chatRepo)
	message, err := uc.Create(context.Background(), 5, 1, domain.CreateChatMessageInput{
		ImageURL: "/uploads/chat-image.jpg",
	})
	if err != nil {
		t.Fatalf("create chat message returned error: %v", err)
	}
	if message.ImageURL == "" {
		t.Fatalf("expected image url to be stored")
	}
}

func TestChatUseCaseRequiresHouseMembership(t *testing.T) {
	chatRepo := &fakeChatRepo{}
	uc := NewChatUseCase(newFakeHouseRepo(), chatRepo)

	_, err := uc.Create(context.Background(), 99, 1, domain.CreateChatMessageInput{
		Content: "Привет, соседи",
	})
	if err != domain.ErrNotFound {
		t.Fatalf("expected not found membership error, got %v", err)
	}
}
