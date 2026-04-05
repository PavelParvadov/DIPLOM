package usecase

import (
	"context"

	"happyhouse/backend/internal/domain"
	"happyhouse/backend/pkg/pagination"
)

type CommentUseCase struct {
	houses   domain.HouseRepository
	posts    domain.PostRepository
	comments domain.CommentRepository
}

func NewCommentUseCase(houses domain.HouseRepository, posts domain.PostRepository, comments domain.CommentRepository) *CommentUseCase {
	return &CommentUseCase{houses: houses, posts: posts, comments: comments}
}

func (uc *CommentUseCase) List(ctx context.Context, userID, houseID, postID int64, filter domain.ListCommentsFilter) ([]domain.Comment, int, error) {
	if _, err := ensureMember(ctx, uc.houses, userID, houseID); err != nil {
		return nil, 0, err
	}
	if _, err := uc.posts.GetByID(ctx, houseID, postID); err != nil {
		return nil, 0, err
	}
	filter.Page, filter.PageSize = pagination.Normalize(filter.Page, filter.PageSize)
	return uc.comments.ListByPost(ctx, houseID, postID, filter)
}

func (uc *CommentUseCase) Create(ctx context.Context, userID, houseID, postID int64, input domain.CreateCommentInput) (*domain.Comment, error) {
	if !validateComment(input.Content) {
		return nil, domain.ErrValidation
	}
	if _, err := ensureMember(ctx, uc.houses, userID, houseID); err != nil {
		return nil, err
	}
	if _, err := uc.posts.GetByID(ctx, houseID, postID); err != nil {
		return nil, err
	}
	return uc.comments.Create(ctx, houseID, postID, userID, input)
}
