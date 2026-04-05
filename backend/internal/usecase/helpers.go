package usecase

import (
	"context"
	"strings"

	"happyhouse/backend/internal/domain"
)

func validateDisplayName(value string) bool {
	return len(strings.TrimSpace(value)) >= 2
}

func validateLogin(value string) bool {
	value = strings.TrimSpace(value)
	return len(value) >= 3 && len(value) <= 32
}

func validatePassword(value string) bool {
	return len(value) >= 6
}

func validateCategoryName(value string) bool {
	return len(strings.TrimSpace(value)) >= 2
}

func validateHouse(name, address string) bool {
	return len(strings.TrimSpace(name)) >= 3 && len(strings.TrimSpace(address)) >= 5
}

func validatePost(title, content string) bool {
	return len(strings.TrimSpace(title)) >= 3 && len(strings.TrimSpace(content)) >= 5
}

func validateComment(content string) bool {
	return len(strings.TrimSpace(content)) >= 1
}

func validateChatMessage(content, imageURL string) bool {
	return len(strings.TrimSpace(content)) >= 1 || strings.TrimSpace(imageURL) != ""
}

func ensureMember(ctx context.Context, houses domain.HouseRepository, userID, houseID int64) (*domain.Membership, error) {
	membership, err := houses.GetMembership(ctx, userID, houseID)
	if err != nil {
		return nil, err
	}
	return membership, nil
}

func ensureAdmin(ctx context.Context, houses domain.HouseRepository, userID, houseID int64) error {
	membership, err := ensureMember(ctx, houses, userID, houseID)
	if err != nil {
		return err
	}
	if membership.Role != domain.RoleAdmin {
		return domain.ErrForbidden
	}
	return nil
}
