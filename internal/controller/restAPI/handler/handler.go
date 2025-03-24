package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"users-app/internal/entity"
	"users-app/pkg/logger"

	"github.com/gofrs/uuid/v5"
)

//go:generate go run go.uber.org/mock/mockgen@latest -source=handler.go -destination=../../../mocks/handler.go -package=mocks -typed
type UserService interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (entity.User, error)
	CreateUser(ctx context.Context, user entity.User) error
	UpdateUser(ctx context.Context, user entity.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type Handler struct {
	log         logger.Logger
	userService UserService
}

func New(log logger.Logger, userService UserService) *Handler {
	return &Handler{
		log:         log,
		userService: userService,
	}
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := r.URL.Query().Get("id")
	if id == "" {
		h.sendErr(w, http.StatusBadRequest, errors.New("id is empty"), "id is empty")
		return
	}

	userID, err := uuid.FromString(id)
	if err != nil {
		h.sendErr(w, http.StatusBadRequest, err, "невалидный id пользователя: "+id)
		return
	}

	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			h.sendErr(w, http.StatusNotFound, err, "user not found")
		}

		h.sendErr(w, http.StatusInternalServerError, err, "failed to get user")
		return
	}

	h.sendJSON(w, http.StatusOK, user)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user entity.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.sendErr(w, http.StatusBadRequest, err, "failed to decode request body")
		return
	}

	if err := h.userService.CreateUser(ctx, user); err != nil {
		if errors.Is(err, entity.ErrAlreadyExists) {
			h.sendErr(w, http.StatusConflict, err, "user with email "+user.Email+" already exists")
			return
		}

		h.sendErr(w, http.StatusInternalServerError, err, "failed to create user")
		return
	}

	h.sendJSON(w, http.StatusCreated, user)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user entity.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.sendErr(w, http.StatusBadRequest, err, "failed to decode request body")
		return
	}

	if err := h.userService.UpdateUser(ctx, user); err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			h.sendErr(w, http.StatusNotFound, err, "user not found")
			return
		}

		h.sendErr(w, http.StatusInternalServerError, err, "failed to update user")
		return
	}

	h.sendJSON(w, http.StatusOK, "user updated")
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := r.URL.Query().Get("id")
	if id == "" {
		h.sendErr(w, http.StatusBadRequest, errors.New("id is empty"), "id is empty")
		return
	}

	userID, err := uuid.FromString(id)
	if err != nil {
		h.sendErr(w, http.StatusBadRequest, err, "невалидный id пользователя: "+id)
		return
	}

	if err := h.userService.DeleteUser(ctx, userID); err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			h.sendErr(w, http.StatusNotFound, err, "user not found")
			return
		}

		h.sendErr(w, http.StatusInternalServerError, err, "failed to delete user")
		return
	}

	h.sendJSON(w, http.StatusOK, "user deleted")
}
