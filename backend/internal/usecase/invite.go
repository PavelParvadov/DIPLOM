package usecase

import (
	"context"

	"happyhouse/backend/internal/domain"
	"happyhouse/backend/pkg/auth"
)

type InviteCodeUseCase struct {
	houses  domain.HouseRepository
	invites domain.InviteCodeRepository
}

func NewInviteCodeUseCase(houses domain.HouseRepository, invites domain.InviteCodeRepository) *InviteCodeUseCase {
	return &InviteCodeUseCase{houses: houses, invites: invites}
}

func (uc *InviteCodeUseCase) List(ctx context.Context, userID, houseID int64) ([]domain.InviteCode, error) {
	if err := ensureAdmin(ctx, uc.houses, userID, houseID); err != nil {
		return nil, err
	}
	return uc.invites.ListByHouse(ctx, houseID)
}

func (uc *InviteCodeUseCase) Create(ctx context.Context, userID, houseID int64, input domain.CreateInviteCodeInput) (*domain.InviteCode, error) {
	if err := ensureAdmin(ctx, uc.houses, userID, houseID); err != nil {
		return nil, err
	}
	code, err := auth.NewRefreshToken()
	if err != nil {
		return nil, err
	}
	if len(code) > 10 {
		code = code[:10]
	}
	return uc.invites.Create(ctx, houseID, userID, input, code)
}

func (uc *InviteCodeUseCase) Deactivate(ctx context.Context, userID, houseID, inviteCodeID int64) error {
	if err := ensureAdmin(ctx, uc.houses, userID, houseID); err != nil {
		return err
	}
	return uc.invites.Deactivate(ctx, houseID, inviteCodeID)
}
