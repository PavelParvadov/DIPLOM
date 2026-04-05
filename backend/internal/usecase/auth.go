package usecase

import (
	"context"
	"time"

	"happyhouse/backend/internal/domain"
	"happyhouse/backend/pkg/auth"
)

type AuthUseCase struct {
	users      domain.UserRepository
	tokens     domain.RefreshTokenRepository
	tokenMaker *auth.TokenManager
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthUseCase(
	users domain.UserRepository,
	tokens domain.RefreshTokenRepository,
	tokenMaker *auth.TokenManager,
	accessTTL, refreshTTL time.Duration,
) *AuthUseCase {
	return &AuthUseCase{
		users:      users,
		tokens:     tokens,
		tokenMaker: tokenMaker,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (uc *AuthUseCase) Register(ctx context.Context, input domain.RegisterInput) (*domain.User, *domain.Tokens, error) {
	if !validateLogin(input.Login) || !validatePassword(input.Password) || !validateDisplayName(input.DisplayName) {
		return nil, nil, domain.ErrValidation
	}

	hash, err := auth.HashPassword(input.Password)
	if err != nil {
		return nil, nil, err
	}

	user, err := uc.users.Create(ctx, input, hash)
	if err != nil {
		return nil, nil, err
	}

	tokens, err := uc.issueTokens(ctx, user.ID, user.Login)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, input domain.LoginInput) (*domain.User, *domain.Tokens, error) {
	user, passwordHash, err := uc.users.GetByLogin(ctx, input.Login)
	if err != nil {
		return nil, nil, err
	}
	if err := auth.CheckPassword(passwordHash, input.Password); err != nil {
		return nil, nil, domain.ErrInvalidCredentials
	}

	tokens, err := uc.issueTokens(ctx, user.ID, user.Login)
	if err != nil {
		return nil, nil, err
	}
	return user, tokens, nil
}

func (uc *AuthUseCase) GetCurrentUser(ctx context.Context, userID int64) (*domain.User, error) {
	return uc.users.GetByID(ctx, userID)
}

func (uc *AuthUseCase) Refresh(ctx context.Context, refreshToken string) (*domain.Tokens, error) {
	session, err := uc.tokens.Get(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	if session.ExpiresAt.Before(time.Now()) {
		_ = uc.tokens.Delete(ctx, refreshToken)
		return nil, domain.ErrUnauthorized
	}

	user, err := uc.users.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	_ = uc.tokens.Delete(ctx, refreshToken)
	return uc.issueTokens(ctx, user.ID, user.Login)
}

func (uc *AuthUseCase) Logout(ctx context.Context, refreshToken string) error {
	return uc.tokens.Delete(ctx, refreshToken)
}

func (uc *AuthUseCase) issueTokens(ctx context.Context, userID int64, login string) (*domain.Tokens, error) {
	accessToken, expiresAt, err := uc.tokenMaker.NewAccessToken(userID, login, uc.accessTTL)
	if err != nil {
		return nil, err
	}

	refreshToken, err := auth.NewRefreshToken()
	if err != nil {
		return nil, err
	}

	if err := uc.tokens.Create(ctx, userID, refreshToken, time.Now().Add(uc.refreshTTL)); err != nil {
		return nil, err
	}

	return &domain.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}
