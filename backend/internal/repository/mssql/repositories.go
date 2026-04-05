package mssql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"happyhouse/backend/internal/domain"
	"happyhouse/backend/pkg/pagination"
)

type UserRepository struct {
	db *sql.DB
}

type HouseRepository struct {
	db *sql.DB
}

type CategoryRepository struct {
	db *sql.DB
}

type PostRepository struct {
	db *sql.DB
}

type CommentRepository struct {
	db *sql.DB
}

type ChatRepository struct {
	db *sql.DB
}

type InviteCodeRepository struct {
	db *sql.DB
}

type RefreshTokenRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository             { return &UserRepository{db: db} }
func NewHouseRepository(db *sql.DB) *HouseRepository           { return &HouseRepository{db: db} }
func NewCategoryRepository(db *sql.DB) *CategoryRepository     { return &CategoryRepository{db: db} }
func NewPostRepository(db *sql.DB) *PostRepository             { return &PostRepository{db: db} }
func NewCommentRepository(db *sql.DB) *CommentRepository       { return &CommentRepository{db: db} }
func NewChatRepository(db *sql.DB) *ChatRepository             { return &ChatRepository{db: db} }
func NewInviteCodeRepository(db *sql.DB) *InviteCodeRepository { return &InviteCodeRepository{db: db} }
func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, input domain.RegisterInput, passwordHash string) (*domain.User, error) {
	query := `
		INSERT INTO users (login, password_hash, display_name)
		OUTPUT INSERTED.id, INSERTED.login, INSERTED.display_name, INSERTED.created_at
		VALUES (@p1, @p2, @p3);
	`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, input.Login, passwordHash, input.DisplayName).
		Scan(&user.ID, &user.Login, &user.DisplayName, &user.CreatedAt)
	if err != nil {
		return nil, normalizeError(err)
	}
	return user, nil
}

func (r *UserRepository) GetByLogin(ctx context.Context, login string) (*domain.User, string, error) {
	query := `
		SELECT id, login, display_name, created_at, password_hash
		FROM users
		WHERE login = @p1;
	`

	user := &domain.User{}
	var passwordHash string
	err := r.db.QueryRowContext(ctx, query, login).
		Scan(&user.ID, &user.Login, &user.DisplayName, &user.CreatedAt, &passwordHash)
	if err != nil {
		return nil, "", normalizeError(err)
	}
	return user, passwordHash, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `
		SELECT id, login, display_name, created_at
		FROM users
		WHERE id = @p1;
	`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&user.ID, &user.Login, &user.DisplayName, &user.CreatedAt)
	if err != nil {
		return nil, normalizeError(err)
	}
	return user, nil
}

func (r *HouseRepository) Create(ctx context.Context, createdBy int64, input domain.CreateHouseInput, defaultCategories []string) (*domain.House, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	query := `
		INSERT INTO houses (name, address, created_by)
		OUTPUT INSERTED.id, INSERTED.name, INSERTED.address, INSERTED.created_by, INSERTED.created_at
		VALUES (@p1, @p2, @p3);
	`

	house := &domain.House{}
	if err = tx.QueryRowContext(ctx, query, input.Name, input.Address, createdBy).
		Scan(&house.ID, &house.Name, &house.Address, &house.CreatedBy, &house.CreatedAt); err != nil {
		return nil, normalizeError(err)
	}

	if _, err = tx.ExecContext(ctx, `INSERT INTO house_members (user_id, house_id, role) VALUES (@p1, @p2, @p3);`, createdBy, house.ID, domain.RoleAdmin); err != nil {
		return nil, normalizeError(err)
	}

	for _, name := range defaultCategories {
		if _, err = tx.ExecContext(ctx, `INSERT INTO categories (house_id, name) VALUES (@p1, @p2);`, house.ID, name); err != nil {
			return nil, normalizeError(err)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return house, nil
}

func (r *HouseRepository) ListByUser(ctx context.Context, userID int64) ([]domain.UserHouse, error) {
	query := `
		SELECT h.id, h.name, h.address, h.created_by, h.created_at, hm.role
		FROM house_members hm
		JOIN houses h ON h.id = hm.house_id
		WHERE hm.user_id = @p1
		ORDER BY h.name;
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.UserHouse, 0)
	for rows.Next() {
		var house domain.UserHouse
		if err := rows.Scan(
			&house.ID,
			&house.Name,
			&house.Address,
			&house.CreatedBy,
			&house.CreatedAt,
			&house.Role,
		); err != nil {
			return nil, err
		}
		items = append(items, house)
	}
	return items, rows.Err()
}

func (r *HouseRepository) GetMembership(ctx context.Context, userID, houseID int64) (*domain.Membership, error) {
	query := `
		SELECT user_id, house_id, role, joined_at
		FROM house_members
		WHERE user_id = @p1 AND house_id = @p2;
	`

	membership := &domain.Membership{}
	err := r.db.QueryRowContext(ctx, query, userID, houseID).
		Scan(&membership.UserID, &membership.HouseID, &membership.Role, &membership.JoinedAt)
	if err != nil {
		return nil, normalizeError(err)
	}
	return membership, nil
}

func (r *HouseRepository) AddMembership(ctx context.Context, userID, houseID int64, role string) error {
	_, err := r.GetMembership(ctx, userID, houseID)
	if err == nil {
		return domain.ErrAlreadyMember
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return err
	}

	query := `INSERT INTO house_members (user_id, house_id, role) VALUES (@p1, @p2, @p3);`
	_, err = r.db.ExecContext(ctx, query, userID, houseID, role)
	return normalizeError(err)
}

func (r *CategoryRepository) ListByHouse(ctx context.Context, houseID int64) ([]domain.Category, error) {
	query := `
		SELECT id, house_id, name, created_at
		FROM categories
		WHERE house_id = @p1
		ORDER BY name;
	`

	rows, err := r.db.QueryContext(ctx, query, houseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.Category, 0)
	for rows.Next() {
		var category domain.Category
		if err := rows.Scan(&category.ID, &category.HouseID, &category.Name, &category.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, category)
	}
	return items, rows.Err()
}

func (r *CategoryRepository) Create(ctx context.Context, houseID int64, input domain.CreateCategoryInput) (*domain.Category, error) {
	query := `
		INSERT INTO categories (house_id, name)
		OUTPUT INSERTED.id, INSERTED.house_id, INSERTED.name, INSERTED.created_at
		VALUES (@p1, @p2);
	`

	category := &domain.Category{}
	err := r.db.QueryRowContext(ctx, query, houseID, input.Name).
		Scan(&category.ID, &category.HouseID, &category.Name, &category.CreatedAt)
	if err != nil {
		return nil, normalizeError(err)
	}
	return category, nil
}

func (r *CategoryRepository) Update(ctx context.Context, houseID, categoryID int64, input domain.UpdateCategoryInput) (*domain.Category, error) {
	query := `
		UPDATE categories
		SET name = @p3
		OUTPUT INSERTED.id, INSERTED.house_id, INSERTED.name, INSERTED.created_at
		WHERE house_id = @p1 AND id = @p2;
	`

	category := &domain.Category{}
	err := r.db.QueryRowContext(ctx, query, houseID, categoryID, input.Name).
		Scan(&category.ID, &category.HouseID, &category.Name, &category.CreatedAt)
	if err != nil {
		return nil, normalizeError(err)
	}
	return category, nil
}

func (r *CategoryRepository) Delete(ctx context.Context, houseID, categoryID int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM categories WHERE house_id = @p1 AND id = @p2;`, houseID, categoryID)
	if err != nil {
		return normalizeError(err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, houseID, categoryID int64) (*domain.Category, error) {
	query := `
		SELECT id, house_id, name, created_at
		FROM categories
		WHERE house_id = @p1 AND id = @p2;
	`

	category := &domain.Category{}
	err := r.db.QueryRowContext(ctx, query, houseID, categoryID).
		Scan(&category.ID, &category.HouseID, &category.Name, &category.CreatedAt)
	if err != nil {
		return nil, normalizeError(err)
	}
	return category, nil
}

func (r *PostRepository) ListByHouse(ctx context.Context, houseID int64, filter domain.ListPostsFilter) ([]domain.Post, int, error) {
	countQuery := `SELECT COUNT(1) FROM posts WHERE house_id = @p1`
	countArgs := []any{houseID}
	if filter.CategoryID != nil {
		countQuery += " AND category_id = @p2"
		countArgs = append(countArgs, *filter.CategoryID)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, normalizeError(err)
	}

	query := `
		SELECT p.id, p.house_id, p.author_id, p.category_id, p.title, p.content, p.image_url,
		       (SELECT COUNT(1) FROM comments cm WHERE cm.post_id = p.id) AS comments_count,
		       p.created_at, p.updated_at,
		       u.display_name, c.name
		FROM posts p
		JOIN users u ON u.id = p.author_id
		JOIN categories c ON c.id = p.category_id
		WHERE p.house_id = @p1`

	args := []any{houseID}
	if filter.CategoryID != nil {
		query += " AND p.category_id = @p2"
		args = append(args, *filter.CategoryID)
	}

	page, pageSize := pagination.Normalize(filter.Page, filter.PageSize)
	offset := pagination.Offset(page, pageSize)
	query += fmt.Sprintf(" ORDER BY p.created_at DESC OFFSET @p%d ROWS FETCH NEXT @p%d ROWS ONLY;", len(args)+1, len(args)+2)
	args = append(args, offset, pageSize)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, normalizeError(err)
	}
	defer rows.Close()

	items := make([]domain.Post, 0)
	for rows.Next() {
		post, err := scanPost(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *post)
	}
	return items, total, rows.Err()
}

func (r *PostRepository) GetByID(ctx context.Context, houseID, postID int64) (*domain.Post, error) {
	query := `
		SELECT p.id, p.house_id, p.author_id, p.category_id, p.title, p.content, p.image_url,
		       (SELECT COUNT(1) FROM comments cm WHERE cm.post_id = p.id) AS comments_count,
		       p.created_at, p.updated_at,
		       u.display_name, c.name
		FROM posts p
		JOIN users u ON u.id = p.author_id
		JOIN categories c ON c.id = p.category_id
		WHERE p.house_id = @p1 AND p.id = @p2;
	`

	row := r.db.QueryRowContext(ctx, query, houseID, postID)
	post, err := scanPost(row)
	if err != nil {
		return nil, normalizeError(err)
	}
	return post, nil
}

func scanPost(scanner interface {
	Scan(dest ...any) error
}) (*domain.Post, error) {
	var post domain.Post
	var imageURL sql.NullString
	if err := scanner.Scan(
		&post.ID,
		&post.HouseID,
		&post.AuthorID,
		&post.CategoryID,
		&post.Title,
		&post.Content,
		&imageURL,
		&post.CommentsCount,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.AuthorName,
		&post.CategoryName,
	); err != nil {
		return nil, err
	}
	post.ImageURL = imageURL.String
	return &post, nil
}

func (r *PostRepository) Create(ctx context.Context, houseID, authorID int64, input domain.CreatePostInput) (*domain.Post, error) {
	query := `
		INSERT INTO posts (house_id, author_id, category_id, title, content, image_url)
		OUTPUT INSERTED.id
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6);
	`

	var id int64
	if err := r.db.QueryRowContext(ctx, query, houseID, authorID, input.CategoryID, input.Title, input.Content, nullableString(input.ImageURL)).Scan(&id); err != nil {
		return nil, normalizeError(err)
	}
	return r.GetByID(ctx, houseID, id)
}

func (r *PostRepository) Update(ctx context.Context, houseID, postID int64, input domain.UpdatePostInput) (*domain.Post, error) {
	query := `
		UPDATE posts
		SET category_id = @p3, title = @p4, content = @p5, image_url = @p6, updated_at = SYSDATETIME()
		WHERE house_id = @p1 AND id = @p2;
	`
	result, err := r.db.ExecContext(ctx, query, houseID, postID, input.CategoryID, input.Title, input.Content, nullableString(input.ImageURL))
	if err != nil {
		return nil, normalizeError(err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, domain.ErrNotFound
	}
	return r.GetByID(ctx, houseID, postID)
}

func (r *PostRepository) Delete(ctx context.Context, houseID, postID int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM posts WHERE house_id = @p1 AND id = @p2;`, houseID, postID)
	if err != nil {
		return normalizeError(err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *CommentRepository) ListByPost(ctx context.Context, houseID, postID int64, filter domain.ListCommentsFilter) ([]domain.Comment, int, error) {
	countQuery := `
		SELECT COUNT(1)
		FROM comments cm
		JOIN posts p ON p.id = cm.post_id
		WHERE p.house_id = @p1 AND cm.post_id = @p2;
	`

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, houseID, postID).Scan(&total); err != nil {
		return nil, 0, normalizeError(err)
	}

	page, pageSize := pagination.Normalize(filter.Page, filter.PageSize)
	offset := pagination.Offset(page, pageSize)
	query := `
		SELECT cm.id, cm.post_id, cm.author_id, cm.content, cm.created_at, u.display_name
		FROM comments cm
		JOIN posts p ON p.id = cm.post_id
		JOIN users u ON u.id = cm.author_id
		WHERE p.house_id = @p1 AND cm.post_id = @p2
		ORDER BY cm.created_at ASC
		OFFSET @p3 ROWS FETCH NEXT @p4 ROWS ONLY;
	`
	rows, err := r.db.QueryContext(ctx, query, houseID, postID, offset, pageSize)
	if err != nil {
		return nil, 0, normalizeError(err)
	}
	defer rows.Close()

	items := make([]domain.Comment, 0)
	for rows.Next() {
		var comment domain.Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.AuthorID, &comment.Content, &comment.CreatedAt, &comment.AuthorName); err != nil {
			return nil, 0, err
		}
		items = append(items, comment)
	}
	return items, total, rows.Err()
}

func (r *CommentRepository) Create(ctx context.Context, houseID, postID, authorID int64, input domain.CreateCommentInput) (*domain.Comment, error) {
	query := `
		INSERT INTO comments (post_id, author_id, content)
		OUTPUT INSERTED.id
		VALUES (@p1, @p2, @p3);
	`
	var id int64
	if err := r.db.QueryRowContext(ctx, query, postID, authorID, input.Content).Scan(&id); err != nil {
		return nil, normalizeError(err)
	}

	row := `
		SELECT cm.id, cm.post_id, cm.author_id, cm.content, cm.created_at, u.display_name
		FROM comments cm
		JOIN posts p ON p.id = cm.post_id
		JOIN users u ON u.id = cm.author_id
		WHERE p.house_id = @p1 AND cm.id = @p2;
	`
	comment := &domain.Comment{}
	err := r.db.QueryRowContext(ctx, row, houseID, id).
		Scan(&comment.ID, &comment.PostID, &comment.AuthorID, &comment.Content, &comment.CreatedAt, &comment.AuthorName)
	if err != nil {
		return nil, normalizeError(err)
	}
	return comment, nil
}

func (r *ChatRepository) ListByHouse(ctx context.Context, houseID int64, limit int) ([]domain.ChatMessage, error) {
	query := `
		SELECT *
		FROM (
			SELECT TOP (@p2) cm.id, cm.house_id, cm.author_id, cm.content, cm.image_url, cm.created_at, u.display_name
			FROM chat_messages cm
			JOIN users u ON u.id = cm.author_id
			WHERE cm.house_id = @p1
			ORDER BY cm.created_at DESC
		) recent_messages
		ORDER BY recent_messages.created_at ASC;
	`

	rows, err := r.db.QueryContext(ctx, query, houseID, limit)
	if err != nil {
		return nil, normalizeError(err)
	}
	defer rows.Close()

	items := make([]domain.ChatMessage, 0)
	for rows.Next() {
		item, err := scanChatMessage(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}
	return items, rows.Err()
}

func (r *ChatRepository) Create(ctx context.Context, houseID, authorID int64, input domain.CreateChatMessageInput) (*domain.ChatMessage, error) {
	query := `
		INSERT INTO chat_messages (house_id, author_id, content, image_url)
		OUTPUT INSERTED.id
		VALUES (@p1, @p2, @p3, @p4);
	`

	var id int64
	if err := r.db.QueryRowContext(ctx, query, houseID, authorID, nullableString(input.Content), nullableString(input.ImageURL)).Scan(&id); err != nil {
		return nil, normalizeError(err)
	}

	row := `
		SELECT cm.id, cm.house_id, cm.author_id, cm.content, cm.image_url, cm.created_at, u.display_name
		FROM chat_messages cm
		JOIN users u ON u.id = cm.author_id
		WHERE cm.id = @p1;
	`

	item, err := scanChatMessage(r.db.QueryRowContext(ctx, row, id))
	if err != nil {
		return nil, normalizeError(err)
	}
	return item, nil
}

func scanChatMessage(scanner interface {
	Scan(dest ...any) error
}) (*domain.ChatMessage, error) {
	var item domain.ChatMessage
	var content sql.NullString
	var imageURL sql.NullString
	if err := scanner.Scan(
		&item.ID,
		&item.HouseID,
		&item.AuthorID,
		&content,
		&imageURL,
		&item.CreatedAt,
		&item.AuthorName,
	); err != nil {
		return nil, err
	}
	item.Content = content.String
	item.ImageURL = imageURL.String
	return &item, nil
}

func (r *InviteCodeRepository) ListByHouse(ctx context.Context, houseID int64) ([]domain.InviteCode, error) {
	query := `
		SELECT ic.id, ic.house_id, ic.code, ic.created_by, ic.is_active, ic.expires_at, ic.created_at, u.login
		FROM invite_codes ic
		JOIN users u ON u.id = ic.created_by
		WHERE ic.house_id = @p1
		ORDER BY ic.created_at DESC;
	`
	rows, err := r.db.QueryContext(ctx, query, houseID)
	if err != nil {
		return nil, normalizeError(err)
	}
	defer rows.Close()

	items := make([]domain.InviteCode, 0)
	for rows.Next() {
		var invite domain.InviteCode
		if err := rows.Scan(
			&invite.ID,
			&invite.HouseID,
			&invite.Code,
			&invite.CreatedBy,
			&invite.IsActive,
			&invite.ExpiresAt,
			&invite.CreatedAt,
			&invite.CreatedByLogin,
		); err != nil {
			return nil, err
		}
		items = append(items, invite)
	}
	return items, rows.Err()
}

func (r *InviteCodeRepository) Create(ctx context.Context, houseID, createdBy int64, input domain.CreateInviteCodeInput, code string) (*domain.InviteCode, error) {
	query := `
		INSERT INTO invite_codes (house_id, code, created_by, expires_at)
		OUTPUT INSERTED.id
		VALUES (@p1, @p2, @p3, @p4);
	`

	var id int64
	if err := r.db.QueryRowContext(ctx, query, houseID, code, createdBy, input.ExpiresAt).Scan(&id); err != nil {
		return nil, normalizeError(err)
	}

	row := `
		SELECT ic.id, ic.house_id, ic.code, ic.created_by, ic.is_active, ic.expires_at, ic.created_at, u.login
		FROM invite_codes ic
		JOIN users u ON u.id = ic.created_by
		WHERE ic.id = @p1;
	`
	invite := &domain.InviteCode{}
	err := r.db.QueryRowContext(ctx, row, id).Scan(
		&invite.ID,
		&invite.HouseID,
		&invite.Code,
		&invite.CreatedBy,
		&invite.IsActive,
		&invite.ExpiresAt,
		&invite.CreatedAt,
		&invite.CreatedByLogin,
	)
	if err != nil {
		return nil, normalizeError(err)
	}
	return invite, nil
}

func (r *InviteCodeRepository) GetActiveByCode(ctx context.Context, code string) (*domain.InviteCode, error) {
	query := `
		SELECT id, house_id, code, created_by, is_active, expires_at, created_at
		FROM invite_codes
		WHERE code = @p1 AND is_active = 1;
	`

	invite := &domain.InviteCode{}
	err := r.db.QueryRowContext(ctx, query, code).
		Scan(&invite.ID, &invite.HouseID, &invite.Code, &invite.CreatedBy, &invite.IsActive, &invite.ExpiresAt, &invite.CreatedAt)
	if err != nil {
		return nil, normalizeError(err)
	}
	return invite, nil
}

func (r *InviteCodeRepository) Deactivate(ctx context.Context, houseID, inviteCodeID int64) error {
	result, err := r.db.ExecContext(ctx, `UPDATE invite_codes SET is_active = 0 WHERE house_id = @p1 AND id = @p2;`, houseID, inviteCodeID)
	if err != nil {
		return normalizeError(err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *RefreshTokenRepository) Create(ctx context.Context, userID int64, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES (@p1, @p2, @p3);
	`
	_, err := r.db.ExecContext(ctx, query, userID, token, expiresAt)
	return normalizeError(err)
}

func (r *RefreshTokenRepository) Get(ctx context.Context, token string) (*domain.RefreshSession, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM refresh_tokens
		WHERE token = @p1;
	`
	session := &domain.RefreshSession{}
	err := r.db.QueryRowContext(ctx, query, token).
		Scan(&session.ID, &session.UserID, &session.Token, &session.ExpiresAt, &session.CreatedAt)
	if err != nil {
		return nil, normalizeError(err)
	}
	return session, nil
}

func (r *RefreshTokenRepository) Delete(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE token = @p1;`, token)
	return normalizeError(err)
}

func (r *RefreshTokenRepository) DeleteByUser(ctx context.Context, userID int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE user_id = @p1;`, userID)
	return normalizeError(err)
}

func normalizeError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrNotFound
	}
	message := strings.ToLower(err.Error())
	if strings.Contains(message, "unique") || strings.Contains(message, "duplicate") || strings.Contains(message, "primary key") {
		return domain.ErrConflict
	}
	return err
}

func nullableString(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}
