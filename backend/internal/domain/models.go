package domain

import "time"

const (
	RoleResident = "resident"
	RoleAdmin    = "admin"
)

type User struct {
	ID          int64     `json:"id"`
	Login       string    `json:"login"`
	DisplayName string    `json:"displayName"`
	CreatedAt   time.Time `json:"createdAt"`
}

type House struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	CreatedBy int64     `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserHouse struct {
	House
	Role string `json:"role"`
}

type Membership struct {
	UserID   int64     `json:"userId"`
	HouseID  int64     `json:"houseId"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joinedAt"`
}

type Category struct {
	ID        int64     `json:"id"`
	HouseID   int64     `json:"houseId"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type Post struct {
	ID            int64     `json:"id"`
	HouseID       int64     `json:"houseId"`
	AuthorID      int64     `json:"authorId"`
	CategoryID    int64     `json:"categoryId"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	ImageURL      string    `json:"imageUrl,omitempty"`
	CommentsCount int       `json:"commentsCount"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	AuthorName    string    `json:"authorName"`
	CategoryName  string    `json:"categoryName"`
}

type Comment struct {
	ID         int64     `json:"id"`
	PostID     int64     `json:"postId"`
	AuthorID   int64     `json:"authorId"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"createdAt"`
	AuthorName string    `json:"authorName"`
}

type ChatMessage struct {
	ID         int64     `json:"id"`
	HouseID    int64     `json:"houseId"`
	AuthorID   int64     `json:"authorId"`
	Content    string    `json:"content"`
	ImageURL   string    `json:"imageUrl,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	AuthorName string    `json:"authorName"`
}

type InviteCode struct {
	ID             int64      `json:"id"`
	HouseID        int64      `json:"houseId"`
	Code           string     `json:"code"`
	CreatedBy      int64      `json:"createdBy"`
	IsActive       bool       `json:"isActive"`
	ExpiresAt      *time.Time `json:"expiresAt,omitempty"`
	CreatedAt      time.Time  `json:"createdAt"`
	CreatedByLogin string     `json:"createdByLogin"`
}

type MediaAsset struct {
	ID          int64     `json:"id"`
	PublicID    string    `json:"publicId"`
	ContentType string    `json:"contentType"`
	Data        []byte    `json:"-"`
	CreatedAt   time.Time `json:"createdAt"`
}

type RefreshSession struct {
	ID        int64
	UserID    int64
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type Tokens struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

type RegisterInput struct {
	Login       string `json:"login"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
}

type LoginInput struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type JoinHouseInput struct {
	Code string `json:"code"`
}

type CreateHouseInput struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type CreateCategoryInput struct {
	Name string `json:"name"`
}

type UpdateCategoryInput struct {
	Name string `json:"name"`
}

type CreatePostInput struct {
	CategoryID int64  `json:"categoryId"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	ImageURL   string `json:"imageUrl"`
}

type UpdatePostInput struct {
	CategoryID int64  `json:"categoryId"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	ImageURL   string `json:"imageUrl"`
}

type CreateCommentInput struct {
	Content string `json:"content"`
}

type CreateChatMessageInput struct {
	Content  string `json:"content"`
	ImageURL string `json:"imageUrl"`
}

type CreateInviteCodeInput struct {
	ExpiresAt *time.Time `json:"expiresAt"`
}

type ListCommentsFilter struct {
	Page     int
	PageSize int
}

type ListPostsFilter struct {
	CategoryID *int64
	Page       int
	PageSize   int
}
