package usecase

import (
	"context"
	"strings"

	"happyhouse/backend/internal/domain"
)

var defaultHouseCategories = []string{"Новости", "Объявления", "Соседи"}

type HouseUseCase struct {
	houses  domain.HouseRepository
	invites domain.InviteCodeRepository
}

func NewHouseUseCase(houses domain.HouseRepository, invites domain.InviteCodeRepository) *HouseUseCase {
	return &HouseUseCase{houses: houses, invites: invites}
}

func (uc *HouseUseCase) List(ctx context.Context, userID int64) ([]domain.UserHouse, error) {
	return uc.houses.ListByUser(ctx, userID)
}

func (uc *HouseUseCase) Create(ctx context.Context, userID int64, input domain.CreateHouseInput) (*domain.House, error) {
	if !validateHouse(input.Name, input.Address) {
		return nil, domain.ErrValidation
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Address = strings.TrimSpace(input.Address)

	return uc.houses.Create(ctx, userID, input, defaultHouseCategories)
}

func (uc *HouseUseCase) JoinByCode(ctx context.Context, userID int64, input domain.JoinHouseInput) error {
	input.Code = strings.TrimSpace(input.Code)
	invite, err := uc.invites.GetActiveByCode(ctx, input.Code)
	if err != nil {
		return err
	}

	if invite.ExpiresAt != nil && invite.ExpiresAt.Before(now()) {
		return domain.ErrExpiredInviteCode
	}

	return uc.houses.AddMembership(ctx, userID, invite.HouseID, domain.RoleResident)
}
