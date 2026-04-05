package usecase

import (
	"context"
	"strings"

	"happyhouse/backend/internal/domain"
	"happyhouse/backend/pkg/pagination"
)

type PostUseCase struct {
	houses     domain.HouseRepository
	posts      domain.PostRepository
	categories domain.CategoryRepository
}

func NewPostUseCase(houses domain.HouseRepository, posts domain.PostRepository, categories domain.CategoryRepository) *PostUseCase {
	return &PostUseCase{
		houses:     houses,
		posts:      posts,
		categories: categories,
	}
}

func (uc *PostUseCase) List(ctx context.Context, userID, houseID int64, filter domain.ListPostsFilter) ([]domain.Post, int, error) {
	if _, err := ensureMember(ctx, uc.houses, userID, houseID); err != nil {
		return nil, 0, err
	}
	filter.Page, filter.PageSize = pagination.Normalize(filter.Page, filter.PageSize)
	return uc.posts.ListByHouse(ctx, houseID, filter)
}

func (uc *PostUseCase) Get(ctx context.Context, userID, houseID, postID int64) (*domain.Post, error) {
	if _, err := ensureMember(ctx, uc.houses, userID, houseID); err != nil {
		return nil, err
	}
	return uc.posts.GetByID(ctx, houseID, postID)
}

func (uc *PostUseCase) Create(ctx context.Context, userID, houseID int64, input domain.CreatePostInput) (*domain.Post, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Content = strings.TrimSpace(input.Content)
	if !validatePost(input.Title, input.Content) {
		return nil, domain.ErrValidation
	}
	if _, err := ensureMember(ctx, uc.houses, userID, houseID); err != nil {
		return nil, err
	}
	if _, err := uc.categories.GetByID(ctx, houseID, input.CategoryID); err != nil {
		return nil, err
	}
	return uc.posts.Create(ctx, houseID, userID, input)
}

func (uc *PostUseCase) Update(ctx context.Context, userID, houseID, postID int64, input domain.UpdatePostInput) (*domain.Post, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Content = strings.TrimSpace(input.Content)
	if !validatePost(input.Title, input.Content) {
		return nil, domain.ErrValidation
	}
	membership, err := ensureMember(ctx, uc.houses, userID, houseID)
	if err != nil {
		return nil, err
	}

	post, err := uc.posts.GetByID(ctx, houseID, postID)
	if err != nil {
		return nil, err
	}
	if membership.Role != domain.RoleAdmin && post.AuthorID != userID {
		return nil, domain.ErrForbidden
	}
	if _, err := uc.categories.GetByID(ctx, houseID, input.CategoryID); err != nil {
		return nil, err
	}
	return uc.posts.Update(ctx, houseID, postID, input)
}

func (uc *PostUseCase) Delete(ctx context.Context, userID, houseID, postID int64) error {
	membership, err := ensureMember(ctx, uc.houses, userID, houseID)
	if err != nil {
		return err
	}

	post, err := uc.posts.GetByID(ctx, houseID, postID)
	if err != nil {
		return err
	}
	if membership.Role != domain.RoleAdmin && post.AuthorID != userID {
		return domain.ErrForbidden
	}
	return uc.posts.Delete(ctx, houseID, postID)
}
