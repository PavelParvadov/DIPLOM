package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"happyhouse/backend/internal/domain"
	transportdto "happyhouse/backend/internal/transport/http/dto"
	authmw "happyhouse/backend/internal/transport/http/middleware"
	"happyhouse/backend/internal/usecase"
	"happyhouse/backend/pkg/auth"
	"happyhouse/backend/pkg/httpx"
)

type Handler struct {
	auth       *usecase.AuthUseCase
	houses     *usecase.HouseUseCase
	categories *usecase.CategoryUseCase
	posts      *usecase.PostUseCase
	comments   *usecase.CommentUseCase
	chats      *usecase.ChatUseCase
	invites    *usecase.InviteCodeUseCase
	media      *usecase.MediaUseCase
	tokenMgr   *auth.TokenManager
	origin     string
	uploadDir  string
}

func New(
	authUC *usecase.AuthUseCase,
	houseUC *usecase.HouseUseCase,
	categoryUC *usecase.CategoryUseCase,
	postUC *usecase.PostUseCase,
	commentUC *usecase.CommentUseCase,
	chatUC *usecase.ChatUseCase,
	inviteUC *usecase.InviteCodeUseCase,
	mediaUC *usecase.MediaUseCase,
	tokenMgr *auth.TokenManager,
	frontendOrigin string,
	uploadDir string,
) *Handler {
	return &Handler{
		auth:       authUC,
		houses:     houseUC,
		categories: categoryUC,
		posts:      postUC,
		comments:   commentUC,
		chats:      chatUC,
		invites:    inviteUC,
		media:      mediaUC,
		tokenMgr:   tokenMgr,
		origin:     frontendOrigin,
		uploadDir:  uploadDir,
	}
}

func (h *Handler) Router() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{h.origin},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"status": "ok"})
	})
	router.Get("/uploads/*", h.handleUploadedAsset)

	router.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", h.handleRegister)
		r.Post("/auth/login", h.handleLogin)
		r.Post("/auth/refresh", h.handleRefresh)
		r.Post("/auth/logout", h.handleLogout)

		r.Group(func(private chi.Router) {
			private.Use(authmw.Authenticate(h.tokenMgr))
			private.Get("/me", h.handleCurrentUser)
			private.Get("/houses", h.handleListHouses)
			private.Post("/houses", h.handleCreateHouse)
			private.Post("/houses/join-by-code", h.handleJoinHouseByCode)

			private.Route("/houses/{houseId}", func(house chi.Router) {
				house.Get("/categories", h.handleListCategories)
				house.Post("/categories", h.handleCreateCategory)
				house.Patch("/categories/{categoryId}", h.handleUpdateCategory)
				house.Delete("/categories/{categoryId}", h.handleDeleteCategory)

				house.Get("/posts", h.handleListPosts)
				house.Post("/posts", h.handleCreatePost)
				house.Get("/posts/{postId}", h.handleGetPost)
				house.Patch("/posts/{postId}", h.handleUpdatePost)
				house.Delete("/posts/{postId}", h.handleDeletePost)

				house.Get("/posts/{postId}/comments", h.handleListComments)
				house.Post("/posts/{postId}/comments", h.handleCreateComment)

				house.Get("/chat/messages", h.handleListChatMessages)
				house.Post("/chat/messages", h.handleCreateChatMessage)

				house.Get("/invite-codes", h.handleListInviteCodes)
				house.Post("/invite-codes", h.handleCreateInviteCode)
				house.Patch("/invite-codes/{inviteCodeId}/deactivate", h.handleDeactivateInviteCode)
			})
		})
	})

	return router
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var input domain.RegisterInput
	if err := httpx.DecodeJSON(r, &input); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	user, tokens, err := h.auth.Register(r.Context(), input)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, transportdto.AuthResponse{User: user, Tokens: tokens})
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var input domain.LoginInput
	if err := httpx.DecodeJSON(r, &input); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	user, tokens, err := h.auth.Login(r.Context(), input)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, transportdto.AuthResponse{User: user, Tokens: tokens})
}

func (h *Handler) handleRefresh(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := httpx.DecodeJSON(r, &payload); err != nil || payload.RefreshToken == "" {
		httpx.Error(w, http.StatusBadRequest, "refreshToken is required")
		return
	}
	tokens, err := h.auth.Refresh(r.Context(), payload.RefreshToken)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"tokens": tokens})
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := httpx.DecodeJSON(r, &payload); err != nil || payload.RefreshToken == "" {
		httpx.Error(w, http.StatusBadRequest, "refreshToken is required")
		return
	}
	if err := h.auth.Logout(r.Context(), payload.RefreshToken); err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"success": true})
}

func (h *Handler) handleCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, err := authmw.UserIDFromContext(r.Context())
	if err != nil {
		h.writeError(w, err)
		return
	}
	user, err := h.auth.GetCurrentUser(r.Context(), userID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"user": user})
}

func (h *Handler) handleListHouses(w http.ResponseWriter, r *http.Request) {
	userID, err := authmw.UserIDFromContext(r.Context())
	if err != nil {
		h.writeError(w, err)
		return
	}
	houses, err := h.houses.List(r.Context(), userID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"items": houses})
}

func (h *Handler) handleCreateHouse(w http.ResponseWriter, r *http.Request) {
	userID, err := authmw.UserIDFromContext(r.Context())
	if err != nil {
		h.writeError(w, err)
		return
	}
	var input domain.CreateHouseInput
	if err := httpx.DecodeJSON(r, &input); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	house, err := h.houses.Create(r.Context(), userID, input)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, httpx.Envelope{"item": house})
}

func (h *Handler) handleJoinHouseByCode(w http.ResponseWriter, r *http.Request) {
	userID, err := authmw.UserIDFromContext(r.Context())
	if err != nil {
		h.writeError(w, err)
		return
	}
	var input domain.JoinHouseInput
	if err := httpx.DecodeJSON(r, &input); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.houses.JoinByCode(r.Context(), userID, input); err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"success": true})
}

func (h *Handler) handleListCategories(w http.ResponseWriter, r *http.Request) {
	userID, houseID, ok := h.userAndHouse(w, r)
	if !ok {
		return
	}
	items, err := h.categories.List(r.Context(), userID, houseID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"items": items})
}

func (h *Handler) handleCreateCategory(w http.ResponseWriter, r *http.Request) {
	userID, houseID, ok := h.userAndHouse(w, r)
	if !ok {
		return
	}
	var input domain.CreateCategoryInput
	if err := httpx.DecodeJSON(r, &input); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	category, err := h.categories.Create(r.Context(), userID, houseID, input)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, httpx.Envelope{"item": category})
}

func (h *Handler) handleUpdateCategory(w http.ResponseWriter, r *http.Request) {
	userID, houseID, ok := h.userAndHouse(w, r)
	if !ok {
		return
	}
	categoryID, err := paramInt64(r, "categoryId")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid categoryId")
		return
	}
	var input domain.UpdateCategoryInput
	if err := httpx.DecodeJSON(r, &input); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	category, err := h.categories.Update(r.Context(), userID, houseID, categoryID, input)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"item": category})
}

func (h *Handler) handleDeleteCategory(w http.ResponseWriter, r *http.Request) {
	userID, houseID, ok := h.userAndHouse(w, r)
	if !ok {
		return
	}
	categoryID, err := paramInt64(r, "categoryId")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid categoryId")
		return
	}
	if err := h.categories.Delete(r.Context(), userID, houseID, categoryID); err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"success": true})
}

func (h *Handler) handleListPosts(w http.ResponseWriter, r *http.Request) {
	userID, houseID, ok := h.userAndHouse(w, r)
	if !ok {
		return
	}

	page := queryInt(r, "page", 1)
	pageSize := queryInt(r, "pageSize", 10)
	filter := domain.ListPostsFilter{Page: page, PageSize: pageSize}
	if categoryID := r.URL.Query().Get("categoryId"); categoryID != "" {
		value, err := strconv.ParseInt(categoryID, 10, 64)
		if err != nil {
			httpx.Error(w, http.StatusBadRequest, "invalid categoryId")
			return
		}
		filter.CategoryID = &value
	}

	items, total, err := h.posts.List(r.Context(), userID, houseID, filter)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, transportdto.ListPostsResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func (h *Handler) handleGetPost(w http.ResponseWriter, r *http.Request) {
	userID, houseID, ok := h.userAndHouse(w, r)
	if !ok {
		return
	}
	postID, err := paramInt64(r, "postId")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid postId")
		return
	}
	post, err := h.posts.Get(r.Context(), userID, houseID, postID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"item": post})
}

func (h *Handler) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	userID, houseID, ok := h.userAndHouse(w, r)
	if !ok {
		return
	}
	input, err := h.decodeCreatePostInput(r)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	post, err := h.posts.Create(r.Context(), userID, houseID, input)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, httpx.Envelope{"item": post})
}

func (h *Handler) handleUpdatePost(w http.ResponseWriter, r *http.Request) {
	userID, houseID, ok := h.userAndHouse(w, r)
	if !ok {
		return
	}
	postID, err := paramInt64(r, "postId")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid postId")
		return
	}
	var input domain.UpdatePostInput
	if err := httpx.DecodeJSON(r, &input); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	post, err := h.posts.Update(r.Context(), userID, houseID, postID, input)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"item": post})
}

func (h *Handler) handleDeletePost(w http.ResponseWriter, r *http.Request) {
	userID, houseID, ok := h.userAndHouse(w, r)
	if !ok {
		return
	}
	postID, err := paramInt64(r, "postId")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid postId")
		return
	}
	if err := h.posts.Delete(r.Context(), userID, houseID, postID); err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"success": true})
}

func (h *Handler) handleListComments(w http.ResponseWriter, r *http.Request) {
	userID, houseID, ok := h.userAndHouse(w, r)
	if !ok {
		return
	}
	postID, err := paramInt64(r, "postId")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid postId")
		return
	}
	page := queryInt(r, "page", 1)
	pageSize := queryInt(r, "pageSize", 4)
	items, total, err := h.comments.List(r.Context(), userID, houseID, postID, domain.ListCommentsFilter{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, transportdto.ListCommentsResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func (h *Handler) handleCreateComment(w http.ResponseWriter, r *http.Request) {
	userID, houseID, ok := h.userAndHouse(w, r)
	if !ok {
		return
	}
	postID, err := paramInt64(r, "postId")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid postId")
		return
	}
	var input domain.CreateCommentInput
	if err := httpx.DecodeJSON(r, &input); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	comment, err := h.comments.Create(r.Context(), userID, houseID, postID, input)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, httpx.Envelope{"item": comment})
}

func (h *Handler) handleListChatMessages(w http.ResponseWriter, r *http.Request) {
	userID, houseID, ok := h.userAndHouse(w, r)
	if !ok {
		return
	}
	items, err := h.chats.List(r.Context(), userID, houseID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"items": items})
}

func (h *Handler) handleCreateChatMessage(w http.ResponseWriter, r *http.Request) {
	userID, houseID, ok := h.userAndHouse(w, r)
	if !ok {
		return
	}
	input, err := h.decodeCreateChatMessageInput(r)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.chats.Create(r.Context(), userID, houseID, input)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, httpx.Envelope{"item": item})
}

func (h *Handler) handleListInviteCodes(w http.ResponseWriter, r *http.Request) {
	userID, houseID, ok := h.userAndHouse(w, r)
	if !ok {
		return
	}
	items, err := h.invites.List(r.Context(), userID, houseID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"items": items})
}

func (h *Handler) handleCreateInviteCode(w http.ResponseWriter, r *http.Request) {
	userID, houseID, ok := h.userAndHouse(w, r)
	if !ok {
		return
	}
	var input transportdto.CreateInviteCodeRequest
	if err := httpx.DecodeJSON(r, &input); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	item, err := h.invites.Create(r.Context(), userID, houseID, domain.CreateInviteCodeInput{ExpiresAt: input.ExpiresAt})
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, httpx.Envelope{"item": item})
}

func (h *Handler) handleDeactivateInviteCode(w http.ResponseWriter, r *http.Request) {
	userID, houseID, ok := h.userAndHouse(w, r)
	if !ok {
		return
	}
	inviteCodeID, err := paramInt64(r, "inviteCodeId")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid inviteCodeId")
		return
	}
	if err := h.invites.Deactivate(r.Context(), userID, houseID, inviteCodeID); err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"success": true})
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		httpx.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrInvalidCredentials), errors.Is(err, domain.ErrUnauthorized):
		httpx.Error(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, domain.ErrForbidden):
		httpx.Error(w, http.StatusForbidden, err.Error())
	case errors.Is(err, domain.ErrNotFound):
		httpx.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrAlreadyMember):
		httpx.Error(w, http.StatusConflict, "Вы уже состоите в этом доме")
	case errors.Is(err, domain.ErrConflict):
		httpx.Error(w, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrExpiredInviteCode):
		httpx.Error(w, http.StatusBadRequest, err.Error())
	default:
		httpx.Error(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *Handler) userAndHouse(w http.ResponseWriter, r *http.Request) (int64, int64, bool) {
	userID, err := authmw.UserIDFromContext(r.Context())
	if err != nil {
		h.writeError(w, err)
		return 0, 0, false
	}
	houseID, err := paramInt64(r, "houseId")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid houseId")
		return 0, 0, false
	}
	return userID, houseID, true
}

func paramInt64(r *http.Request, key string) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, key), 10, 64)
}

func queryInt(r *http.Request, key string, fallback int) int {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}
