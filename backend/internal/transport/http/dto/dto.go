package dto

import (
	"time"

	"happyhouse/backend/internal/domain"
)

type AuthResponse struct {
	User   *domain.User   `json:"user"`
	Tokens *domain.Tokens `json:"tokens"`
}

type ListPostsResponse struct {
	Items    []domain.Post `json:"items"`
	Total    int           `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"pageSize"`
}

type ListCommentsResponse struct {
	Items    []domain.Comment `json:"items"`
	Total    int              `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"pageSize"`
}

type CreateInviteCodeRequest struct {
	ExpiresAt *time.Time `json:"expiresAt"`
}
