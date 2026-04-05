package middleware

import (
	"context"
	"net/http"
	"strings"

	"happyhouse/backend/internal/domain"
	"happyhouse/backend/pkg/auth"
	"happyhouse/backend/pkg/httpx"
)

type contextKey string

const userIDContextKey contextKey = "userID"

func Authenticate(tokenManager *auth.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				httpx.Error(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				httpx.Error(w, http.StatusUnauthorized, "invalid authorization header")
				return
			}

			claims, err := tokenManager.ParseAccessToken(parts[1])
			if err != nil {
				httpx.Error(w, http.StatusUnauthorized, "invalid access token")
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(userIDContextKey).(int64)
	if !ok || userID == 0 {
		return 0, domain.ErrUnauthorized
	}
	return userID, nil
}
