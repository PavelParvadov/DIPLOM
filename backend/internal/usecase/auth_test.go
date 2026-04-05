package usecase

import (
	"context"
	"testing"
	"time"

	"happyhouse/backend/internal/domain"
	"happyhouse/backend/pkg/auth"
)

type fakeUserRepo struct {
	createdHash  string
	usersByID    map[int64]*domain.User
	usersByLogin map[string]*domain.User
	hashes       map[string]string
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{
		usersByID:    make(map[int64]*domain.User),
		usersByLogin: make(map[string]*domain.User),
		hashes:       make(map[string]string),
	}
}

func (r *fakeUserRepo) Create(_ context.Context, input domain.RegisterInput, passwordHash string) (*domain.User, error) {
	r.createdHash = passwordHash
	user := &domain.User{
		ID:          int64(len(r.usersByID) + 1),
		Login:       input.Login,
		DisplayName: input.DisplayName,
		CreatedAt:   time.Now(),
	}
	r.usersByID[user.ID] = user
	r.usersByLogin[user.Login] = user
	r.hashes[user.Login] = passwordHash
	return user, nil
}

func (r *fakeUserRepo) GetByLogin(_ context.Context, login string) (*domain.User, string, error) {
	user, ok := r.usersByLogin[login]
	if !ok {
		return nil, "", domain.ErrNotFound
	}
	return user, r.hashes[login], nil
}

func (r *fakeUserRepo) GetByID(_ context.Context, id int64) (*domain.User, error) {
	user, ok := r.usersByID[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return user, nil
}

type fakeRefreshRepo struct {
	tokens map[string]domain.RefreshSession
}

func newFakeRefreshRepo() *fakeRefreshRepo {
	return &fakeRefreshRepo{tokens: make(map[string]domain.RefreshSession)}
}

func (r *fakeRefreshRepo) Create(_ context.Context, userID int64, token string, expiresAt time.Time) error {
	r.tokens[token] = domain.RefreshSession{UserID: userID, Token: token, ExpiresAt: expiresAt}
	return nil
}

func (r *fakeRefreshRepo) Get(_ context.Context, token string) (*domain.RefreshSession, error) {
	session, ok := r.tokens[token]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return &session, nil
}

func (r *fakeRefreshRepo) Delete(_ context.Context, token string) error {
	delete(r.tokens, token)
	return nil
}

func (r *fakeRefreshRepo) DeleteByUser(_ context.Context, userID int64) error {
	for token, session := range r.tokens {
		if session.UserID == userID {
			delete(r.tokens, token)
		}
	}
	return nil
}

func TestAuthUseCaseRegisterHashesPasswordAndIssuesTokens(t *testing.T) {
	userRepo := newFakeUserRepo()
	refreshRepo := newFakeRefreshRepo()
	uc := NewAuthUseCase(userRepo, refreshRepo, auth.NewTokenManager("secret"), time.Minute, time.Hour)

	user, tokens, err := uc.Register(context.Background(), domain.RegisterInput{
		Login:       "pavel_admin",
		Password:    "demo1234",
		DisplayName: "Павел",
	})
	if err != nil {
		t.Fatalf("register returned error: %v", err)
	}
	if user == nil || tokens == nil {
		t.Fatalf("expected user and tokens to be returned")
	}
	if userRepo.createdHash == "" || userRepo.createdHash == "demo1234" {
		t.Fatalf("expected password to be hashed, got %q", userRepo.createdHash)
	}
	if tokens.AccessToken == "" || tokens.RefreshToken == "" {
		t.Fatalf("expected both access and refresh tokens")
	}
}

func TestAuthUseCaseLoginValidatesCredentials(t *testing.T) {
	userRepo := newFakeUserRepo()
	refreshRepo := newFakeRefreshRepo()
	hash, err := auth.HashPassword("demo1234")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	userRepo.usersByID[1] = &domain.User{ID: 1, Login: "anna", DisplayName: "Анна", CreatedAt: time.Now()}
	userRepo.usersByLogin["anna"] = userRepo.usersByID[1]
	userRepo.hashes["anna"] = hash

	uc := NewAuthUseCase(userRepo, refreshRepo, auth.NewTokenManager("secret"), time.Minute, time.Hour)
	_, _, err = uc.Login(context.Background(), domain.LoginInput{
		Login:    "anna",
		Password: "wrong-password",
	})
	if err != domain.ErrInvalidCredentials {
		t.Fatalf("expected invalid credentials error, got %v", err)
	}
}
