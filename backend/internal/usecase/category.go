package usecase

import (
	"context"

	"happyhouse/backend/internal/domain"
)

type CategoryUseCase struct {
	houses     domain.HouseRepository
	categories domain.CategoryRepository
}

func NewCategoryUseCase(houses domain.HouseRepository, categories domain.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{houses: houses, categories: categories}
}

func (uc *CategoryUseCase) List(ctx context.Context, userID, houseID int64) ([]domain.Category, error) {
	if _, err := ensureMember(ctx, uc.houses, userID, houseID); err != nil {
		return nil, err
	}
	return uc.categories.ListByHouse(ctx, houseID)
}

func (uc *CategoryUseCase) Create(ctx context.Context, userID, houseID int64, input domain.CreateCategoryInput) (*domain.Category, error) {
	if !validateCategoryName(input.Name) {
		return nil, domain.ErrValidation
	}
	if err := ensureAdmin(ctx, uc.houses, userID, houseID); err != nil {
		return nil, err
	}
	return uc.categories.Create(ctx, houseID, input)
}

func (uc *CategoryUseCase) Update(ctx context.Context, userID, houseID, categoryID int64, input domain.UpdateCategoryInput) (*domain.Category, error) {
	if !validateCategoryName(input.Name) {
		return nil, domain.ErrValidation
	}
	if err := ensureAdmin(ctx, uc.houses, userID, houseID); err != nil {
		return nil, err
	}
	return uc.categories.Update(ctx, houseID, categoryID, input)
}

func (uc *CategoryUseCase) Delete(ctx context.Context, userID, houseID, categoryID int64) error {
	if err := ensureAdmin(ctx, uc.houses, userID, houseID); err != nil {
		return err
	}
	return uc.categories.Delete(ctx, houseID, categoryID)
}
