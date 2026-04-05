package domain

import (
	"context"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, input RegisterInput, passwordHash string) (*User, error)
	GetByLogin(ctx context.Context, login string) (*User, string, error)
	GetByID(ctx context.Context, id int64) (*User, error)
}

type HouseRepository interface {
	Create(ctx context.Context, createdBy int64, input CreateHouseInput, defaultCategories []string) (*House, error)
	ListByUser(ctx context.Context, userID int64) ([]UserHouse, error)
	GetMembership(ctx context.Context, userID, houseID int64) (*Membership, error)
	AddMembership(ctx context.Context, userID, houseID int64, role string) error
}

type CategoryRepository interface {
	ListByHouse(ctx context.Context, houseID int64) ([]Category, error)
	Create(ctx context.Context, houseID int64, input CreateCategoryInput) (*Category, error)
	Update(ctx context.Context, houseID, categoryID int64, input UpdateCategoryInput) (*Category, error)
	Delete(ctx context.Context, houseID, categoryID int64) error
	GetByID(ctx context.Context, houseID, categoryID int64) (*Category, error)
}

type PostRepository interface {
	ListByHouse(ctx context.Context, houseID int64, filter ListPostsFilter) ([]Post, int, error)
	GetByID(ctx context.Context, houseID, postID int64) (*Post, error)
	Create(ctx context.Context, houseID, authorID int64, input CreatePostInput) (*Post, error)
	Update(ctx context.Context, houseID, postID int64, input UpdatePostInput) (*Post, error)
	Delete(ctx context.Context, houseID, postID int64) error
}

type CommentRepository interface {
	ListByPost(ctx context.Context, houseID, postID int64, filter ListCommentsFilter) ([]Comment, int, error)
	Create(ctx context.Context, houseID, postID, authorID int64, input CreateCommentInput) (*Comment, error)
}

type ChatRepository interface {
	ListByHouse(ctx context.Context, houseID int64, limit int) ([]ChatMessage, error)
	Create(ctx context.Context, houseID, authorID int64, input CreateChatMessageInput) (*ChatMessage, error)
}

type InviteCodeRepository interface {
	ListByHouse(ctx context.Context, houseID int64) ([]InviteCode, error)
	Create(ctx context.Context, houseID, createdBy int64, input CreateInviteCodeInput, code string) (*InviteCode, error)
	GetActiveByCode(ctx context.Context, code string) (*InviteCode, error)
	Deactivate(ctx context.Context, houseID, inviteCodeID int64) error
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, userID int64, token string, expiresAt time.Time) error
	Get(ctx context.Context, token string) (*RefreshSession, error)
	Delete(ctx context.Context, token string) error
	DeleteByUser(ctx context.Context, userID int64) error
}
